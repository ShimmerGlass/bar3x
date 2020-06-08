package ui

import (
	"image/draw"
)

type Drawable interface {
	Init() error
	SetContext(ctx Context)

	Width() int
	OnWidthChange(func(int))

	Height() int
	OnHeightChange(func(int))

	Visible() bool
	SetVisible(v bool)
	OnVisibleChange(func(bool))

	Draw(x, y int, im draw.Image)
	Notify()

	Context() Context
}

type ParentDrawable interface {
	Drawable

	Add(Drawable)
	Children() []Drawable
	ChildContext(index int) Context
}

type TextDrawable interface {
	Drawable
	SetText(string)
}

type notifier struct {
	parent Drawable
}

func (n notifier) Notify() {
	n.parent.Notify()
}
