package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUserWithEmail(dbCollection *mongo.Collection, email string) (map[string]interface{}, error) {
	var err error

	insertUser := bson.D{
		{Key: "email", Value: email},
	}

	insertResult, err := dbCollection.InsertOne(context.Background(), insertUser)

	if err != nil {
		return nil, err
	}

	newUser := map[string]interface{}{
		"email": email,
		"_id":   insertResult.InsertedID,
	}

	return newUser, err
}
