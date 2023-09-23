package api

import (
	"log"
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	custJWT "github.com/adamkoro/adventcalendar-backend/lib/jwt"
	"github.com/adamkoro/adventcalendar-backend/lib/model"
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(c.Request.Body.Read([]byte{}))
		errormessage := "error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp := model.ErrorResponse{Error: "invalid request body"}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	if validationErr := validate.Struct(&data); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	err := Db.Ping()
	if err != nil {
		errormessage := "error connecting to database"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "database connection error"}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	err = Db.Login(data.Username, data.Password)
	if err != nil {
		errormessage := "username or password incorrect"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}
	token, err := custJWT.GenerateJWT(data.Username, env.GetSecretKey())
	if err != nil {
		errormessage := "error generating JWT token"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "generating token"}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	loginresp := model.SuccessResponse{Status: "login successful"}
	log.Println(loginresp.Status)
	c.SetCookie("token", token, 86400, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &loginresp)
}

// @Summary Admin user logout
// @Description Logout admin user via cookie and get empty cookie
// @Tags auth
// @Produce json
// @Success 200 {object} model.SuccessResponse "logout successful"
// @Failure 400 {object} model.ErrorResponse "Cookie not found or JWT validation error"
// @Router /api/admin/logout [post]
func Logout(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		errormessage := "error getting cookie"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "cookie not found, please login again"}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	_, err = custJWT.ValidateJWT(cookie, env.GetSecretKey())
	if err != nil {
		errormessage := "error validating JWT token"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "invalid token, please login again"}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	logoutresp := model.SuccessResponse{Status: "logout successful"}
	log.Println(logoutresp.Status)
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &logoutresp)
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
