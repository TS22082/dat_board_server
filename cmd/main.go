package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TS22082/dat_board_server/api/auth"
	"github.com/TS22082/dat_board_server/api/test"
	"github.com/TS22082/dat_board_server/db"
	utils "github.com/TS22082/dat_board_server/scripts/utilities"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Failed to load .env file: %v", err)
	}

	client := db.Connect()

	collectionToGet := client.Database("dat_board").Collection("test")

	itemsInTestCollection, err := collectionToGet.Find(context.Background(), bson.D{})

	var results []bson.M

	if err != nil {
		fmt.Printf("Failed to get items from test collection: %v", err)
	}

	if err = itemsInTestCollection.All(context.Background(), &results); err != nil {
		fmt.Printf("Failed to decode items from test collection: %v", err)
	}

	http.HandleFunc("/api", test.HelloHandler)
	http.HandleFunc("/api/2", test.HelloHandler2)

	http.HandleFunc("/api/github/gh_login", utils.PassDbToClient(auth.HandleGhLogin, client))

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Server failed to start: %v", err)
	}
}
