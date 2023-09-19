package api

import (
	"log"
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	custJWT "github.com/adamkoro/adventcalendar-backend/lib/jwt"
	"github.com/adamkoro/adventcalendar-backend/lib/model"
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/gin-gonic/gin"
)

var Db pg.Repository

func Login(c *gin.Context) {
	var data pg.LoginRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println(c.Request.Body.Read([]byte{}))
		errormessage := "Error binding JSON: " + err.Error()
		log.Println(errormessage)
		errorresp := model.ErrorResponse{Error: "invalid request body"}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}

	err := Db.Login(data.Username, data.Password)
	if err != nil {
		errormessage := "Username or password incorrect"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: errormessage}
		c.JSON(http.StatusUnauthorized, &errorresp)
		return
	}
	token, err := custJWT.GenerateJWT(data.Username, env.GetSecretKey())
	if err != nil {
		errormessage := "Error generating JWT token"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "generating token"}
		c.JSON(http.StatusInternalServerError, &errorresp)
		return
	}
	loginresp := model.SuccessResponse{Status: "Login successful"}
	log.Println(loginresp.Status)
	c.SetCookie("token", token, 86400, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &loginresp)
}

func Logout(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		errormessage := "Error getting cookie"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "cookie not found, please login again"}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	_, err = custJWT.ValidateJWT(cookie, env.GetSecretKey())
	if err != nil {
		errormessage := "Error validating JWT token"
		log.Println(errormessage+" : ", err.Error())
		errorresp := model.ErrorResponse{Error: "invalid token, please login again"}
		c.JSON(http.StatusBadRequest, &errorresp)
		return
	}
	logoutresp := model.SuccessResponse{Status: "Logout successful"}
	log.Println(logoutresp.Status)
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, &logoutresp)
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func ApiLogin(c *gin.Context) {

}
