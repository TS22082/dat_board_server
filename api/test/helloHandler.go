package test

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

// HelloHandler is a simple http handler that writes "Get request 1!" to the response writer.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Message: "Get request 1!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
