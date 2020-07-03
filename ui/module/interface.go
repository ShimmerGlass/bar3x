package module

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shimmerglass/bar3x/lib/bandwidth"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type Interface struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock
	bw    *bandwidth.BW

	iface     string
	showLabel bool

	Transfer   *transfer
	Label      *base.Text
	LabelSizer *base.Sizer
}

func NewInterface(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *Interface {
	return &Interface{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),
		showLabel:  true,
	}
}

func (m *Interface) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Row ref="Root">
			<Sizer ref="LabelSizer" PaddingRight="{h_padding}">
				<Text ref="Label" Color="{accent_color}" />
			</Sizer>
			<Transfer ref="Transfer" />
		</Row>
	`)
	if err != nil {
		return err
	}

	m.clock.Add(m, time.Second)
	return nil
}

func (m *Interface) Update(context.Context) {
	if m.bw == nil {
		return
	}

	rw, tx, err := m.bw.Read()
	if err != nil {
		log.Println(err)
		return
	}

	m.Transfer.Set(rw, tx)

	if m.showLabel {
		m.Label.SetText(m.iface)
	} else {
		m.LabelSizer.SetVisible(false)
	}
}

// parameters

func (m *Interface) Iface() string {
	return m.iface
}
func (m *Interface) SetIface(s string) {
	m.iface = s
	m.bw = bandwidth.New(s)
}

func (m *Interface) ShowLabel() bool {
	return m.showLabel
}
func (m *Interface) SetShowLabel(v bool) {
	m.showLabel = v
}
