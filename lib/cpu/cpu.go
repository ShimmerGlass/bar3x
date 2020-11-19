package cpu

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

var usage []float64
var lock sync.Mutex
var once sync.Once

func Start() {
	once.Do(func() {
		go func() {
			for {
				percents, _ := cpu.Percent(time.Second, true)

				lock.Lock()
				usage = nil
				for _, p := range percents {
					usage = append(usage, p/100)
				}
				lock.Unlock()
			}
		}()
	})
}

func Read() []float64 {
	lock.Lock()
	defer lock.Unlock()
	return usage
}
