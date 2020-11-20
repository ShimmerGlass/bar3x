package debug

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/shimmerglass/bar3x/bar"
)

func Run(addr string, bars *bar.Bars) error {
	http.Handle("/", devtoolsHandler(bars))
	return http.ListenAndServe(addr, nil)
}

func devtoolsHandler(bars *bar.Bars) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		view := createView(bars)
		rw.Write([]byte(view))
	})
}
