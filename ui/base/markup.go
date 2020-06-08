package base

import (
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/markup"
)

func RegisterMarkup(mk *markup.Markup) {
	mk.Register("Bar", func(p ui.ParentDrawable) ui.Drawable {
		return NewBar(p)
	})
	mk.Register("Col", func(p ui.ParentDrawable) ui.Drawable {
		return NewCol(p)
	})
	mk.Register("Graph", func(p ui.ParentDrawable) ui.Drawable {
		return NewGraph(p)
	})
	mk.Register("Image", func(p ui.ParentDrawable) ui.Drawable {
		return NewImage(p)
	})
	mk.Register("Layers", func(p ui.ParentDrawable) ui.Drawable {
		return NewLayers(p)
	})
	mk.Register("Rect", func(p ui.ParentDrawable) ui.Drawable {
		return NewRect(p)
	})
	mk.Register("Row", func(p ui.ParentDrawable) ui.Drawable {
		return NewRow(p)
	})
	mk.Register("Sizer", func(p ui.ParentDrawable) ui.Drawable {
		return NewSizer(p)
	})
	mk.Register("Text", func(p ui.ParentDrawable) ui.Drawable {
		return NewText(p)
	})
	mk.Register("Icon", func(p ui.ParentDrawable) ui.Drawable {
		return NewIcon(p)
	})
	mk.Register("SeparatorArrow", func(p ui.ParentDrawable) ui.Drawable {
		return NewSeparatorArrow(p)
	})
}
