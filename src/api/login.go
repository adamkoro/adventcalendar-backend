package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adamkoro/adventcalendar-backend/postgres"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	SecretKey string
	Rd        *redis.Client
)

func Login(c *gin.Context) {
	var data LoginRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := postgres.Login(Db, data.Username, data.Password)
	if err != nil {
		errormessage := "Username or password incorrect"
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}

	token, err := generateJWT(data.Username)
	if err != nil {
		errormessage := "Error generating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	err = createSession(Rd, data.Username, token, c.ClientIP())
	if err != nil {
		errormessage := "Error creating session: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	loginresp := SuccessResponse{Status: "Login successful"}
	log.Println(loginresp.Status)
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
	cookie, err := c.Cookie("token")
	if err != nil {
		errormessage := "Error getting cookie: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	claims, err := validateJWT(cookie, c)
	if err != nil {
		errormessage := "Error validating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	if claims["authorized"] == true {
		c.Next()
	} else {
		errormessage := "Unauthorized"
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
}

func createSession(rd *redis.Client, username string, token string, sourceIp string) error {
	session := Session{
		Username: username,
		Token:    token,
		SourceIP: sourceIp,
		LoginAt:  time.Now().String(),
	}
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return rd.Set(context.Background(), uuid.New().String(), data, 86400*time.Second).Err()
}
