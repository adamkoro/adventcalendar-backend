package api

import (
	"net/http"

	db "github.com/adamkoro/adventcalendar-backend/lib/mariadb"
	"github.com/adamkoro/adventcalendar-backend/lib/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

var (
	MqChannel *amqp.Channel
	MqQUeue   amqp.Queue
	Db        db.Repository
	validate  = validator.New()
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// ///////////////////////
// Email send
// ///////////////////////
func EmailSend(c *gin.Context) {
	var rMail db.EmailRequest
	log.Debug().Msg("binding request body...")
	if err := c.ShouldBindJSON(&rMail); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&rMail); validationErr != nil {
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
	log.Debug().Msg("getting email from the database...")
	mail, err := Db.GetEmailByName(&rMail)
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid email name"})
		return
	}
	log.Debug().Msg("getting email from the database successful")
	log.Debug().Msg("sending message to RabbitMQ...")
	err = SendMessage(MqChannel, MqQUeue.Name, mail.To, mail.Subject, mail.Body)
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "sending message to RabbitMQ"})
		return
	}
	log.Debug().Msg("sending message to RabbitMQ successful")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "message sent to RabbitMQ."})
}

// ///////////////////////
// Custom email send
// ///////////////////////
func CustomEmailSend(c *gin.Context) {
	var message db.MQMessage
	log.Debug().Msg("binding request body...")
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&message); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	log.Debug().Msg("validating request body successful")
	log.Debug().Msg("sending message to RabbitMQ...")
	err := SendMessage(MqChannel, MqQUeue.Name, message.EmailTo, message.Subject, message.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "sending message to RabbitMQ"})
		return
	}
	log.Debug().Msg("sending message to RabbitMQ successful")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "message sent to RabbitMQ."})
}

// ///////////////////////
// Get all emails
// ///////////////////////
func GetEmails(c *gin.Context) {
	log.Debug().Msg("establishing connection to the database..")
	err := Db.Ping()
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "can not establish connection to the database"})
		return
	}
	log.Debug().Msg("establishing connection to the database successful")
	log.Debug().Msg("getting emails from the database...")
	emails, err := Db.GetAllEmails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "getting emails"})
		return
	}
	log.Debug().Msg("getting emails from the database successful")
	c.JSON(http.StatusOK, emails)
}

// ///////////////////////
// Create email
// ///////////////////////
func CreateEmail(c *gin.Context) {
	var email db.Email
	log.Debug().Msg("binding request body...")
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&email); validationErr != nil {
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
	log.Debug().Msg("creating email in the database...")
	err = Db.CreateEmail(&email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "creating email"})
		return
	}
	log.Debug().Msg("creating email in the database successful")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "email created"})
}

// ///////////////////////
// Update email
// ///////////////////////
func UpdateEmail(c *gin.Context) {
	var email db.UpdateEmailRequest
	log.Debug().Msg("binding request body...")
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&email); validationErr != nil {
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
	log.Debug().Msg("updating email in the database...")
	err = Db.UpdateEmail(&email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "updating email"})
		return
	}
	log.Debug().Msg("updating email in the database successful")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "email updated"})
}

// ///////////////////////
// Delete email
// ///////////////////////
func DeleteEmail(c *gin.Context) {
	var rMail db.DeleteEmailRequest
	log.Debug().Msg("binding request body...")
	if err := c.ShouldBindJSON(&rMail); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	log.Debug().Msg("binding request body successful")
	log.Debug().Msg("checking if email is default...")
	if rMail.Name == "default" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "cannot delete default email"})
		return
	}
	log.Debug().Msg("checking if email is default successful")
	log.Debug().Msg("validating request body...")
	if validationErr := validate.Struct(&rMail); validationErr != nil {
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
	log.Debug().Msg("deleting email from the database...")
	err = Db.DeleteEmailByName(&rMail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "deleting email"})
		return
	}
	log.Debug().Msg("deleting email from the database successful")
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "email deleted"})
}
