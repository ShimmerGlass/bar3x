package module

import (
	"context"
	"math"
	"runtime"
	"strconv"
	"time"

	"github.com/shimmerglass/bar3x/lib/cpu"
	"github.com/shimmerglass/bar3x/lib/process"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type CPU struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock
	pw    *process.ProcessWatcher

	PcTxt     *TextUnit
	ProcTxt   *base.Text
	ProcSizer *base.Sizer
	TextRow   *base.Row
	Bar       *base.Bar
	BarSizer  *base.Sizer
	BarsSizer *base.Sizer
	BarsRow   *base.Row
	Bars      []*base.Bar

	perCoreBars bool
	maxProcess  bool
	avgBar      bool
}

func NewCPU(p ui.ParentDrawable, mk *markup.Markup, clock *Clock, pw *process.ProcessWatcher) *CPU {
	return &CPU{
		mk:         mk,
		clock:      clock,
		pw:         pw,
		moduleBase: newBase(p),

		perCoreBars: true,
		maxProcess:  true,
		avgBar:      true,
	}
}

func (m *CPU) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Row ref="Root">
			<Sizer PaddingRight="{h_padding}">
				<Icon>{icons.chip}</Icon>
			</Sizer>
			<Col>
				<Row ref="TextRow">
					<Sizer Width="30" HAlign="left">
						<TxtUnit ref="PcTxt" />
					</Sizer>
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
						ref="Bar"
						Height="2"
						Width="{$TextRow.Width}"
						Direction="left-right"
						BgColor="{neutral_color}"
						FgColor="{accent_color}"
					/>
				</Sizer>
			</Col>
			<Sizer ref="BarsSizer" PaddingLeft="{h_padding}">
				<Row ref="BarsRow" />
			</Sizer>
		</Row>
	`)
	if err != nil {
		return err
	}

	type barRefs struct {
		Bar *base.Bar
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		refs := &barRefs{}
		bar := m.mk.MustParse(m.BarsRow, refs, `
			<Sizer PaddingLeft="1">
				<Bar
					ref="Bar"
					BgColor="{neutral_color}"
					FgColor="{accent_color}"
					Width="3"
					Height="{bar_height - v_padding * 2}"
				/>
			</Sizer>
		`)
		m.Bars = append(m.Bars, refs.Bar)
		m.BarsRow.Add(bar)
	}

	m.clock.Add(m, time.Second)
	cpu.Start()
	return nil
}

func (m *CPU) Update(context.Context) {
	usage := cpu.Read()

	avg := 0.0
	if len(usage) > 0 {
		for _, v := range usage {
			avg += v
		}
		avg /= float64(len(usage))
	}

	m.PcTxt.Set(m.formatCPU(avg*100), "%")

	if m.maxProcess {
		m.ProcTxt.SetText(m.pw.MaxCPU)
	} else {
		m.ProcSizer.SetVisible(false)
	}

	if m.perCoreBars {
		for i, v := range usage {
			m.Bars[i].SetValue(v)
		}
	} else {
		m.BarsSizer.SetVisible(false)
	}

	if m.avgBar {
		m.Bar.SetValue(avg)
	} else {
		m.BarSizer.SetVisible(false)
	}
}

func (m *CPU) formatCPU(v float64) string {
	usage := int(math.Round(v))
	usageStr := strconv.Itoa(usage)
	switch len(usageStr) {
	case 3:
		usageStr = "00"
	case 1:
		usageStr = " " + usageStr
	}

	return usageStr
}

// parameters

func (m *CPU) ShowPerCoreBars() bool {
	return m.perCoreBars
}
func (m *CPU) SetShowPerCoreBars(v bool) {
	m.perCoreBars = v
}

func (m *CPU) ShowMaxProcess() bool {
	return m.maxProcess
}
func (m *CPU) SetShowMaxProcess(v bool) {
	m.maxProcess = v
}

func (m *CPU) ShowAvgBar() bool {
	return m.avgBar
}
func (m *CPU) SetShowAvgBar(v bool) {
	m.avgBar = v
}
