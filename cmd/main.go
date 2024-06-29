package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TS22082/dat_board_server/api/auth"
	"github.com/TS22082/dat_board_server/api/user"
	"github.com/TS22082/dat_board_server/db"
	utils "github.com/TS22082/dat_board_server/scripts/utilities"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to load .env file: %v\n", err)
	}

	dbClient, dbCtx := db.Connect()

	defer func() {
		if err = dbClient.Disconnect(dbCtx); err != nil {
			fmt.Printf("Failed to disconnect from MongoDB: %v\n", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/github/gh_login", utils.PassDbToClient(auth.HandleGhLogin, dbClient))
	mux.HandleFunc("GET /api/verify_jwt", utils.PassDbToClient(auth.VerifyJWTHandler, dbClient))
	mux.HandleFunc("GET /api/user", utils.PassDbToClient(user.GetUserByTokenHandler, dbClient))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serverErrors := make(chan error, 1)

	go func() {
		fmt.Println("Server is starting on :8080")
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		fmt.Printf("Error starting server: %v\n", err)
	case <-shutdown:
		fmt.Println("Shutdown signal received, stopping server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Error during server shutdown: %v\n", err)
		}
	}

	fmt.Println("Server stopped")
}
