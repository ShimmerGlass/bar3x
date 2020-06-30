package ui

import (
	"image"
	"image/draw"
)

type EventType int

const (
	EventTypeLeftClick = iota
	EventTypeRightClick
)

type Event struct {
	Type EventType
	At   image.Point
}

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

	SendEvent(Event) bool

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
