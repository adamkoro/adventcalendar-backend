package env_test

import (
	"os"
	"testing"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
)

func TestGetLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("LOG_LEVEL")

	expected := "debug"
	actual := env.GetLogLevel()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetHttpPort(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	expected := 8080
	actual := env.GetHttpPort()

	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TestGetMetricsPort(t *testing.T) {
	os.Setenv("METRICS_PORT", "8081")
	defer os.Unsetenv("METRICS_PORT")

	expected := 8081
	actual := env.GetMetricsPort()

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

func TestGetRedisHost(t *testing.T) {
	os.Setenv("REDIS_HOST", "localhost")
	defer os.Unsetenv("REDIS_HOST")

	expected := "localhost"
	actual := env.GetRedisHost()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetRedisPort(t *testing.T) {
	os.Setenv("REDIS_PORT", "6379")
	defer os.Unsetenv("REDIS_PORT")

	expected := 6379
	actual := env.GetRedisPort()

	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TestGetRedisPassword(t *testing.T) {
	os.Setenv("REDIS_PASSWORD", "password")
	defer os.Unsetenv("REDIS_PASSWORD")

	expected := "password"
	actual := env.GetRedisPassword()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetRedisDb(t *testing.T) {
	os.Setenv("REDIS_DB", "0")
	defer os.Unsetenv("REDIS_DB")

	expected := 0
	actual := env.GetRedisDb()

	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TestGetRabbitmqHost(t *testing.T) {
	os.Setenv("RABBITMQ_HOST", "localhost")
	defer os.Unsetenv("RABBITMQ_HOST")

	expected := "localhost"
	actual := env.GetRabbitmqHost()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetRabbitmqPort(t *testing.T) {
	os.Setenv("RABBITMQ_PORT", "5672")
	defer os.Unsetenv("RABBITMQ_PORT")

	expected := 5672
	actual := env.GetRabbitmqPort()

	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TesGetRabbitmqUser(t *testing.T) {
	os.Setenv("RABBITMQ_USER", "guest")
	defer os.Unsetenv("RABBITMQ_USER")

	expected := "guest"
	actual := env.GetRabbitmqUser()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetRabbitmqPassword(t *testing.T) {
	os.Setenv("RABBITMQ_PASSWORD", "guest")
	defer os.Unsetenv("RABBITMQ_PASSWORD")

	expected := "guest"
	actual := env.GetRabbitmqPassword()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetRabbitmqVhost(t *testing.T) {
	os.Setenv("RABBITMQ_VHOST", "/")
	defer os.Unsetenv("RABBITMQ_VHOST")

	expected := "/"
	actual := env.GetRabbitmqVhost()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
