package base

// #cgo pkg-config: cairo
// #cgo pkg-config: freetype2
// #include <cairo.h>
// #include <cairo-ft.h>
// #include <ft2build.h>
// #include FT_SFNT_NAMES_H
// #include FT_FREETYPE_H
// #include FT_GLYPH_H
// #include FT_OUTLINE_H
// #include FT_BBOX_H
// #include FT_TYPE1_TABLES_H
import "C"

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"sync"
	"unicode/utf8"
	"unsafe"
)

type textOptions struct {
	font     string
	fontKey  uint64
	fontSize float64
	maxWidth int
	color    color.Color
}

type faceCacheEntry struct {
	backingArray []byte
	face         C.FT_Face
}

var faceCacheLock sync.Mutex
var faceCache = map[uint64]faceCacheEntry{}

const heightChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func textSize(text string, opt textOptions) (int, int, string, error) {
	surface := C.cairo_image_surface_create(C.cairo_format_t(0), C.int(0), C.int(0))
	defer C.cairo_surface_destroy(surface)
	ctx := C.cairo_create(surface)
	defer C.cairo_destroy(ctx)

	err := setCtxFont(ctx, opt)
	if err != nil {
		return 0, 0, "", err
	}
	C.cairo_set_font_size(ctx, C.double(opt.fontSize))

	drawnText := text
	w, _ := measureText(ctx, text, opt)
	_, h := measureText(ctx, heightChars, opt)

	max := float64(opt.maxWidth)
	txt := text
	for opt.maxWidth > 0 && w > max {
		_, s := utf8.DecodeLastRuneInString(txt)
		txt = strings.TrimSpace(txt[:len(txt)-s])
		drawnText = txt + "…"
		w, _ = measureText(ctx, drawnText, opt)
	}

	return int(math.Ceil(w)), int(math.Ceil(h)), drawnText, nil
}

func renderText(text string, w, h int, opt textOptions) (image.Image, error) {
	rgba := color.RGBAModel.Convert(opt.color).(color.RGBA)
	surface := C.cairo_image_surface_create(C.cairo_format_t(0), C.int(w), C.int(h))
	defer C.cairo_surface_destroy(surface)
	ctx := C.cairo_create(surface)
	defer C.cairo_destroy(ctx)

	err := setCtxFont(ctx, opt)
	if err != nil {
		return nil, err
	}

	C.cairo_set_font_size(ctx, C.double(opt.fontSize))
	C.cairo_set_source_rgba(ctx,
		C.double(float64(rgba.R)/255),
		C.double(float64(rgba.G)/255),
		C.double(float64(rgba.B)/255),
		C.double(float64(rgba.A)/255),
	)

	C.cairo_move_to(ctx, 0, C.double(opt.fontSize))
	cs := C.CString(text)
	C.cairo_show_text(ctx, cs)
	C.free(unsafe.Pointer(cs))

	C.cairo_surface_flush(surface)
	dataPtr := C.cairo_image_surface_get_data(surface)
	if dataPtr == nil {
		return nil, fmt.Errorf("cairo: can't access surface pixel data")
	}
	stride := C.cairo_image_surface_get_stride(surface)
	height := C.cairo_image_surface_get_height(surface)
	width := C.cairo_image_surface_get_width(surface)
	data := C.GoBytes(unsafe.Pointer(dataPtr), stride*height)

	return &BGRA{
		Pix:    data,
		Stride: int(stride),
		Rect:   image.Rect(0, 0, int(width), int(height)),
	}, nil
}

func setCtxFont(ctx *C.cairo_t, opt textOptions) error {
	if !strings.HasPrefix(opt.font, "base64:") {
		s := C.CString(opt.font)
		defer C.free(unsafe.Pointer(s))
		C.cairo_select_font_face(ctx, s, C.cairo_font_slant_t(0), C.cairo_font_weight_t(0))
		return nil
	}

	fontB64 := opt.font[7:]
	faceCacheLock.Lock()
	defer faceCacheLock.Unlock()

	var cf *C.cairo_font_face_t
	if e, ok := faceCache[opt.fontKey]; ok {
		cf = C.cairo_ft_font_face_create_for_ft_face(e.face, 0)
	} else {
		dec, err := base64.StdEncoding.DecodeString(fontB64)
		if err != nil {
			return fmt.Errorf("text: could not decode base64 font")
		}

		var lib C.FT_Library
		if status := C.FT_Init_FreeType(&lib); status != 0 {
			return fmt.Errorf("text: could not open FreeType library: %d", status)
		}

		var face C.FT_Face
		if status := C.FT_New_Memory_Face(lib, (*C.FT_Byte)(&dec[0]), C.FT_Long(len(dec)), 0, &face); status != 0 {
			return fmt.Errorf("text: error creating face: %d", status)
		}
		faceCache[opt.fontKey] = faceCacheEntry{
			backingArray: dec, // keep a ref to font data so it doesnt get garbage collected by go
			face:         face,
		}
		cf = C.cairo_ft_font_face_create_for_ft_face(face, 0)
	}

	C.cairo_set_font_face(ctx, cf)
	return nil
}

func measureText(ctx *C.cairo_t, text string, opt textOptions) (float64, float64) {
	cte := C.cairo_text_extents_t{}
	cs := C.CString(text)
	defer C.free(unsafe.Pointer(cs))
	C.cairo_text_extents(ctx, cs, &cte)

	return float64(cte.x_advance), opt.fontSize + float64(cte.height) + float64(cte.y_bearing)
}

var (
	tst          uint32 = 1
	littleEndian bool   = (*[4]byte)(unsafe.Pointer(&tst))[0] == 1
)

// NewBGRA returns a new BGRA with the given bounds.
func NewBGRA(r image.Rectangle) *BGRA {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 4*w*h)
	return &BGRA{Pix: buf, Stride: 4 * w, Rect: r}
}

// BGRA is an in-memory image whose At method returns BGRAColor values.
type BGRA struct {
	// Pix holds the image's pixels, in B, G, R, A order on small endian systems
	// and A, R, G, B on big endian systems.
	// See http://cairographics.org/manual/cairo-Image-Surfaces.html#cairo-format-t
	// The pixel at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (self *BGRA) ColorModel() color.Model {
	return BGRAColorModel
}

func (self *BGRA) Bounds() image.Rectangle {
	return self.Rect
}

func (self *BGRA) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(self.Rect)) {
		return BGRAColor{}
	}
	i := self.PixOffset(x, y)
	if littleEndian {
		return BGRAColor{B: self.Pix[i+0], G: self.Pix[i+1], R: self.Pix[i+2], A: self.Pix[i+3]}
	} else {
		return BGRAColor{A: self.Pix[i+0], R: self.Pix[i+1], G: self.Pix[i+2], B: self.Pix[i+3]}
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (self *BGRA) PixOffset(x, y int) int {
	return (y-self.Rect.Min.Y)*self.Stride + (x-self.Rect.Min.X)*4
}

func (self *BGRA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(self.Rect)) {
		return
	}
	i := self.PixOffset(x, y)
	c1 := BGRAColorModel.Convert(c).(BGRAColor)
	if littleEndian {
		self.Pix[i+0] = c1.B
		self.Pix[i+1] = c1.G
		self.Pix[i+2] = c1.R
		self.Pix[i+3] = c1.A
	} else {
		self.Pix[i+0] = c1.A
		self.Pix[i+1] = c1.R
		self.Pix[i+2] = c1.G
		self.Pix[i+3] = c1.B
	}
}

var BGRAColorModel = color.ModelFunc(
	func(c color.Color) color.Color {
		if _, ok := c.(BGRAColor); ok {
			return c
		}
		r, g, b, a := c.RGBA()
		return BGRAColor{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	},
)

// BGRAColor represents a traditional 32-bit alpha-premultiplied color,
// having 8 bits for each of alpha, red, green and blue.
type BGRAColor struct {
	B, G, R, A uint8
}

func (c BGRAColor) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8

	g = uint32(c.G)
	g |= g << 8

	b = uint32(c.B)
	b |= b << 8

	a = uint32(c.A)
	a |= a << 8

	return r, g, b, a
}
