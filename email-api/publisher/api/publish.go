package api

import (
	"encoding/json"

	db "github.com/adamkoro/adventcalendar-backend/lib/mariadb"
	rabbitMQ "github.com/adamkoro/adventcalendar-backend/lib/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendMessage(channel *amqp.Channel, queueName, emailTo, subject, message string) error {
	messageJson, err := json.Marshal(db.MQMessage{
		EmailTo: emailTo,
		Subject: subject,
		Message: message,
	})
	if err != nil {
		return err
	}
	err = rabbitMQ.Publish(channel, queueName, messageJson)
	if err != nil {
		return err
	}
	return nil
}
