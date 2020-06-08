package tray

import (
	"image/color"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
)

func createTrayWindow(x *xgbutil.XUtil, screen *xproto.ScreenInfo, bgColor color.Color) (xproto.Window, error) {
	wid, _ := xproto.NewWindowId(x.Conn())

	selMask := xproto.CwBackPixel | xproto.CwBorderPixel | xproto.CwOverrideRedirect | xproto.CwColormap
	selVal := []uint32{
		screen.BlackPixel,
		screen.BlackPixel,
		1,
		uint32(screen.DefaultColormap),
	}

	// CreateWindow takes a boatload of parameters.
	xproto.CreateWindow(x.Conn(), screen.RootDepth, wid, screen.Root,
		-1, -1, 1, 1, 0,
		xproto.WindowClassInputOutput, screen.RootVisual, uint32(selMask), selVal)

	trayOrAtom, err := xprop.Atom(x, "_NET_SYSTEM_TRAY_ORIENTATION", false)
	if err != nil {
		return 0, err
	}

	payload := make([]byte, 4)
	xgb.Put32(payload, 0) // Horizontal
	xproto.ChangeProperty(x.Conn(),
		xproto.PropModeReplace,
		wid,
		trayOrAtom,
		xproto.AtomCardinal,
		32,
		1,
		payload,
	)

	trayVizAtom, err := xprop.Atom(x, "_NET_SYSTEM_TRAY_VISUAL", false)
	if err != nil {
		return 0, err
	}

	payload = make([]byte, 4)
	xgb.Put32(payload, uint32(screen.RootVisual))
	xproto.ChangeProperty(x.Conn(),
		xproto.PropModeReplace,
		wid,
		trayVizAtom,
		xproto.AtomVisualid,
		32,
		1,
		payload,
	)

	trayColorsAtom, err := xprop.Atom(x, "_NET_SYSTEM_TRAY_COLORS", false)
	if err != nil {
		return 0, err
	}

	r, g, b, _ := bgColor.RGBA()
	data := make([]byte, 3*4*4)
	off := 0
	for i := 0; i < 4; i++ {
		xgb.Put32(data[off:], uint32(r))
		off += 4
		xgb.Put32(data[off:], uint32(g))
		off += 4
		xgb.Put32(data[off:], uint32(b))
		off += 4
	}
	xproto.ChangeProperty(x.Conn(),
		xproto.PropModeReplace,
		wid,
		trayColorsAtom,
		xproto.AtomCardinal,
		32, 12,
		data,
	)

	return wid, nil
}
