package mirror

import (
	"bytes"
	"image"
	"net"
	sync "sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedMirrorServer

	addr string

	lock      sync.Mutex
	modules   map[string]*Image
	listeners map[string]map[chan *Image]struct{}

	throttle *time.Timer

	start sync.Once
}

func NewServer(addr string) *Server {
	return &Server{
		addr:      addr,
		modules:   map[string]*Image{},
		listeners: map[string]map[chan *Image]struct{}{},
		throttle:  time.NewTimer(time.Second),
	}
}

func (s *Server) Send(name string, img *image.RGBA) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.start.Do(func() {
		srv := grpc.NewServer()
		lis, err := net.Listen("tcp", s.addr)
		if err != nil {
			log.Fatal(err)
		}

		RegisterMirrorServer(srv, s)

		go func() {
			err = srv.Serve(lis)
			if err != nil {
				log.Error(err)
			}
		}()
	})

	imgpb := &Image{
		Pixels: img.Pix,
		Stride: int32(img.Stride),
		Width:  int32(img.Rect.Dx()),
		Height: int32(img.Rect.Dy()),
	}

	if s.modules[name] != nil && bytes.Equal(s.modules[name].Pixels, imgpb.Pixels) {
		return
	}

	s.modules[name] = imgpb

	for c := range s.listeners[name] {
		select {
		case c <- imgpb:
		case <-time.After(time.Second):
		}
	}
}

func (s *Server) Subscribe(req *SubscribeRequest, srv Mirror_SubscribeServer) error {
	c := make(chan *Image, 1)

	s.lock.Lock()
	if s.listeners[req.Name] == nil {
		s.listeners[req.Name] = map[chan *Image]struct{}{}
	}
	s.listeners[req.Name][c] = struct{}{}
	img, ok := s.modules[req.Name]
	if ok {
		c <- img
	}
	s.lock.Unlock()
	defer func() {
		s.lock.Lock()
		delete(s.listeners[req.Name], c)
		s.lock.Unlock()
	}()

	for img := range c {
		log.Infof("mirror: %s: sending img", req.Name)
		err := srv.Send(img)
		if err != nil {
			log.Errorf("mirror server: send: %s", err)
			return err
		}
	}

	return nil
}
