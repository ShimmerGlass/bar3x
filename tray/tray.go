package tray

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/shimmerglass/bar3x/x"
	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Config struct {
	IconPadding     int
	IconSize        int
	BarWidth        int
	BarHeight       int
	BackgroundColor color.Color
}

type Client struct {
	Class     string
	Win       *xwindow.Window
	Mapped    bool
	ShouldMap bool
}

type Tray struct {
	x      *xgbutil.XUtil
	barWin xproto.Window
	cfg    Config

	clients    map[xproto.Window]*Client
	clientLock sync.Mutex

	updateCb func(State)
}

func New(x *xgbutil.XUtil, barWin xproto.Window, cfg Config, updateCb func(State)) *Tray {
	return &Tray{
		updateCb: updateCb,
		clients:  map[xproto.Window]*Client{},
		x:        x,
		barWin:   barWin,
		cfg:      cfg,
	}
}

func (t *Tray) Init() error {
	// xproto.Setup retrieves the Setup information from the setup bytes
	// gathered during connection.
	setup := xproto.Setup(t.x.Conn())

	// This is the default screen with all its associated info.
	screen := setup.DefaultScreen(t.x.Conn())

	win, err := t.createTrayWindow(screen, t.cfg.BackgroundColor)
	if err != nil {
		return err
	}

	trayAtomName := fmt.Sprintf("_NET_SYSTEM_TRAY_S%d", 0)
	err = xproto.SetSelectionOwnerChecked(
		t.x.Conn(),
		win,
		x.MustAtom(t.x, trayAtomName),
		xproto.TimeCurrentTime,
	).Check()
	if err != nil {
		return err
	}

	reply, err := xproto.GetSelectionOwner(t.x.Conn(), x.MustAtom(t.x, trayAtomName)).Reply()
	if err != nil {
		return err
	}

	if reply.Owner != win {
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
		Type:   x.MustAtom(t.x, "MANAGER"),
		Format: 32,
		Data: xproto.ClientMessageDataUnionData32New([]uint32{
			xproto.TimeCurrentTime,
			uint32(x.MustAtom(t.x, trayAtomName)),
			uint32(win),
			0, 0,
		}),
	}

	// notify clients to dock with us
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
		return t.addClient(clientWin)
	}

	return nil
}
