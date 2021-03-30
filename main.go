// Example draw-text shows how to draw text to an xgraphics.Image type.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheCreeper/go-notify"
	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/shimmerglass/bar3x/bar"
	debugsrv "github.com/shimmerglass/bar3x/debug"
)

const childEnv = "BAR3X_CHILD"

func main() {
	if os.Getenv(childEnv) != "" {
		runChild()
		return
	}

	var cmd *exec.Cmd

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	go func() {
		for range sigs {
			if cmd != nil {
				syscall.Kill(cmd.Process.Pid, syscall.SIGUSR1)
			}
		}
	}()

	for {
		cmd = exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=1", childEnv))
		err := cmd.Run()
		if err != nil {
			log.Error(err)
		}
		code := cmd.ProcessState.ExitCode()

		if code > 0 {
			ntf := notify.NewNotification("bar3x", fmt.Sprintf("bar3x: exited with status %d", cmd.ProcessState.ExitCode()))
			ntf.Show()
			time.Sleep(time.Second)
		}
	}
}

func runChild() {
	cfgPath := flag.String("config", "", "YAML Config file path")
	themePath := flag.String("theme", "", "YAML Theme file path")
	debug := flag.Bool("debug", false, "Enable profile server on port 6060")
	debugAddr := flag.String("debug-addr", "127.0.0.1:6060", "Debug server addr")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGUSR1)
	go func() {
		for range sigs {
			os.Exit(0)
		}
	}()

	cfg, err := getConfig(*cfgPath, *themePath)
	if err != nil {
		log.Fatal(err)
	}

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	err = randr.Init(X.Conn())
	if err != nil {
		log.Fatal(err)
	}

	err = randr.SelectInputChecked(X.Conn(), X.RootWin(),
		randr.NotifyMaskScreenChange|
			randr.NotifyMaskCrtcChange|
			randr.NotifyMaskOutputChange|
			randr.NotifyMaskOutputProperty).Check()
	if err != nil {
		log.Fatal(err)
	}

	xevent.HookFun(func(xu *xgbutil.XUtil, event interface{}) bool {
		switch event.(type) {
		case randr.ScreenChangeNotifyEvent:
			os.Exit(0)
		case randr.NotifyEvent:
			os.Exit(0)
		}

		return true
	}).Connect(X)

	xevent.ErrorHandlerSet(X, func(err xgb.Error) {
		// we sometimes get BadWindow errors from the tray, I'm not sure why
		// silence them to avoid flooding the output
		if _, ok := err.(xproto.WindowError); ok {
			return
		}

		log.Errorf("x handler error: %s", err)
	})

	bars, err := bar.CreateBars(cfg, X)
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		go func() {
			err := debugsrv.Run(*debugAddr, bars)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	xevent.Main(X)
}
