package module

import (
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
	log "github.com/sirupsen/logrus"
)

type sepRowEl struct {
	moduleBase
	Content ui.ParentDrawable
	el      ui.Drawable
}

func (s *sepRowEl) SetContext(ctx ui.Context) {
	s.moduleBase.SetContext(ctx)
	if s.Root != nil {
		s.Root.SetContext(ctx)
	}
}

func (s *sepRowEl) Visible() bool {
	return s.el.Visible()
}

func (s *sepRowEl) OnVisibleChange(c func(bool)) {
	s.el.OnVisibleChange(c)
}

type ModuleRow struct {
	*base.Row
	mk *markup.Markup
}

func NewSepRow(parent ui.ParentDrawable, mk *markup.Markup) *ModuleRow {
	return &ModuleRow{
		Row: base.NewRow(parent),
		mk:  mk,
	}
}

func (c *ModuleRow) Add(d ui.Drawable) {
	ctx := c.Parent().ChildContext(0)

	if !ctx.Has("module") {
		c.Row.Add(d)
		return
	}

	el := &sepRowEl{moduleBase: newBase(c), el: d}
	el.SetContext(c.ChildContext(len(c.Row.Children())))

	var err error
	el.Root, err = c.mk.Parse(el, el, ctx.MustString("module"))
	if err != nil {
		log.Fatalf("config: module: %s", err)
	}

	if el.Content == nil {
		log.Fatal(`config: module: missing container with ref="Content"`)
	}

	nctx := el.Content.ChildContext(0)
	el.el.SetContext(nctx)

	el.Content.Add(d)
	c.Row.Add(el)
}
