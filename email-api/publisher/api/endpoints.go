package api

import (
	"net/http"

	db "github.com/adamkoro/adventcalendar-backend/lib/mariadb"
	"github.com/adamkoro/adventcalendar-backend/lib/model"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	MqChannel *amqp.Channel
	MqQUeue   amqp.Queue
	Db        db.Repository
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func EmailSend(c *gin.Context) {
	var rMail model.EmailRequest
	if err := c.ShouldBindJSON(&rMail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mail, err := Db.GetEmailByName(rMail.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = SendMessage(MqChannel, MqQUeue.Name, mail.To, mail.Subject, mail.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent to RabbitMQ."})
}

func CustomEmailSend(c *gin.Context) {
	var message model.MQMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := SendMessage(MqChannel, MqQUeue.Name, message.EmailTo, message.Subject, message.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent to RabbitMQ."})
}

func DeleteEmail(c *gin.Context) {
	var rMail model.EmailRequest
	if err := c.ShouldBindJSON(&rMail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if rMail.Name == "default" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can not delete default email."})
		return
	}
	err := Db.DeleteEmailByName(rMail.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email deleted."})
}

func GetEmails(c *gin.Context) {
	emails, err := Db.GetAllEmails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emails)
}

func CreateEmail(c *gin.Context) {
	var email db.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := Db.CreateEmail(&email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email created."})
}
