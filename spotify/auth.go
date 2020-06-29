package spotify

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// GetCodeURL returns the URL to get the spotify code
func GetCodeURL(spotifyClientID, spotifySecret, scope, redirectURL string) string {
	return fmt.Sprintf("%s/authorize/?response_type=code&client_id=%s&scope=%s&redirect_uri=%s", authURL, spotifyClientID, url.QueryEscape(scope), redirectURL)
}

// GetTokens returns both access and refresh token for spotify
func GetTokens(spotifyClientID, spotifySecret, code, redirectURL string) (map[string]interface{}, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURL)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/token", authURL), strings.NewReader(data.Encode()))
	req.Header.Add("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(spotifyClientID+":"+spotifySecret)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return make(map[string]interface{}), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return make(map[string]interface{}), errors.New("Invalid Request")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var tokens map[string]interface{}
	if err = json.Unmarshal(body, &tokens); err != nil {
		return make(map[string]interface{}), err
	}

	return tokens, nil
}

// refreshToken gets a new accessToken and sets it to the SPOTIFY_ACCESS_TOKEN environment variable
func refreshToken(spotifyClientID, spotifySecret, refreshToken string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/token", authURL), strings.NewReader(data.Encode()))
	req.Header.Add("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(spotifyClientID+":"+spotifySecret)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	token := make(map[string]interface{})
	if err = json.Unmarshal(body, &token); err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", token["access_token"]), nil
}
