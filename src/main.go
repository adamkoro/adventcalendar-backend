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
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	httpPort     int
	postgresConn *gorm.DB
)

func init() {
	httpPort = env.GetHttpPort()
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
}

func main() {
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

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(httpPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server...")
	log.Println("Listening on port " + strconv.Itoa(httpPort))
	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
