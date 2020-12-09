package module

import (
	"context"
	"path/filepath"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	gonm "github.com/Wifx/gonetworkmanager"
	"github.com/shimmerglass/bar3x/ui"
	"github.com/shimmerglass/bar3x/ui/base"
	"github.com/shimmerglass/bar3x/ui/markup"
)

type connState struct {
	isVPN  bool
	name   string
	ip     string
	status gonm.NmActiveConnectionState
}

type connIndicator struct {
	VPN    *base.Sizer
	Root   *base.Sizer
	Name   *base.Text
	Status *base.Circle
}

type Connections struct {
	moduleBase

	Row *base.Row

	els []*connIndicator

	typeFilter map[string]bool
	nameFilter map[string]bool

	mk    *markup.Markup
	clock *Clock
	nm    gonm.NetworkManager
}

func NewConnections(p ui.ParentDrawable, mk *markup.Markup, clock *Clock) *Connections {
	return &Connections{
		clock:      clock,
		mk:         mk,
		moduleBase: newBase(p),
		typeFilter: map[string]bool{
			"vpn":            true,
			"802-3-ethernet": true,
		},
	}
}

func (m *Connections) Init() error {
	_, err := m.mk.Parse(m, m, `
		<Sizer ref="Root">
			<Row ref="Row" />
		</Sizer>
	`)
	if err != nil {
		return err
	}

	m.clock.Add(m, time.Second)
	return nil
}

func (m *Connections) SetTypeFilter(v []string) {
	f := map[string]bool{}
	for _, v := range v {
		f[v] = true
	}
	m.typeFilter = f
}

func (m *Connections) TypeFilter() []string {
	r := []string{}
	for f := range m.typeFilter {
		r = append(r, f)
	}
	return r
}

func (m *Connections) SetNameFilter(v []string) {
	f := map[string]bool{}
	for _, v := range v {
		f[v] = true
	}
	m.nameFilter = f
}

func (m *Connections) NameFilter() []string {
	r := []string{}
	for f := range m.nameFilter {
		r = append(r, f)
	}
	return r
}

func (m *Connections) Update(context.Context) {
	connections := m.getConnections()

	if len(connections) == 0 {
		m.SetVisible(false)
		return
	}

	for i := len(m.els); i < len(connections); i++ {
		m.addIndicator()
	}

	for i := range connections {
		indicator := m.els[i]
		conn := connections[i]
		currentTxt := indicator.Name.Text()

		if currentTxt != conn.name && currentTxt != conn.ip {
			indicator.Name.SetText(conn.name)
		}
		switch conn.status {
		case gonm.NmActiveConnectionStateUnknown, gonm.NmActiveConnectionStateActivating, gonm.NmActiveConnectionStateDeactivating:
			indicator.Status.SetColor(m.Context().MustColor("warning_color"))
		case gonm.NmActiveConnectionStateActivated:
			indicator.Status.SetColor(m.Context().MustColor("success_color"))
		case gonm.NmActiveConnectionStateDeactivated:
			indicator.Status.SetColor(m.Context().MustColor("danger_color"))
		}
		func(i int) {
			if conn.ip == "" {
				return
			}
			indicator.Root.SetOnPointerEnter(func(ui.Event) bool {
				indicator.Name.SetText(conn.ip)
				indicator.Name.Notify()
				return true
			})
			indicator.Root.SetOnPointerLeave(func(ui.Event) bool {
				indicator.Name.SetText(conn.name)
				indicator.Name.Notify()
				return true
			})
		}(i)

		indicator.Root.SetVisible(true)
		indicator.VPN.SetVisible(connections[i].isVPN)
	}
	for i := len(connections); i < len(m.els); i++ {
		m.els[i].Root.SetVisible(false)
	}

	m.SetVisible(true)
}

func (m *Connections) addIndicator() {
	el := &connIndicator{}
	root := m.mk.MustParse(m.Row, el, `
		<Sizer
			ref="Root"
			PaddingLeft="{h_padding}"
		>
			<Rect
				Radius="1"
				Color="{neutral_color}">
				<Sizer
					Height="{bar_height - v_padding * 3}"
					PaddingLeft="{h_padding}"
					PaddingRight="{h_padding}"
				>
					<Row>
						<Sizer ref="VPN" PaddingRight="{h_padding}">
							<Icon FontSize="{text_small_font_size}">{icons.lock}</Icon>
						</Sizer>
						<Text ref="Name" FontSize="{text_small_font_size}" />
						<Sizer PaddingLeft="{h_padding}">
							<Circle ref="Status" Radius="3" />
						</Sizer>
					</Row>
				</Sizer>
			</Rect>
		</Sizer>
	`)
	m.Row.Add(root)
	m.els = append(m.els, el)
}

func (m *Connections) getConnections() []connState {
	if m.nm == nil {
		var err error
		m.nm, err = gonm.NewNetworkManager()
		if err != nil {
			log.Println(err)
			m.SetVisible(false)
			return nil
		}
	}

	conns, err := m.nm.GetPropertyActiveConnections()
	if err != nil {
		log.Printf("vpn: %s", err)
		m.SetVisible(false)
		return nil
	}

	connections := []connState{}

	for _, c := range conns {
		t, _ := c.GetPropertyType()
		if len(m.typeFilter) > 0 && !m.typeFilter[t] {
			continue
		}
		isVPN, err := c.GetPropertyVPN()
		if err != nil {
			log.Printf("vpn: %s", err)
			continue
		}
		conn, err := c.GetPropertyConnection()
		if err != nil {
			continue
		}
		status, err := c.GetPropertyState()
		if err != nil {
			continue
		}
		var ip string
		ipv6, err := c.GetPropertyIP6Config()
		if err != nil {
			continue
		}
		if ipv6 != nil {
			addr, err := ipv6.GetPropertyAddressData()
			if err != nil {
				continue
			}
			if len(addr) > 0 {
				ip = addr[0].Address
			}
		}
		ipv4, err := c.GetPropertyIP4Config()
		if err != nil {
			continue
		}
		if ipv4 != nil {
			addr, err := ipv4.GetPropertyAddressData()
			if err != nil {
				continue
			}
			if len(addr) > 0 {
				ip = addr[0].Address
			}
		}
		path, err := conn.GetPropertyFilename()
		if err != nil {
			continue
		}
		file := filepath.Base(path)
		name := strings.TrimSuffix(file, filepath.Ext(file))

		if len(m.nameFilter) > 0 && !m.nameFilter[name] {
			continue
		}

		connections = append(connections, connState{
			isVPN:  isVPN,
			ip:     ip,
			name:   name,
			status: status,
		})
	}

	sort.Slice(connections, func(i, j int) bool {
		return connections[i].name < connections[j].name
	})

	return connections
}
