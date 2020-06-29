package base

import (
	"github.com/shimmerglass/bar3x/ui"
)

type WatchInt struct {
	V   int
	cbs []func(int)
}

func (w *WatchInt) Add(c func(int)) {
	w.cbs = append(w.cbs, c)
	c(w.V)
}

func (w *WatchInt) Set(v int) {
	if w.V == v {
		return
	}
	w.V = v
	for _, c := range w.cbs {
		c(v)
	}
}

type WatchBool struct {
	V   bool
	cbs []func(bool)
}

func (w *WatchBool) Add(c func(bool)) {
	w.cbs = append(w.cbs, c)
	c(w.V)
}

func (w *WatchBool) Set(v bool) {
	if w.V == v {
		return
	}
	w.V = v
	for _, c := range w.cbs {
		c(v)
	}
}

type Base struct {
	parent  ui.ParentDrawable
	visible *WatchBool
	width   *WatchInt
	height  *WatchInt
	ctx     ui.Context
}

func NewBase(p ui.ParentDrawable) Base {
	return Base{
		parent:  p,
		visible: &WatchBool{V: true},
		width:   &WatchInt{},
		height:  &WatchInt{},
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
	return b.width.V
}
func (b *Base) OnWidthChange(c func(int)) {
	b.width.Add(c)
}

func (b *Base) Height() int {
	return b.height.V
}
func (b *Base) OnHeightChange(c func(int)) {
	b.height.Add(c)
}

func (b Base) Visible() bool {
	return b.visible.V
}
func (b *Base) SetVisible(v bool) {
	b.visible.Set(v)
}
func (b *Base) OnVisibleChange(c func(bool)) {
	b.visible.Add(c)
}
