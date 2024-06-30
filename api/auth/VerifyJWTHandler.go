package auth

import (
	"encoding/json"
	"net/http"

	utils "github.com/TS22082/dat_board_server/scripts/utilities"
	"go.mongodb.org/mongo-driver/mongo"
)

var TokenMissingError = map[string]interface{}{
	"error":   true,
	"message": "Token is missing",
}

var TokenInvalidError = map[string]interface{}{
	"error":   true,
	"message": "Token is invalid",
}

func VerifyJWTHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}

	token := r.Header.Get("Authorization")

	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TokenMissingError)
		return
	}

	verified, err := utils.VerifyJWT(token)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TokenInvalidError)
		return
	}

	if !verified {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TokenInvalidError)
		return
	}

	response := map[string]interface{}{
		"verified": true,
		"error":    false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
