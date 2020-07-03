package module

import (
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/shirou/gopsutil/disk"
)

type DiskBandwidth struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock

	devs      []string
	showLabel bool
	unit      string

	lastT time.Time
	last  map[string]disk.IOCountersStat

	Transfer   *transfer
	Label      *base.Text
	LabelSizer *base.Sizer
}

func NewDiskBW(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *DiskBandwidth {
	return &DiskBandwidth{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),
		last:       map[string]disk.IOCountersStat{},
		showLabel:  true,
	}
}

func (m *DiskBandwidth) Init() error {
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

func (m *DiskBandwidth) Update(context.Context) {
	allStats, err := disk.IOCounters(m.devs...)
	if err != nil {
		log.Println(err)
		return
	}
	dt := time.Since(m.lastT)

	var readPs, writePs float64
	for name, stats := range allStats {
		lastStats, ok := m.last[name]
		if !ok {
			m.last[name] = stats
		} else {
			readD := stats.ReadBytes - lastStats.ReadBytes
			writeD := stats.WriteBytes - lastStats.WriteBytes
			m.last[name] = stats

			readPs += float64(readD) / dt.Seconds()
			writePs += float64(writeD) / dt.Seconds()
		}
	}

	if m.unit == "bits" {
		readPs *= 8
		writePs *= 8
	}

	m.lastT = time.Now()
	m.Transfer.Set(int(readPs), int(writePs))

	if m.showLabel {
		m.Label.SetText(strings.Join(m.devs, ","))
	} else {
		m.LabelSizer.SetVisible(false)
	}
}

// parameters

func (m *DiskBandwidth) ShowLabel() bool {
	return m.showLabel
}
func (m *DiskBandwidth) SetShowLabel(v bool) {
	m.showLabel = v
}

func (m *DiskBandwidth) Unit() string {
	return m.unit
}
func (m *DiskBandwidth) SetUnit(v string) {
	m.unit = v
}

func (m *DiskBandwidth) Devs() []string {
	return m.devs
}
func (m *DiskBandwidth) SetDevs(v []string) {
	m.devs = v
}
