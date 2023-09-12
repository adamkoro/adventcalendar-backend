package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	endpoints "github.com/adamkoro/adventcalendar-backend/admin-api/api"
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

func TestLogin(t *testing.T) {
	router := gin.New()
	router.POST("/login", endpoints.Login)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLogout(t *testing.T) {
	router := gin.New()
	router.POST("/logout", endpoints.Logout)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetAllUsers(t *testing.T) {
	router := gin.New()
	router.GET("/users", endpoints.GetAllUsers)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetUser(t *testing.T) {
	router := gin.New()
	router.GET("/user", endpoints.GetUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	router := gin.New()
	router.PUT("/user", endpoints.UpdateUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/user", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateUser(t *testing.T) {
	router := gin.New()
	router.POST("/user", endpoints.CreateUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	router := gin.New()
	router.DELETE("/user", endpoints.DeleteUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/user", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
