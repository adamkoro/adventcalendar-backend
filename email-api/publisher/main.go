package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	endpoints "github.com/adamkoro/adventcalendar-backend/email-api/publisher/api"
	"github.com/adamkoro/adventcalendar-backend/lib/env"
	md "github.com/adamkoro/adventcalendar-backend/lib/mariadb"
	rabbitMQ "github.com/adamkoro/adventcalendar-backend/lib/rabbitmq"
	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	httpPort    int
	metricsPort int
	rabbitConn  *amqp.Connection
	db          *md.Repository
)

func main() {
	figure.NewFigure("AdventCalendar Email Publisher", "big", false).Print()
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
			endpoints.MqChannel, err = rabbitMQ.CreateChannel(rabbitConn)
			if err != nil {
				log.Error().Msg(err.Error())
			}
			log.Debug().Msg("channel created")
			endpoints.MqQUeue, err = rabbitMQ.DeclareQueue(endpoints.MqChannel, "email")
			if err != nil {
				log.Error().Msg(err.Error())
			}
			log.Debug().Msg("queue declared")
		}
		for {
			if rabbitConn == nil || rabbitConn.IsClosed() {
				log.Error().Msg("lost connection to the rabbitmq, reconnecting...")
				rabbitConn, err = createRabbitMqConnection()
				if err != nil {
					log.Error().Msg(err.Error())
				} else {
					log.Info().Msg("connected to the rabbitmq")
					endpoints.MqChannel, err = rabbitMQ.CreateChannel(rabbitConn)
					if err != nil {
						log.Error().Msg(err.Error())
					}
					log.Debug().Msg("channel created")
					endpoints.MqQUeue, err = rabbitMQ.DeclareQueue(endpoints.MqChannel, "email")
					if err != nil {
						log.Error().Msg(err.Error())
					}
					log.Debug().Msg("queue declared")
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
	// MariaDB connection check
	/////////////////////////
	go func() {
		var isConnected bool
		log.Debug().Msg("establishing connection to the mariadb...")
		mariadbConn, err := createMariaDbConnection()
		if err != nil {
			log.Error().Msg(err.Error())
		} else {
			log.Info().Msg("connected to the mariadb")
		}
		ctx := context.Background()
		db := md.NewRepository(mariadbConn, &ctx)
		isConnected = true
		for {
			log.Debug().Msg("pinging the mariadb...")
			err = db.Ping()
			if err != nil {
				log.Info().Msg("lost connection to the mariadb, reconnecting...")
				log.Error().Msg(err.Error())
				isConnected = false
			} else {
				log.Debug().Msg("pinging the mariadb successful")
				if !isConnected {
					log.Info().Msg("reconnected to the mariadb")
					isConnected = true
				}
			}
			endpoints.Db = *db
			time.Sleep(5 * time.Second)
		}
	}()
	/////////////////////////
	// Api server setup
	/////////////////////////
	router := gin.New()
	router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		router.Use(endpoints.CORSMiddleware())
	}
	router.Use(endpoints.SetJsonLogger())
	log.Debug().Msg("setting up endpoints...")
	api := router.Group("/api")
	{
		// Public endpoints
		// Health check
		api.GET("/ping", endpoints.Ping)
		// Admin endpoints (require authentication)
		admin := api.Group("/admin")
		admin.Use(endpoints.AuthRequired)
		{
			admin.GET("/emailmanage/email", endpoints.GetEmails)
			admin.POST("/emailmanage/email", endpoints.CreateEmail)
			admin.PUT("/emailmanage/email", nil)
			admin.DELETE("/emailmanage/email", endpoints.DeleteEmail)
			admin.POST("/emailmanage/sendemail", endpoints.EmailSend)
			admin.POST("/emailmanage/customemail", endpoints.CustomEmailSend)
		}
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
	metrics.Use(endpoints.SetJsonLogger())
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

func createMariaDbConnection() (*gorm.DB, error) {
	return db.Connect(env.GetDbUser(), env.GetDbPassword(), env.GetDbHost(), env.GetDbName(), env.GetDbPort())
}
