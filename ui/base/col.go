package base

import (
	"image"
	"image/draw"

	"github.com/shimmerglass/bar3x/ui"
)

var _ ui.ParentDrawable = (*Col)(nil)

type Col struct {
	Base
	inner []ui.Drawable
	align string
}

func NewCol(p ui.ParentDrawable) *Col {
	return &Col{
		Base:  NewBase(p),
		align: HAlighCenter,
	}
}

func (c *Col) Add(d ui.Drawable) {
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

func (c *Col) Align() string {
	return c.align
}
func (c *Col) SetAlign(v string) {
	c.align = v
}

func (c *Col) SetContext(ctx ui.Context) {
	c.ctx = ctx
	c.updateChildrenCtx()
}

func (c *Col) ChildContext(i int) ui.Context {
	return c.ctx.New(ui.Context{
		"index":            i,
		"visible_index":    i,
		"is_first":         false,
		"is_last":          false,
		"is_first_visible": false,
		"is_last_visible":  false,
	})
}

func (c *Col) Children() []ui.Drawable {
	return c.inner
}

func (c *Col) SendEvent(ev ui.Event) bool {
	if !c.Base.SendEvent(ev) {
		return false
	}

	y := 0
	for _, i := range c.inner {
		if !i.Visible() {
			continue
		}
		w, h := i.Width(), i.Height()
		x := (c.width.V - w) / 2

		if ev.At.In(image.Rect(x, y, x+w, y+h)) {
			iev := ev
			iev.At = ev.At.Sub(image.Pt(x, y))
			i.SendEvent(iev)
		}

		y += h
	}

	return true
}

func (c *Col) updateSize() {
	var w, h int
	for _, i := range c.inner {
		if !i.Visible() {
			continue
		}
		ew, eh := i.Width(), i.Height()
		h += eh
		if ew > w {
			w = ew
		}
	}
	c.width.Set(w)
	c.height.Set(h)
}

func (r *Col) updateChildrenCtx() {
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

func (c *Col) Draw(x, y int, im draw.Image) {
	xOff, yOff := 0, 0
	for _, i := range c.inner {
		if !i.Visible() {
			continue
		}
		switch c.align {
		case HAlignLeft:
			xOff = 0
		case HAlighCenter:
			xOff = (c.width.V - i.Width()) / 2
		case VAlignBottom:
			xOff = c.width.V - i.Width()
		}
		i.Draw(x+xOff, y+yOff, im)
		yOff += i.Height()
	}
}
