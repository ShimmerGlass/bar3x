package module

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/shirou/gopsutil/disk"
)

type DiskBandwidth struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock

	devs  []string
	lastT time.Time
	last  map[string]disk.IOCountersStat

	Transfer *transfer
}

func NewDiskBW(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *DiskBandwidth {
	return &DiskBandwidth{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),
		last:       map[string]disk.IOCountersStat{},
	}
}

func (m *DiskBandwidth) Init() error {
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

func (m *DiskBandwidth) Devs() []string {
	return m.devs
}

func (m *DiskBandwidth) SetDevs(v []string) {
	m.devs = v
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

	m.lastT = time.Now()
	m.Transfer.Set(int(readPs), int(writePs))
}
