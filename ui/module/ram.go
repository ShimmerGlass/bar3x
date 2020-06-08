package module

import (
	"time"

	"github.com/shimmerglass/bar3x/lib/process"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/shirou/gopsutil/mem"
)

type RAM struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock
	pw    *process.ProcessWatcher

	GbTxt   *TextUnit
	ProcTxt *base.Text
	TextRow *base.Row
	Bar     *base.Bar
}

func NewRAM(p ui.ParentDrawable, mk *markup.Markup, clock *Clock, pw *process.ProcessWatcher) *RAM {
	return &RAM{
		mk:         mk,
		clock:      clock,
		pw:         pw,
		moduleBase: newBase(p),
	}
}

func (m *RAM) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Row ref="Root">
			<Sizer PaddingRight="{h_padding}">
				<Icon>{icons["chip2"]}</Icon>
			</Sizer>
			<Col>
				<Row ref="TextRow">
					<TxtUnit ref="GbTxt" />
					<Sizer
						PaddingLeft="{h_padding}"
						Width="60"
						HAlign="right"
					>
						<Text
							ref="ProcTxt"
							MaxWidth="60"
							Color="{inactive_light_color}"
						/>
					</Sizer>
				</Row>
				<Sizer PaddingTop="3">
					<Bar
						ref="Bar"
						Height="2"
						Width="{$TextRow.Width}"
						Direction="left-right"
						BgColor="{inactive_color}"
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

func (m *RAM) Update() {
	v, _ := mem.VirtualMemory()
	free := v.Free + v.Cached
	m.GbTxt.Set(humanateBytes(free))
	m.ProcTxt.SetText(m.pw.MaxRAM)

	pc := float64(v.Total-free) / float64(v.Total)
	m.Bar.SetValue(pc)
}
