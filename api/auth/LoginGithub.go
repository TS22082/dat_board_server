// Package auth contains authentication related functions.
package auth

import (
	"encoding/json" // Package json implements encoding and decoding of JSON objects.
	"net/http"      // Package http provides HTTP client and server implementations.
	"net/url"       // Package url parses URLs and implements query escaping.
	"os"            // Package os provides a platform-independent interface to operating system functionality.
	"strings"       // Package strings implements simple functions to manipulate UTF-8 encoded strings.
)

// LoginGithub is a function that handles the callback from GitHub's OAuth flow.
func LoginGithub(w http.ResponseWriter, r *http.Request) {
	// Get the code from the query parameters.
	code := r.URL.Query().Get("code")

	// Get the GitHub client ID and client secret from the environment variables.
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	// Create a url.Values object and set the client ID, client secret, and code.
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)

	// Create a new HTTP request to exchange the code for an access token.
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		// If there's an error, respond with a 500 Internal Server Error status code and return.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Accept and Content-Type headers on the request.
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a new HTTP client and send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// If there's an error, respond with a 500 Internal Server Error status code and return.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Close the response body when the function returns.
	defer resp.Body.Close()

	// Define a struct to hold the access token from the response.
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	// Decode the JSON response into the tokenResponse struct.
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)

	if err != nil {
		// If there's an error, respond with a 500 Internal Server Error status code and return.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a new HTTP request to get the user's GitHub username.
	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		// If there's an error, respond with a 500 Internal Server Error status code and return.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Authorization header on the request.
	req.Header.Set("Authorization", "token "+tokenResponse.AccessToken)

	// Send the request.
	resp, err = client.Do(req)
	if err != nil {
		// If there's an error, respond with a 500 Internal Server Error status code and return.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Close the response body when the function returns.
	defer resp.Body.Close()

	// Define a struct to hold the username from the response.
	var userResponse struct {
		Login string `json:"login"`
	}

	// Decode the JSON response into the userResponse struct.
	err = json.NewDecoder(resp.Body).Decode(&userResponse)

	if err != nil {
		// If there's an error, respond with a 500 Internal Server Error status code and return.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// At this point, you have the user's GitHub username in userResponse.Login.
	// You can use this to look up the user in your database, create a new user, etc.
}
