package x

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

type Screen struct {
	X, Y          int
	Width, Height int
	Outputs       []string
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
NextCrtc:
	for _, crtc := range resources.Crtcs {
		info, err := randr.GetCrtcInfo(X, crtc, 0).Reply()
		if err != nil {
			return nil, fmt.Errorf("x get crtc: %w", err)
		}
		if len(info.Outputs) == 0 {
			continue NextCrtc
		}

		outputNames := []string{}
		for _, out := range info.Outputs {
			output, err := randr.GetOutputInfo(X, out, 0).Reply()
			if err != nil {
				return nil, fmt.Errorf("x get output info: %w", err)
			}
			outputNames = append(outputNames, string(output.Name))
		}

		res = append(res, Screen{
			X:       int(info.X),
			Y:       int(info.Y),
			Width:   int(info.Width),
			Height:  int(info.Height),
			Outputs: outputNames,
		})
	}

	return res, nil
}
