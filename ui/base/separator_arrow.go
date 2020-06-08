package base

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/fogleman/gg"
	"github.com/shimmerglass/bar3x/ui"
)

type SeparatorArrow struct {
	Base
	color color.Color
}

func NewSeparatorArrow(p ui.ParentDrawable) *SeparatorArrow {
	return &SeparatorArrow{Base: NewBase(p)}
}

func (b *SeparatorArrow) SetWidth(v int) {
	b.width.set(v)
}
func (b *SeparatorArrow) SetHeight(v int) {
	b.height.set(v)
}

func (b *SeparatorArrow) Color() color.Color {
	return b.color
}
func (b *SeparatorArrow) SetColor(v color.Color) {
	b.color = v
}

func (s *SeparatorArrow) Draw(x, y int, im draw.Image) {
	w, h := s.width.v, s.height.v

	dc := gg.NewContext(w, h)
	dc.SetColor(s.color)
	dc.DrawLine(float64(w), 0, 0, float64(h)/2)
	dc.DrawLine(0, float64(h)/2, float64(w), float64(h))
	dc.Stroke()

	draw.Draw(im, im.Bounds().Add(image.Pt(x, y)), dc.Image(), image.ZP, draw.Over)
}
