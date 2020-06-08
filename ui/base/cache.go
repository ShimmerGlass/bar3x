package base

import (
	"image"
	"image/draw"

	"github.com/shimmerglass/bar3x/ui"
)

var cache = map[interface{}]draw.Image{}

var CacheImageFn = func(s image.Rectangle) draw.Image {
	return image.NewRGBA(s)
}

var CacheCopyFn = func(dst draw.Image, r image.Rectangle, src draw.Image) {
	draw.Draw(dst, r, src, image.ZP, draw.Over)
}

type cached interface {
	key() interface{}
	ui.Drawable
}

type cacheEl struct {
	cached cached
}

func (c *cacheEl) Draw(x, y int, im draw.Image) {
	key := c.cached.key()
	i, ok := cache[key]
	if !ok {
		i = CacheImageFn(image.Rect(0, 0, c.cached.Width(), c.cached.Height()))
		c.cached.Draw(0, 0, i)
		cache[key] = i
	}

	CacheCopyFn(im, i.Bounds().Add(image.Pt(x, y)), i)
}

func (c *cacheEl) Notify() {
	c.cached.Notify()
}

func (c *cacheEl) Width() int {
	return c.cached.Width()
}
func (c *cacheEl) OnWidthChange(cb func(int)) {
	c.cached.OnWidthChange(cb)
}

func (c *cacheEl) Height() int {
	return c.cached.Height()
}

func (c *cacheEl) OnHeightChange(cb func(int)) {
	c.cached.OnHeightChange(cb)
}

func (c *cacheEl) Visible() bool {
	return c.cached.Visible()
}
func (c *cacheEl) SetVisible(v bool) {
	c.cached.SetVisible(v)
}
func (c *cacheEl) OnVisibleChange(cb func(bool)) {
	c.cached.OnVisibleChange(cb)
}
