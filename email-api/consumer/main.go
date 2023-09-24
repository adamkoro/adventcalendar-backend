package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	rabbitMQ "github.com/adamkoro/adventcalendar-backend/lib/rabbitmq"
	"github.com/common-nighthawk/go-figure"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	rabbitConn  *amqp.Connection
	channel     *amqp.Channel
	queue       amqp.Queue
	consume     <-chan amqp.Delivery
	httpPort    int
	metricsPort int
)

func main() {
	figure.NewFigure("AdventCalendar Email Consumer", "big", false).Print()
	/////////////////////////
	// Environment variables
	/////////////////////////
	httpPort = env.GetHttpPort()
	metricsPort = env.GetMetricsPort()
	/////////////////////////
	// Logger setup
	/////////////////////////
	logLevel := env.GetLogLevel()
	switch logLevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	/////////////////////////
	// RabbitMQ connection check
	/////////////////////////
	go func() {
		var err error
		var wait time.Duration = 5 * time.Second
		rabbitConn, err = createRabbitMqConnection()
		if err != nil {
			log.Error().Msg(err.Error())
		} else {
			log.Info().Msg("connected to the rabbitmq")
			channel, err = rabbitMQ.CreateChannel(rabbitConn)
			if err != nil {
				log.Error().Msg(err.Error())
			}
			log.Debug().Msg("channel created")
			queue, err = rabbitMQ.DeclareQueue(channel, "email")
			if err != nil {
				log.Error().Msg(err.Error())
			}
			log.Debug().Msg("queue declared")
			consume, err = rabbitMQ.Consume(channel, queue.Name)
			if err != nil {
				log.Error().Msg(err.Error())
			}
		}
		for {
			if rabbitConn == nil || rabbitConn.IsClosed() {
				log.Error().Msg("connection closed, reconnecting...")
				rabbitConn, err = createRabbitMqConnection()
				if err != nil {
					log.Error().Msg(err.Error())
				} else {
					log.Info().Msg("connected to the rabbitmq")
					channel, err = rabbitMQ.CreateChannel(rabbitConn)
					if err != nil {
						log.Error().Msg(err.Error())
					}
					log.Debug().Msg("channel created")
					queue, err = rabbitMQ.DeclareQueue(channel, "email")
					if err != nil {
						log.Error().Msg(err.Error())
					}
					log.Debug().Msg("queue declared")
					consume, err = rabbitMQ.Consume(channel, queue.Name)
					if err != nil {
						log.Error().Msg(err.Error())
					}
				}
			}
			if rabbitConn != nil {
				notify := rabbitConn.NotifyClose(make(chan *amqp.Error))
				select {
				case err = <-notify:
					log.Error().Msg(err.Error())
					rabbitConn = nil
				case <-time.After(5 * time.Second):
					// Check the connection every 5 seconds
				}
			}
			if rabbitConn == nil {
				// Wait before attempting to reconnect
				time.Sleep(wait)
			}
		}
	}()
	/////////////////////////
	// RabbitMQ consume messages
	/////////////////////////
	go func() {
		for d := range consume {
			var message rabbitMQ.MQMessage
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Error().Msg(err.Error())
			}
			log.Info().Msg("email received")
			err = sendMail(env.GetSmtpAuth(), env.GetSmtpHost(), env.GetSmtpPort(), env.GetSmtpUser(), env.GetSmtpPassword(), env.GetSmtpFrom(), message.EmailTo, message.Subject, message.Message)
			if err != nil {
				log.Error().Msg(err.Error())
			}
			log.Info().Msg("email sent")
		}
	}()
	/////////////////////////
	// Api server setup
	/////////////////////////
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		})))
	log.Debug().Msg("setting up endpoints...")
	api := router.Group("/api")
	{
		// Public endpoints
		// Health check
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})
	}
	api_server := &http.Server{
		Addr:         ":" + strconv.Itoa(httpPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Debug().Msg("setting up endpoints successful")
	/////////////////////////
	// Metrics server setup
	/////////////////////////
	metrics := gin.New()
	metrics.Use(gin.Recovery())
	metrics.Use(logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		})))
	log.Debug().Msg("setting up metrics endpoints...")
	metrics.GET("/metrics", gin.WrapH(promhttp.Handler()))
	metrics_server := &http.Server{
		Addr:         ":" + strconv.Itoa(metricsPort),
		Handler:      metrics,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Debug().Msg("setting up metrics endpoints successful")
	/////////////////////////
	// Server info
	/////////////////////////
	log.Info().Msg("starting server...")
	log.Info().Msg("listening api on port " + strconv.Itoa(httpPort))
	log.Info().Msg("listening metrics on port " + strconv.Itoa(metricsPort))
	/////////////////////////
	// Api server start
	/////////////////////////
	go func() {
		if err := api_server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic().Msg("api listen: " + err.Error())
		}
	}()
	/////////////////////////
	// Metrics server start
	/////////////////////////
	go func() {
		if err := metrics_server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic().Msg("metrics listen: " + err.Error())
		}
	}()
	/////////////////////////
	// Graceful shutdown
	/////////////////////////
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := api_server.Shutdown(ctx); err != nil {
		log.Error().Msg("server shutdown: " + err.Error())
	}
	select {
	case <-ctx.Done():
		log.Info().Msg("timeout of 5 seconds")
	}
	log.Info().Msg("server exiting")
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
