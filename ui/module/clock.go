package module

import (
	"context"
	"reflect"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const sequentialThreshold = 100 * time.Millisecond

type Updatable interface {
	Update(context.Context)
}

type Clock struct {
	paint func()

	lock     sync.Mutex
	els      map[int][]Updatable
	totalEls int

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
	c.totalEls++
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

		done := make([]reflect.SelectCase, 0, c.totalEls+1)
		done = append(done, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(time.After(sequentialThreshold)),
		})

		for interval, modules := range c.els {
			if counter%interval != 0 {
				continue
			}
			needPaint = true
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, time.Duration(interval)*time.Millisecond)
			defer cancel()

			for _, m := range modules {
				c := make(chan struct{})
				done = append(done, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(c),
				})
				go func(m Updatable, c chan struct{}) {
					m.Update(ctx)
					close(c)
				}(m, c)
			}
		}

		for len(done) > 1 {
			chosen, _, _ := reflect.Select(done)

			// idx zero is the timeout chan
			if chosen == 0 {
				// we are timing out, wait for late modules and paint afterwards
				go func() {
					// we  ignore the timeout now
					done = done[1:]
					log.Warnf("clock: %d modules took more than %s to update", len(done), sequentialThreshold)
					for len(done) > 0 {
						chosen, _, _ := reflect.Select(done)
						done[chosen] = done[len(done)-1]
						done = done[:len(done)-1]
					}
					c.paint()
				}()

				// render the modules that updated in time now
				break
			}

			// remove done module an wait for others
			done[chosen] = done[len(done)-1]
			done = done[:len(done)-1]
		}

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
