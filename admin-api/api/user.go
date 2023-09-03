package api

import (
	"log"
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Db *gorm.DB

func CreateUser(c *gin.Context) {
	var data CreateUserRequest
	var errorresp ErrorResponse
	var createuserresp SuccessResponse

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
	var data UserRequest
	var errorresp ErrorResponse
	var getuserresp UserResponse

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
	getuserresp.Created = user.CreatedAt.String()
	getuserresp.Modified = user.ModifiedAt.String()
	c.JSON(http.StatusOK, &getuserresp)
}

func GetAllUsers(c *gin.Context) {
	var errorresp ErrorResponse
	var getallusersresp []UserResponse

	users, err := postgres.GetAllUsers(Db)
	if err != nil {
		errormessage := "Error while getting all users: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}

	for _, user := range users {
		var userresp UserResponse
		userresp.Username = user.Username
		userresp.Email = user.Email
		userresp.Created = user.CreatedAt.String()
		userresp.Modified = user.ModifiedAt.String()
		getallusersresp = append(getallusersresp, userresp)
	}
	c.JSON(http.StatusOK, &getallusersresp)
}

func UpdateUser(c *gin.Context) {
	var data CreateUserRequest
	var errorresp ErrorResponse
	var updateuserresp SuccessResponse

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := postgres.UpdateUser(Db, data.Username, data.Email, data.Password)
	if err != nil {
		errormessage := "Error while updating user: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	updateuserresp.Status = "User updated"
	log.Println(updateuserresp.Status)
	c.JSON(http.StatusOK, &updateuserresp)
}

func DeleteUser(c *gin.Context) {
	var data UserRequest
	var errorresp ErrorResponse
	var deleteuserresp SuccessResponse

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := postgres.DeleteUser(Db, data.Username)
	if err != nil {
		errormessage := "Error while deleting user: " + err.Error()
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	deleteuserresp.Status = "User deleted"
	log.Println(deleteuserresp.Status)
	c.JSON(http.StatusOK, &deleteuserresp)
}
