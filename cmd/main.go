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
	"github.com/TS22082/dat_board_server/scripts/middleware"
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

	router := http.NewServeMux()
	protectedRouter := http.NewServeMux()

	router.HandleFunc("GET /api/github/gh_login", utils.PassDbTohandler(auth.HandleGhLogin, dbClient))
	router.HandleFunc("GET /api/verify_jwt", utils.PassDbTohandler(auth.VerifyJWTHandler, dbClient))

	protectedRouter.HandleFunc("GET /api/user", utils.PassDbTohandler(user.GetUserByTokenHandler, dbClient))

	router.Handle("/", middleware.EnsureAuth(protectedRouter, dbClient))

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.CorsRules,
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: stack(router),
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
