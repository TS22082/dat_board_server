package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TS22082/dat_board_server/api/test"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/dat_board")

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		fmt.Printf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	http.HandleFunc("/", test.HelloHandler)
	http.HandleFunc("/2", test.HelloHandler2)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Server failed to start: %v", err)
	}
}
