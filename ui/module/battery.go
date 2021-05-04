package module

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	log "github.com/sirupsen/logrus"
)

type Battery struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock

	name        string
	showPercent string

	Charging *base.Sizer
	Bar      *base.Rect
	Value    *TextUnit
}

func NewBattery(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *Battery {
	return &Battery{
		mk:          mk,
		clock:       clock,
		moduleBase:  newBase(p),
		showPercent: "true",
	}
}

func (m *Battery) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Row ref="Root">
			<Col>
				<Sizer Width="10">
					<Rect
						Color="{accent_color}"
						Width="4"
						Height="2"
					/>
				</Sizer>
				<Layers>
					<Rect
						Color="{accent_color}"
						Width="10"
						Height="14"
					/>
					<Sizer Width="10" Height="14">
						<Rect
							Color="{bg_color}"
							Width="8"
							Height="12"
						/>
					</Sizer>
					<Sizer Width="10" Height="14" PaddingBottom="2" VAlign="bottom">
						<Rect
							ref="Bar"
							Color="{accent_color}"
							Width="6"
							Height="5"
						/>
					</Sizer>
					<Sizer ref="Charging" Width="10" Height="14">
						<Icon Color="{text_color}" FontSize="9">`+"\u26a1"+`</Icon>
					</Sizer>
				</Layers>
			</Col>
			<Sizer PaddingLeft="{h_padding}">
				<TxtUnit ref="Value" />
			</Sizer>
		</Row>
	`)
	if err != nil {
		return err
	}

	m.clock.Add(m, time.Second)
	return nil
}

func (m *Battery) Update(ctx context.Context) {
	value := 0.0
	charging := false

	if c, err := ioutil.ReadFile(filepath.Join("/sys/class/power_supply/", m.name, "capacity")); err == nil {
		v, err := strconv.Atoi(strings.TrimSpace(string(c)))
		if err == nil {
			value = float64(v) / 100
		} else {
			log.Errorf("error parsing battery %s capacity: %s", m.Name(), err)
		}
	} else {
		log.Error("error reading battery %s capacity: %s", m.Name(), err)
	}

	if c, err := ioutil.ReadFile(filepath.Join("/sys/class/power_supply/", m.name, "status")); err == nil {
		charging = !strings.Contains(string(c), "Discharging")
	} else {
		log.Error("error reading battery %s status: %s", m.Name(), err)
	}

	m.Bar.SetHeight(int(value * 10))

	if value < 0.25 {
		m.Bar.SetColor(m.Context().MustColor("danger_color"))
	} else {
		m.Bar.SetColor(m.Context().MustColor("accent_color"))
	}

	if m.showPercent == "true" || (m.showPercent == "low" && value < 0.25) {
		m.Value.Unit.SetText("%")
		m.Value.Value.SetText(fmt.Sprintf("%.0f", value*100))
		m.Value.SetVisible(true)
	} else {
		m.Value.SetVisible(false)
	}

	m.Charging.SetVisible(charging)
}

// parameters

func (m *Battery) Name() string {
	return m.name
}
func (m *Battery) SetName(v string) {
	m.name = v
}

func (m *Battery) ShowPercent() string {
	return m.showPercent
}

func (m *Battery) SetShowPercent(v string) {
	m.showPercent = v
}
