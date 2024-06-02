package test

import (
	"encoding/json"
	"net/http"
)

type Response2 struct {
	Message string `json:"message"`
}

// HelloHandler2 is a simple http handler that writes "Get request 2!" to the response writer.
func HelloHandler2(w http.ResponseWriter, r *http.Request) {
	res := Response2{
		Message: "Get request 2!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
