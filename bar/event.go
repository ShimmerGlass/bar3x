package bar

import (
	"image"

	"github.com/shimmerglass/bar3x/ui"
)

func (b *Bar) displatchEvent(ev ui.Event) {
	{
		h := b.LeftRoot.Height()
		x, y := b.padding, (b.height-h)/2

		nev := ev
		nev.At = ev.At.Sub(image.Pt(x, y))
		b.LeftRoot.SendEvent(nev)
	}

	{
		w, h := b.CenterRoot.Width(), b.CenterRoot.Height()
		x, y := (b.screen.Width-w)/2, (b.height-h)/2

		nev := ev
		nev.At = ev.At.Sub(image.Pt(x, y))
		b.CenterRoot.SendEvent(nev)
	}

	{
		w, h := b.RightRoot.Width(), b.RightRoot.Height()
		x, y := b.screen.Width-w-b.padding, (b.height-h)/2

		nev := ev
		nev.At = ev.At.Sub(image.Pt(x, y))
		b.RightRoot.SendEvent(nev)
	}
}
