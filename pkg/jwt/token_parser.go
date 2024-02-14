package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	ErrInvalidAccessToken = fmt.Errorf("invalid access token")
	ErrInvalidClaims      = fmt.Errorf("invalid claims")
)

func ParseToken(tokenString string, secret []byte) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidAccessToken
		}

		return secret, nil
	})

	if err != nil || !token.Valid {
		return Claims{}, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return *claims, nil
	}

	return Claims{}, ErrInvalidClaims
}
