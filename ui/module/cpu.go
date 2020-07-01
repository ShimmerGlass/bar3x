package module

import (
	"context"
	"math"
	"runtime"
	"strconv"
	"time"

	"github.com/shimmerglass/bar3x/lib/process"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	moduleBase

	mk      *markup.Markup
	clock   *Clock
	pw      *process.ProcessWatcher
	PcTxt   *TextUnit
	ProcTxt *base.Text
	BarsRow *base.Row
	TextRow *base.Row
	Bar     *base.Bar
	Bars    []*base.Bar
}

func NewCPU(p ui.ParentDrawable, mk *markup.Markup, clock *Clock, pw *process.ProcessWatcher) *CPU {
	return &CPU{
		mk:         mk,
		clock:      clock,
		pw:         pw,
		moduleBase: newBase(p),
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
				<Sizer PaddingTop="3">
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
			<Sizer PaddingLeft="{h_padding}">
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
	return nil
}

func (m *CPU) Update(context.Context) {
	vals, _ := cpu.Percent(0, true)
	avg := 0.0
	for _, v := range vals {
		avg += v
	}
	avg /= float64(len(vals))

	m.PcTxt.Set(m.formatCPU(avg), "%")
	m.ProcTxt.SetText(m.pw.MaxCPU)

	for i, b := range m.Bars {
		b.SetValue(vals[i] / 100)
	}

	m.Bar.SetValue(avg / 100)
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
