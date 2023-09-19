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

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	mdb "github.com/adamkoro/adventcalendar-backend/lib/mongo"
	endpoints "github.com/adamkoro/adventcalendar-backend/public-api/api"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	httpPort    int
	metricsPort int
	db          *mdb.Repository
)

func main() {
	httpPort = env.GetHttpPort()
	metricsPort = env.GetMetricsPort()
	go func() {
		var isConnected bool
		mongoDbConn, mongoDbContext, err := createMongoDbConnection()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Connected to the mongodb.")
		}
		db = mdb.NewRepository(mongoDbConn, mongoDbContext)
		isConnected = true
		for {
			err := db.PingDb()
			if err != nil {
				log.Println("Lost connection to the mongodb, trying to reconnect...")
				isConnected = false
			} else {
				if !isConnected {
					log.Println("Reconnected to the mongodb.")
					isConnected = true
				}
			}
			endpoints.Db = db
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
		admin := api.Group("/admin")
		admin.Use(endpoints.AuthRequired)
		{
			admin.GET("/public/dbs", endpoints.GetAllDatabase)
			admin.POST("/public/day", endpoints.CreateDay)
			admin.PUT("/public/day", endpoints.UpdateDay)
			admin.DELETE("/public/day", endpoints.DeleteDay)
		}
		public := api.Group("/public")
		{
			public.GET("/day", endpoints.GetDay)
			public.GET("/days", endpoints.GetAllDay)
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

func createMongoDbConnection() (*mongo.Client, *context.Context, error) {
	return db.Connect(env.GetDbUser(), env.GetDbPassword(), env.GetDbHost(), env.GetDbPort())
}
