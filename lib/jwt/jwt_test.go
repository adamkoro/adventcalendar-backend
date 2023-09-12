package jwt_test

import (
	"testing"

	"github.com/adamkoro/adventcalendar-backend/lib/jwt"
)

var (
	secretKey string
	token     string
	err       error
)

func init() {
	secretKey = "test"
}

func TestGenerateJWT(t *testing.T) {
	token, err = jwt.GenerateJWT("test", secretKey)
	if err != nil {
		t.Errorf("Failed to generate JWT: %v", err)
	}
	if token == "" {
		t.Errorf("Failed to generate JWT: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	token, _ = jwt.GenerateJWT("test", secretKey)
	claims, err := jwt.ValidateJWT(token, secretKey)
	if err != nil {
		t.Errorf("Failed to validate JWT: %v", err)
	}
	if claims["user"] != "test" {
		t.Errorf("Failed to validate JWT: %v", err)
	}
	if claims["session"] != "test" {
		t.Errorf("Failed to validate JWT: %v", err)
	}
}
