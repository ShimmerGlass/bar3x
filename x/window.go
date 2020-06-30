package x

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
)

func WindowClass(x *xgbutil.XUtil, win xproto.Window) (string, error) {
	classRes, err := xproto.GetProperty(
		x.Conn(),
		false,
		win,
		MustAtom(x, "WM_CLASS"),
		xproto.AtomString,
		0,
		32,
	).Reply()
	if err != nil {
		return "", fmt.Errorf("x: get class: %w", err)
	}

	return string(classRes.Value), nil
}
