package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	custJWT "github.com/adamkoro/adventcalendar-backend/lib/jwt"
	custModel "github.com/adamkoro/adventcalendar-backend/lib/model"
	"github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	SecretKey string
	Rd        *redis.Client
)

func Login(c *gin.Context) {
	var data custModel.LoginRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := postgres.Login(Db, data.Username, data.Password)
	if err != nil {
		errormessage := "Username or password incorrect"
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}
	session_uuid := uuid.New().String()
	token, err := custJWT.GenerateJWT(data.Username, session_uuid)
	if err != nil {
		errormessage := "Error generating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	err = createSession(Rd, data.Username, token, c.ClientIP(), session_uuid)
	if err != nil {
		errormessage := "Error creating session: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	loginresp := custModel.SuccessResponse{Status: "Login successful"}
	log.Println(loginresp.Status)
	c.SetCookie("token", token, 86400, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &loginresp)
}

func Logout(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		errormessage := "Error getting cookie: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	claims, err := custJWT.ValidateJWT(cookie, c)
	if err != nil {
		errormessage := "Error validating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	if claims["authorized"] != true {
		errormessage := "Unauthorized"
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}
	if claims["session"] == "" {
		errormessage := "Session not found"
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}
	err = Rd.Del(context.Background(), claims["session"].(string)).Err()
	if err != nil {
		errormessage := "Error deleting session from redis: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	logoutresp := custModel.SuccessResponse{Status: "Logout successful"}
	log.Println(logoutresp.Status)
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &logoutresp)
}

func AuthRequired(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		errormessage := "Error getting cookie: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: "cookie not found"}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	claims, err := custJWT.ValidateJWT(cookie, c)
	if err != nil {
		errormessage := "Error validating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	if claims["authorized"] != true {
		errormessage := "Unauthorized"
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	if claims["session"] == "" {
		errormessage := "Session not found"
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	backend_session, err := Rd.Get(context.Background(), claims["session"].(string)).Bytes()
	if err != nil {
		errormessage := "Error getting session from redis: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusInternalServerError, &errorresp)
		return
	}
	var tokenSession custModel.Session
	err = json.Unmarshal(backend_session, &tokenSession)
	if err != nil {
		errormessage := "Error unmarshalling session: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusInternalServerError, &errorresp)
		return
	}
	if cookie != tokenSession.Token {
		errormessage := "Session invalid"
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	c.Next()
}

func createSession(rd *redis.Client, username string, token string, sourceIp string, session_uuid string) error {
	session := custModel.Session{
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
