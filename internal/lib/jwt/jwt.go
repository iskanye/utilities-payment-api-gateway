package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type claims struct {
	UserID  int64  `json:"uid"`
	IsAdmin bool   `json:"is_admin"`
	Email   string `json:"email"`
	jwt.StandardClaims
}

type TokenPayload struct {
	UserID  int64 `json:"uid"`
	IsAdmin bool  `json:"is_admin"`
}

type TokenSaver interface {
	Set(val string) error
}

type TokenProvider interface {
	Get(key string) error
}

func ValidateToken(tokenStr string, secret string) (TokenPayload, error) {
	var claims claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return TokenPayload{}, err
	}

	if !token.Valid {
		return TokenPayload{}, fmt.Errorf("invalid token")
	}

	return TokenPayload{UserID: claims.UserID, IsAdmin: claims.IsAdmin}, nil
}
