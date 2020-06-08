package spotify

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/godbus/dbus"
	"github.com/zmb3/spotify"
)

func (s *Spotify) track() (id spotify.ID, title string, artist string, err error) {
	res := dbus.Variant{}

	err = s.conn.Object(
		"org.mpris.MediaPlayer2.spotify",
		"/org/mpris/MediaPlayer2",
	).Call(
		"org.freedesktop.DBus.Properties.Get",
		0,
		"org.mpris.MediaPlayer2.Player",
		"Metadata",
	).Store(&res)
	if err != nil {
		return
	}

	val, ok := res.Value().(map[string]dbus.Variant)
	if !ok {
		return
	}

	artistVar, ok := val["xesam:artist"].Value().([]string)
	if ok {
		artist = strings.Join(artistVar, ", ")
	}

	title, _ = val["xesam:title"].Value().(string)

	rid, _ := val["mpris:trackid"].Value().(string)
	if len(rid) > 0 {
		parts := strings.Split(rid, ":")
		id = spotify.ID(parts[len(parts)-1])
	}

	return
}

func (s *Spotify) playStatus() (playing bool, err error) {
	res := dbus.Variant{}

	err = s.conn.Object(
		"org.mpris.MediaPlayer2.spotify",
		"/org/mpris/MediaPlayer2",
	).Call(
		"org.freedesktop.DBus.Properties.Get",
		0,
		"org.mpris.MediaPlayer2.Player",
		"PlaybackStatus",
	).Store(&res)
	if err != nil {
		return
	}

	return res.Value().(string) == "Playing", nil
}

func (s *Spotify) watchStatus() {
	s.update()

	for {
		err := s.conn.BusObject().Call(
			"org.freedesktop.DBus.AddMatch",
			0,
			"type='signal',path='/org/mpris/MediaPlayer2',interface='org.freedesktop.DBus.Properties',sender='org.mpris.MediaPlayer2.spotify'",
		).Err
		if err != nil {
			log.Println(err)
		}
		c := make(chan *dbus.Signal, 10)
		s.conn.Signal(c)

		for range c {
			s.update()
		}

		time.Sleep(time.Second)
	}
}
