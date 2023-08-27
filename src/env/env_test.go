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
