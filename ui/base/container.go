package base

import (
	"image"
	"image/draw"

	"github.com/shimmerglass/bar3x/ui"
)

type drawableAt struct {
	D    ui.Drawable
	Rect image.Rectangle
}

type containerLayout []drawableAt

func (l *containerLayout) Add(d ui.Drawable, x, y, w, h int) {
	*l = append(*l, drawableAt{
		D:    d,
		Rect: image.Rect(x, y, x+w, y+h),
	})
}

func (l containerLayout) SendEvent(ev ui.Event, last ui.Event) {
	for _, d := range l {
		if ev.At.In(d.Rect) {
			iev := ev
			iev.At = ev.At.Sub(d.Rect.Min)
			d.D.SendEvent(iev)

			if !last.At.In(d.Rect) {
				d.D.SendEvent(ui.Event{
					Type: ui.EventPointerEnter,
					At:   ev.At.Sub(d.Rect.Min),
				})
			}
		}

		if last.At.In(d.Rect) && !ev.At.In(d.Rect) {
			d.D.SendEvent(ui.Event{
				Type: ui.EventPointerLeave,
				At:   last.At.Sub(d.Rect.Min),
			})
		}
	}
}

func (l containerLayout) Draw(x, y int, im draw.Image) {
	for _, d := range l {
		d.D.Draw(x+d.Rect.Min.X, y+d.Rect.Min.Y, im)
	}
}
