package module

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"

	owm "github.com/briandowns/openweathermap"
)

type Weather struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock

	TempTxt *TextUnit
	Icon    *base.Text

	location string
	apiKey   string
	unit     string

	w *owm.CurrentWeatherData
}

func NewWeather(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *Weather {
	return &Weather{
		mk:         mk,
		clock:      clock,
		moduleBase: newBase(p),
		unit:       "C",
	}
}

func (m *Weather) Init() error {
	w, err := owm.NewCurrent(m.unit, "en", m.apiKey)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	m.w = w

	_, err = m.mk.Parse(m, m, `
		<Row ref="Root">
			<Sizer PaddingRight="{h_padding}">
				<Icon ref="Icon" />
			</Sizer>
			<TxtUnit ref="TempTxt" />
		</Row>
	`)
	if err != nil {
		return err
	}

	m.clock.Add(m, 30*time.Minute)
	return nil
}

func (m *Weather) Update(context.Context) {
	if m.location == "" {
		return
	}

	err := m.w.CurrentByName(m.location)
	if err != nil {
		log.Println(err)
		m.SetVisible(false)
		return
	}

	var icon string
	if len(m.w.Weather) > 0 {
		switch m.w.Weather[0].Icon {
		case "01d":
			icon = "\ue30d"
		case "02d":
			icon = "\ue302"
		case "03d", "03n":
			icon = "\ue33d"
		case "04d", "04n":
			icon = "\ue312"
		case "09d", "09n":
			icon = "\ue318"
		case "10d":
			icon = "\ue309"
		case "11d", "11n":
			icon = "\ue31d"
		case "13d":
			icon = "\ue30a"
		case "50d":
			icon = "\ue3ae"
		case "01n":
			icon = "\ue32b"
		case "02n":
			icon = "\ue32e"
		case "10n":
			icon = "\ue334"
		case "13n":
			icon = "\ue327"
		case "50n":
			icon = "\ue346"
		}
	}

	m.Icon.SetText(icon)
	m.TempTxt.Set(fmt.Sprintf("%.1f", m.w.Main.Temp), "°")
	m.SetVisible(true)
}

// parameters

func (m *Weather) Location() string {
	return m.location
}
func (m *Weather) SetLocation(v string) {
	m.location = v
}

func (m *Weather) ApiKey() string {
	return m.apiKey
}
func (m *Weather) SetApiKey(v string) {
	m.apiKey = v
}

func (m *Weather) Unit() string {
	return m.unit
}
func (m *Weather) SetUnit(v string) {
	m.unit = v
}
