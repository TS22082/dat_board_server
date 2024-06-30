package user

import (
	"encoding/json"
	"net/http"

	"github.com/TS22082/dat_board_server/scripts/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByTokenHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	userFromAuth, ok := r.Context().Value(middleware.AuthenticatedUser).(map[string]interface{})

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Problems getting user from middleware",
		})
		return
	}

	response := map[string]interface{}{
		"error": false,
		"user":  userFromAuth,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
