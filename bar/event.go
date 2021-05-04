package bar

import (
	"image"

	"github.com/shimmerglass/bar3x/ui"
)

func (b *Bar) displatchEvent(ev ui.Event) {
	if b.LeftRoot != nil {
		h := b.LeftRoot.Height()
		x, y := b.padding, (b.height-h)/2

		nev := ev
		nev.At = ev.At.Sub(image.Pt(x, y))
		b.LeftRoot.SendEvent(nev)
	}

	if b.CenterRoot != nil {
		w, h := b.CenterRoot.Width(), b.CenterRoot.Height()
		x, y := (b.screen.Width-w)/2, (b.height-h)/2

		nev := ev
		nev.At = ev.At.Sub(image.Pt(x, y))
		b.CenterRoot.SendEvent(nev)
	}

	if b.RightRoot != nil {
		w, h := b.RightRoot.Width(), b.RightRoot.Height()
		x, y := b.screen.Width-w-b.padding-b.TrayWidth, (b.height-h)/2

		nev := ev
		nev.At = ev.At.Sub(image.Pt(x, y))
		b.RightRoot.SendEvent(nev)
	}
}
