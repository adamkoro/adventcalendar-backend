package api

import (
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/model"
	mdb "github.com/adamkoro/adventcalendar-backend/lib/mongo"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var (
	Db       *mdb.Repository
	validate = validator.New()
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func GetAllDatabase(c *gin.Context) {
	log.Debug().Msg("getting all databases...")
	dbs, err := Db.GetAllDatabase()
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	log.Debug().Msg("databases retrieved successfully")
	c.JSON(http.StatusOK, dbs)
}

func CreateDay(c *gin.Context) {
	var day mdb.AdventCalendarDay
	log.Debug().Msg("binding request body...")
	if err := c.BindJSON(&day); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request payload"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&day); validationErr != nil {
		log.Error().Msg(validationErr.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	log.Debug().Msg("validation successful")
	log.Debug().Msg("creating day...")
	if err := Db.CreateDay(&day, "adventcalendar", "days"); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	log.Debug().Msg("day created successfully")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "day created successfully"})
}

func GetDay(c *gin.Context) {
	var day mdb.AdventCalendarDay
	log.Debug().Msg("binding request body...")
	if err := c.BindJSON(&day); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request payload"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&day); validationErr != nil {
		log.Error().Msg(validationErr.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	log.Debug().Msg("validation successful")
	log.Debug().Msg("getting day...")
	respDay, err := Db.GetDay(day.Day, "adventcalendar", "days")
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	log.Debug().Msg("day retrieved successfully")
	c.JSON(http.StatusOK, respDay)
}

func GetAllDay(c *gin.Context) {
	log.Debug().Msg("getting all days...")
	days, err := Db.GetAllDay("adventcalendar", "days")
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	log.Debug().Msg("all days retrieved successfully")
	c.JSON(http.StatusOK, days)
}

func UpdateDay(c *gin.Context) {
	var day mdb.AdventCalendarDayUpdate
	log.Debug().Msg("binding request body...")
	if err := c.BindJSON(&day); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request payload"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&day); validationErr != nil {
		log.Error().Msg(validationErr.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	log.Debug().Msg("validation successful")
	log.Debug().Msg("updating day...")
	if err := Db.UpdateDay(&day, "adventcalendar", "days"); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	log.Debug().Msg("day updated successfully")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "day updated successfully"})
}

func DeleteDay(c *gin.Context) {
	var dayID mdb.DayIDRequest
	log.Debug().Msg("binding request body...")
	if err := c.BindJSON(&dayID); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request payload"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&dayID); validationErr != nil {
		log.Error().Msg(validationErr.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	log.Debug().Msg("validation successful")
	log.Debug().Msg("deleting day...")
	if err := Db.DeleteDay(&dayID, "adventcalendar", "days"); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	log.Debug().Msg("day deleted successfully")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "Day deleted successfully."})
}
