package spotify

import (
	log "github.com/sirupsen/logrus"

	"github.com/zmb3/spotify"
)

func (s *Spotify) toggleLibrary(id spotify.ID) error {
	s.clientReady.Wait()

	res, err := s.client.UserHasTracks(id)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		return nil
	}

	if !res[0] {
		err := s.client.AddTracksToLibrary(id)
		if err != nil {
			return err
		}
	} else {
		err := s.client.RemoveTracksFromLibrary(id)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	s.userTracks.Delete(string(id))

	return nil
}
