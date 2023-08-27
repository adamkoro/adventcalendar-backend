package env

import (
	"os"
	"strconv"
)

func GetHttpPort() int {
	port := 8080
	if os.Getenv("PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}
	return port
}
