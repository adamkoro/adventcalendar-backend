package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	endpoints "github.com/adamkoro/adventcalendar-backend/auth-api/api"
	"github.com/adamkoro/adventcalendar-backend/lib/env"
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	httpPort    int
	metricsPort int
	db          pg.Repository
)

func main() {
	figure.NewFigure("AdventCalendar Auth Api", "big", false).Print()
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
	// Postgres connection check
	/////////////////////////
	go func() {
		var isConnected bool
		log.Debug().Msg("establishing connection to the postgres...")
		postgresConn, err := createPostgresConnection()
		if err != nil {
			log.Error().Msg(err.Error())
		} else {
			log.Info().Msg("connected to the postgres")
		}
		ctx := context.Background()
		db := pg.NewRepository(postgresConn, &ctx)
		isConnected = true
		for {
			log.Debug().Msg("pinging the postgres...")
			err := db.Ping()
			if err != nil {
				log.Info().Msg("lost connection to the postgres, reconnecting...")
				log.Error().Msg(err.Error())
				isConnected = false
			} else {
				log.Debug().Msg("pinging the postgres successful")
				if !isConnected {
					log.Info().Msg("reconnected to the postgres")
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
	// Swagger
	log.Debug().Msg("setting up endpoints...")
	router.StaticFile("/swagger/doc.json", "./docs/swagger.json")
	api := router.Group("/api")
	{
		// Public endpoints
		// Health check
		api.GET("/ping", endpoints.Ping)
		// Login
		api.POST("/auth/login", endpoints.Login)
		// Basic auth
		api.POST("/auth/api_login", endpoints.ApiLogin)
		// Logout
		api.POST("/auth/logout", endpoints.Logout)
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

func createPostgresConnection() (*gorm.DB, error) {
	return db.Connect(env.GetDbHost(), env.GetDbUser(), env.GetDbPassword(), env.GetDbName(), env.GetDbPort(), env.GetDbSslMode())
}
