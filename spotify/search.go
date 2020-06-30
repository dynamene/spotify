package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// FindTrack find a track based on the given info
func FindTrack(track SingleTrack) string {
	searchURL := fmt.Sprintf("%s/v1/search/?q=%s&type=track", baseURL, url.QueryEscape(track.Title+" "+track.Artist))

	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SPOTIFY_ACCESS_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		accessToken, err := refreshToken(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"), os.Getenv("SPOTIFY_REFRESH_TOKEN"))
		if err != nil {
			return ""
		}

		os.Setenv("SPOTIFY_ACCESS_TOKEN", accessToken)
		return FindTrack(track)
	}

	// If we make too may requests then we wait and make te request again
	if resp.StatusCode == http.StatusTooManyRequests {
		delay, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		time.Sleep(time.Duration(delay) * time.Second)
		return FindTrack(track)
	}

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)
	searchResult := SearchTrack{}
	if err = json.Unmarshal(body, &searchResult); err != nil {
		fmt.Println(err)
		return ""
	}

	trackID := ""

	for _, v := range searchResult.Tracks.Items {
		score := 0

		artists := []string{}
		for _, i := range v.Artists {
			artists = append(artists, i.Name)
		}

		if v.Name == track.Title && v.Album.Artists[0].Name == track.Artist {
			score += 2
		}

		if v.Album.Name == track.Album {
			score++
		}

		if sameStringSlice(track.Contributors, artists) {
			score++
		}

		if int(math.Round(float64(v.DurationMs/1000))) == track.Duration {
			score++
		}

		if score >= 4 {
			trackID = v.URI
			break
		}
	}

	return trackID
}
