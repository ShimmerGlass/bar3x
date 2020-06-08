package base

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/fogleman/gg"
	"github.com/shimmerglass/bar3x/ui"
)

const (
	BarDirectionBottomTop = "bottom-top"
	BarDirectionLeftRight = "left-right"
)

type Bar struct {
	Base

	radius    int
	direction string
	fgColor   color.Color
	bgColor   color.Color

	value float64
}

func NewBar(p ui.ParentDrawable) *Bar {
	return &Bar{
		Base:      NewBase(p),
		direction: BarDirectionBottomTop,
	}
}

func (b *Bar) Value() float64 {
	return b.value
}
func (b *Bar) SetValue(v float64) {
	b.value = v
}

func (b *Bar) SetWidth(v int) {
	b.width.set(v)
}
func (b *Bar) SetHeight(v int) {
	b.height.set(v)
}

func (b *Bar) Radius() int {
	return b.radius
}
func (b *Bar) SetRadius(v int) {
	b.radius = v
}

func (b *Bar) Direction() string {
	return b.direction
}
func (b *Bar) SetDirection(v string) {
	b.direction = v
}

func (b *Bar) FgColor() color.Color {
	return b.fgColor
}
func (b *Bar) SetFgColor(v color.Color) {
	b.fgColor = v
}

func (b *Bar) BgColor() color.Color {
	return b.bgColor
}
func (b *Bar) SetBgColor(v color.Color) {
	b.bgColor = v
}

func (b *Bar) Draw(dx, dy int, im draw.Image) {
	dc := gg.NewContext(b.width.v, b.height.v)

	r, w, h, z := float64(b.radius), float64(b.width.v), float64(b.height.v), float64(0)

	var x10, x11, x12,
		y10, y11, y12, y13,
		x20, x21, x22,
		y20, y21, y22, y23 float64

	switch b.direction {
	case BarDirectionLeftRight:
		s := math.Round(b.value * float64(b.width.v))

		x10, x11, x12 = z, r, s
		y10, y11, y12, y13 = z, r, h, h-r

		x20, x21, x22 = s, w-r, w
		y20, y21, y22, y23 = z, r, h-r, h

		dc.NewSubPath()
		dc.MoveTo(x11, y10)
		dc.LineTo(x12, y10)
		dc.LineTo(x12, y12)
		dc.LineTo(x11, y12)
		dc.DrawArc(x11, y13, r, gg.Radians(90), gg.Radians(180))
		dc.LineTo(x10, y11)
		dc.DrawArc(x11, y11, r, gg.Radians(180), gg.Radians(270))
		dc.ClosePath()
		dc.SetColor(b.fgColor)
		dc.Fill()

		dc.NewSubPath()
		dc.MoveTo(x20, y20)
		dc.LineTo(x21, y20)
		dc.DrawArc(x21, y21, r, gg.Radians(270), gg.Radians(360))
		dc.LineTo(x22, y22)
		dc.DrawArc(x21, y22, r, gg.Radians(0), gg.Radians(90))
		dc.LineTo(x20, y23)
		dc.LineTo(x20, y20)
		dc.ClosePath()
		dc.SetColor(b.bgColor)
		dc.Fill()

	case BarDirectionBottomTop:
		s := math.Round((1 - b.value) * float64(b.height.v))
		x0, x1, x2, x3 := z, r, w-r, w

		y10, y11, y12 := z, r, h
		y20, y21, y22 := s, s+h-r, s+h

		dc.NewSubPath()
		dc.MoveTo(x1, y10)
		dc.LineTo(x2, y10)
		dc.DrawArc(x2, y11, r, gg.Radians(270), gg.Radians(360))
		dc.LineTo(x3, y12)
		dc.LineTo(x0, y12)
		dc.LineTo(x0, y11)
		dc.DrawArc(x1, y11, r, gg.Radians(180), gg.Radians(270))
		dc.ClosePath()
		dc.SetColor(b.bgColor)
		dc.Fill()

		dc.NewSubPath()
		dc.MoveTo(x0, y20)
		dc.LineTo(x3, y20)
		dc.LineTo(x3, y21)
		dc.DrawArc(x2, y21, r, gg.Radians(0), gg.Radians(90))
		dc.LineTo(x1, y22)
		dc.DrawArc(x1, y21, r, gg.Radians(90), gg.Radians(180))
		dc.LineTo(x0, y20)
		dc.ClosePath()
		dc.SetColor(b.fgColor)
		dc.Fill()
	}

	draw.Draw(
		im,
		image.Rect(dx, dy, dx+b.width.v, dy+b.height.v),
		dc.Image(),
		image.ZP,
		draw.Over,
	)
}
