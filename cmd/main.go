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

	client := db.Connect()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/github/gh_login", utils.PassDbToClient(auth.HandleGhLogin, client))
	mux.HandleFunc("GET /api/verify_jwt", utils.PassDbToClient(auth.VerifyJWTHandler, client))
	mux.HandleFunc("GET /api/user", utils.PassDbToClient(user.GetUserByTokenHandler, client))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		fmt.Println("Server is starting on :8080")
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		fmt.Printf("Error starting server: %v\n", err)
	case <-shutdown:
		fmt.Println("Shutdown signal received, stopping server...")

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		server.Shutdown(ctx)
	}

	fmt.Println("Server stopped")
	os.Exit(0)
}
