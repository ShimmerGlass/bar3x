package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/skratchdot/open-golang/open"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const spotifyRedirectURI = "http://localhost:8080/callback"

var spotifyScopes = []string{
	spotify.ScopeUserReadPrivate,
	spotify.ScopePlaylistReadPrivate,
	spotify.ScopeUserLibraryRead,
	spotify.ScopeUserLibraryModify,
	spotify.ScopeUserReadCurrentlyPlaying,
	spotify.ScopeUserReadPlaybackState,
}

func spotifyClient(keyID, keySecret string) (*spotify.Client, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	tokenFilePath := path.Join(usr.HomeDir, ".spotify")

	client, err := spotifyClientSaved(tokenFilePath, keyID, keySecret)
	if err == nil && client != nil {
		return client, nil
	}
	if err != nil {
		log.Println(err)
	}

	client, err = spotifyClientAcquire(keyID, keySecret)
	if err != nil {
		return nil, err
	}

	tokenFile, err := os.OpenFile(tokenFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		log.Println(err)
		return client, nil
	}
	defer tokenFile.Close()

	token, err := client.Token()
	if err != nil {
		log.Println(err)
	} else {
		err = json.NewEncoder(tokenFile).Encode(token)
		if err != nil {
			log.Println(err)
			return client, nil
		}
	}

	return client, nil
}

func spotifyClientSaved(path string, keyID, keySecret string) (*spotify.Client, error) {
	spotifyAuthenticator := spotify.NewAuthenticator(spotifyRedirectURI)
	spotifyAuthenticator.SetAuthInfo(keyID, keySecret)

	tokenFile, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return nil, nil
	}

	token := &oauth2.Token{}
	err = json.NewDecoder(tokenFile).Decode(token)
	if err != nil {
		return nil, err
	}

	client := spotifyAuthenticator.NewClient(token)
	return &client, nil
}

func spotifyClientAcquire(keyID, keySecret string) (*spotify.Client, error) {
	spotifyAuthenticator := spotify.NewAuthenticator(spotifyRedirectURI)
	spotifyAuthenticator.SetAuthInfo(keyID, keySecret)

	ch := make(chan spotify.Client)
	state := "abc123"

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		tok, err := spotifyAuthenticator.Token(state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}

		// use the token to get an authenticated client
		client := spotifyAuthenticator.NewClient(tok)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "Login Completed!<script>window.close();</script>")

		ch <- client
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go http.ListenAndServe(":8080", nil)

	url := spotifyAuthenticator.AuthURL(state)
	time.Sleep(2 * time.Second)
	open.Run(url)

	// wait for auth to complete
	client := <-ch
	return &client, nil
}
