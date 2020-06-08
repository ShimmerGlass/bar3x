package tray

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
)

func xembedData(X *xgbutil.XUtil, win xproto.Window) (int, bool, error) {
	embedInfoAtom, err := xprop.Atom(X, "_XEMBED_INFO", false)
	if err != nil {
		return 0, false, err
	}

	xembedc := xproto.GetProperty(X.Conn(),
		false,
		win,
		embedInfoAtom,
		xproto.GetPropertyTypeAny,
		0,
		2*32,
	)
	xembedr, err := xembedc.Reply()
	if err != nil {
		return 0, false, err
	}

	if xembedr != nil && xembedr.Length > 0 {
		xembedVersion := xgb.Get32(xembedr.Value[:4])
		xembedFlags := xgb.Get32(xembedr.Value[4:])
		if xembedVersion > 1 {
			xembedVersion = 1
		}

		shouldMap := xembedFlags&1 > 0
		return int(xembedVersion), shouldMap, nil
	}

	return 0, false, nil
}
