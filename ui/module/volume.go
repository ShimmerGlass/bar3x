package module

import (
	"github.com/shimmerglass/bar3x/lib/pulse"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	"github.com/sirupsen/logrus"
)

const (
	volumeMuteIcon = "\uf026"
	volumeIcon     = "\uf028"
)

type Volume struct {
	moduleBase

	mk *markup.Markup

	Icon *base.Text
	Bar  *base.Bar
}

func NewVolume(p ui.ParentDrawable, mk *markup.Markup) *Volume {
	return &Volume{
		mk:         mk,
		moduleBase: newBase(p),
	}
}

func (m *Volume) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Row ref="Root">
			<Sizer PaddingRight="{h_padding}">
				<Icon ref="Icon" />
			</Sizer>
			<Bar
				ref="Bar"
				Width="100"
				Height="6"
				Radius="2"
				Direction="left-right"
				FgColor="{accent_color}"
				BgColor="{neutral_color}"
				Hover:BgColor="{neutral_light_color}"
			/>
		</Row>
	`)

	m.Bar.SetOnLeftClick(func(ev ui.Event) bool {
		v := float64(ev.At.X) / float64(m.Bar.Width())
		err := pulse.SetVolume(v)
		if err != nil {
			logrus.Error(err)
		}
		return true
	})

	c := make(chan struct{})
	pulse.Watch(c)

	go func() {
		for range c {
			vol := pulse.Volume()

			if vol > 1 {
				vol = 1
			}
			var i string
			if vol > 0 {
				i = volumeIcon
			} else {
				i = volumeMuteIcon
			}
			m.Icon.SetText(i)
			m.Bar.SetValue(vol)

			m.Notify()
		}
	}()

	return err
}
