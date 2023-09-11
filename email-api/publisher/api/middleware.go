package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	custJWT "github.com/adamkoro/adventcalendar-backend/lib/jwt"
	custModel "github.com/adamkoro/adventcalendar-backend/lib/model"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	SecretKey string
	Rd        *redis.Client
)

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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("CORS enabled")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3030")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Access-Control-Allow-Credentials")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
