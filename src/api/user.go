package api

import (
	"log"
	"net/http"

	postgres "github.com/adamkoro/adventcalendar-backend/postgres"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Db *gorm.DB

func CreateUser(c *gin.Context) {
	var data createUserRequest
	var errorresp errorResponse
	var createuserresp successResponse

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := postgres.CreateUser(Db, data.Username, data.Email, data.Password)
	if err != nil {
		errormessage := "Error while creating user: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	createuserresp.Status = "User created"
	log.Println(createuserresp.Status)
	c.JSON(http.StatusOK, &createuserresp)
}

func GetUser(c *gin.Context) {
	var data getUserRequest
	var errorresp errorResponse
	var getuserresp getUserResponse

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	user, err := postgres.GetUser(Db, data.Username)
	if err != nil {
		errormessage := "Error while getting user: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	getuserresp.Username = user.Username
	getuserresp.Email = user.Email
	c.JSON(http.StatusOK, &getuserresp)
}
