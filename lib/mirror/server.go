package mirror

import (
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

	start sync.Once
}

func NewServer(addr string) *Server {
	return &Server{
		addr:      addr,
		modules:   map[string]*Image{},
		listeners: map[string]map[chan *Image]struct{}{},
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
			return err
		}
	}

	return nil
}
