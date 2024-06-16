package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func VerifyJWT(token string) (bool, error) {

	// remove the "Bearer " prefix from the token
	token = token[7:]

	decriptedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println("Error verifying token: ", err)
		return false, err
	}

	claims, ok := decriptedToken.Claims.(jwt.MapClaims)

	if !ok || !decriptedToken.Valid {
		return false, nil
	}

	exp := claims["exp"]

	if int64(exp.(float64)) < time.Now().Unix() {
		return false, nil
	}

	if decriptedToken.Valid {
		return true, nil
	}

	return false, nil
}
