package module

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/shirou/gopsutil/disk"
)

type DiskUsage struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock

	Txt *TextUnit

	mountPoint string
}

func NewDiskUsage(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *DiskUsage {
	return &DiskUsage{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),
		mountPoint: "/",
	}
}

func (m *DiskUsage) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Sizer ref="Root" Height="{bar_height}">
			<Row>
				<Icon>{icons.disk}</Icon>
				<Sizer PaddingLeft="{h_padding}">
					<TxtUnit ref="Txt" />
				</Sizer>
			</Row>
		</Sizer>
	`)
	if err != nil {
		return err
	}

	m.clock.Add(m, 10*time.Second)
	return nil
}

func (m *DiskUsage) Update(context.Context) {
	stat, err := disk.Usage("/")
	var p float64
	if err != nil {
		log.Println(err)
		p = 0
	} else {
		p = stat.UsedPercent
	}

	m.Txt.Set(fmt.Sprintf("%.0f", p), "%")
}

// parameters

func (m *DiskUsage) MountPoint() string {
	return m.mountPoint
}
func (m *DiskUsage) SetMountPoint(v string) {
	m.mountPoint = v
}
