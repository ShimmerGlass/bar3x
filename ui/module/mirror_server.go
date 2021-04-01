package module

import (
	"image"
	"image/draw"

	"github.com/shimmerglass/bar3x/lib/mirror"
	"github.com/shimmerglass/bar3x/ui"
)

type MirrorServer struct {
	moduleBase

	srv  *mirror.Server
	name string
}

func NewMirrorServer(p ui.ParentDrawable, srv *mirror.Server) *MirrorServer {
	return &MirrorServer{
		moduleBase: newBase(p),
		srv:        srv,
	}
}

func (m *MirrorServer) Init() error {
	return nil
}

func (m *MirrorServer) Add(child ui.Drawable) {
	m.Root = child
}

func (m *MirrorServer) Draw(x, y int, im draw.Image) {
	m.moduleBase.Draw(x, y, im)

	w, h := m.Width(), m.Height()
	sim := image.NewRGBA(image.Rect(0, 0, w, h))
	for xx := 0; xx < w; xx++ {
		for yy := 0; yy < h; yy++ {
			sim.Set(xx, yy, im.At(x+xx, y+yy))
		}
	}

	m.srv.Send(m.name, sim)
}

// parameters

func (m *MirrorServer) Name() string {
	return m.name
}

func (m *MirrorServer) SetName(v string) {
	m.name = v
}
