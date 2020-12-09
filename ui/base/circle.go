package base

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/fogleman/gg"
	"github.com/shimmerglass/bar3x/ui"
)

var _ ui.Drawable = (*Circle)(nil)

type Circle struct {
	Base

	color  color.Color
	radius int
}

func NewCircle(p ui.ParentDrawable) *Circle {
	return &Circle{
		Base: NewBase(p),
	}
}

func (c *Circle) SetContext(ctx ui.Context) {
	c.Base.SetContext(ctx)
}

func (b *Circle) Color() color.Color {
	return b.color
}
func (b *Circle) SetColor(v color.Color) {
	b.color = v
}

func (b *Circle) Radius() int {
	return b.radius
}
func (b *Circle) SetRadius(v int) {
	b.radius = v
	b.width.Set(v * 2)
	b.height.Set(v * 2)
}

func (s *Circle) SendEvent(ev ui.Event) bool {
	if !s.Base.SendEvent(ev) {
		return false
	}

	return true
}

func (r *Circle) Draw(x, y int, im draw.Image) {
	w, h := r.Width(), r.Height()
	dc := gg.NewContext(w, h)
	dc.DrawCircle(
		float64(r.radius),
		float64(r.radius),
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
}
