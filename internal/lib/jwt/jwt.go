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

func ValidateToken(tokenStr string, secret string) (int64, bool, error) {
	var claims claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, false, err
	}

	if !token.Valid {
		return 0, false, fmt.Errorf("invalid token")
	}

	return claims.UserID, claims.IsAdmin, nil
}
