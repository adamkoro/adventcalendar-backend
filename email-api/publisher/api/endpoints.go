package api

import (
	"net/http"

	db "github.com/adamkoro/adventcalendar-backend/lib/mariadb"
	"github.com/adamkoro/adventcalendar-backend/lib/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
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

func EmailSend(c *gin.Context) {
	var rMail db.EmailRequest
	if err := c.ShouldBindJSON(&rMail); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	if validationErr := validate.Struct(&rMail); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	mail, err := Db.GetEmailByName(rMail.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid email name"})
		return
	}
	err = SendMessage(MqChannel, MqQUeue.Name, mail.To, mail.Subject, mail.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "sending message to RabbitMQ"})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "Message sent to RabbitMQ."})
}

func CustomEmailSend(c *gin.Context) {
	var message db.MQMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	if validationErr := validate.Struct(&message); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	err := SendMessage(MqChannel, MqQUeue.Name, message.EmailTo, message.Subject, message.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "sending message to RabbitMQ"})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "Message sent to RabbitMQ."})
}

func DeleteEmail(c *gin.Context) {
	var rMail db.EmailRequest
	if err := c.ShouldBindJSON(&rMail); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	if rMail.Name == "default" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "cannot delete default email"})
		return
	}
	if validationErr := validate.Struct(&rMail); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	err := Db.DeleteEmailByName(rMail.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "deleting email"})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "Email deleted."})
}

func GetEmails(c *gin.Context) {
	emails, err := Db.GetAllEmails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "getting emails"})
		return
	}
	c.JSON(http.StatusOK, emails)
}

func CreateEmail(c *gin.Context) {
	var email db.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	if validationErr := validate.Struct(&email); validationErr != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: validationErr.Error()})
		return
	}
	err := Db.CreateEmail(&email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "creating email"})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Status: "Email created."})
}
