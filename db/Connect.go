package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
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

	return client
}
