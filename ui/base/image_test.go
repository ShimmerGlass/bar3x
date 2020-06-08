package base

import (
	"image"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageDefault(t *testing.T) {
	img := NewImage(nil)
	img.SetImage(image.NewRGBA(image.Rect(0, 0, 10, 20)))
	require.Equal(t, 10, img.Width())
	require.Equal(t, 20, img.Height())
}

func TestImageSetSize(t *testing.T) {
	img := NewImage(nil)
	img.SetWidth(30)
	img.SetHeight(40)
	img.SetImage(image.NewRGBA(image.Rect(0, 0, 10, 20)))
	require.Equal(t, 30, img.Width())
	require.Equal(t, 40, img.Height())
}
