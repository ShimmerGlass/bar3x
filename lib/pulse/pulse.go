package pulse

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/godbus/dbus"
	"github.com/sqp/pulseaudio"
)

type Pulse struct {
	pulse *pulseaudio.Client

	Up chan float64
}

func New() *Pulse {
	m := &Pulse{
		Up: make(chan float64),
	}

	go func() {
		for m.pulse == nil {
			var err error
			m.pulse, err = pulseaudio.New()
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second)
				continue
			}

			go func() {
				for {
					errs := m.pulse.Register(m)
					if len(errs) > 0 {
						log.Println(err)
						time.Sleep(time.Second)
						continue
					}

					m.pulse.Listen()
				}
			}()

			m.update()
		}
	}()

	return m
}

func (m *Pulse) DeviceVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	m.update()
}

func (m *Pulse) DeviceMuteUpdated(path dbus.ObjectPath, state bool) {
	m.update()
}

func (m *Pulse) update() {
	vol, err := m.volume()
	if err != nil {
		log.Println(err)
		return
	}

	m.Up <- vol
}

func (m *Pulse) volume() (float64, error) {
	sinks, err := m.pulse.Core().ListPath("Sinks")
	if err != nil {
		return 0, err
	}

	if len(sinks) == 0 {
		return 0, nil
	}

	var muted bool
	err = m.pulse.Device(sinks[0]).Get("Mute", &muted)
	if err != nil {
		return 0, err
	}
	if muted {
		return 0, nil
	}
	var volumes []uint32
	err = m.pulse.Device(sinks[0]).Get("Volume", &volumes)
	if err != nil {
		return 0, err
	}

	var volumeSteps uint32
	err = m.pulse.Device(sinks[0]).Get("VolumeSteps", &volumeSteps)
	if err != nil {
		return 0, err
	}

	var volTotal uint32
	for _, v := range volumes {
		volTotal += v
	}

	return float64(volTotal) / float64(len(volumes)) / float64(volumeSteps), nil
}
