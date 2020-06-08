package ui

import (
	"fmt"
	"strings"
)

func Print(d Drawable, prop string) {
	print(0, d, prop)
}

func print(indent int, d Drawable, prop string) {
	if d == nil {
		return
	}
	fmt.Printf("%s%T: %+v\n", strings.Repeat("\t", indent), d, d.Context()[prop])

	if p, ok := d.(ParentDrawable); ok {
		for _, c := range p.Children() {
			print(indent+1, c, prop)
		}
	}
}
