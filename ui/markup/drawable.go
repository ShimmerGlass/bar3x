package markup

import (
	"context"
	"image/draw"

	"github.com/PaesslerAG/gval"
	"github.com/shimmerglass/bar3x/ui"
	log "github.com/sirupsen/logrus"
)

type property struct {
	field *field
	expr  gval.Evaluable
}

type ctxProp struct {
	name string
	expr gval.Evaluable
}

type MarkupDrawable struct {
	parent     ui.Drawable
	inner      ui.Drawable
	properties []property
	ctxProp    []ctxProp
}

func newMarkupDrawable(p ui.Drawable, inner ui.Drawable) *MarkupDrawable {
	return &MarkupDrawable{
		parent: p,
		inner:  inner,
	}
}

func (b *MarkupDrawable) Init() error {
	return b.inner.Init()
}

func (b *MarkupDrawable) SetContext(ctx ui.Context) {
	for _, prop := range b.ctxProp {
		v, err := prop.expr(context.Background(), ctx)
		if err != nil {
			log.Fatal(err)
		}

		ctx = ctx.New(ui.Context{
			prop.name: v,
		})
	}

	b.inner.SetContext(ctx)

	for _, p := range b.properties {
		v, err := p.expr(context.Background(), ctx)
		if err != nil {
			log.Fatal(err)
		}

		err = p.field.Set(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (b *MarkupDrawable) Context() ui.Context {
	return b.inner.Context()
}

func (b *MarkupDrawable) Add(d ui.Drawable) {}

func (b *MarkupDrawable) ChildContext(int) ui.Context {
	return b.inner.Context()
}
func (b *MarkupDrawable) Children() []ui.Drawable {
	return []ui.Drawable{b.inner}
}

func (b *MarkupDrawable) Notify() {
	b.parent.Notify()
}

func (b *MarkupDrawable) Width() int {
	return b.inner.Width()
}
func (b *MarkupDrawable) OnWidthChange(c func(int)) {
	b.inner.OnWidthChange(c)
}

func (b *MarkupDrawable) Height() int {
	return b.inner.Height()
}
func (b *MarkupDrawable) OnHeightChange(c func(int)) {
	b.inner.OnHeightChange(c)
}

func (b *MarkupDrawable) Visible() bool {
	return b.inner.Visible()
}
func (b *MarkupDrawable) SetVisible(v bool) {
	b.inner.SetVisible(v)
}
func (b *MarkupDrawable) OnVisibleChange(c func(bool)) {
	b.inner.OnVisibleChange(c)
}

func (b *MarkupDrawable) Draw(x, y int, im draw.Image) {
	b.inner.Draw(x, y, im)
}

func (b *MarkupDrawable) SendEvent(ev ui.Event) bool {
	return b.inner.SendEvent(ev)
}
