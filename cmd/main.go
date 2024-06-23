package main

import (
	"fmt"
	"net/http"

	"github.com/TS22082/dat_board_server/api/auth"
	"github.com/TS22082/dat_board_server/api/test"
	"github.com/TS22082/dat_board_server/api/user"
	"github.com/TS22082/dat_board_server/db"
	utils "github.com/TS22082/dat_board_server/scripts/utilities"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Failed to load .env file: %v", err)
	}

	client := db.Connect()

	http.HandleFunc("/api", test.HelloHandler)
	http.HandleFunc("/api/2", test.HelloHandler2)

	http.HandleFunc("GET /api/github/gh_login", utils.PassDbToClient(auth.HandleGhLogin, client))
	http.HandleFunc("GET /api/verify_jwt", utils.PassDbToClient(auth.VerifyJWTHandler, client))

	http.HandleFunc("GET /api/user", utils.PassDbToClient(user.GetUserByTokenHandler, client))

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Server failed to start: %v", err)
	}
}
