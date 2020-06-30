package tray

import (
	"image/color"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/shimmerglass/bar3x/x"
)

func (t *Tray) createTrayWindow(screen *xproto.ScreenInfo, bgColor color.Color) (xproto.Window, error) {
	wid, _ := xproto.NewWindowId(t.x.Conn())

	selMask := xproto.CwBackPixel | xproto.CwBorderPixel | xproto.CwOverrideRedirect | xproto.CwColormap
	selVal := []uint32{
		screen.BlackPixel,
		screen.BlackPixel,
		1,
		uint32(screen.DefaultColormap),
	}

	xproto.CreateWindow(
		t.x.Conn(),
		screen.RootDepth,
		wid,
		screen.Root,
		-1, -1, 1, 1, 0,
		xproto.WindowClassInputOutput,
		screen.RootVisual,
		uint32(selMask),
		selVal,
	)

	payload := make([]byte, 4)
	xgb.Put32(payload, 0) // Horizontal
	xproto.ChangeProperty(t.x.Conn(),
		xproto.PropModeReplace,
		wid,
		x.MustAtom(t.x, "_NET_SYSTEM_TRAY_ORIENTATION"),
		xproto.AtomCardinal,
		32,
		1,
		payload,
	)

	payload = make([]byte, 4)
	xgb.Put32(payload, uint32(screen.RootVisual))
	xproto.ChangeProperty(t.x.Conn(),
		xproto.PropModeReplace,
		wid,
		x.MustAtom(t.x, "_NET_SYSTEM_TRAY_VISUAL"),
		xproto.AtomVisualid,
		32,
		1,
		payload,
	)

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
	xproto.ChangeProperty(t.x.Conn(),
		xproto.PropModeReplace,
		wid,
		x.MustAtom(t.x, "_NET_SYSTEM_TRAY_COLORS"),
		xproto.AtomCardinal,
		32, 12,
		data,
	)

	return wid, nil
}
