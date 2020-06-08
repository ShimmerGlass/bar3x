package module

import (
	"sync"
	"time"
)

type Updatable interface {
	Update()
}

type Clock struct {
	paint func()

	lock sync.Mutex
	els  map[int][]Updatable

	interval int
	reset    int
}

func NewClock(paint func()) *Clock {
	c := &Clock{
		paint: paint,
		els:   map[int][]Updatable{},
	}

	return c
}

func (c *Clock) Add(module Updatable, interval time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	intMs := int(interval.Milliseconds())

	if c.interval == 0 {
		c.interval = intMs
		c.reset = intMs
	} else {
		c.interval = GCD(c.interval, intMs)
		c.reset = c.reset * intMs / GCD(c.reset, intMs)
	}
	c.els[intMs] = append(c.els[intMs], module)
}

func (c *Clock) Run() {
	if c.interval == 0 {
		return
	}

	counter := 0
	for {
		c.lock.Lock()
		start := time.Now()
		needPaint := false

		var wg sync.WaitGroup
		for interval, modules := range c.els {
			if counter%interval != 0 {
				continue
			}
			needPaint = true

			wg.Add(len(modules))
			for _, m := range modules {
				go func(m Updatable) {
					m.Update()
					wg.Done()
				}(m)
			}
		}
		wg.Wait()

		if needPaint {
			c.paint()
		}

		elapsed := time.Since(start)
		time.Sleep(time.Duration(c.interval)*time.Millisecond - elapsed)
		counter += c.interval
		if counter > c.reset {
			counter = 0
		}
		c.lock.Unlock()
	}
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
