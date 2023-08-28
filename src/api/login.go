package api

import (
	"net/http"

	"log"

	postgres "github.com/adamkoro/adventcalendar-backend/postgres"
	"github.com/gin-gonic/gin"
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
	loginresp.Status = "Login successful"
	log.Println(loginresp.Status)
	c.JSON(http.StatusOK, &loginresp)
}
