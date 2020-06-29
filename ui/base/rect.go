package base

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/shimmerglass/bar3x/ui"
)

type Rect struct {
	Base
	color color.Color
}

func NewRect(p ui.ParentDrawable) *Rect {
	return &Rect{Base: NewBase(p)}
}

func (b *Rect) SetWidth(v int) {
	b.width.Set(v)
}
func (b *Rect) SetHeight(v int) {
	b.height.Set(v)
}

func (b *Rect) Color() color.Color {
	return b.color
}
func (b *Rect) SetColor(v color.Color) {
	b.color = v
}

func (r *Rect) Draw(x, y int, im draw.Image) {
	_, _, _, a := r.color.RGBA()
	if a == 0 {
		return
	}

	draw.Draw(
		im,
		image.Rect(x, y, x+r.width.V, y+r.height.V),
		image.NewUniform(r.color),
		image.ZP,
		draw.Over,
	)
}
