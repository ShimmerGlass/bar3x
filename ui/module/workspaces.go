package module

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	log "github.com/sirupsen/logrus"
	"go.i3wm.org/i3/v4"
)

type workspaceIndicator struct {
	Root      *base.Sizer
	Rect      *base.Rect
	Text      *base.Text
	Content   *base.Sizer
	Indicator *base.Rect
}

type Workspaces struct {
	moduleBase

	mk          *markup.Markup
	display     string
	maxWidth    int
	onlyCurrent bool

	Row *base.Row
	els []*workspaceIndicator
}

func NewWorkspaces(p ui.ParentDrawable, mk *markup.Markup) *Workspaces {
	return &Workspaces{
		mk:         mk,
		moduleBase: newBase(p),
		maxWidth:   -1,
	}
}

func (m *Workspaces) Init() error {
	root, err := m.mk.Parse(m, m, `
		<Row ref="Row" />
	`)
	if err != nil {
		return err
	}
	m.Root = root
	m.display = m.Context().MustString("display")

	go func() {
		m.update()
		er := i3.Subscribe(i3.WorkspaceEventType)
		for er.Next() {
			m.update()
		}
	}()

	return err
}

func (m *Workspaces) update() {
	wks, err := i3.GetWorkspaces()
	if err != nil {
		log.Errorf("Workspaces: %s", err)
		return
	}

	sort.Slice(wks, func(i, j int) bool {
		return strings.Compare(wks[i].Name, wks[j].Name) < 1
	})

	j := 0
	for _, wk := range wks {
		if wk.Output != m.display {
			continue
		}
		if m.onlyCurrent && !wk.Visible {
			continue
		}

		if j > len(m.els)-1 {
			m.addIndicator()
		}

		el := m.els[j]

		switch {
		case !wk.Visible && !wk.Focused:
			el.Rect.SetColor(m.Context().MustColor("neutral_color"))
			el.Text.SetColor(m.Context().MustColor("neutral_light_color"))
			el.Indicator.SetColor(m.Context().MustColor("neutral_color"))
		case wk.Visible && !wk.Focused:
			el.Rect.SetColor(m.Context().MustColor("neutral_color"))
			el.Text.SetColor(m.Context().MustColor("text_color"))
			el.Indicator.SetColor(m.Context().MustColor("neutral_light_color"))
		case wk.Visible && wk.Focused:
			el.Rect.SetColor(m.Context().MustColor("neutral_color"))
			el.Text.SetColor(m.Context().MustColor("text_color"))
			el.Indicator.SetColor(m.Context().MustColor("accent_color"))
		}

		el.Text.SetText(wk.Name)
		el.Root.SetVisible(true)

		func(wk i3.Workspace) {
			el.Root.SetOnLeftClick(func(ui.Event) bool {
				i3.RunCommand(fmt.Sprintf("workspace %s", wk.Name))
				return false
			})

			el.Root.SetOnPointerEnter(func(ui.Event) bool {
				el.Indicator.SetColor(m.Context().MustColor("accent_color"))
				el.Indicator.Notify()
				return true
			})

			el.Root.SetOnPointerLeave(func(ui.Event) bool {
				switch {
				case wk.Focused:
					el.Indicator.SetColor(m.Context().MustColor("accent_color"))
				case wk.Visible:
					el.Indicator.SetColor(m.Context().MustColor("neutral_light_color"))
				default:
					el.Indicator.SetColor(m.Context().MustColor("neutral_color"))
				}
				el.Indicator.Notify()
				return true
			})
		}(wk)

		j++
	}

	for i := j; i < len(m.els); i++ {
		m.els[i].Root.SetVisible(false)
	}

	m.Notify()
}

func (m *Workspaces) addIndicator() {
	el := &workspaceIndicator{}
	root := m.mk.MustParse(m.Row, el, `
		<Sizer
			ref="Root"
			PaddingRight="{is_last_visible ? 0 : h_padding}"
		>
			<Rect ref="Rect" Radius="1">
				<Col>
					<Sizer
						ref="Content"
						Height="{bar_height - v_padding * 2}"
						PaddingLeft="{h_padding}"
						PaddingRight="{h_padding}"
					>
						<Text ref="Text" />
					</Sizer>
					<Rect
						ref="Indicator"
						Height="2"
						Width="{$Content.Width}"
					/>
				</Col>
			</Rect>
		</Sizer>
	`)
	if m.maxWidth != -1 {
		el.Text.SetMaxWidth(m.maxWidth)
	}
	m.Row.Add(root)
	m.els = append(m.els, el)
}

// parameters

func (m *Workspaces) MaxWidth() int {
	return m.maxWidth
}
func (m *Workspaces) SetMaxWidth(v int) {
	m.maxWidth = v
}

func (m *Workspaces) OnlyCurrent() bool {
	return m.onlyCurrent
}
func (m *Workspaces) SetOnlyCurrent(v bool) {
	m.onlyCurrent = v
}
