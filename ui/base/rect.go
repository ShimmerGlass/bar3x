package base

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/fogleman/gg"
	"github.com/shimmerglass/bar3x/ui"
)

var _ ui.ParentDrawable = (*Rect)(nil)

type Rect struct {
	Base

	inner ui.Drawable

	color     color.Color
	setWidth  int
	setHeight int
	radius    int
}

func NewRect(p ui.ParentDrawable) *Rect {
	return &Rect{
		Base:      NewBase(p),
		setWidth:  -1,
		setHeight: -1,
	}
}

func (s *Rect) Add(d ui.Drawable) {
	s.inner = d
	d.OnWidthChange(func(int) {
		s.updateSize()
	})
	d.OnHeightChange(func(int) {
		s.updateSize()
	})
	d.OnVisibleChange(func(bool) {
		s.updateSize()
	})
}

func (c *Rect) ChildContext(i int) ui.Context {
	return c.ctx
}
func (r *Rect) Children() []ui.Drawable {
	if r.inner == nil {
		return nil
	}
	return []ui.Drawable{r.inner}
}
func (c *Rect) SetContext(ctx ui.Context) {
	c.Base.SetContext(ctx)
	if c.inner != nil {
		c.inner.SetContext(ctx)
	}
}

func (b *Rect) SetWidth(v int) {
	b.setWidth = v
	b.updateSize()
}
func (b *Rect) SetHeight(v int) {
	b.setHeight = v
	b.updateSize()
}

func (b *Rect) Color() color.Color {
	return b.color
}
func (b *Rect) SetColor(v color.Color) {
	b.color = v
}

func (b *Rect) Radius() int {
	return b.radius
}
func (b *Rect) SetRadius(v int) {
	b.radius = v
}

func (s *Rect) SendEvent(ev ui.Event) bool {
	if !s.Base.SendEvent(ev) {
		return false
	}

	if s.inner != nil {
		if !s.inner.Visible() {
			return true
		}

		s.inner.SendEvent(ev)
	}

	return true
}

func (r *Rect) updateSize() {
	w, h := r.setWidth, r.setHeight
	if w == -1 {
		if r.inner != nil {
			w = r.inner.Width()
		} else {
			w = 0
		}
	}
	if h == -1 {
		if r.inner != nil {
			h = r.inner.Height()
		} else {
			h = 0
		}
	}

	r.width.Set(w)
	r.height.Set(h)
}

func (r *Rect) Draw(x, y int, im draw.Image) {
	w, h := r.Width(), r.Height()
	dc := gg.NewContext(w, h)
	dc.DrawRoundedRectangle(
		0,
		0,
		float64(w),
		float64(h),
		float64(r.radius),
	)

	dc.SetColor(r.color)
	dc.Fill()

	draw.Draw(
		im,
		image.Rect(x, y, x+w, y+h),
		dc.Image(),
		image.Pt(0, 0),
		draw.Over,
	)

	if r.inner != nil {
		r.inner.Draw(x, y, im)
	}
}
