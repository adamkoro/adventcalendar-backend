package api

import (
	"log"
	"net/http"

	custModel "github.com/adamkoro/adventcalendar-backend/lib/model"
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/gin-gonic/gin"
)

var Db pg.Repository

func CreateUser(c *gin.Context) {
	var data custModel.CreateUserRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON"
		log.Println(errormessage+" : ", err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := Db.CreateUser(data.Username, data.Email, data.Password)
	if err != nil {
		errormessage := "Error while creating user"
		log.Println(errormessage+" : ", err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	createuserresp := custModel.SuccessResponse{Status: "User created"}
	log.Println(createuserresp.Status)
	c.JSON(http.StatusOK, &createuserresp)
}

func GetUser(c *gin.Context) {
	var data custModel.UserRequest
	var getuserresp custModel.UserResponse

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON"
		log.Println(errormessage+" : ", err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	user, err := Db.GetUser(data.Username)
	if err != nil {
		errormessage := "Error while getting user"
		log.Println(errormessage+" : ", err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}

	getuserresp.Id = int(user.Key)
	getuserresp.Email = user.Email
	getuserresp.Created = user.CreatedAt.String()
	getuserresp.Modified = user.ModifiedAt.String()
	c.JSON(http.StatusOK, &getuserresp)
}

func GetAllUsers(c *gin.Context) {
	var getallusersresp []custModel.UserResponse

	users, err := Db.GetAllUsers()
	if err != nil {
		errormessage := "Error while getting all users"
		log.Println(errormessage+" : ", err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}

	for _, user := range users {
		var userresp custModel.UserResponse
		userresp.Id = int(user.Key)
		userresp.Username = user.Username
		userresp.Email = user.Email
		userresp.Created = user.CreatedAt.String()
		userresp.Modified = user.ModifiedAt.String()
		getallusersresp = append(getallusersresp, userresp)
	}
	c.JSON(http.StatusOK, &getallusersresp)
}

func UpdateUser(c *gin.Context) {
	var data custModel.UpdateUserRequest
	var updateuserresp custModel.SuccessResponse

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON to struct"
		log.Println(errormessage + ":" + err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := Db.UpdateUser(data.Username, data.Email, data.Password)
	if err != nil {
		errormessage := "Error while updating user"
		log.Println(errormessage + ":" + err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	updateuserresp.Status = "User updated"
	log.Println(updateuserresp.Status)
	c.JSON(http.StatusOK, &updateuserresp)
}

func DeleteUser(c *gin.Context) {
	var data custModel.UserRequest
	var errorresp custModel.ErrorResponse
	var deleteuserresp custModel.SuccessResponse

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON to struct"
		log.Println(errormessage + ":" + err.Error())
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	if data.Username == "admin" {
		errormessage := "Cannot delete admin user"
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	err := Db.DeleteUser(data.Username)
	if err != nil {
		errormessage := "Error while deleting user"
		log.Println(errormessage + ":" + err.Error())
		errorresp.Error = errormessage
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	deleteuserresp.Status = "User deleted"
	log.Println(deleteuserresp.Status)
	c.JSON(http.StatusOK, &deleteuserresp)
}
