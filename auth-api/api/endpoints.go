package api

import (
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	custJWT "github.com/adamkoro/adventcalendar-backend/lib/jwt"
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

// @Summary Admin user login
// @Description Login admin user via username and password and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body pg.LoginRequest true "Login"
// @Success 200 {object} model.SuccessResponse "login successful"
// @Failure 400 {object} model.ErrorResponse "Invalid json request or validation error"
// @Failure 401 {object} model.ErrorResponse "Username or password incorrect"
// @Failure 500 {object} model.ErrorResponse "Error generating JWT token or database connection error"
// @Router /api/admin/login [post]
func Login(c *gin.Context) {
	var data pg.LoginRequest
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
	log.Debug().Msg("validating request body successful")
	log.Debug().Msg("establishing connection to the database..")
	err := Db.Ping()
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "can not establish connection to the database"})
		return
	}
	log.Debug().Msg("establishing connection to the database successful")
	log.Debug().Msg("login user...")
	err = Db.Login(data.Username, data.Password)
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "username or password incorrect"})
		return
	}
	log.Debug().Msg("login user successful from database")
	log.Debug().Msg("generating token...")
	token, err := custJWT.GenerateJWT(data.Username, env.GetSecretKey())
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "generating token"})
		return
	}
	log.Debug().Msg("generating token successful")
	log.Debug().Msg("setting cookie...")
	c.SetCookie("token", token, 86400, "/", "localhost", false, true)
	log.Debug().Msg("setting cookie successful")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "login successful"})
	log.Debug().Msg("login successful")
}

// @Summary Admin user logout
// @Description Logout admin user via cookie and get empty cookie
// @Tags auth
// @Produce json
// @Success 200 {object} model.SuccessResponse "logout successful"
// @Failure 400 {object} model.ErrorResponse "Cookie not found or JWT validation error"
// @Router /api/admin/logout [post]
func Logout(c *gin.Context) {
	log.Debug().Msg("getting cookie")
	cookie, err := c.Cookie("token")
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "cookie not found, please login again"})
		return
	}
	log.Debug().Msg("getting cookie successful")
	log.Debug().Msg("validating token")
	_, err = custJWT.ValidateJWT(cookie, env.GetSecretKey())
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid token, please login again"})
		return
	}
	log.Debug().Msg("validating token successful")
	log.Debug().Msg("deleting cookie")
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	log.Debug().Msg("deleting cookie successful")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "logout successful"})
	log.Debug().Msg("logout successful")
}

// @Summary Ping
// @Description Ping
// @Tags auth
// Produce string
// @Success 200 {string} string "pong"
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func ApiLogin(c *gin.Context) {

}
