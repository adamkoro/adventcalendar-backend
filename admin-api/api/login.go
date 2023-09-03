package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adamkoro/adventcalendar-backend/lib/postgres"
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
	session_uuid := uuid.New().String()
	token, err := generateJWT(data.Username, session_uuid)
	if err != nil {
		errormessage := "Error generating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	err = createSession(Rd, data.Username, token, c.ClientIP(), session_uuid)
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

func Logout(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		errormessage := "Error getting cookie: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	claims, err := validateJWT(cookie, c)
	if err != nil {
		errormessage := "Error validating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	if claims["authorized"] != true {
		errormessage := "Unauthorized"
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}
	if claims["session"] == "" {
		errormessage := "Session not found"
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}
	err = Rd.Del(context.Background(), claims["session"].(string)).Err()
	if err != nil {
		errormessage := "Error deleting session from redis: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	logoutresp := SuccessResponse{Status: "Logout successful"}
	log.Println(logoutresp.Status)
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &logoutresp)
}

func generateJWT(username string, session string) (string, error) {
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
	if claims["authorized"] != true {
		errormessage := "Unauthorized"
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	if claims["session"] == "" {
		errormessage := "Session not found"
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	backend_session, err := Rd.Get(context.Background(), claims["session"].(string)).Bytes()
	if err != nil {
		errormessage := "Error getting session from redis: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusInternalServerError, &errorresp)
		return
	}
	var tokenSession Session
	err = json.Unmarshal(backend_session, &tokenSession)
	if err != nil {
		errormessage := "Error unmarshalling session: " + err.Error()
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusInternalServerError, &errorresp)
		return
	}
	if cookie != tokenSession.Token {
		errormessage := "Session invalid"
		log.Println(errormessage)
		errorresp := ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	c.Next()
}

func createSession(rd *redis.Client, username string, token string, sourceIp string, session_uuid string) error {
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
	return rd.Set(context.Background(), session_uuid, data, 86400*time.Second).Err()
}
