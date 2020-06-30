package module

import (
	"context"
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
	d := &Date{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),
	}
	d.SetOnLeftClick(func(ui.Event) bool {
		ui.StartCommand("gnome-calendar")
		return true
	})
	return d
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

func (m *Date) Update(context.Context) {
	m.Txt.SetText(time.Now().Format("Mon, _2 Jan 15:04:05"))
}
