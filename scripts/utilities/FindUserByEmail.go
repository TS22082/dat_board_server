package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUserByEmail(dbCollection *mongo.Collection, email string) (*mongo.SingleResult, error) {
	var err error

	userFound := dbCollection.FindOne(context.Background(), bson.D{{Key: "email", Value: email}})

	if userFound.Err() != nil {
		err = userFound.Err()
	}

	return userFound, err

}
