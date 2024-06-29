package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("Failed to ping MongoDB: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected to MongoDB!")

	return client, ctx
}
