package module

import (
	"github.com/shimmerglass/bar3x/lib/spotify"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type Music struct {
	Img         *base.Image
	Placeholder *base.Rect
	TextRow     *base.Row
	Progress    *base.Bar
	Artist      *base.Text
	Sep         *base.Sizer
	Title       *base.Text

	mk *markup.Markup

	state spotify.State

	moduleBase
}

func NewMusic(p ui.ParentDrawable, mk *markup.Markup) *Music {
	return &Music{
		mk:         mk,
		moduleBase: newBase(p),
	}
}

func (m *Music) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Row ref="Root">
			<Sizer PaddingRight="{h_padding}">
				<Layers>
					<Rect
						ref="Placeholder"
						Width="{height - v_padding * 2}"
						Height="{height - v_padding * 2}"
						Color="{inactive_color}"
						Visible="{!$Img.Visible}"
					/>
					<Image
						ref="Img"
						Width="{height - v_padding * 2}"
						Height="{height - v_padding * 2}"
					/>
				</Layers>
			</Sizer>
			<Col>
				<Row ref="TextRow">
					<Text ref="Title" MaxWidth="120" />
					<Sizer ref="Sep" PaddingLeft="{h_padding}" PaddingRight="{h_padding}">
						<Icon>{icons["dot"]}</Icon>
					</Sizer>
					<Text ref="Artist" MaxWidth="120" />
				</Row>
				<Sizer PaddingTop="3">
					<Bar
						ref="Progress"
						FgColor="{accent_color}"
						BgColor="{inactive_color}"
						Height="2"
						Width="{$TextRow.Width}"
						Direction="left-right"
					/>
				</Sizer>
			</Col>
		</Row>
	`)
	if err != nil {
		return err
	}

	go func() {
		spot := spotify.New()
		for s := range spot.Up {
			if !s.Playing {
				m.SetVisible(false)
				m.Notify()
				continue
			} else {
				m.SetVisible(true)
			}
			m.state = s

			m.Artist.SetText(m.state.Artist)
			m.Title.SetText(m.state.Title)

			if m.state.Artist == "" {
				m.Artist.SetVisible(false)
				m.Sep.SetVisible(false)
			} else {
				m.Artist.SetVisible(true)
				m.Sep.SetVisible(true)
			}

			m.Img.SetImage(m.state.Image)
			if s.Image == nil {
				m.Img.SetVisible(false)
			} else {
				m.Img.SetVisible(true)
			}

			m.Progress.SetValue(s.Progress)

			m.Notify()
		}
	}()

	return nil
}
