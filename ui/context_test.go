package ui

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseColor(t *testing.T) {
	cases := map[string]color.Color{
		"#012345":   color.RGBA{0x01, 0x23, 0x45, 0xff},
		"#01234567": color.RGBA{0x01, 0x23, 0x45, 0x67},
	}

	for in, out := range cases {
		v, err := ParseColor(in)
		require.NoError(t, err)
		require.Equal(t, out, v)
	}
}
