package x

import (
	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
)

func MustAtom(x *xgbutil.XUtil, name string) xproto.Atom {
	atom, err := xprop.Atom(x, name, false)
	if err != nil {
		log.Fatalf("cannot get atom %s: %s", name, err)
	}

	return atom
}
