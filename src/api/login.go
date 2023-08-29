package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adamkoro/adventcalendar-backend/postgres"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)

var (
	SecretKey string
	Rd        *redis.Client
)

func Login(c *gin.Context) {
	var data LoginRequest
	var errorresp ErrorResponse
	var loginresp SuccessResponse

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := postgres.Login(Db, data.Username, data.Password)
	if err != nil {
		errormessage := "Username or/and password incorrect"
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}

	token, err := generateJWT(data.Username)
	if err != nil {
		errormessage := "Error generating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	loginresp.Status = "Login successful"
	log.Println(loginresp.Status)
	createSession(Rd, data.Username)
	c.SetCookie("token", token, 86400, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &loginresp)
}

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(86400 * time.Second).Unix(),
		"authorized": true,
		"user":       username,
	})
	signToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return signToken, nil
}

func validateJWT(tokenString string, c *gin.Context) (jwt.MapClaims, error) {
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

func AuthRequired(c *gin.Context) {
	var errorresp ErrorResponse
	cookie, err := c.Cookie("token")
	if err != nil {
		errormessage := "Error getting cookie: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	claims, err := validateJWT(cookie, c)
	if err != nil {
		errormessage := "Error validating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	if claims["authorized"] == true {
		c.Next()
	} else {
		errormessage := "Unauthorized"
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
}

func createSession(rd *redis.Client, username string) error {
	return rd.Set(context.Background(), username, true, 86400*time.Second).Err()
}
