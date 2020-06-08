package module

import (
	"time"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type Date struct {
	moduleBase

	clock *Clock
	mk    *markup.Markup
	Txt   *base.Text
}

func NewDate(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *Date {
	return &Date{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),
	}
}

func (m *Date) Init() error {
	m.mk.MustParse(m, m, `
		<Row ref="Root">
			<Sizer PaddingRight="{h_padding}">
				<Icon>{icons.calendar}</Icon>
			</Sizer>
			<Text ref="Txt" />
		</Row>
	`)

	m.clock.Add(m, time.Second)
	return nil
}

func (m *Date) Update() {
	m.Txt.SetText(time.Now().Format("Mon, _2 Jan 15:04:05"))
}
