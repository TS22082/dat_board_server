package utils

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByToken(dbCollection *mongo.Collection, token string) (map[string]interface{}, error) {
	// remove the "Bearer " prefix from the token
	var err error

	token = token[7:]

	decriptedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := decriptedToken.Claims.(jwt.MapClaims)

	if !ok || !decriptedToken.Valid {
		return nil, nil
	}

	formattedToken := map[string]interface{}{
		"email": claims["email"],
		"id":    claims["id"],
		"exp":   claims["exp"],
	}

	return formattedToken, nil
}
