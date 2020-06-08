package x

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

type Screen struct {
	X, Y          int
	Width, Height int
}

func Screens(X *xgb.Conn) ([]Screen, error) {
	// Every extension must be initialized before it can be used.
	err := randr.Init(X)
	if err != nil {
		return nil, err
	}
	// Get the root window on the default screen.
	root := xproto.Setup(X).DefaultScreen(X).Root

	// Gets the current screen resources. Screen resources contains a list
	// of names, crtcs, outputs and modes, among other things.
	resources, err := randr.GetScreenResources(X, root).Reply()
	if err != nil {
		return nil, err
	}

	res := []Screen{}

	// Iterate through all of the crtcs and show some of their info.
	for _, crtc := range resources.Crtcs {
		info, err := randr.GetCrtcInfo(X, crtc, 0).Reply()
		if err != nil {
			return nil, err
		}

		res = append(res, Screen{
			X:      int(info.X),
			Y:      int(info.Y),
			Width:  int(info.Width),
			Height: int(info.Height),
		})
	}

	return res, nil
}
