package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamkoro/adventcalendar-backend/admin-api/api"
	postgres "github.com/adamkoro/adventcalendar-backend/admin-api/postgres"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", api.Login)

	t.Run("Login with correct credentials", func(t *testing.T) {
		db, _ := postgres.Connect()
		postgres.Migrate(db)
		postgres.CreateUser(db, "testuser", "test@domain.test", "testpassword")
		defer postgres.Close(db)

		loginRequest := api.LoginRequest{
			Username: "testuser",
			Password: "testpassword",
		}
		loginRequestJson, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginRequestJson))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Login with incorrect credentials", func(t *testing.T) {

		loginRequest := api.LoginRequest{
			Username: "testuser",
			Password: "testpassword",
		}
		loginRequestJson, _ := json.Marshal(loginRequest)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginRequestJson))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 401, resp.Code)
	})
}
