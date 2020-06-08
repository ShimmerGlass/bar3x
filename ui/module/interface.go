package module

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shimmerglass/bar3x/lib/bandwidth"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type Interface struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock

	iface string
	bw    *bandwidth.BW

	Transfer *transfer
}

func NewInterface(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *Interface {
	return &Interface{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),
	}
}

func (m *Interface) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Row ref="Root">
			<Transfer ref="Transfer" />
		</Row>
	`)
	if err != nil {
		return err
	}

	m.clock.Add(m, time.Second)
	return nil
}

func (m *Interface) Iface() string {
	return m.iface
}
func (m *Interface) SetIface(s string) {
	m.iface = s
	m.bw = bandwidth.New(s)
}

func (m *Interface) Update() {
	if m.bw == nil {
		return
	}

	rw, tx, err := m.bw.Read()
	if err != nil {
		log.Println(err)
		return
	}

	m.Transfer.Set(rw, tx)
}
