package spotify

import (
	"bytes"
	"image"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	_ "image/jpeg"
	_ "image/png"

	"github.com/godbus/dbus"
	lru "github.com/hashicorp/golang-lru"
	"github.com/patrickmn/go-cache"
	"github.com/zmb3/spotify"
)

type State struct {
	Playing  bool
	Title    string
	Artist   string
	Image    image.Image
	ImageURL string
	Progress float64
}

type Spotify struct {
	conn *dbus.Conn

	client      *spotify.Client
	clientReady sync.WaitGroup
	imgCache    *lru.Cache

	userTracks *cache.Cache

	state *State

	Up chan State
}

func New(keyID, keySecret string) *Spotify {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal("Failed to connect to session bus:", err)
	}

	imgCache, _ := lru.New(5)
	s := &Spotify{
		conn:       conn,
		userTracks: cache.New(10*time.Minute, time.Minute),
		Up:         make(chan State),

		imgCache: imgCache,
	}
	go s.watchStatus()
	go s.watchPosition()

	s.clientReady.Add(1)
	go func() {
		client, err := spotifyClient(keyID, keySecret)
		if err != nil {
			log.Println(err)
			return
		}

		s.client = client
		s.clientReady.Done()
	}()

	return s
}

func (s *Spotify) update() {
	playing, err := s.playStatus()
	if err != nil {
		log.Println(err)
		return
	}

	id, title, artist, err := s.track()
	if err != nil {
		log.Println(err)
		return
	}

	state := &State{
		Playing: playing,
		Title:   title,
		Artist:  artist,
	}

	if s.state != nil && *s.state == *state {
		return
	}

	if playing {
		go func() {
			if err := s.updateImage(id, state); err != nil {
				log.Println(err)
			}
		}()
	}

	s.state = state
	s.sendUpdate()
}

func (s *Spotify) updateImage(id spotify.ID, state *State) error {
	s.clientReady.Wait()
	track, err := s.client.GetTrack(id)
	if err != nil {
		return err
	}

	if len(track.Album.Images) == 0 {
		return nil
	}

	imgBuf := &bytes.Buffer{}
	imgProp := track.Album.Images[len(track.Album.Images)-1] // last one is smallest

	var img image.Image
	if imgi, ok := s.imgCache.Get(imgProp.URL); ok {
		img = imgi.(image.Image)
	} else {
		err = imgProp.Download(imgBuf)
		if err != nil {
			return err
		}

		img, _, err = image.Decode(imgBuf)
		if err != nil {
			return err
		}
		s.imgCache.Add(imgProp.URL, img)
	}

	state.Image = img
	state.ImageURL = imgProp.URL
	s.sendUpdate()

	return nil
}

func (s *Spotify) sendUpdate() {
	s.Up <- *s.state
}

func (s *Spotify) watchPosition() {
	s.clientReady.Wait()
	for range time.Tick(2 * time.Second) {
		if s.state == nil {
			continue
		}
		if !s.state.Playing {
			continue
		}

		state, err := s.client.PlayerState()
		if err != nil {
			log.Println(err)
			continue
		}

		if s.state != nil && state.Item != nil {
			s.state.Progress = float64(state.Progress) / float64(state.Item.Duration)
			s.sendUpdate()
		} else {
			s.state.Progress = 0
			s.sendUpdate()
		}
	}
}
