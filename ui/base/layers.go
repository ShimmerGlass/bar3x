package base

import (
	"image"
	"image/draw"

	"github.com/shimmerglass/bar3x/ui"
)

var _ ui.ParentDrawable = (*Layers)(nil)

type Layers struct {
	Base
	inner  []ui.Drawable
	hAlign string
	vAlign string
}

func NewLayers(p ui.ParentDrawable) *Layers {
	return &Layers{
		Base:   NewBase(p),
		hAlign: HAlighCenter,
		vAlign: VAlighMiddle,
	}
}

func (l *Layers) Add(d ui.Drawable) {
	l.inner = append(l.inner, d)
	d.OnWidthChange(func(int) {
		l.updateSize()
	})
	d.OnHeightChange(func(int) {
		l.updateSize()
	})
	d.OnVisibleChange(func(bool) {
		l.updateSize()
	})
}

func (l *Layers) SetContext(ctx ui.Context) {
	l.Base.SetContext(ctx)
	for _, i := range l.inner {
		i.SetContext(ctx)
	}
}

func (l *Layers) HAlign() string {
	return l.hAlign
}
func (l *Layers) SetHAlign(v string) {
	l.hAlign = v
}

func (l *Layers) VAlign() string {
	return l.vAlign
}
func (l *Layers) SetVAlign(v string) {
	l.vAlign = v
}

func (c *Layers) ChildContext(i int) ui.Context {
	return c.ctx
}

func (c *Layers) Children() []ui.Drawable {
	return c.inner
}

func (c *Layers) SendEvent(ev ui.Event) bool {
	if !c.Base.SendEvent(ev) {
		return false
	}

	for _, i := range c.inner {
		if !i.Visible() {
			continue
		}
		w, h := i.Width(), i.Height()
		x := (c.width.V - w) / 2
		y := (c.height.V - h) / 2

		if ev.At.In(image.Rect(x, y, x+w, y+h)) {
			iev := ev
			iev.At = ev.At.Sub(image.Pt(x, y))
			i.SendEvent(iev)
		}

		y += h
	}

	return true
}

func (l *Layers) updateSize() {
	var w, h int
	for _, i := range l.inner {
		if !i.Visible() {
			continue
		}
		ew, eh := i.Width(), i.Height()
		if ew > w {
			w = ew
		}
		if eh > h {
			h = eh
		}
	}
	l.width.Set(w)
	l.height.Set(h)
}

func (l *Layers) Draw(x, y int, im draw.Image) {
	for _, i := range l.inner {
		if !i.Visible() {
			continue
		}
		var xOff, yOff int
		switch l.hAlign {
		case HAlignLeft:
			xOff = 0
		case HAlighCenter:
			xOff = (l.width.V - i.Width()) / 2
		case HAlignRight:
			xOff = l.width.V - i.Width()
		}
		switch l.vAlign {
		case VAlignTop:
			yOff = 0
		case VAlighMiddle:
			yOff = (l.height.V - i.Height()) / 2
		case VAlignBottom:
			yOff = l.height.V - i.Height()
		}
		i.Draw(x+xOff, y+yOff, im)
	}
}
