package pulse

// #cgo pkg-config: libpulse
// #include <pulse/pulseaudio.h>
// #include "./pulse.h"
import "C"
import (
	"fmt"
	"os/exec"
	"sync"

	log "github.com/sirupsen/logrus"
)

var currentVolume float64

var lock sync.Mutex
var watchers []chan struct{}

func init() {
	log.Info("connecting to pulse")
	r := C.initialize()
	if r != 0 {
		log.Error("could not connect to pulse")
	}
	go C.run()
}

func Volume() float64 {
	return currentVolume
}

func SetVolume(v float64) error {
	percent := int(v * 100)
	return exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", fmt.Sprintf("%d%%", percent)).Run()
}

func Watch(c chan struct{}) {
	lock.Lock()
	defer lock.Unlock()
	watchers = append(watchers, c)
	go func() {
		c <- struct{}{}
	}()
}
