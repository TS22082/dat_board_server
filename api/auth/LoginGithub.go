package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func LoginGithub(w http.ResponseWriter, r *http.Request) {
	// Get the code from the request
	code := r.URL.Query().Get("code")

	// Get the client id and secret from the environment
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	// Create a form data
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)

	// Create a new request to get the access token
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Get the access token from the response

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "token "+tokenResponse.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var userResponse struct {
		Login string `json:"login"`
	}

	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
