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

	endpoints "github.com/adamkoro/adventcalendar-backend/api"
	"github.com/adamkoro/adventcalendar-backend/env"
	"github.com/adamkoro/adventcalendar-backend/postgres"
	rd "github.com/adamkoro/adventcalendar-backend/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	httpPort     int
	metricsPort  int
	postgresConn *gorm.DB
	redirConn    *redis.Client
)

func init() {
	httpPort = env.GetHttpPort()
	metricsPort = env.GetMetricsPort()
	postgresConn, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	err = postgres.Migrate(postgresConn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrated database")
	endpoints.Db = postgresConn
	redirConn = rd.Connect()
	err = rd.Ping(redirConn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to redis")
	endpoints.Rd = redirConn
}

func main() {
	// Api server
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
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
