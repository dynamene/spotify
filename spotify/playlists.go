package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetPlaylist returns playlist information from playlist link
func GetPlaylist(playlistLink string) (map[string]interface{}, error) {
	s := strings.Split(playlistLink, "/")
	playlistID := strings.Split(s[len(s)-1], "?")[0]
	playlistURL := fmt.Sprintf("%s/v1/playlists/%s", baseURL, playlistID)

	req, _ := http.NewRequest("GET", playlistURL, nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SPOTIFY_ACCESS_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		accessToken, err := refreshToken(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"), os.Getenv("SPOTIFY_REFRESH_TOKEN"))
		if err != nil {
			return map[string]interface{}{}, err
		}

		os.Setenv("SPOTIFY_ACCESS_TOKEN", accessToken)
		return GetPlaylist(playlistLink)
	}

	// If we make too may requests then we wait and make te request again
	if resp.StatusCode == http.StatusTooManyRequests {
		delay, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		time.Sleep(time.Duration(delay) * time.Second)
		return GetPlaylist(playlistLink)
	}

	if resp.StatusCode != http.StatusOK {
		return map[string]interface{}{}, errors.New("Invalid playlist link")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	playlist := Playlist{}
	tracks := []SingleTrack{}
	if err = json.Unmarshal(body, &playlist); err != nil {
		return map[string]interface{}{}, errors.New("Invalid playlist link")
	}

	if len(playlist.Tracks.Items) > 20 {
		playlist.Tracks.Items = playlist.Tracks.Items[:20]
	}

	for _, v := range playlist.Tracks.Items {
		track := SingleTrack{
			Album:        v.Track.Album.Name,
			Artist:       v.Track.Album.Artists[0].Name,
			Contributors: []string{},
			Duration:     int(math.Round(float64(v.Track.DurationMs / 1000))),
			Title:        v.Track.Name,
			TrackCover:   v.Track.Album.Images[1].URL,
		}

		for _, i := range v.Track.Artists {
			track.Contributors = append(track.Contributors, i.Name)
		}

		tracks = append(tracks, track)
	}

	playlistInfo := map[string]interface{}{
		"name":          playlist.Name,
		"description":   playlist.Description,
		"playlistCover": playlist.Images[0].URL,
		"tracks":        tracks,
		"numTracks":     len(tracks),
	}

	return playlistInfo, nil
}

// CreatePlaylist creates a playlist
func CreatePlaylist(name string, description string, tracks []SingleTrack) map[string]interface{} {
	createURL := fmt.Sprintf("%s/v1/me/playlists", baseURL)
	reqBody := []byte(fmt.Sprintf(`{ "name": "%s", "description": "%s" }`, name, description))

	req, _ := http.NewRequest("POST", createURL, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SPOTIFY_ACCESS_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return map[string]interface{}{}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		accessToken, err := refreshToken(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"), os.Getenv("SPOTIFY_REFRESH_TOKEN"))
		if err != nil {
			return map[string]interface{}{}
		}

		os.Setenv("SPOTIFY_ACCESS_TOKEN", accessToken)
		return CreatePlaylist(name, description, tracks)
	}

	// If we make too may requests then we wait and make te request again
	if resp.StatusCode == http.StatusTooManyRequests {
		delay, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		time.Sleep(time.Duration(delay) * time.Second)
		return CreatePlaylist(name, description, tracks)
	}

	if resp.StatusCode != http.StatusCreated {
		return map[string]interface{}{}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	playlist := CreatedPlaylist{}
	if err = json.Unmarshal(body, &playlist); err != nil {
		return map[string]interface{}{}
	}

	trackIDs := []string{}
	missingTracks := []SingleTrack{}
	for _, v := range tracks {
		id := FindTrack(v)
		if id == "" {
			missingTracks = append(missingTracks, v)
		} else {
			trackIDs = append(trackIDs, id)
		}
	}

	// Add tracks to playlist
	addTracksURL := fmt.Sprintf("%s/v1/playlists/%s/tracks", baseURL, playlist.ID)
	ids, _ := json.Marshal(trackIDs)
	reqBody = []byte(fmt.Sprintf(`{ "uris": %s }`, ids))

	req, _ = http.NewRequest("POST", addTracksURL, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SPOTIFY_ACCESS_TOKEN"))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return map[string]interface{}{}
	}

	if resp.StatusCode == http.StatusUnauthorized {
		accessToken, err := refreshToken(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"), os.Getenv("SPOTIFY_REFRESH_TOKEN"))
		if err != nil {
			return map[string]interface{}{}
		}

		os.Setenv("SPOTIFY_ACCESS_TOKEN", accessToken)
		return CreatePlaylist(name, description, tracks)
	}

	// If we make too may requests then we wait and make te request again
	if resp.StatusCode == http.StatusTooManyRequests {
		delay, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		time.Sleep(time.Duration(delay) * time.Second)
		return CreatePlaylist(name, description, tracks)
	}

	playlistInfo := map[string]interface{}{
		"link":             playlist.ExternalUrls.Spotify,
		"missingTracks":    missingTracks,
		"numMissingTracks": len(missingTracks),
	}

	return playlistInfo
}
