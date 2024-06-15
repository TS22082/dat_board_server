package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/TS22082/dat_board_server/scripts/middleware"
	utils "github.com/TS22082/dat_board_server/scripts/utilities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response struct {
	AccessToken  string `json:"access_token"`
	PrimaryEmail string `json:"primary_email"`
	ID           string `json:"id"`
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

	access_err_message := map[string]interface{}{
		"error":   true,
		"message": "Failed to get access token",
	}

	email_err_message := map[string]interface{}{
		"error":   true,
		"message": "Failed to get user email from token",
	}

	if err != nil {
		json.NewEncoder(w).Encode(access_err_message)
		return
	}

	if statusCode != http.StatusOK {
		json.NewEncoder(w).Encode(access_err_message)
		return
	}

	if ghAuthResults["access_token"] == nil {
		json.NewEncoder(w).Encode(access_err_message)
		return
	}

	email, err := utils.GetUserEmail(ghAuthResults["access_token"].(string))

	if err != nil {
		json.NewEncoder(w).Encode(email_err_message)
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

		jwt_token, err := utils.CreateJWT(email, newUser["_id"].(primitive.ObjectID).Hex())

		fmt.Println("jwt_token: ", jwt_token)

		if err != nil {
			message := map[string]interface{}{
				"error": true,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(message)
		}

		fmt.Println("happened here: ", newUser)

		response := Response{
			AccessToken:  jwt_token,
			PrimaryEmail: email,
			ID:           newUser["_id"].(primitive.ObjectID).Hex(),
			Error:        false,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
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
		message := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
		return
	}

	jwt_token, err := utils.CreateJWT(email, user.ID)

	if err != nil {
		message := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
	}

	response := Response{
		AccessToken:  jwt_token,
		PrimaryEmail: email,
		ID:           user.ID,
		Error:        false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
