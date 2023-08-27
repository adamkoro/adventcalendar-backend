package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adamkoro/adventcalendar-backend/endpoints"
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
	router.GET("/ping", endpoints.Ping)

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(httpPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
