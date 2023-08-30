package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	endpoints "github.com/adamkoro/adventcalendar-backend/api"
	"github.com/gin-gonic/gin"
)

func TestPing(t *testing.T) {
	router := gin.New()
	router.GET("/ping", endpoints.Ping)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "pong" {
		t.Errorf("expected %q, got %q", "pong", w.Body.String())
	}
}
