package api

import (
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	custJWT "github.com/adamkoro/adventcalendar-backend/lib/jwt"
	custModel "github.com/adamkoro/adventcalendar-backend/lib/model"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func AuthRequired(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		log.Error().Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, custModel.ErrorResponse{Error: "cookie not found, please login again"})
		return
	}
	_, err = custJWT.ValidateJWT(cookie, env.GetSecretKey())
	if err != nil {
		log.Error().Msg(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, custModel.ErrorResponse{Error: "invalid token, please login again"})
		return
	}
	c.Next()
}

func SetJsonLogger() gin.HandlerFunc {
	return logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	)
}

func CORSMiddleware() gin.HandlerFunc {
	log.Info().Msg("setting up CORS middleware")
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
