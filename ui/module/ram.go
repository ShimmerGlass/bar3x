package module

import (
	"context"
	"fmt"
	"time"

	"github.com/shimmerglass/bar3x/lib/process"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

type RAM struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock
	pw    *process.ProcessWatcher

	GbTxt     *TextUnit
	ProcTxt   *base.Text
	ProcSizer *base.Sizer
	TextRow   *base.Row
	BarEl     *base.Bar
	BarSizer  *base.Sizer

	showBar  bool
	showProc bool
	format   string
}

func NewRAM(p ui.ParentDrawable, mk *markup.Markup, clock *Clock, pw *process.ProcessWatcher) *RAM {
	return &RAM{
		mk:         mk,
		clock:      clock,
		pw:         pw,
		moduleBase: newBase(p),
		format:     "free",
		showBar:    true,
		showProc:   true,
	}
}

func (m *RAM) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Row ref="Root">
			<Sizer PaddingRight="{h_padding}">
				<Icon>{icons.chip2}</Icon>
			</Sizer>
			<Col>
				<Row ref="TextRow">
					<TxtUnit ref="GbTxt" />
					<Sizer
						ref="ProcSizer"
						PaddingLeft="{h_padding}"
						Width="60"
						HAlign="right"
					>
						<Text
							ref="ProcTxt"
							MaxWidth="60"
							Color="{neutral_light_color}"
						/>
					</Sizer>
				</Row>
				<Sizer ref="BarSizer" PaddingTop="3">
					<Bar
						ref="BarEl"
						Height="2"
						Width="{$TextRow.Width}"
						Direction="left-right"
						BgColor="{neutral_color}"
						FgColor="{accent_color}"
					/>
				</Sizer>
			</Col>
		</Row>
	`)
	if err != nil {
		return err
	}

	m.clock.Add(m, time.Second)
	return nil
}

func (m *RAM) Update(context.Context) {
	v, _ := mem.VirtualMemory()
	free := v.Free + v.Cached

	switch m.format {
	case "free":
		m.GbTxt.Set(humanateBytes(free))
	case "used":
		m.GbTxt.Set(humanateBytes(v.Total - free))
	case "used-percent":
		m.GbTxt.Set(fmt.Sprintf("%.0f", float64(v.Total-free)/float64(v.Total)), "%")
	default:
		log.Fatalf("RAM: unknown format %q, possible values are free, free, used-percent")
	}

	if m.showProc {
		m.ProcTxt.SetText(m.pw.MaxRAM)
	} else {
		m.ProcSizer.SetVisible(false)
	}

	pc := float64(v.Total-free) / float64(v.Total)
	m.BarEl.SetValue(pc)
}

// parameters

func (m *RAM) Format() string {
	return m.format
}
func (m *RAM) SetFormat(v string) {
	m.format = v
}

func (m *RAM) ShowBar() bool {
	return m.showBar
}
func (m *RAM) SetShowBar(b bool) {
	m.showBar = b
}

func (m *RAM) ShowMaxProcess() bool {
	return m.showProc
}
func (m *RAM) SetShowMaxProcess(v bool) {
	m.showProc = v
}
