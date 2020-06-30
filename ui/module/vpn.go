package module

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Wifx/gonetworkmanager"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type VPN struct {
	moduleBase

	mk    *markup.Markup
	clock *Clock
	nm    gonetworkmanager.NetworkManager
}

func NewVPN(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *VPN {
	return &VPN{
		clock:      clock,
		mk:         mk,
		moduleBase: newBase(p),
	}
}

func (m *VPN) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Icon ref="Root">{icons.lock}</Icon>
	`)
	if err != nil {
		return err
	}

	m.clock.Add(m, time.Second)
	return nil
}

func (m *VPN) Update(context.Context) {
	if m.nm == nil {
		var err error
		m.nm, err = gonetworkmanager.NewNetworkManager()
		if err != nil {
			log.Println(err)
			m.SetVisible(false)
			return
		}
	}
	conns, err := m.nm.GetPropertyActiveConnections()
	if err != nil {
		log.Printf("vpn: %s", err)
		m.SetVisible(false)
		return
	}
	for _, c := range conns {
		isVPN, err := c.GetPropertyVPN()
		if err != nil {
			log.Printf("vpn: %s", err)
			continue
		}
		if isVPN {
			m.SetVisible(true)
			return
		}
	}

	m.SetVisible(false)
}
