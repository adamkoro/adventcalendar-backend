package api

import (
	"log"
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/model"
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	Db       pg.Repository
	validate = validator.New()
)

func CreateUser(c *gin.Context) {
	var data pg.CreateUserRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "invalid request body"}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	if validationErr := validate.Struct(&data); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	err := Db.CreateUser(data.Username, data.Email, data.Password)
	if err != nil {
		errormessage := "Error while creating user"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "error while creating user"}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	createuserresp := model.SuccessResponse{Status: "User created"}
	c.JSON(http.StatusOK, &createuserresp)
}

func GetUser(c *gin.Context) {
	var data pg.UserRequest
	var getuserresp pg.UserResponse
	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "invalid request body"}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	if validationErr := validate.Struct(&data); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	user, err := Db.GetUser(data.Username)
	if err != nil {
		errormessage := "Error while getting user"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "error while getting user"}
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
	var getallusersresp []pg.UserResponse

	users, err := Db.GetAllUsers()
	if err != nil {
		errormessage := "Error while getting all users"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "error while getting all users"}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	for _, user := range users {
		var userresp pg.UserResponse
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
	var data pg.UpdateUserRequest
	var updateuserresp model.SuccessResponse
	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON to struct"
		log.Println(errormessage + ":" + err.Error())
		errorresp := model.ErrorResponse{Error: "invalid request body"}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	if validationErr := validate.Struct(&data); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	err := Db.UpdateUser(data.Username, data.Email, data.Password)
	if err != nil {
		errormessage := "Error while updating user"
		log.Println(errormessage + ":" + err.Error())
		errorresp := model.ErrorResponse{Error: "error while updating user"}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	updateuserresp.Status = "User updated"
	c.JSON(http.StatusOK, &updateuserresp)
}

func DeleteUser(c *gin.Context) {
	var data pg.UserRequest
	var errorresp model.ErrorResponse
	var deleteuserresp model.SuccessResponse

	if err := c.ShouldBindJSON(&data); err != nil {
		errormessage := "Error binding JSON to struct"
		log.Println(errormessage + ":" + err.Error())
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, "invalid request body")
		return
	}
	if validationErr := validate.Struct(&data); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	if data.Username == "admin" {
		errormessage := "Cannot delete admin user"
		log.Println(errormessage)
		errorresp.Error = errormessage
		c.JSON(http.StatusBadRequest, "cannot delete admin user")
		return
	}
	err := Db.DeleteUser(data.Username)
	if err != nil {
		errormessage := "Error while deleting user"
		log.Println(errormessage + ":" + err.Error())
		errorresp.Error = errormessage
		c.JSON(http.StatusInternalServerError, "error while deleting user")
		return
	}
	deleteuserresp.Status = "User deleted"
	c.JSON(http.StatusOK, &deleteuserresp)
}
