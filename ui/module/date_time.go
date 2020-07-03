package module

import (
	"context"
	"time"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	log "github.com/sirupsen/logrus"
)

type DateTime struct {
	moduleBase

	clock *Clock
	mk    *markup.Markup
	Txt   *base.Text

	format   string
	tz       *time.Location
	showIcon bool
}

func NewDateTime(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *DateTime {
	d := &DateTime{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),
		tz:         time.Local,
		format:     "Mon Jan 2 3:04 PM",
	}
	d.SetOnLeftClick(func(ui.Event) bool {
		ui.StartCommand("gnome-calendar")
		return true
	})
	return d
}

func (m *DateTime) Init() error {
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

func (m *DateTime) Update(context.Context) {
	m.Txt.SetText(time.Now().In(m.tz).Format(m.format))
}

// parameters

func (m *DateTime) Format() string {
	return m.format
}
func (m *DateTime) SetFormat(v string) {
	m.format = v
}

func (m *DateTime) Timezone() string {
	return m.tz.String()
}
func (m *DateTime) SetTimezone(v string) {
	tz, err := time.LoadLocation(v)
	if err != nil {
		log.Fatalf("Data: Timezone: %s", err)
	}
	m.tz = tz
}
