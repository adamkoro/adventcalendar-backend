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
	rabbitMQ "github.com/adamkoro/adventcalendar-backend/lib/rabbitmq"
	rd "github.com/adamkoro/adventcalendar-backend/lib/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var (
	httpPort    int
	metricsPort int
	redisConn   *redis.Client
	rabbitConn  *amqp.Connection
)

func main() {
	httpPort = env.GetHttpPort()
	metricsPort = env.GetMetricsPort()
	// RabbitMQ connection check
	go func() {
		var isConnected bool
		var channel *amqp.Channel
		rabbitConn, err := createRabbitMqConnection()
		if err != nil {
			log.Println(err)
		}
		isConnected = true
		channel, err = rabbitMQ.CreateChannel(rabbitConn)
		if err != nil {
			log.Println(err)
		}
		err = rabbitMQ.CreateExchange(channel, "email")
		if err != nil {
			log.Println(err)
		}
		log.Println("Connected to the rabbitmq.")
		for {
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
			admin.GET("/email", nil)
			admin.POST("/email", endpoints.SendMessageToRabbitMQ)
			admin.PUT("/email", nil)
			admin.DELETE("/email", nil)
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

func createRedisConnection() *redis.Client {
	return rd.Connect(env.GetRedisHost(), env.GetRedisPort(), env.GetRedisPassword(), env.GetRedisDb())
}

func createRabbitMqConnection() (*amqp.Connection, error) {
	return rabbitMQ.Connect(env.GetRabbitmqUser(), env.GetRabbitmqPassword(), env.GetRabbitmqHost(), env.GetRabbitmqVhost(), env.GetRabbitmqPort())
}
