package middleware

import (
	"log"
	"net/http"
)

func CorsRules(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		(w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		log.Println("allowing all cors")
		next.ServeHTTP(w, r)
	})
}
