package module

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Wifx/gonetworkmanager"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type VPN struct {
	moduleBase

	Txt *base.Text

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
		<Row ref="Root">
			<Sizer PaddingRight="{h_padding}">
				<Icon>{icons.lock}</Icon>
			</Sizer>
			<Text ref="Txt" />
		</Row>
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

	activeVPNs := []string{}

	for _, c := range conns {
		isVPN, err := c.GetPropertyVPN()
		if err != nil {
			log.Printf("vpn: %s", err)
			continue
		}
		if isVPN {
			conn, err := c.GetPropertyConnection()
			if err != nil {
				continue
			}
			path, err := conn.GetPropertyFilename()
			if err != nil {
				continue
			}
			file := filepath.Base(path)
			name := strings.TrimSuffix(file, filepath.Ext(file))

			activeVPNs = append(activeVPNs, name)
		}
	}

	if len(activeVPNs) == 0 {
		m.SetVisible(false)
		return
	}

	m.SetVisible(true)
	m.Txt.SetText(strings.Join(activeVPNs, ", "))
}
