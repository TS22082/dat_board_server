package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateJWT(email string, id string) (string, error) {

	fmt.Printf("email: %s\n", email)
	fmt.Printf("id: %s\n", id)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    id,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("super_secret_key"))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
