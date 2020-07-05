package base

import (
	"image/color"
	"image/draw"

	"github.com/shimmerglass/bar3x/ui"
)

const (
	GraphDirUp   = "up"
	GraphDirDown = "down"
)

type Graph struct {
	Base
	data      []float64
	direction string
	color     color.Color
}

func NewGraph(p ui.ParentDrawable) *Graph {
	return &Graph{
		Base:      NewBase(p),
		direction: GraphDirUp,
	}
}

func (b *Graph) SetWidth(v int) {
	b.width.Set(v)
}
func (b *Graph) SetHeight(v int) {
	b.height.Set(v)
}

func (b *Graph) Color() color.Color {
	return b.color
}
func (b *Graph) SetColor(v color.Color) {
	b.color = v
}

func (b *Graph) Data() []float64 {
	return b.data
}
func (b *Graph) SetData(v []float64) {
	b.data = v
}

func (b *Graph) Direction() string {
	return b.direction
}
func (b *Graph) SetDirection(v string) {
	b.direction = v
}

func (g Graph) Draw(xt, yt int, im draw.Image) {
	off := g.width.V - len(g.data)
	for i, d := range g.data {
		var y1, y2 int
		x := i + off
		switch g.direction {
		case GraphDirUp:
			y1 = int((1 - d) * float64(g.height.V))
			y2 = g.height.V
		case GraphDirDown:
			y1 = 0
			y2 = int(d * float64(g.height.V))
		}

		for i := y1; i <= y2; i++ {
			im.Set(
				x+xt, yt+i,
				g.color,
			)
		}
	}
}
