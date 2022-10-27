package mirror

import (
	"context"
	"image"
	"io"
	"sync"
	"time"

	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

type Client struct {
	addr string

	c     MirrorClient
	start sync.Once
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Subscribe(name string) (<-chan *image.RGBA, error) {
	c.start.Do(func() {
		log.Infof("mirror client: subscribing to %s", c.addr)
		conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
			return
		}

		c.c = NewMirrorClient(conn)
	})

	imgs := make(chan *image.RGBA)

	first := make(chan struct{}, 1)
	first <- struct{}{}
	throttle := time.NewTicker(100 * time.Millisecond)

	go func() {
	Connect:
		for {
			srv, err := c.c.Subscribe(context.Background(), &SubscribeRequest{Name: name})
			if err != nil {
				log.Error(err)
				time.Sleep(5 * time.Second)
				continue
			}

			for {
				img, err := srv.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Errorf("mirror client: receive: %s", err)
					time.Sleep(5 * time.Second)
					continue Connect
				}

				select {
				case <-throttle.C:
				case <-first:
				default:
					continue
				}

				log.Infof("mirror: %s: received img", name)
				imgs <- &image.RGBA{
					Pix:    img.Pixels,
					Stride: int(img.Stride),
					Rect:   image.Rect(0, 0, int(img.Width), int(img.Height)),
				}
			}
		}
	}()

	return imgs, nil
}
