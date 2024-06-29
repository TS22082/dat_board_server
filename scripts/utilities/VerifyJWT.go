package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func VerifyJWT(token string) (bool, error) {
	decriptedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return false, err
	}

	claims, ok := decriptedToken.Claims.(jwt.MapClaims)

	if !ok || !decriptedToken.Valid {
		return false, errors.New("token is not valid")
	}

	exp := claims["exp"]

	if int64(exp.(float64)) < time.Now().Unix() {
		return false, errors.New("token is expired")
	}

	if decriptedToken.Valid {
		return true, nil
	}

	return false, errors.New("token did not pass validation check")
}
