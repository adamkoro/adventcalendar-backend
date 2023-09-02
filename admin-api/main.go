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
	"github.com/adamkoro/adventcalendar-backend/admin-api/env"
	"github.com/adamkoro/adventcalendar-backend/admin-api/postgres"
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

func main() {
	httpPort = env.GetHttpPort()
	metricsPort = env.GetMetricsPort()
	postgresChannel := make(chan *gorm.DB, 1)
	postgresConn, _ = postgres.Connect()
	dbctx, dbcancel := context.WithCancel(context.Background())
	defer dbcancel()
	go monitorDbConnection(dbctx, postgresConn, postgresChannel)

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

func monitorDbConnection(ctx context.Context, db *gorm.DB, dbchannel chan *gorm.DB) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := postgres.Ping(postgresConn); err != nil {
				log.Printf("DB connection lost. Retrying...")
				db, err = postgres.Connect()
				if err != nil {
					log.Println(err)
				}
				dbchannel <- db
			}
		}
	}
}

/*func dada() {
	// Redis
	redirConn = rd.Connect()
	err = rd.Ping(redirConn)
	if err != nil {
		log.Println(err)
	}
	log.Println("Connected to redis")
	endpoints.Rd = redirConn
	// Admin user create and update
	err = postgres.CreateUser(postgresConn, env.GetAdminUsername(), env.GetAdminEmail(), env.GetAdminPassword())
	if err != nil {
		log.Println(err)
	}
	isAdminExists, err := postgres.GetUser(postgresConn, env.GetAdminUsername())
	if err != nil {
		log.Println(err)
	}
	if isAdminExists.Username != "" {
		err = postgres.UpdateUser(postgresConn, env.GetAdminUsername(), env.GetAdminEmail(), env.GetAdminPassword())
		if err != nil {
			log.Println(err)
		}
	}
}*/
