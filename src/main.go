package main

import (
	"strconv"

	"github.com/adamkoro/adventcalendar-backend/env"
	"github.com/gin-gonic/gin"
)

var httpPort int

func init() {
	httpPort = env.GetHttpPort()
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})
	router.Run(":" + strconv.Itoa(httpPort))
}
