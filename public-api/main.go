package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	mdb "github.com/adamkoro/adventcalendar-backend/lib/mongo"
	endpoints "github.com/adamkoro/adventcalendar-backend/public-api/api"
	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	httpPort    int
	metricsPort int
	db          *mdb.Repository
)

func main() {
	figure.NewFigure("AdventCalendar Public Api", "big", false).Print()
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
	// MongoDb connection check
	/////////////////////////
	go func() {
		var isConnected bool
		log.Debug().Msg("establishing connection to the mongodb...")
		mongoDbConn, mongoDbContext, err := createMongoDbConnection()
		if err != nil {
			log.Error().Msg(err.Error())
		} else {
			log.Info().Msg("connected to the mongodb")
		}
		db = mdb.NewRepository(mongoDbConn, mongoDbContext)
		isConnected = true
		for {
			log.Debug().Msg("pinging the mongodb...")
			err := db.PingDb()
			if err != nil {
				log.Error().Msg(err.Error())
				isConnected = false
			} else {
				log.Debug().Msg("pinging the mongodb successful")
				if !isConnected {
					log.Info().Msg("reconnected to the mongodb")
					isConnected = true
				}
			}
			endpoints.Db = db
			time.Sleep(5 * time.Second)
		}
	}()
	/////////////////////////
	// Api server
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
		admin := api.Group("/admin")
		admin.Use(endpoints.AuthRequired)
		{
			admin.GET("/public/dbs", endpoints.GetAllDatabase)
			admin.POST("/public/day", endpoints.CreateDay)
			admin.PATCH("/public/day", endpoints.UpdateDay)
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
	log.Debug().Msg("setting up endpoints successful")
	/////////////////////////
	// Metrics server
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

func createMongoDbConnection() (*mongo.Client, *context.Context, error) {
	return db.Connect(env.GetDbUser(), env.GetDbPassword(), env.GetDbHost(), env.GetDbPort())
}
