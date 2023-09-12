package api

import (
	"log"
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	custJWT "github.com/adamkoro/adventcalendar-backend/lib/jwt"
	custModel "github.com/adamkoro/adventcalendar-backend/lib/model"
	"github.com/gin-gonic/gin"
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
	_, err = custJWT.ValidateJWT(cookie, env.GetSecretKey())
	if err != nil {
		errormessage := "Error validating JWT: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.AbortWithStatusJSON(http.StatusUnauthorized, &errorresp)
		return
	}
	c.Next()
}
