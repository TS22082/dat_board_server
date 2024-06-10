package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/TS22082/dat_board_server/scripts/middleware"
	utils "github.com/TS22082/dat_board_server/scripts/utilities"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response struct {
	StatusCode int                    `json:"status_code"`
	Body       map[string]interface{} `json:"body"`
}

func HandleGhLogin(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
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

	if ghAuthResults["access_token"] == nil {
		http.Error(w, fmt.Sprintf("Failed to get access token: %v", ghAuthResults), http.StatusInternalServerError)
		return
	}

	emails, err := utils.GetUserEmails(ghAuthResults["access_token"].(string))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user emails: %v", err), http.StatusInternalServerError)
		return
	}

	primaryEmail := string("")

	for _, email := range emails {
		if email["primary"] == true {
			primaryEmail = email["email"].(string)
			break
		}
	}

	fmt.Printf("Primary Email ==> %v\n", primaryEmail)

	response.Body["emails"] = emails

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
