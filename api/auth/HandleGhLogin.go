package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/TS22082/dat_board_server/scripts/middleware"
	utils "github.com/TS22082/dat_board_server/scripts/utilities"
)

type Response struct {
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
	if code == "" || code == "null" {
		http.Error(w, "Code parameter is missing or empty", http.StatusBadRequest)
		return
	}

	url := "https://github.com/login/oauth/access_token"
	ghAuthPayload := map[string]string{
		"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
		"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		"code":          code,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	ghAuthParams := utils.HTTPRequestParams{
		URL:     url,
		Method:  "POST",
		Headers: headers,
		Body:    ghAuthPayload,
	}

	ghAuthResults, statusCode, err := utils.MakeHTTPRequest(ghAuthParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get access token: %v", err), http.StatusInternalServerError)
		return
	}

	response := Response{
		StatusCode: statusCode,
		Body:       ghAuthResults,
	}

	ghUserParams := utils.HTTPRequestParams{
		URL:    "https://api.github.com/user",
		Method: "GET",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("token %v", ghAuthResults["access_token"]),
		},
	}

	ghUserResults, _, err := utils.MakeHTTPRequest(ghUserParams)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user data: %v", err), http.StatusInternalServerError)
		return
	}

	response.Body["user"] = ghUserResults

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
