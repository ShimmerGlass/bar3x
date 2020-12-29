package rgb

import (
	"context"

	"github.com/shimmerglass/bar3x/lib/pulse"
	"github.com/shimmerglass/rgbx/rgbx"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type RGB struct {
	client rgbx.RGBizerClient
}

func New(addr string) *RGB {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := rgbx.NewRGBizerClient(conn)

	return &RGB{client}
}

func (r *RGB) Run() {
	go r.volume()
}

func (r *RGB) volume() {
	c := make(chan struct{})
	pulse.Watch(c)

	last := -1.0
	for range c {
		vol := pulse.Volume()
		if vol == last {
			continue
		}
		last = vol

		_, err := r.client.Set(context.Background(), &rgbx.SetRequest{
			Priority:   100,
			DurationMs: 1000,
			Effect: &rgbx.SetRequest_Static{
				Static: &rgbx.EffectStatic{
					Color: &rgbx.Color{},
				},
			},
		})
		if err != nil {
			log.Error(err)
		}
		_, err = r.client.Set(context.Background(), &rgbx.SetRequest{
			Priority:   101,
			DurationMs: 1000,
			Effect: &rgbx.SetRequest_Progress{
				Progress: &rgbx.EffectProgress{
					Color: &rgbx.Color{R: 0xff, G: 0x00, B: 0xFF},
					Value: vol,
					Rows:  []int32{0},
				},
			},
		})
		if err != nil {
			log.Error(err)
		}
	}
}
