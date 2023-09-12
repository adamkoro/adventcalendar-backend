package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(username string, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(86400 * time.Second).Unix(),
		"user": username,
	})
	signToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signToken, nil
}

func ValidateJWT(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
