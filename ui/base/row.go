package base

import (
	"image/draw"

	"github.com/shimmerglass/bar3x/ui"
)

type Row struct {
	Base
	inner []ui.Drawable
	align string
}

func NewRow(p ui.ParentDrawable) *Row {
	return &Row{
		Base:  NewBase(p),
		align: VAlighMiddle,
	}
}

func (c *Row) Add(d ui.Drawable) {
	c.inner = append(c.inner, d)
	d.OnWidthChange(func(int) {
		c.updateSize()
	})
	d.OnHeightChange(func(int) {
		c.updateSize()
	})
	d.OnVisibleChange(func(bool) {
		c.updateSize()
		c.updateChildrenCtx()
	})
}

func (c *Row) ChildContext(i int) ui.Context {
	return c.ctx.New(ui.Context{
		"index":            i,
		"visible_index":    i,
		"is_first":         false,
		"is_last":          false,
		"is_first_visible": false,
		"is_last_visible":  false,
	})
}

func (c *Row) SetContext(ctx ui.Context) {
	c.ctx = ctx
	c.updateChildrenCtx()
}

func (c *Row) Children() []ui.Drawable {
	return c.inner
}

func (c *Row) Align() string {
	return c.align
}
func (c *Row) SetAlign(v string) {
	c.align = v
}

func (r *Row) updateSize() {
	var w, h int
	for _, i := range r.inner {
		if !i.Visible() {
			continue
		}
		ew, eh := i.Width(), i.Height()
		w += ew
		if eh > h {
			h = eh
		}
	}

	r.width.Set(w)
	r.height.Set(h)
}

func (r *Row) updateChildrenCtx() {
	firstVisible := len(r.inner)
	lastVisible := 0

	for i, el := range r.inner {
		if !el.Visible() {
			continue
		}
		if i < firstVisible {
			firstVisible = i
		}
		if i > lastVisible {
			lastVisible = i
		}
	}

	vi := 0
	for i, el := range r.inner {
		el.SetContext(r.ctx.New(ui.Context{
			"index":            i,
			"visible_index":    vi,
			"is_first":         i == 0,
			"is_last":          i == len(r.inner)-1,
			"is_first_visible": i == firstVisible,
			"is_last_visible":  i == lastVisible,
		}))
		if el.Visible() {
			vi++
		}
	}
}

func (r Row) Draw(x, y int, im draw.Image) {
	xOff, yOff := 0, 0
	for _, i := range r.inner {
		if !i.Visible() {
			continue
		}
		switch r.align {
		case VAlignTop:
			yOff = 0
		case VAlighMiddle:
			yOff = (r.height.V - i.Height()) / 2
		case VAlignBottom:
			yOff = r.height.V - i.Height()
		}
		i.Draw(x+xOff, y+yOff, im)
		xOff += i.Width()
	}
}
