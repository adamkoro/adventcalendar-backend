package api

import (
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/model"
	mdb "github.com/adamkoro/adventcalendar-backend/lib/mongo"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	Db       *mdb.Repository
	validate = validator.New()
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func GetAllDatabase(c *gin.Context) {
	dbs, err := Db.GetAllDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal server error."})
		return
	}
	c.JSON(http.StatusOK, dbs)
}

func CreateDay(c *gin.Context) {
	var day mdb.AdventCalendarDay
	if err := c.BindJSON(&day); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request payload."})
		return
	}
	if validationErr := validate.Struct(&day); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	if err := Db.CreateDay(&day, "adventcalendar", "days"); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal server error."})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "Day created successfully."})
}

func GetDay(c *gin.Context) {
	var day mdb.AdventCalendarDay
	if err := c.BindJSON(&day); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request payload."})
		return
	}
	if validationErr := validate.Struct(&day); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	respDay, err := Db.GetDay(day.Day, "adventcalendar", "days")
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal server error."})
		return
	}
	c.JSON(http.StatusOK, respDay)
}

func GetAllDay(c *gin.Context) {
	days, err := Db.GetAllDay("adventcalendar", "days")
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal server error."})
		return
	}
	c.JSON(http.StatusOK, days)
}

func UpdateDay(c *gin.Context) {
	var day mdb.AdventCalendarDayUpdate
	if err := c.BindJSON(&day); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request payload."})
		return
	}
	if validationErr := validate.Struct(&day); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	if err := Db.UpdateDay(&day, "adventcalendar", "days"); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal server error."})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "Day updated successfully."})
}

func DeleteDay(c *gin.Context) {
	var dayID mdb.DayIDRequest
	if err := c.BindJSON(&dayID); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request payload."})
		return
	}
	if validationErr := validate.Struct(&dayID); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	if err := Db.DeleteDay(&dayID, "adventcalendar", "days"); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal server error."})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "Day deleted successfully."})
}
