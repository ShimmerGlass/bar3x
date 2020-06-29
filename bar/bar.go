package bar

import (
	"fmt"
	"image"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/shimmerglass/bar3x/ui/module"
	"github.com/shimmerglass/bar3x/x"
)

type Bar struct {
	Win xproto.Window
	Buf *xgraphics.Image

	lock sync.Mutex

	screen     x.Screen
	ctx        ui.Context
	background *ui.Root

	TrayWidth int
	padding   int

	lastLeft   *image.RGBA
	lastCenter *image.RGBA
	lastRight  *image.RGBA
}

func NewBar(ctx ui.Context, X *xgbutil.XUtil, screen x.Screen) (*Bar, error) {
	w, h := screen.Width, ctx.MustInt("height")

	win, err := xwindow.Generate(X)
	if err != nil {
		return nil, err
	}

	win.Create(X.RootWin(), screen.X, screen.Y, w, h, 0)

	// Make this window close gracefully.
	win.WMGracefulClose(func(w *xwindow.Window) {
		xevent.Detach(w.X, w.Id)
		keybind.Detach(w.X, w.Id)
		mousebind.Detach(w.X, w.Id)
		w.Destroy()
		xevent.Quit(w.X)
	})

	// Set WM_STATE so it is interpreted as a top-level window.
	err = icccm.WmStateSet(X, win.Id, &icccm.WmState{
		State: icccm.StateNormal,
	})
	if err != nil { // not a fatal error
		xgbutil.Logger.Printf("Could not set WM_STATE: %s", err)
	}

	// Set WM_NORMAL_HINTS so the window can't be resized.
	err = icccm.WmNormalHintsSet(X, win.Id, &icccm.NormalHints{
		Flags:     icccm.SizeHintPMinSize | icccm.SizeHintPMaxSize,
		MinWidth:  uint(w),
		MinHeight: uint(h),
		MaxWidth:  uint(w),
		MaxHeight: uint(h),
	})
	if err != nil { // not a fatal error
		xgbutil.Logger.Printf("Could not set WM_NORMAL_HINTS: %s", err)
	}

	// Set _NET_WM_NAME so it looks nice.
	err = ewmh.WmNameSet(X, win.Id, "bar3x")
	if err != nil { // not a fatal error
		xgbutil.Logger.Printf("Could not set _NET_WM_NAME: %s", err)
	}

	// Now we can map, since we've set all our properties.
	// (The initial map is when the window manager starts managing.)
	err = ewmh.WmWindowTypeSet(X, win.Id, []string{"_NET_WM_WINDOW_TYPE_DOCK"})
	if err != nil {
		return nil, err
	}

	strutAtom, err := xprop.Atom(X, "_NET_WM_STRUT_PARTIAL", false)
	if err != nil {
		return nil, err
	}

	strutVals := make([]byte, 4*12)
	xgb.Put32(strutVals[8:], uint32(h))
	xgb.Put32(strutVals[32:], uint32(screen.X))
	xgb.Put32(strutVals[36:], uint32(screen.X+screen.Width))

	xproto.ChangeProperty(X.Conn(),
		xproto.PropModeReplace,
		win.Id,
		strutAtom,
		xproto.AtomCardinal,
		32, 12,
		strutVals,
	)

	ximg := xgraphics.New(X, image.Rect(0, 0, w, h))

	ximg.XSurfaceSet(win.Id)
	ximg.XDraw()
	ximg.XPaint(win.Id)

	win.Map()

	b := &Bar{
		Win:     win.Id,
		Buf:     ximg,
		padding: ctx.MustInt("h_padding"),
		ctx:     ctx,
		screen:  screen,
	}
	err = b.initBackground()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bar) initBackground() error {
	clock := module.NewClock(func() {})
	mk := markup.New()
	base.RegisterMarkup(mk)
	module.RegisterMarkup(mk, clock)

	root := ui.NewRoot(b.ctx.New(ui.Context{"bar_width": b.screen.Width}), b.updateAll)
	el, err := mk.Parse(root, nil, b.ctx.MustString("bar_background"))
	if err != nil {
		return fmt.Errorf("config: bar_background: %w", err)
	}

	root.Inner = el
	b.background = root

	root.Notify()
	go clock.Run()

	return nil
}

func (b *Bar) updateAll() {
	b.background.Paint()
	DrawCopySrcRGBAToBGRA(b.Buf, b.Buf.Rect, b.background.Image(), image.ZP)
	b.Buf.XDraw()
	if b.lastLeft != nil {
		b.PaintLeft(b.lastLeft)
	}
	if b.lastCenter != nil {
		b.PaintCenter(b.lastCenter)
	}
	if b.lastRight != nil {
		b.PaintRight(b.lastRight)
	}
}

func (b *Bar) PaintLeft(im *image.RGBA) {
	maxW := im.Rect.Dx()
	if b.lastLeft != nil && b.lastLeft.Rect.Dx() > maxW {
		maxW = b.lastLeft.Rect.Dx()
	}
	maxH := im.Rect.Dy()
	if b.lastLeft != nil && b.lastLeft.Rect.Dy() > maxH {
		maxH = b.lastLeft.Rect.Dy()
	}

	source := image.Pt(b.padding, (b.Buf.Rect.Dy()-maxH)/2)
	maxBounds := image.Rect(0, 0, maxW, maxH).Add(source)

	DrawCopySrcRGBAToBGRA(
		b.Buf,
		maxBounds,
		b.background.Image(),
		source,
	)

	DrawCopyOverRGBAToBGRA(
		b.Buf,
		im.Rect.Add(image.Pt(
			b.padding,
			(b.Buf.Rect.Dy()-im.Rect.Dy())/2,
		)),
		im,
		image.ZP,
	)
	b.lastLeft = im

	b.paintPart(maxBounds)
}

func (b *Bar) PaintCenter(im *image.RGBA) {
	maxW := im.Rect.Dx()
	if b.lastCenter != nil && b.lastCenter.Rect.Dx() > maxW {
		maxW = b.lastCenter.Rect.Dx()
	}
	maxH := im.Rect.Dy()
	if b.lastCenter != nil && b.lastCenter.Rect.Dy() > maxH {
		maxH = b.lastCenter.Rect.Dy()
	}

	source := image.Pt((b.Buf.Rect.Dx()-maxW)/2, (b.Buf.Rect.Dy()-maxH)/2)
	maxBounds := image.Rect(0, 0, maxW, maxH).Add(source)

	DrawCopySrcRGBAToBGRA(
		b.Buf,
		maxBounds,
		b.background.Image(),
		source,
	)

	DrawCopyOverRGBAToBGRA(
		b.Buf,
		im.Rect.Add(image.Pt(
			(b.Buf.Rect.Dx()-im.Rect.Dx())/2,
			(b.Buf.Rect.Dy()-im.Rect.Dy())/2,
		)),
		im,
		image.ZP,
	)
	b.lastCenter = im

	b.paintPart(maxBounds)
}

func (b *Bar) PaintRight(im *image.RGBA) {
	maxW := im.Rect.Dx()
	if b.lastRight != nil && b.lastRight.Rect.Dx() > maxW {
		maxW = b.lastRight.Rect.Dx()
	}
	maxH := im.Rect.Dy()
	if b.lastRight != nil && b.lastRight.Rect.Dy() > maxH {
		maxH = b.lastRight.Rect.Dy()
	}

	source := image.Pt(b.Buf.Rect.Dx()-maxW-b.padding-b.TrayWidth, (b.Buf.Rect.Dy()-maxH)/2)
	maxBounds := image.Rect(0, 0, maxW, maxH).Add(source)

	DrawCopySrcRGBAToBGRA(
		b.Buf,
		maxBounds,
		b.background.Image(),
		source,
	)

	DrawCopyOverRGBAToBGRA(
		b.Buf,
		im.Rect.Add(image.Pt(
			b.Buf.Rect.Dx()-im.Rect.Dx()-b.padding-b.TrayWidth,
			(b.Buf.Rect.Dy()-im.Rect.Dy())/2,
		)),
		im,
		image.ZP,
	)
	b.lastRight = im

	b.paintPart(maxBounds)
}

func (b *Bar) paintPart(bounds image.Rectangle) {
	b.lock.Lock()
	b.Buf.SubImage(bounds).(*xgraphics.Image).XDraw()
	b.Buf.XPaint(b.Win)
	b.lock.Unlock()
}

func (b *Bar) SetTrayWidth(w int) {
	b.TrayWidth = w
	b.updateAll()
}
