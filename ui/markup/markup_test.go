package markup

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"github.com/shimmerglass/bar3x/ui"
	"github.com/stretchr/testify/require"
)

type TestDrawable struct {
	text      string
	Str       string
	someVal   string
	someValCb func(string)
	Color     color.Color
	Inner     []ui.Drawable
}

func (t *TestDrawable) SetSomeVal(i string) {
	t.someVal = i
	if t.someValCb != nil {
		t.someValCb(i)
	}
}

func (t *TestDrawable) SomeVal() string {
	return t.someVal
}

func (t *TestDrawable) OnSomeValChange(cb func(string)) {
	t.someValCb = cb
	cb(t.someVal)
}

func (t *TestDrawable) Size() image.Rectangle {
	return image.Rectangle{}
}

func (t *TestDrawable) Draw(x, y int, im draw.Image) {}
func (t *TestDrawable) Notify()                      {}
func (t *TestDrawable) Visible() bool {
	return true
}
func (t *TestDrawable) Show() {}
func (t *TestDrawable) Hide() {}
func (t *TestDrawable) Add(d ui.Drawable) {
	t.Inner = append(t.Inner, d)
}
func (t *TestDrawable) SetText(s string) {
	t.text = s
}

type refs struct {
	Root *TestDrawable
	Test *TestDrawable
}

func TestSimple(t *testing.T) {
	m := New()
	m.Register("test", func(parent ui.ParentDrawable) ui.Drawable {
		return &TestDrawable{}
	})

	r := &refs{}
	el, err := m.Parse(
		nil,
		map[string]interface{}{
			"color":  color.RGBA{0xff, 0, 0, 0xff},
			"number": 2,
		},
		r,
		`
			<test
				ref="Root"
				SomeVal="test_someval"
			>
				<test SomeVal='{$Root.SomeVal + "_hey"}'>Hey</test>
				<test ref="Test" />
			</test>
		`,
	)
	require.Nil(t, err)

	time.Sleep(100 * time.Millisecond)

	require.Equal(t, el, &TestDrawable{
		someVal: "test_someval",
		Inner: []ui.Drawable{
			&TestDrawable{text: "Hey", someVal: "test_someval_hey"},
			&TestDrawable{},
		},
	})

	require.Equal(t, r, &refs{
		Test: &TestDrawable{Str: "lul"},
	})
}
