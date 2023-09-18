package main

import (
	"encoding/json"
	"log"
	"net/smtp"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	"github.com/adamkoro/adventcalendar-backend/lib/model"
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
			var message model.MQMessage
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Println(err)
			}
			log.Println("Message received")
			err = sendMail(env.GetSmtpAuth(), env.GetSmtpHost(), env.GetSmtpPort(), env.GetSmtpUser(), env.GetSmtpPassword(), env.GetSmtpFrom(), message.EmailTo, message.Subject, message.Message)
			if err != nil {
				log.Println(err)
			}
			log.Println("Email sent")
		}
	}()
	<-forever
}

func createRabbitMqConnection() (*amqp.Connection, error) {
	return rabbitMQ.Connect(env.GetRabbitmqUser(), env.GetRabbitmqPassword(), env.GetRabbitmqHost(), env.GetRabbitmqVhost(), env.GetRabbitmqPort())
}

func sendMail(smtpAuth bool, smtpHost, smtpPort, smtpUser, smtpPassword, from, to, subject, body string) error {
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)
	if smtpAuth {
		auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
		return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	}
	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, msg)
}
