package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/TS22082/dat_board_server/scripts/middleware"
)

type GithubResponse struct {
	StatusCode int                    `json:"status_code"`
	Body       map[string]interface{} `json:"body"`
}

func HandleGhLogin(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	hasCodeParam := r.URL.Query().Has("code")

	if !hasCodeParam {
		http.Error(w, "No code parameter in query string", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "Code parameter is empty", http.StatusBadRequest)
		return
	}

	url := "https://github.com/login/oauth/access_token"
	payload := map[string]string{
		"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
		"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		"code":          code,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
		return
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		http.Error(w, "Failed to create request to GitHub", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get access token from GitHub", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		http.Error(w, "Failed to unmarshal response body", http.StatusInternalServerError)
		return
	}

	response := GithubResponse{
		StatusCode: resp.StatusCode,
		Body:       result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
