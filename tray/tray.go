package tray

import (
	"fmt"
	"image/color"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
)

type Client struct {
	Class     string
	Win       xproto.Window
	Mapped    bool
	ShouldMap bool
}

type Tray struct {
	x        *xgbutil.XUtil
	barWin   xproto.Window
	iconSize int

	barHeight int
	barWidth  int

	clients map[xproto.Window]*Client

	updateCb func(State)
}

func New(x *xgbutil.XUtil, barWin xproto.Window, updateCb func(State)) *Tray {
	return &Tray{
		updateCb:  updateCb,
		clients:   map[xproto.Window]*Client{},
		x:         x,
		barWin:    barWin,
		iconSize:  20,
		barHeight: 30,
		barWidth:  1920,
	}
}

func (t *Tray) Init(background color.Color) error {
	// xproto.Setup retrieves the Setup information from the setup bytes
	// gathered during connection.
	setup := xproto.Setup(t.x.Conn())

	// This is the default screen with all its associated info.
	screen := setup.DefaultScreen(t.x.Conn())

	win, err := createTrayWindow(t.x, screen, background)
	if err != nil {
		return err
	}

	trayAtom, err := xprop.Atom(t.x, fmt.Sprintf("_NET_SYSTEM_TRAY_S%d", 0), false)
	if err != nil {
		return err
	}

	managerAtom, err := xprop.Atom(t.x, "MANAGER", false)
	if err != nil {
		return err
	}

	xproto.SetSelectionOwner(t.x.Conn(), win, trayAtom, xproto.TimeCurrentTime)
	cookie := xproto.GetSelectionOwner(t.x.Conn(), trayAtom)
	reply, err := cookie.Reply()
	if err != nil {
		return err
	}

	if reply.Owner == win {
	} else {
		log.Warnf("no selection, owned by %v, disabling tray", reply.Owner)
		return nil
	}

	xevent.ClientMessageFun(func(X *xgbutil.XUtil, ev xevent.ClientMessageEvent) {
		atmName, err := xprop.AtomName(X, ev.Type)
		if err != nil {
			log.Errorf("tray: %s", err)
			return
		}

		switch atmName {
		case "_NET_SYSTEM_TRAY_OPCODE":
			err := t.handleTrayMessage(t.x, ev, t.barWin)
			if err != nil {
				log.Errorf("tray: %s", err)
			}

		}
	}).Connect(t.x, win)

	ev := xproto.ClientMessageEvent{
		Window: screen.Root,
		Type:   managerAtom,
		Format: 32,
		Data: xproto.ClientMessageDataUnionData32New([]uint32{
			xproto.TimeCurrentTime,
			uint32(trayAtom),
			uint32(win),
			0, 0,
		}),
	}

	err = xproto.SendEventChecked(t.x.Conn(), false, screen.Root, 0xFFFFFF, string(ev.Bytes())).Check()
	if err != nil {
		return err
	}

	return nil
}

func (t *Tray) handleTrayMessage(x *xgbutil.XUtil, ev xevent.ClientMessageEvent, barWin xproto.Window) error {
	op := ev.Data.Data32[1]
	switch op {
	case 0: // dock
		clientWin := xproto.Window(ev.Data.Data32[2])

		xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
			log.Debugf("map event for %d", ev.Window)
			client, ok := t.clients[ev.Window]
			if !ok {
				log.Warnf("map event for unknown client %d", ev.Window)
				return
			}
			client.Mapped = true
			client.ShouldMap = true
			t.configureClients()
		}).Connect(t.x, clientWin)

		xevent.UnmapNotifyFun(func(X *xgbutil.XUtil, ev xevent.UnmapNotifyEvent) {
			log.Debugf("unmap event for %d", ev.Window)
			client, ok := t.clients[ev.Window]
			if !ok {
				log.Warnf("unmap event for unknown client %d", ev.Window)
				return
			}
			client.Mapped = false
			client.ShouldMap = false
			t.configureClients()
		}).Connect(t.x, clientWin)

		xevent.PropertyNotifyFun(func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
			log.Debugf("unmap event for %d", ev.Window)
			client, ok := t.clients[ev.Window]
			if !ok {
				log.Warnf("unmap event for unknown client %d", ev.Window)
				return
			}
			_, shouldMap, err := xembedData(x, clientWin)
			if err != nil {
				log.Errorf("tray: %s", err)
				return
			}

			client.ShouldMap = shouldMap
			t.configureClients()
			t.configureClientsMap()
		}).Connect(t.x, clientWin)

		xevent.DestroyNotifyFun(func(X *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
			log.Debugf("destroy event for %d", ev.Window)
			_, ok := t.clients[ev.Window]
			if !ok {
				log.Warnf("destroy event for unknown client %d", ev.Window)
				return
			}
			delete(t.clients, ev.Window)
			t.configureClients()
			xevent.Detach(X, ev.Window)
		}).Connect(t.x, clientWin)

		xproto.ChangeWindowAttributes(x.Conn(),
			clientWin,
			xproto.CwEventMask,
			[]uint32{
				xproto.EventMaskPropertyChange | xproto.EventMaskStructureNotify | xproto.EventMaskResizeRedirect,
			},
		)

		xembedVersion, shouldMap, err := xembedData(x, clientWin)
		if err != nil {
			return err
		}

		xproto.ReparentWindow(x.Conn(),
			clientWin,
			barWin,
			0, 0,
		)

		cwCookie := xproto.ConfigureWindowChecked(x.Conn(),
			clientWin,
			xproto.ConfigWindowWidth|xproto.ConfigWindowHeight,
			[]uint32{uint32(t.iconSize), uint32(t.iconSize)},
		)
		err = cwCookie.Check()
		if err != nil {
			return err
		}

		embedAtom, err := xprop.Atom(x, "_XEMBED", false)
		if err != nil {
			return err
		}

		ev := xproto.ClientMessageEvent{
			Window: clientWin,
			Type:   embedAtom,
			Format: 32,
			Data: xproto.ClientMessageDataUnionData32New([]uint32{
				xproto.TimeCurrentTime,
				0, // XEMBED_EMBEDDED_NOTIFY
				uint32(barWin),
				uint32(xembedVersion),
				0, // unused
			}),
		}

		err = xproto.SendEventChecked(x.Conn(), false, clientWin, xproto.EventMaskNoEvent, string(ev.Bytes())).Check()
		if err != nil {
			return err
		}

		classAtom, err := xprop.Atom(x, "WM_CLASS", false)
		if err != nil {
			return err
		}

		classc := xproto.GetProperty(
			x.Conn(),
			false,
			clientWin,
			classAtom,
			xproto.AtomString,
			0,
			32,
		)

		classRes, err := classc.Reply()
		if err != nil {
			return err
		}

		xproto.ChangeSaveSet(x.Conn(), xproto.SetModeInsert, clientWin)
		t.clients[clientWin] = &Client{
			Class:     string(classRes.Value),
			Win:       clientWin,
			Mapped:    false,
			ShouldMap: shouldMap,
		}

		t.configureClients()
		t.configureClientsMap()
	}

	return nil
}

func (t *Tray) configureClients() error {
	padding := (t.barHeight - t.iconSize) / 2

	clients := make([]*Client, 0, len(t.clients))
	for _, c := range t.clients {
		clients = append(clients, c)
	}

	sort.Slice(clients, func(i, j int) bool {
		return strings.Compare(clients[i].Class, clients[j].Class) == -1
	})

	size := padding
	for _, c := range clients {
		if !c.ShouldMap {
			continue
		}
		size += t.iconSize + 2
	}

	t.updateCb(State{Width: size})

	off := padding
	for _, c := range clients {
		if !c.ShouldMap {
			continue
		}
		err := xproto.ConfigureWindowChecked(t.x.Conn(),
			c.Win,
			xproto.ConfigWindowX|xproto.ConfigWindowY,
			[]uint32{
				uint32(t.barWidth - off - t.iconSize),
				uint32(padding),
			},
		).Check()
		if err != nil {
			return err
		}
		off += t.iconSize + 2
	}

	return nil
}

func (t *Tray) configureClientsMap() error {
	for _, c := range t.clients {
		if c.Mapped && !c.ShouldMap {
			xproto.UnmapWindow(t.x.Conn(), c.Win)
			c.Mapped = false
			c.ShouldMap = false
		}
		if !c.Mapped && c.ShouldMap {
			xproto.MapWindow(t.x.Conn(), c.Win)
			c.Mapped = true
			c.ShouldMap = true
		}
	}

	return nil
}
