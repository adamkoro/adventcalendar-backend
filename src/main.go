package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	endpoints "github.com/adamkoro/adventcalendar-backend/api"
	"github.com/adamkoro/adventcalendar-backend/env"
	"github.com/gin-gonic/gin"
)

var (
	httpPort int
)

func init() {
	httpPort = env.GetHttpPort()
}

func main() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	api := router.Group("/api")
	{
		api.GET("/ping", endpoints.Ping)
	}

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(httpPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Starting server...")
	log.Println("Listening on port " + strconv.Itoa(httpPort))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
