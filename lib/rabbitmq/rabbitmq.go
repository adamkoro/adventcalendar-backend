package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToMq(username, password, host, port, vhost string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s%s", username, password, host, port, vhost))
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
