package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int64, secret string, expiry int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}