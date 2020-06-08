package ui

import (
	"fmt"
	"image/color"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	keyNotFoundFmt      = "config: key %q not found, please add it to your config file"
	keyBadTypeFmt       = "config: bad type for key %q: %T, expected %s"
	colorParseFormatFmt = "config: bad format for color %q, expected #RRGGBB(AA)"
)

type Context map[string]interface{}

func (c Context) New(vals Context) Context {
	n := Context{}
	for k, v := range c {
		n[k] = v
	}
	if vals != nil {
		for k, v := range vals {
			n[k] = v
		}
	}
	return n
}

func (c Context) Has(k string) bool {
	_, ok := c[k]
	return ok
}

func (c Context) MustString(n string) string {
	if !c.Has(n) {
		log.Fatalf(keyNotFoundFmt, n)
	}
	v, ok := c[n].(string)
	if !ok {
		log.Fatalf(keyBadTypeFmt, n, c[n], "string")
	}
	return v
}

func (c Context) MustInt(n string) int {
	if !c.Has(n) {
		log.Fatalf(keyNotFoundFmt, n)
	}
	v, ok := c[n].(int)
	if !ok {
		log.Fatalf(keyBadTypeFmt, n, c[n], "int")
	}
	return v
}

func (c Context) MustFloat(n string) float64 {
	if !c.Has(n) {
		log.Fatalf(keyNotFoundFmt, n)
	}
	v, ok := c[n].(float64)
	if !ok {
		log.Fatalf(keyBadTypeFmt, n, c[n], "float64")
	}
	return v
}

func (c Context) MustColor(n string) color.Color {
	if !c.Has(n) {
		log.Fatalf(keyNotFoundFmt, n)
	}

	v := c[n]

	switch val := v.(type) {
	case color.Color:
		return val
	case string:
		c, err := ParseColor(val)
		if err != nil {
			log.Fatalf("could not parse key %q (%q) as color: %s", n, v, err)
		}

		return c
	default:
		log.Fatalf(keyBadTypeFmt, n, v, "color")
	}

	return nil
}

func ParseColor(s string) (color.Color, error) {
	if len(s) != 7 && len(s) != 9 {
		return nil, fmt.Errorf(colorParseFormatFmt, s)
	}

	if s[0] != '#' {
		return nil, fmt.Errorf(colorParseFormatFmt, s)
	}

	r, err := strconv.ParseInt(s[1:3], 16, 64)
	if err != nil {
		return nil, fmt.Errorf(colorParseFormatFmt+":%w", s, err)
	}
	g, err := strconv.ParseInt(s[3:5], 16, 64)
	if err != nil {
		return nil, fmt.Errorf(colorParseFormatFmt+":%w", s, err)
	}
	b, err := strconv.ParseInt(s[5:7], 16, 64)
	if err != nil {
		return nil, fmt.Errorf(colorParseFormatFmt+":%w", s, err)
	}
	a := int64(255)
	if len(s) == 9 {
		a, err = strconv.ParseInt(s[7:9], 16, 64)
		if err != nil {
			return nil, fmt.Errorf(colorParseFormatFmt+":%w", s, err)
		}
	}

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}, nil
}
