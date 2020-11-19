// Example draw-text shows how to draw text to an xgraphics.Image type.
package main

import (
	"flag"
	"fmt"
	"io"
	"log/syslog"
	"os"
	"os/exec"
	"time"

	"github.com/TheCreeper/go-notify"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/smallnest/ringbuffer"

	"net/http"
	_ "net/http/pprof"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/shimmerglass/bar3x/bar"
)

const childEnv = "BAR3X_CHILD"

func main() {
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	if err != nil {
		log.Error(err)
	}
	if err == nil {
		log.AddHook(hook)
	}

	if os.Getenv(childEnv) != "" {
		runChild()
		return
	}

	for {
		errBuf := ringbuffer.New(1024)
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Stderr = io.MultiWriter(errBuf, os.Stderr)
		cmd.Stdout = os.Stdout
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=1", childEnv))
		err := cmd.Run()
		if err != nil {
			log.Error(string(errBuf.Bytes()))
		}
		ntf := notify.NewNotification("bar3x", fmt.Sprintf("bar3x: exited with status %d", cmd.ProcessState.ExitCode()))
		ntf.Show()

		time.Sleep(time.Second)
	}
}

func runChild() {
	cfgPath := flag.String("config", "", "YAML Config file path")
	themePath := flag.String("theme", "", "YAML Theme file path")
	debug := flag.Bool("debug", false, "Enable profile server on port 6060")
	flag.Parse()

	cfg, err := getConfig(*cfgPath, *themePath)
	if err != nil {
		log.Fatal(err)
	}

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	xevent.ErrorHandlerSet(X, func(err xgb.Error) {
		// we sometimes get BadWindow errors from the tray, I'm not sure why
		// silence them to avoid flooding the output
		if _, ok := err.(xproto.WindowError); ok {
			return
		}

		log.Errorf("x handler error: %s", err)
	})

	_, err = bar.CreateBars(cfg, X)
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	xevent.Main(X)
}
