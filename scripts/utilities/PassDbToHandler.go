package utils

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func PassDbTohandler(handler func(http.ResponseWriter, *http.Request, *mongo.Client), client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, client)
	}
}
