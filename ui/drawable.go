package ui

import (
	"image"
	"image/draw"
)

type EventType int

const (
	EventTypeLeftClick = iota
	EventTypeRightClick
	EventPointerMove
	EventPointerEnter
	EventPointerLeave
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

	OnLeftClick() func(Event) bool
	OnRightClick() func(Event) bool
	SetOnLeftClick(cb func(Event) bool)
	SetOnRightClick(cb func(Event) bool)
	OnPointerMove() func(Event) bool
	SetOnPointerMove(cb func(Event) bool)
	OnPointerEnter() func(Event) bool
	SetOnPointerEnter(cb func(Event) bool)
	OnPointerLeave() func(Event) bool
	SetOnPointerLeave(cb func(Event) bool)

	Draw(x, y int, im draw.Image)
	Notify()

	Context() Context
	Children() []Drawable
}

type ParentDrawable interface {
	Drawable

	Add(Drawable)
	Children() []Drawable
	ChildContext(index int) Context
}
