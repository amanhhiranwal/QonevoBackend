package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// 🔐 Generate token
func GenerateToken(userID int64, secret string, expiry int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Hour * time.Duration(expiry)),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// 🔍 Parse token
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}


func ValidateToken(
	tokenString string,
) (*jwt.Token, error) {

	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, errors.New("invalid token")
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}