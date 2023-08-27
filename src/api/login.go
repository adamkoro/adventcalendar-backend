package api

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var data loginRequest
	var errorresp errorResponse
	var loginresp successResponse
	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	loginresp.Status = "Logged in"
	c.JSON(http.StatusOK, &loginresp)
}
