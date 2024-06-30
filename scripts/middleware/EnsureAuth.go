package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	utils "github.com/TS22082/dat_board_server/scripts/utilities"
	"go.mongodb.org/mongo-driver/mongo"
)

type contextKey string

const AuthenticatedUser contextKey = "middleware.auth.authUser"

func EnsureAuth(next http.Handler, client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
				"message": "error getting user",
			})

			return
		}

		authedUser := map[string]interface{}{
			"id":    user["id"],
			"email": user["email"],
		}

		ctx := context.WithValue(r.Context(), AuthenticatedUser, authedUser)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
