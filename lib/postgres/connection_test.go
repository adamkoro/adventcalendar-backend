package postgres_test

import (
	"testing"

	"github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	db, err := postgres.Connect()
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func TestMigrate(t *testing.T) {
	db, _ := postgres.Connect()
	err := postgres.Migrate(db)
	assert.Nil(t, err)
}

func TestClose(t *testing.T) {
	db, _ := postgres.Connect()
	err := postgres.Close(db)
	assert.Nil(t, err)
}

func TestCreateUser(t *testing.T) {

	db, _ := postgres.Connect()
	postgres.Migrate(db)
	defer postgres.Close(db)

	err := postgres.CreateUser(db, "testuser", "test@domain.test", "testpassword")
	assert.Nil(t, err)
}
