package playmo

import (
	"encoding/json"
	"net/http"

	"github.com/KengoWada/playmo/spotify"
	"gopkg.in/validator.v2"
)

// Spotify PlayMo Spotify Cloud Function
func Spotify(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "GET":
		link, _ := req.URL.Query()["link"]
		playlist, err := spotify.GetPlaylist(link[0])
		resp.WriteHeader(http.StatusOK)
		if err != nil {
			data := map[string]interface{}{
				"playlist": map[string]string{},
				"isValid":  false,
			}
			json.NewEncoder(resp).Encode(data)
			return
		}

		data := map[string]interface{}{
			"playlist": playlist,
			"isValid":  true,
		}
		json.NewEncoder(resp).Encode(data)
		return
	case "POST":
		var data struct {
			Description string                `json:"description"`
			Name        string                `json:"name" validate:"nonzero"`
			Tracks      []spotify.SingleTrack `json:"tracks"`
		}

		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			data := map[string]interface{}{
				"link":             "",
				"missingTracks":    []string{},
				"numMissingTracks": 0,
			}
			resp.WriteHeader(http.StatusCreated)
			json.NewEncoder(resp).Encode(data)
			return
		}

		if err := validator.Validate(data); err != nil {
			data := map[string]interface{}{
				"link":             "",
				"missingTracks":    []string{},
				"numMissingTracks": 0,
			}
			resp.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(resp).Encode(data)
			return
		}

		response := spotify.CreatePlaylist(data.Name, data.Description, data.Tracks)
		if len(response) == 0 {
			response = map[string]interface{}{
				"link":             "",
				"missingTracks":    []string{},
				"numMissingTracks": 0,
			}
		}
		resp.WriteHeader(http.StatusCreated)
		json.NewEncoder(resp).Encode(response)
		return
	default:
		response := map[string]interface{}{
			"message": "Invalid request method",
			"status":  http.StatusBadRequest,
		}
		resp.WriteHeader(http.StatusCreated)
		json.NewEncoder(resp).Encode(response)
		return
	}
}
