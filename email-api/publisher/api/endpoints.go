package api

import (
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/model"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	MqChannel *amqp.Channel
	MqQUeue   amqp.Queue
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func SendMessageToRabbitMQ(c *gin.Context) {
	var message model.MQMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := SendMessage(MqChannel, MqQUeue.Name, message.EmailTo, message.Subject, message.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent to RabbitMQ."})
}
