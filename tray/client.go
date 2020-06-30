package tray

import (
	"fmt"
	"sort"
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/shimmerglass/bar3x/x"
	log "github.com/sirupsen/logrus"
)

func (t *Tray) addClient(clientWin xproto.Window) error {
	client := &Client{
		Win:    xwindow.New(t.x, clientWin),
		Mapped: false,
	}

	class, err := x.WindowClass(t.x, clientWin)
	if err != nil {
		return err
	}
	client.Class = class

	t.clientLock.Lock()
	t.clients[clientWin] = client
	t.clientLock.Unlock()

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		log.Debugf("map event for %d", ev.Window)
		client, ok := t.clients[ev.Window]
		if !ok {
			log.Warnf("map event for unknown client %d", ev.Window)
			return
		}
		client.Mapped = true
		client.ShouldMap = true
		t.configureClientsPosition()
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
		t.configureClientsPosition()
	}).Connect(t.x, clientWin)

	xevent.PropertyNotifyFun(func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		log.Debugf("unmap event for %d", ev.Window)
		client, ok := t.clients[ev.Window]
		if !ok {
			log.Warnf("unmap event for unknown client %d", ev.Window)
			return
		}
		_, shouldMap, err := xembedData(t.x, ev.Window)
		if err == nil {
			client.ShouldMap = shouldMap
			t.configureClientsPosition()
			t.configureClientsMap()
		}
	}).Connect(t.x, clientWin)

	xevent.DestroyNotifyFun(func(X *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		log.Debugf("destroy event for %d", ev.Window)
		_, ok := t.clients[ev.Window]
		if !ok {
			log.Warnf("destroy event for unknown client %d", ev.Window)
			return
		}

		t.removeClient(ev.Window)
	}).Connect(t.x, clientWin)

	err = client.Win.Listen(
		xproto.EventMaskPropertyChange,
		xproto.EventMaskStructureNotify,
		xproto.EventMaskResizeRedirect,
	)
	if err != nil {
		return err
	}

	xembedVersion, shouldMap, err := xembedData(t.x, clientWin)
	if err != nil {
		return fmt.Errorf("xembed get: %w", err)
	}
	client.ShouldMap = shouldMap

	xproto.ReparentWindow(t.x.Conn(),
		clientWin,
		t.barWin,
		0, 0,
	)

	client.Win.Resize(t.cfg.IconSize, t.cfg.IconSize)

	ev := xproto.ClientMessageEvent{
		Window: clientWin,
		Type:   x.MustAtom(t.x, "_XEMBED"),
		Format: 32,
		Data: xproto.ClientMessageDataUnionData32New([]uint32{
			xproto.TimeCurrentTime,
			0, // XEMBED_EMBEDDED_NOTIFY
			uint32(t.barWin),
			uint32(xembedVersion),
			0, // unused
		}),
	}

	xproto.SendEvent(t.x.Conn(), false, clientWin, xproto.EventMaskNoEvent, string(ev.Bytes()))

	xproto.ChangeSaveSet(t.x.Conn(), xproto.SetModeInsert, clientWin)

	t.configureClientsPosition()
	t.configureClientsMap()

	return nil
}

func (t *Tray) removeClient(win xproto.Window) {
	t.clientLock.Lock()
	delete(t.clients, win)
	t.clientLock.Unlock()

	t.configureClientsPosition()
	xevent.Detach(t.x, win)
}

func (t *Tray) configureClientsPosition() error {
	padding := (t.cfg.BarHeight - t.cfg.IconSize) / 2

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
		size += t.cfg.IconSize + t.cfg.IconPadding
	}

	// we notify the underlying bar before moving the windows
	// this is to ensure new windows are drawn on a clean surface
	// overwise the window's background is poluted with the bar's
	// graphics
	t.updateCb(State{Width: size})

	off := t.cfg.BarWidth
	for _, c := range clients {
		if !c.ShouldMap {
			continue
		}
		off -= t.cfg.IconSize + t.cfg.IconPadding
		c.Win.Move(off, padding)
	}

	return nil
}

func (t *Tray) configureClientsMap() error {
	for _, c := range t.clients {
		if c.Mapped && !c.ShouldMap {
			c.Win.Unmap()
			c.Mapped = false
			c.ShouldMap = false
		}
		if !c.Mapped && c.ShouldMap {
			c.Win.Map()
			c.Mapped = true
			c.ShouldMap = true
		}
	}

	return nil
}
