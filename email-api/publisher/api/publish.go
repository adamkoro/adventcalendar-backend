package api

import (
	"encoding/json"

	"github.com/adamkoro/adventcalendar-backend/lib/model"
	rabbitMQ "github.com/adamkoro/adventcalendar-backend/lib/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendMessage(channel *amqp.Channel, emailTo, subject, message string) error {
	messageJson, err := json.Marshal(model.MQMessage{
		EmailTo: emailTo,
		Subject: subject,
		Message: message,
	})
	if err != nil {
		return err
	}
	err = rabbitMQ.Publish(channel, "email", messageJson)
	if err != nil {
		return err
	}
	return nil
}
