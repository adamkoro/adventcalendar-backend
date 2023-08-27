package env_test

import (
	"os"
	"testing"

	"github.com/adamkoro/adventcalendar-backend/env"
)

func TestGetHttpPort(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	expected := 8080
	actual := env.GetHttpPort()

	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TestGetDbHost(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	defer os.Unsetenv("DB_HOST")

	expected := "localhost"
	actual := env.GetDbHost()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetDbPort(t *testing.T) {
	os.Setenv("DB_PORT", "5432")
	defer os.Unsetenv("DB_PORT")

	expected := 5432
	actual := env.GetDbPort()

	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TestGetDbUser(t *testing.T) {
	os.Setenv("DB_USER", "postgres")
	defer os.Unsetenv("DB_USER")

	expected := "postgres"
	actual := env.GetDbUser()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetDbPassword(t *testing.T) {
	os.Setenv("DB_PASSWORD", "postgres")
	defer os.Unsetenv("DB_PASSWORD")

	expected := "postgres"
	actual := env.GetDbPassword()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetDbName(t *testing.T) {
	os.Setenv("DB_NAME", "postgres")
	defer os.Unsetenv("DB_NAME")

	expected := "postgres"
	actual := env.GetDbName()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetDbSslMode(t *testing.T) {
	os.Setenv("DB_SSL_MODE", "disable")
	defer os.Unsetenv("DB_SSL_MODE")

	expected := "disable"
	actual := env.GetDbSslMode()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetDbMaxIdleConns(t *testing.T) {
	os.Setenv("DB_MAX_IDLE_CONNS", "10")
	defer os.Unsetenv("DB_MAX_IDLE_CONNS")

	expected := 10
	actual := env.GetDbMaxIdleConns()

	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TestGetDbMaxOpenConns(t *testing.T) {
	os.Setenv("DB_MAX_OPEN_CONNS", "10")
	defer os.Unsetenv("DB_MAX_OPEN_CONNS")

	expected := 10
	actual := env.GetDbMaxOpenConns()

	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}
