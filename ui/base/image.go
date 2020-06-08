package base

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
	"github.com/shimmerglass/bar3x/ui"
)

type Image struct {
	Base

	setWidth  int
	setHeight int

	path string

	img  image.Image
	rimg image.Image
}

func NewImage(p ui.ParentDrawable) *Image {
	return &Image{
		Base:      NewBase(p),
		setHeight: -1,
		setWidth:  -1,
	}
}

func (i *Image) SetWidth(v int) {
	i.setWidth = v
	i.width.set(v)
	i.updateImg()
}

func (i *Image) SetHeight(v int) {
	i.setHeight = v
	i.height.set(v)
	i.updateImg()
}

func (i *Image) Image() image.Image {
	return i.img
}
func (i *Image) SetImage(v image.Image) {
	if v == i.img {
		return
	}

	i.img = v
	i.rimg = nil
	r := image.ZR
	if v != nil {
		r = v.Bounds()
	}

	if i.setWidth == -1 {
		i.width.set(r.Dx())
	}
	if i.setHeight == -1 {
		i.height.set(r.Dy())
	}

	i.updateImg()
}

func (i *Image) Path() string {
	return i.path
}
func (i *Image) SetPath(p string) {
	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		panic(fmt.Errorf("cannot decode image at %q: %w", p, err))
	}

	i.SetImage(img)
}

func (i *Image) updateImg() {
	if i.img == nil {
		return
	}

	isize := i.img.Bounds()
	ww, wh := isize.Dx(), isize.Dy()
	if i.setWidth != -1 {
		ww = i.setWidth
	}
	if i.setHeight != -1 {
		wh = i.setHeight
	}

	if i.rimg != nil && ww == isize.Dx() && wh == isize.Dy() {
		return
	}

	i.rimg = resize.Resize(uint(ww), uint(wh), i.img, resize.Lanczos3)
}

func (i *Image) Draw(x, y int, im draw.Image) {
	if i.rimg == nil {
		return
	}
	draw.Draw(
		im,
		image.Rect(
			x, y,
			im.Bounds().Dx(),
			im.Bounds().Dy(),
		),
		i.rimg,
		image.ZP,
		draw.Over,
	)
}
