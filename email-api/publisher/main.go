package main

import (
	"context"
	"log"
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
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

var (
	httpPort    int
	metricsPort int
	rabbitConn  *amqp.Connection
	db          *md.Repository
)

func main() {
	httpPort = env.GetHttpPort()
	metricsPort = env.GetMetricsPort()
	// RabbitMQ connection check
	go func() {
		//var isConnected bool
		var channel *amqp.Channel
		var queue amqp.Queue
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
		endpoints.MqChannel = channel
		endpoints.MqQUeue = queue
		/*for {
			pingConn, err := createRabbitMqConnection()
			if err != nil {
				log.Println("Lost connection to the rabbitmq, reconnecting...")
				rabbitConn, err = createRabbitMqConnection()
				if err != nil {
					isConnected = false
					log.Println("Failed to reconnect to the rabbitmq.")
				}
			} else {
				if !isConnected {
					log.Println("Reconnected to the rabbitmq.")
					isConnected = true
					channel, err := rabbitMQ.CreateChannel(rabbitConn)
					if err != nil {
						log.Println(err)
					}
					endpoints.MqChannel = channel
				}
			}
			rabbitMQ.CloseConnection(pingConn)
			endpoints.MqChannel = channel
			endpoints.MqQUeue = queue
			time.Sleep(5 * time.Second)
		}*/
	}()
	// MariaDB connection check
	go func() {
		var isConnected bool
		mariadbConn, err := createMariaDbConnection()
		if err != nil {
			log.Println(err)
		}
		ctx := context.Background()
		db := md.NewRepository(mariadbConn, &ctx)
		isConnected = true
		log.Println("Connected to the mariadb.")
		for {
			err := db.Ping()
			if err != nil {
				log.Println("Lost connection to the mariadb, reconnecting...")
				isConnected = false
			} else {
				if !isConnected {
					log.Println("Reconnected to the mariadb.")
					isConnected = true
				}
			}
			endpoints.Db = *db
			time.Sleep(5 * time.Second)
		}
	}()
	// Api server
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	if gin.Mode() == gin.DebugMode {
		router.Use(endpoints.CORSMiddleware())
	}
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

	quit := make(chan os.Signal)
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

func createMariaDbConnection() (*gorm.DB, error) {
	return db.Connect(env.GetDbUser(), env.GetDbPassword(), env.GetDbHost(), env.GetDbName(), env.GetDbPort())
}
