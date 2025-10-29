package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type claims struct {
	UserID int64  `json:"uid"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func ValidateToken(tokenStr string, secret string) (int64, error) {
	var claims claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	/*userID, ok := claims["uid"].(int)
	if !ok {
		return 0, fmt.Errorf("cannot convert uid to int")
	}*/

	return claims.UserID, nil
}
