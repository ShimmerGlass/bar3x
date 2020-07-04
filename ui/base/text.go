package base

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/segmentio/fasthash/fnv1a"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/cache"
)

type textCacheKey struct {
	font     uint64
	fontSize float64
	color    color.Color
	text     string
}

type Text struct {
	Base

	defaultColorKey    string
	defaultFontKey     string
	defaultFontSizeKey string

	fontKey     uint64
	setFont     string
	setFontSize float64
	setColor    color.Color
	text        string
	drawnText   string
	maxWidth    int
}

func NewText(p ui.ParentDrawable) *Text {
	return &Text{
		Base:               NewBase(p),
		defaultColorKey:    "text_color",
		defaultFontKey:     "text_font",
		defaultFontSizeKey: "text_font_size",
	}
}

func NewIcon(p ui.ParentDrawable) *Text {
	return &Text{
		Base:               NewBase(p),
		defaultColorKey:    "accent_color",
		defaultFontKey:     "icon_font",
		defaultFontSizeKey: "icon_font_size",
	}
}

func (t *Text) Init() error {
	// this prevents to change to font after the text element is created
	// but this is an unlikely usage and allows for cache optimizations
	font := t.Font()

	// we spend a long time hashing font base64 for modules such as Cmd
	// that create Text element often. Hashing the first 512 bytes seems
	// to work
	if len(font) > 512 {
		font = font[:512]
	}
	t.fontKey = fnv1a.HashString64(font)

	return nil
}

func (t *Text) Text() string {
	return t.text
}
func (t *Text) SetText(s string) {
	old := t.text
	t.text = s
	if s != old {
		t.updateSize()
	}
}

func (t *Text) Font() string {
	if t.setFont != "" {
		return t.setFont
	}
	return t.ctx.MustString(t.defaultFontKey)
}
func (t *Text) SetFont(f string) {
	t.setFont = f
	t.updateSize()
}

func (t *Text) FontSize() float64 {
	if t.setFontSize > 0 {
		return t.setFontSize
	}
	return t.ctx.MustFloat(t.defaultFontSizeKey)
}
func (t *Text) SetFontSize(f float64) {
	t.setFontSize = f
	t.updateSize()
}

func (t *Text) Color() color.Color {
	if t.setColor != nil {
		return t.setColor
	}
	return t.ctx.MustColor(t.defaultColorKey)
}
func (t *Text) SetColor(c color.Color) {
	t.setColor = c
}

func (t *Text) MaxWidth() int {
	return t.maxWidth
}
func (t *Text) SetMaxWidth(i int) {
	t.maxWidth = i
	t.updateSize()
}

func (t *Text) updateSize() {
	fontSize := t.FontSize()

	if t.text == "" || fontSize == 0 {
		t.width.Set(0)
		t.height.Set(0)
		return
	}

	w, h, drawnText, err := textSize(t.text, textOptions{
		font:     t.Font(),
		fontKey:  t.fontKey,
		fontSize: fontSize,
		maxWidth: t.maxWidth,
	})
	if err != nil {
		log.Fatal(err)
	}
	t.width.Set(w)
	t.height.Set(h)
	t.drawnText = drawnText
}

func (t *Text) Draw(x, y int, im draw.Image) {
	k := textCacheKey{
		font:     t.fontKey,
		fontSize: t.FontSize(),
		color:    t.Color(),
		text:     t.drawnText,
	}
	w, h := t.width.V, t.height.V

	cache.Draw(k, w, h, x, y, im, t.draw)
}

func (t *Text) draw(im draw.Image) {
	w, h := t.width.V, t.height.V
	if w == 0 || h == 0 {
		return
	}

	img, err := renderText(t.drawnText, w, h, textOptions{
		font:     t.Font(),
		fontKey:  t.fontKey,
		fontSize: t.FontSize(),
		color:    t.Color(),
	})
	if err != nil {
		log.Fatal(err)
	}

	draw.Draw(
		im,
		image.Rect(0, 0, w, h),
		img,
		image.Point{},
		draw.Over,
	)
}
