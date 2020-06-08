package base

import (
	"github.com/shimmerglass/bar3x/ui"
)

type watchInt struct {
	v   int
	cbs []func(int)
}

func (w *watchInt) add(c func(int)) {
	w.cbs = append(w.cbs, c)
	c(w.v)
}

func (w *watchInt) set(v int) {
	if w.v == v {
		return
	}
	w.v = v
	for _, c := range w.cbs {
		c(v)
	}
}

type watchBool struct {
	v   bool
	cbs []func(bool)
}

func (w *watchBool) add(c func(bool)) {
	w.cbs = append(w.cbs, c)
	c(w.v)
}

func (w *watchBool) set(v bool) {
	if w.v == v {
		return
	}
	w.v = v
	for _, c := range w.cbs {
		c(v)
	}
}

type Base struct {
	parent  ui.ParentDrawable
	visible *watchBool
	width   *watchInt
	height  *watchInt
	ctx     ui.Context
}

func NewBase(p ui.ParentDrawable) Base {
	return Base{
		parent:  p,
		visible: &watchBool{v: true},
		width:   &watchInt{},
		height:  &watchInt{},
	}
}

func (b *Base) Init() error {
	return nil
}
func (b *Base) SetContext(ctx ui.Context) {
	b.ctx = ctx
}

func (b *Base) Parent() ui.ParentDrawable {
	return b.parent
}
func (b *Base) Context() ui.Context {
	return b.ctx
}

func (b Base) Notify() {
	b.parent.Notify()
}

func (b *Base) Width() int {
	return b.width.v
}
func (b *Base) OnWidthChange(c func(int)) {
	b.width.add(c)
}

func (b *Base) Height() int {
	return b.height.v
}
func (b *Base) OnHeightChange(c func(int)) {
	b.height.add(c)
}

func (b Base) Visible() bool {
	return b.visible.v
}
func (b *Base) SetVisible(v bool) {
	b.visible.set(v)
}
func (b *Base) OnVisibleChange(c func(bool)) {
	b.visible.add(c)
}
