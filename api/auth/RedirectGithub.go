package auth

import (
	"net/http"
	"os"
)

func RedirectGithub(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")

	url := "https://github.com/login/oauth/authorize?client_id=" + clientID

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
