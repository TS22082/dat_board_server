package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/TS22082/dat_board_server/scripts/middleware"
	utils "github.com/TS22082/dat_board_server/scripts/utilities"
)

type GithubResponse struct {
	StatusCode int                    `json:"status_code"`
	Body       map[string]interface{} `json:"body"`
}

func HandleGhLogin(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code parameter is missing or empty", http.StatusBadRequest)
		return
	}

	url := "https://github.com/login/oauth/access_token"
	payload := map[string]string{
		"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
		"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		"code":          code,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	params := utils.HTTPRequestParams{
		URL:     url,
		Method:  "POST",
		Headers: headers,
		Body:    payload,
	}

	result, statusCode, err := utils.MakeHTTPRequest(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get access token: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		StatusCode int                    `json:"status_code"`
		Body       map[string]interface{} `json:"body"`
	}{
		StatusCode: statusCode,
		Body:       result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
