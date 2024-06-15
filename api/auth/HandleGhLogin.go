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
	AccessToken  string `json:"access_token"`
	PrimaryEmail string `json:"primary_email"`
	Error        bool   `json:"error"`
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

	if statusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to get access token: %v", ghAuthResults), http.StatusInternalServerError)
		return
	}

	if ghAuthResults["access_token"] == nil {
		http.Error(w, fmt.Sprintf("Failed to get access token: %v", ghAuthResults), http.StatusInternalServerError)
		return
	}

	email, err := utils.GetUserEmail(ghAuthResults["access_token"].(string))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user emails: %v", err), http.StatusInternalServerError)
		return
	}

	userCollection := client.Database("dat_board").Collection("Users")
	userFound, err := utils.FindUserByEmail(userCollection, email)

	if err == mongo.ErrNoDocuments {
		newUser, err := utils.CreateUserWithEmail(userCollection, email)

		if err != nil {
			message := map[string]interface{}{
				"error":   true,
				"message": err.Error(),
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(message)
		}

		message := map[string]interface{}{
			"error":     false,
			"message":   "User created successfully",
			"newUser":   newUser,
			"authToken": ghAuthResults["access_token"].(string),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
		return
	}

	if err != nil {
		message := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
		return
	}

	type userModel struct {
		Email string `bson:"email"`
		ID    string `bson:"_id"`
	}

	var user userModel

	err = userFound.Decode(&user)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode user: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("{\n  Email: %s\n  ID: %s\n}", user.Email, user.ID)

	response := Response{
		AccessToken:  ghAuthResults["access_token"].(string),
		PrimaryEmail: email,
		Error:        false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
