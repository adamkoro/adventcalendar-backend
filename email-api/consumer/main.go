package main

import (
	"log"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	rabbitMQ "github.com/adamkoro/adventcalendar-backend/lib/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitConn *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	forever    chan struct{}
)

func main() {
	rabbitConn, err := createRabbitMqConnection()
	if err != nil {
		log.Println(err)
	}
	//isConnected = true
	channel, err = rabbitMQ.CreateChannel(rabbitConn)
	if err != nil {
		log.Println(err)
	}
	log.Println("Connected to the rabbitmq.")
	log.Println("Channel created.")
	queue, err = rabbitMQ.DeclareQueue(channel, "email")
	if err != nil {
		log.Println(err)
	}
	log.Println("Queue declared.")
	consume, err := rabbitMQ.Consume(channel, queue.Name)
	if err != nil {
		log.Println(err)
	}
	go func() {
		for d := range consume {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	<-forever
}

func createRabbitMqConnection() (*amqp.Connection, error) {
	return rabbitMQ.Connect(env.GetRabbitmqUser(), env.GetRabbitmqPassword(), env.GetRabbitmqHost(), env.GetRabbitmqVhost(), env.GetRabbitmqPort())
}
