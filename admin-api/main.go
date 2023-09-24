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
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

var (
	httpPort    int
	metricsPort int
	db          pg.Repository
	//redisConn   *redis.Client
)

func main() {
	figure.NewFigure("AdventCalendar Admin Api", "big", false).Print()
	httpPort = env.GetHttpPort()
	metricsPort = env.GetMetricsPort()
	// Postgres connection check
	go func() {
		var isConnected bool
		postgresConn, err := createPostgresConnection()
		ctx := context.Background()
		db := pg.NewRepository(postgresConn, &ctx)
		if err != nil {
			log.Println(err)
		}
		isConnected = true
		log.Println("Connected to the postgres.")
		for {
			err := db.Ping()
			if err != nil {
				log.Println("Lost connection to the postgres, reconnecting...")
				isConnected = false
			} else {
				if !isConnected {
					log.Println("Reconnected to the postgres.")
					isConnected = true
				}
			}
			endpoints.Db = *db
			time.Sleep(5 * time.Second)
		}
	}()

	// Redis connection check
	/*go func() {
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
			time.Sleep(5 * time.Second)
		}
	}()*/

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
			// Single user
			admin.GET("/usermanage/user", endpoints.GetUser)
			admin.POST("/usermanage/user", endpoints.CreateUser)
			admin.PUT("/usermanage/user", endpoints.UpdateUser)
			admin.DELETE("/usermanage/user", endpoints.DeleteUser)
			// All users
			admin.GET("/usermanage/users", endpoints.GetAllUsers)
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

func createPostgresConnection() (*gorm.DB, error) {
	return db.Connect(env.GetDbHost(), env.GetDbUser(), env.GetDbPassword(), env.GetDbName(), env.GetDbPort(), env.GetDbSslMode())
}

func createRedisConnection() {

}
