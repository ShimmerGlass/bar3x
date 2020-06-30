package cache

import (
	"image"
	"image/draw"

	lru "github.com/hashicorp/golang-lru"
)

func init() {
	cache, _ = lru.New(2048)
}

var cache *lru.Cache

func Draw(key interface{}, w, h, x, y int, im draw.Image, id func(draw.Image)) {
	bounds := image.Rect(0, 0, w, h)
	v, ok := cache.Get(key)
	if !ok {
		iim := image.NewRGBA(bounds)
		id(iim)
		cache.Add(key, iim)
		v = iim
	}
	cim := v.(*image.RGBA)
	draw.Draw(im, bounds.Add(image.Pt(x, y)), cim, image.Point{}, draw.Over)
}
