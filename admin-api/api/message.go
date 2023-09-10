package api

import (
	mq "github.com/adamkoro/adventcalendar-backend/lib/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	MqConn *amqp.Connection
)

// Message is a struct that represents a message to be sent to the queue

func SendMessage(mqConn *amqp.Connection, msg MqMessage) error {
	ch, err := mq.CreateChannel(mqConn)
	if err != nil {
		return err
	}
	defer mq.CloseChannel(ch)

	q, err := mq.DeclareQueue(ch, msg.QueueName)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg.Body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
