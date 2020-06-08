package module

import (
	"fmt"
	"math"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type TextUnit struct {
	moduleBase

	Value *base.Text
	Unit  *base.Text

	mk *markup.Markup
}

func NewTextUnit(p ui.ParentDrawable, mk *markup.Markup) *TextUnit {
	return &TextUnit{
		mk:         mk,
		moduleBase: newBase(p),
	}
}

func (t *TextUnit) Init() error {
	_, err := t.mk.Parse(t, t, `
		<Row ref="Root">
			<Text
				ref="Value"
				Color="{text_color}"
				Font="{text_font}"
				FontSize="{text_font_size}"
			/>
			<Text
				ref="Unit"
				Color="{inactive_light_color}"
				Font="{text_font}"
				FontSize="{text_font_size}"
			/>
		</Row>
	`)
	return err
}

func (t *TextUnit) Set(v, u string) {
	t.Value.SetText(v)
	t.Unit.SetText(u)
}

func humanateBytes(s uint64) (string, string) {
	base := 1024.0
	sizes := []string{"", "k", "M", "G", "T", "P", "E"}

	if s < 10 {
		return fmt.Sprintf("%d", s), ""
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f, val), suffix
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
