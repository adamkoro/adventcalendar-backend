package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(username, password, host, vhost string, port int) (*amqp.Connection, error) {
	stringport := fmt.Sprintf("%d", port)
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s%s", username, password, host, stringport, vhost))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func CloseChannel(ch *amqp.Channel) {
	ch.Close()
}

func CloseConnection(conn *amqp.Connection) {
	defer conn.Close()
}

func DeclareQueue(ch *amqp.Channel, chName string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		chName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return q, err
	}
	return q, nil
}

func CreateExchange(ch *amqp.Channel, exName string) error {
	err := ch.ExchangeDeclare(
		exName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func Publish(ch *amqp.Channel, exchangeName string, body []byte) error {
	err := ch.PublishWithContext(context.Background(), "exchangeName", "direct", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}
	return nil
}
