package api

import (
	"net/http"

	"github.com/adamkoro/adventcalendar-backend/lib/model"
	mdb "github.com/adamkoro/adventcalendar-backend/lib/mongo"
	"github.com/gin-gonic/gin"
)

var (
	Db *mdb.Repository
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func GetAllDatabase(c *gin.Context) {
	dbs, err := Db.GetAllDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Internal server error."})
		return
	}
	c.JSON(http.StatusOK, dbs)
}
