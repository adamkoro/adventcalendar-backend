package jwt

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	SecretKey string
)

func GenerateJWT(username string, session string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(86400 * time.Second).Unix(),
		"authorized": true,
		"user":       username,
		"session":    session,
	})
	signToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return signToken, nil
}

func ValidateJWT(tokenString string, c *gin.Context) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
