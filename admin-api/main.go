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

	endpoints "github.com/adamkoro/adventcalendar-backend/admin-api/api"
	"github.com/adamkoro/adventcalendar-backend/lib/env"
	"github.com/adamkoro/adventcalendar-backend/lib/postgres"
	mq "github.com/adamkoro/adventcalendar-backend/lib/rabbitmq"
	rd "github.com/adamkoro/adventcalendar-backend/lib/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	httpPort     int
	metricsPort  int
	postgresConn *gorm.DB
	redisConn    *redis.Client
	rabbitmqConn *amqp.Connection
)

func main() {
	httpPort = env.GetHttpPort()
	metricsPort = env.GetMetricsPort()
	endpoints.SecretKey = env.GetSecretKey()
	// Postgres connection check
	go func() {
		var isConnected bool
		postgresConn, err := createPostgresConnection()
		if err != nil {
			log.Println(err)
		}
		isConnected = true
		log.Println("Connected to the postgres.")
		for {
			err := postgres.Ping(postgresConn)
			if err != nil {
				log.Println("Lost connection to the postgres, reconnecting...")
				postgresConn, err = createPostgresConnection()
				if err != nil {
					isConnected = false
					log.Println("Failed to reconnect to the postgres.")
				}
			} else {
				if !isConnected {
					log.Println("Reconnected to the postgres.")
					isConnected = true
				}
			}
			endpoints.Db = postgresConn
			time.Sleep(5 * time.Second)
		}
	}()

	// Redis connection check
	go func() {
		var isConnected bool
		redisConn = createRedisConnection()
		if redisConn != nil {
			isConnected = true
			log.Println("Connected to the redis.")
		}
		for {
			err := rd.Ping(redisConn)
			if err != nil {
				log.Println("Lost connection to the redis, reconnecting...")
				redisConn = createRedisConnection()
				if err != nil {
					isConnected = false
					log.Println("Failed to reconnect to the redis.")
				}
			} else {
				if !isConnected {
					log.Println("Reconnected to the redis.")
					isConnected = true
				}
			}
			endpoints.Rd = redisConn
			time.Sleep(5 * time.Second)
		}
	}()

	// RabbitMQ connection check
	go func() {
		var isConnected bool
		rabbitmqConn, err := createRabbitMqConnection()
		if err != nil {
			log.Println(err)
		}
		isConnected = true
		log.Println("Connected to the rabbitmq.")
		for {
			pingConn, err := createRabbitMqConnection()
			if err != nil {
				log.Println("Lost connection to the rabbitmq, reconnecting...")
				rabbitmqConn, err = createRabbitMqConnection()
				if err != nil {
					isConnected = false
					log.Println("Failed to reconnect to the rabbitmq.")
				}
			} else {
				if !isConnected {
					log.Println("Reconnected to the rabbitmq.")
					isConnected = true
				}
			}
			mq.CloseConnection(pingConn)
			endpoints.MqConn = rabbitmqConn
			time.Sleep(5 * time.Second)
		}
	}()

	// Api server
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(CORSMiddleware())
	api := router.Group("/api")
	{
		// Public endpoints
		// Health check
		api.GET("/ping", endpoints.Ping)
		// Login
		api.POST("/login", endpoints.Login)
		// Logout
		api.POST("/logout", endpoints.Logout)
		// Admin endpoints (require authentication)
		admin := api.Group("/admin")
		admin.Use(endpoints.AuthRequired)
		{
			// Single user
			admin.GET("/user", endpoints.GetUser)
			admin.POST("/user", endpoints.CreateUser)
			admin.PUT("/user", endpoints.UpdateUser)
			admin.DELETE("/user", endpoints.DeleteUser)
			// All users
			admin.GET("/users", endpoints.GetAllUsers)
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3030")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Access-Control-Allow-Credentials")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func createPostgresConnection() (*gorm.DB, error) {
	return postgres.Connect(env.GetDbHost(), env.GetDbUser(), env.GetDbPassword(), env.GetDbName(), env.GetDbPort(), env.GetDbSslMode())
}

func createRedisConnection() *redis.Client {
	return rd.Connect(env.GetRedisHost(), env.GetRedisPort(), env.GetRedisPassword(), env.GetRedisDb())
}

func createRabbitMqConnection() (*amqp.Connection, error) {
	return mq.Connect(env.GetRabbitmqUser(), env.GetRabbitmqPassword(), env.GetRabbitmqHost(), env.GetRabbitmqVhost(), env.GetRabbitmqPort())
}
