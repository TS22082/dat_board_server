package user

import (
	"encoding/json"
	"net/http"

	"github.com/TS22082/dat_board_server/scripts/middleware"
	utils "github.com/TS22082/dat_board_server/scripts/utilities"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByTokenHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	middleware.EnableCors(&w)

	token := r.Header.Get("Authorization")

	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Token is missing",
		})
		return
	}

	verified, err := utils.VerifyJWT(token)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Token is invalid",
		})
		return
	}

	if !verified {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Token is invalid",
		})
		return
	}

	user, err := utils.GetUserByToken(client.Database("dat_board").Collection("users"), token)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Failed to get user",
		})
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	response := map[string]interface{}{
		"error": false,
		"uesr":  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
