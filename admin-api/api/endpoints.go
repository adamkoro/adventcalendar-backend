package api

import (
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/model"
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var (
	Db       pg.Repository
	validate = validator.New()
)

func CreateUser(c *gin.Context) {
	var data pg.CreateUserRequest
	log.Debug().Msg("binding request body...")
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&data); validationErr != nil {
		log.Error().Msg(validationErr.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	log.Debug().Msg("validation successful")
	log.Debug().Msg("establishing connection to the database..")
	err := Db.Ping()
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "can not establish connection to the database"})
		return
	}
	log.Debug().Msg("establishing connection to the database successful")
	log.Debug().Msg("creating user...")
	err = Db.CreateUser(data.Username, data.Email, data.Password)
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "while creating user"})
		return
	}
	log.Debug().Msg("user created successfully")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "user created successfully"})
}

func GetUser(c *gin.Context) {
	var data pg.UserRequest
	var getuserresp pg.UserResponse
	log.Debug().Msg("binding request body...")
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&data); validationErr != nil {
		log.Error().Msg(validationErr.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	log.Debug().Msg("validation successful")
	log.Debug().Msg("establishing connection to the database..")
	err := Db.Ping()
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "can not establish connection to the database"})
		return
	}
	log.Debug().Msg("establishing connection to the database successful")
	log.Debug().Msg("getting user...")
	user, err := Db.GetUser(data.Username)
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "while getting requested user"})
		return
	}
	log.Debug().Msg("user retrieved successfully")
	getuserresp.Id = int(user.Key)
	getuserresp.Email = user.Email
	getuserresp.Created = user.CreatedAt.String()
	getuserresp.Modified = user.ModifiedAt.String()
	c.JSON(http.StatusOK, &getuserresp)
}

func GetAllUsers(c *gin.Context) {
	var getallusersresp []pg.UserResponse
	log.Debug().Msg("establishing connection to the database..")
	err := Db.Ping()
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "can not establish connection to the database"})
		return
	}
	log.Debug().Msg("establishing connection to the database successful")
	log.Debug().Msg("getting all users...")
	users, err := Db.GetAllUsers()
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "while getting all users"})
		return
	}
	log.Debug().Msg("all users retrieved successfully")
	log.Debug().Msg("creating response...")
	for _, user := range users {
		var userresp pg.UserResponse
		userresp.Id = int(user.Key)
		userresp.Username = user.Username
		userresp.Email = user.Email
		userresp.Created = user.CreatedAt.String()
		userresp.Modified = user.ModifiedAt.String()
		getallusersresp = append(getallusersresp, userresp)
	}
	log.Debug().Msg("response created successfully")
	c.JSON(http.StatusOK, &getallusersresp)
}

func UpdateUser(c *gin.Context) {
	var data pg.UpdateUserRequest
	log.Debug().Msg("binding request body...")
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&data); validationErr != nil {
		log.Error().Msg(validationErr.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	log.Debug().Msg("validation successful")
	log.Debug().Msg("establishing connection to the database..")
	err := Db.Ping()
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "can not establish connection to the database"})
		return
	}
	log.Debug().Msg("establishing connection to the database successful")
	log.Debug().Msg("updating user...")
	err = Db.UpdateUser(data.Username, data.Email, data.Password)
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "while updating user"})
		return
	}
	log.Debug().Msg("user updated successfully")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "user updated successfully"})
}

func DeleteUser(c *gin.Context) {
	var data pg.UserRequest
	log.Debug().Msg("binding request body...")
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&data); validationErr != nil {
		log.Error().Msg(validationErr.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	log.Debug().Msg("validation successful")
	log.Debug().Msg("establishing connection to the database..")
	err := Db.Ping()
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "can not establish connection to the database"})
		return
	}
	log.Debug().Msg("establishing connection to the database successful")
	log.Debug().Msg("checking if user is not admin...")
	if data.Username == "admin" {
		log.Error().Msg("can not delete admin user")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "can not delete admin user"})
		return
	}
	log.Debug().Msg("user is not admin")
	log.Debug().Msg("deleting user...")
	err = Db.DeleteUser(data.Username)
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "while deleting user"})
		return
	}
	log.Debug().Msg("user deleted successfully")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "user deleted successfully"})
}
