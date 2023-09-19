package api

import (
	"log"
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	custJWT "github.com/adamkoro/adventcalendar-backend/lib/jwt"
	custModel "github.com/adamkoro/adventcalendar-backend/lib/model"
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/gin-gonic/gin"
)

var Db pg.Repository

func Login(c *gin.Context) {
	var data custModel.LoginRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(c.Request.Body.Read([]byte{}))
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := Db.Login(data.Username, data.Password)
	if err != nil {
		errormessage := "Username or password incorrect"
		log.Println(errormessage+" : ", err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}
	token, err := custJWT.GenerateJWT(data.Username, env.GetSecretKey())
	if err != nil {
		errormessage := "Error generating JWT token"
		log.Println(errormessage+" : ", err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	loginresp := custModel.SuccessResponse{Status: "Login successful"}
	log.Println(loginresp.Status)
	c.SetCookie("token", token, 86400, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &loginresp)
}

func Logout(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		errormessage := "Error getting cookie"
		log.Println(errormessage+" : ", err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	_, err = custJWT.ValidateJWT(cookie, env.GetSecretKey())
	if err != nil {
		errormessage := "Error validating JWT token"
		log.Println(errormessage+" : ", err.Error())
		errorresp := custModel.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	logoutresp := custModel.SuccessResponse{Status: "Logout successful"}
	log.Println(logoutresp.Status)
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &logoutresp)
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func ApiLogin(c *gin.Context) {

}
