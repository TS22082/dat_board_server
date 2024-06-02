// Package auth contains authentication related functions.
package auth

import (
	"net/http" // Package http provides HTTP client and server implementations.
	"os"       // Package os provides a platform-independent interface to operating system functionality.
)

// RedirectGithub is a function that redirects the user to GitHub for OAuth authentication.
func RedirectGithub(w http.ResponseWriter, r *http.Request) {
	// Get the GitHub client ID from the environment variables.
	clientID := os.Getenv("GITHUB_CLIENT_ID")

	// Construct the URL for GitHub's OAuth authorization endpoint.
	// The client ID is included as a query parameter.
	url := "https://github.com/login/oauth/authorize?client_id=" + clientID

	// Redirect the user to the constructed URL.
	// http.StatusTemporaryRedirect is used as the status code, which means the redirect is only temporary.
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
