package base

import (
	"image"
	"image/draw"

	"github.com/shimmerglass/bar3x/ui"
)

type Pattern struct {
	Base

	inner ui.Drawable
}

func NewPattern(p ui.ParentDrawable) *Pattern {
	return &Pattern{
		Base: NewBase(p),
	}
}

func (p *Pattern) SetWidth(v int) {
	p.width.Set(v)
}
func (p *Pattern) SetHeight(v int) {
	p.height.Set(v)
}

func (p *Pattern) ChildContext(i int) ui.Context {
	return p.ctx
}
func (p *Pattern) Children() []ui.Drawable {
	return []ui.Drawable{p.inner}
}
func (p *Pattern) SetContext(ctx ui.Context) {
	p.Base.SetContext(ctx)
	if p.inner != nil {
		p.inner.SetContext(ctx)
	}
}

func (p *Pattern) Add(d ui.Drawable) {
	p.inner = d
}
func (p *Pattern) SendEvent(ev ui.Event) bool {
	if !p.Base.SendEvent(ev) {
		return false
	}

	if !p.inner.Visible() {
		return true
	}

	iev := ev
	iev.At = image.Pt(ev.At.X%p.inner.Width(), ev.At.Y%p.inner.Height())
	p.inner.SendEvent(iev)

	return true
}

func (p *Pattern) Draw(dx, dy int, im draw.Image) {
	for y := 0; y < p.Height(); y += p.inner.Height() {
		for x := 0; x < p.Width(); x += p.inner.Width() {
			p.inner.Draw(dx+x, dy+y, im)
		}
	}
}
