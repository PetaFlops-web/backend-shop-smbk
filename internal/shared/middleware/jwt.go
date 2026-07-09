package middleware

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateJWT(jwtSecret string, id string, username string) (string, error) {
	if jwtSecret == "" {
		return "", errors.New("JWT secret not configured")
	}
	
    claims := jwt.MapClaims{
        "id":       id,
        "username": username,
        "exp":      time.Now().Add(time.Hour * 72).Unix(),
        "iat":      time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}