package bar

import (
	"image"

	"github.com/BurntSushi/xgbutil/xgraphics"
)

const m = 1<<16 - 1

func DrawCopySrcRGBAToBGRA(dst *xgraphics.Image, r image.Rectangle, src *image.RGBA, sp image.Point) {
	r = r.Intersect(dst.Rect)

	dx, dy := r.Dx(), r.Dy()
	d0 := dst.PixOffset(r.Min.X, r.Min.Y)
	s0 := src.PixOffset(sp.X, sp.Y)
	var (
		ddelta, sdelta int
		i0, i1, idelta int
	)
	if r.Min.Y < sp.Y || r.Min.Y == sp.Y && r.Min.X <= sp.X {
		ddelta = dst.Stride
		sdelta = src.Stride
		i0, i1, idelta = 0, dx*4, +4
	} else {
		// If the source start point is higher than the destination start point, or equal height but to the left,
		// then we compose the rows in right-to-left, bottom-up order instead of left-to-right, top-down.
		d0 += (dy - 1) * dst.Stride
		s0 += (dy - 1) * src.Stride
		ddelta = -dst.Stride
		sdelta = -src.Stride
		i0, i1, idelta = (dx-1)*4, -4, -4
	}
	for ; dy > 0; dy-- {
		dpix := dst.Pix[d0:]
		spix := src.Pix[s0:]
		for i := i0; i != i1; i += idelta {
			s := spix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
			sr := s[0]
			sg := s[1]
			sb := s[2]
			sa := s[3]

			d := dpix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
			d[0] = sb
			d[1] = sg
			d[2] = sr
			d[3] = sa
		}
		d0 += ddelta
		s0 += sdelta
	}
}

func DrawCopyOverRGBAToBGRA(dst *xgraphics.Image, r image.Rectangle, src *image.RGBA, sp image.Point) {
	r = r.Intersect(dst.Rect)

	dx, dy := r.Dx(), r.Dy()
	d0 := dst.PixOffset(r.Min.X, r.Min.Y)
	s0 := src.PixOffset(sp.X, sp.Y)
	var (
		ddelta, sdelta int
		i0, i1, idelta int
	)
	if r.Min.Y < sp.Y || r.Min.Y == sp.Y && r.Min.X <= sp.X {
		ddelta = dst.Stride
		sdelta = src.Stride
		i0, i1, idelta = 0, dx*4, +4
	} else {
		// If the source start point is higher than the destination start point, or equal height but to the left,
		// then we compose the rows in right-to-left, bottom-up order instead of left-to-right, top-down.
		d0 += (dy - 1) * dst.Stride
		s0 += (dy - 1) * src.Stride
		ddelta = -dst.Stride
		sdelta = -src.Stride
		i0, i1, idelta = (dx-1)*4, -4, -4
	}
	for ; dy > 0; dy-- {
		dpix := dst.Pix[d0:]
		spix := src.Pix[s0:]
		for i := i0; i != i1; i += idelta {
			s := spix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
			sr := uint32(s[0]) * 0x101
			sg := uint32(s[1]) * 0x101
			sb := uint32(s[2]) * 0x101
			sa := uint32(s[3]) * 0x101

			// The 0x101 is here for the same reason as in drawRGBA.
			a := (m - sa) * 0x101

			d := dpix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
			d[0] = uint8((uint32(d[0])*a/m + sb) >> 8)
			d[1] = uint8((uint32(d[1])*a/m + sg) >> 8)
			d[2] = uint8((uint32(d[2])*a/m + sr) >> 8)
			d[3] = uint8((uint32(d[3])*a/m + sa) >> 8)
		}
		d0 += ddelta
		s0 += sdelta
	}
}
