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

	keyID     string
	keySecret string

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
						Width="{bar_height - v_padding * 2}"
						Height="{bar_height - v_padding * 2}"
						Color="{neutral_color}"
						Visible="{!$Img.Visible}"
					/>
					<Image
						ref="Img"
						Width="{bar_height - v_padding * 2}"
						Height="{bar_height - v_padding * 2}"
					/>
				</Layers>
			</Sizer>
			<Col>
				<Row ref="TextRow">
					<Text ref="Title" MaxWidth="120" />
					<Sizer ref="Sep" PaddingLeft="{h_padding}" PaddingRight="{h_padding}">
						<Icon>{icons.dot}</Icon>
					</Sizer>
					<Text ref="Artist" MaxWidth="120" />
				</Row>
				<Sizer PaddingTop="3">
					<Bar
						ref="Progress"
						FgColor="{accent_color}"
						BgColor="{neutral_color}"
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
		spot := spotify.New(m.keyID, m.keySecret)
		for s := range spot.Up {
			if !s.Playing {
				m.SetVisible(false)
				m.Notify()
				continue
			} else {
				m.SetVisible(true)
			}

			m.Artist.SetText(s.Artist)
			m.Title.SetText(s.Title)

			if s.Artist == "" {
				m.Artist.SetVisible(false)
				m.Sep.SetVisible(false)
			} else {
				m.Artist.SetVisible(true)
				m.Sep.SetVisible(true)
			}

			m.Img.SetImage(s.Image)
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

// parameters

func (m *Music) SpotifyKeyID() string {
	return m.keyID
}
func (m *Music) SetSpotifyKeyID(v string) {
	m.keyID = v
}

func (m *Music) SpotifyKeySecret() string {
	return m.keySecret
}
func (m *Music) SetSpotifyKeySecret(v string) {
	m.keySecret = v
}
