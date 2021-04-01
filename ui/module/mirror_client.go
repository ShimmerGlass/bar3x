package module

import (
	"github.com/shimmerglass/bar3x/lib/mirror"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
)

type MirrorClient struct {
	moduleBase

	addr string
	name string
}

func NewMirrorClient(p ui.ParentDrawable) *MirrorClient {
	return &MirrorClient{
		moduleBase: newBase(p),
	}
}

func (m *MirrorClient) Init() error {
	root := base.NewImage(m)
	m.Root = root

	client := mirror.NewClient(m.addr)
	c, err := client.Subscribe(m.name)
	if err != nil {
		return err
	}

	go func() {
		for img := range c {
			root.SetImage(img)
			m.Notify()
		}
	}()

	return nil
}

// parameters

func (m *MirrorClient) Name() string {
	return m.name
}

func (m *MirrorClient) SetName(v string) {
	m.name = v
}

func (m *MirrorClient) Addr() string {
	return m.addr
}

func (m *MirrorClient) SetAddr(v string) {
	m.addr = v
}
