package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	rabbitMQ "github.com/adamkoro/adventcalendar-backend/lib/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitConn  *amqp.Connection
	channel     *amqp.Channel
	queue       amqp.Queue
	forever     chan struct{}
	httpPort    int
	metricsPort int
)

func main() {
	httpPort = env.GetHttpPort()
	metricsPort = env.GetMetricsPort()
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
			var message rabbitMQ.MQMessage
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
	// Api server
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	api := router.Group("/api")
	{
		// Public endpoints
		// Health check
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
		api.GET("/health", func(c *gin.Context) {
			if env.GetSmtpAuth() {
			}
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})
	}
	api_server := &http.Server{
		Addr:         ":" + strconv.Itoa(httpPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Metrics server
	metrics := gin.New()
	metrics.GET("/metrics", gin.WrapH(promhttp.Handler()))
	metrics_server := &http.Server{
		Addr:         ":" + strconv.Itoa(metricsPort),
		Handler:      metrics,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server...")
	log.Println("Listening api on port " + strconv.Itoa(httpPort))
	log.Println("Listening metrics on port " + strconv.Itoa(metricsPort))
	go func() {
		// api
		if err := api_server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	go func() {
		// metrics
		if err := metrics_server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := api_server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func createRabbitMqConnection() (*amqp.Connection, error) {
	return rabbitMQ.Connect(env.GetRabbitmqUser(), env.GetRabbitmqPassword(), env.GetRabbitmqHost(), env.GetRabbitmqVhost(), env.GetRabbitmqPort())
}

func sendMail(smtpAuth bool, smtpHost, smtpPort, smtpUser, smtpPassword, from, to, subject, body string) error {
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"\r\n" +
		body)
	if smtpAuth {
		auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
		return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	}
	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, msg)
}
