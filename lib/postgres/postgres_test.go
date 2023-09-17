package postgres_test

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pg "github.com/adamkoro/adventcalendar-backend/lib/postgres"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *pg.Repository
	err  error
)

func TestMain(m *testing.M) {
	// Create a new mock database
	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("can't create sqlmock: %s", err)
	}
	// Open a new GORM database connection with the mock database
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf("can't open gorm connection: %s", err)
	}
	// Create a new repository with the GORM database
	repo = pg.NewRepository(gormDB)
	if repo == nil {
		log.Fatal("repository is nil")
	}
	// Run the tests
	code := m.Run()
	// Close the mock database after the tests have run
	db.Close()
	// Exit with the code returned from the tests
	os.Exit(code)
}

func TestHashPassword(t *testing.T) {
	hash, err := pg.HashPassword("test")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	hash, err := pg.HashPassword("test")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	err = pg.CheckPasswordHash("test", hash)
	assert.NoError(t, err)
}

func TestMigrate(t *testing.T) {
	mock.ExpectExec("^CREATE TABLE (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
	repo.Migrate()
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO (.+)").WithArgs("test", "test@test.test", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"key"}).AddRow(1))
	mock.ExpectCommit()
	err = repo.CreateUser("test", "test@test.test", "test")
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUser(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "0001-01-01 00:00:00 +0000 UTC")
	mock.ExpectQuery("^SELECT (.+) FROM (.+)").WithArgs("test").WillReturnRows(sqlmock.NewRows([]string{"key", "username", "email", "password", "created_at", "modified_at"}).AddRow(1, "test", "test@test.test", "", createdAt, createdAt))
	user, err := repo.GetUser("test")
	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.Key)
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, "test@test.test", user.Email)
	assert.Equal(t, "", user.Password)
	assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", user.CreatedAt.String())
	assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", user.ModifiedAt.String())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllUser(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "0001-01-01 00:00:00 +0000 UTC")
	mock.ExpectQuery("^SELECT (.+) FROM (.+)").WillReturnRows(sqlmock.NewRows([]string{"key", "username", "email", "password", "created_at", "modified_at"}).AddRow(1, "test", "test@test.test", "", createdAt, createdAt))
	users, err := repo.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, uint(1), users[0].Key)
	assert.Equal(t, "test", users[0].Username)
	assert.Equal(t, "test@test.test", users[0].Email)
	assert.Equal(t, "", users[0].Password)
	assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", users[0].CreatedAt.String())
	assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", users[0].ModifiedAt.String())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TODO: Fix this test
/*func TestUpdateUserEmail(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "0001-01-01 00:00:00 +0000 UTC")
	mock.ExpectBegin()
	mock.ExpectQuery("^UPDATE (.+)").WithArgs("testnewemail@test.test", "test").WillReturnRows(sqlmock.NewRows([]string{"key", "username", "email", "password", "created_at", "modified_at"}).AddRow(1, "test", "test@test.test", "hashed_password", createdAt, createdAt))
	mock.ExpectCommit()

	err = repo.UpdateUser("test", "testnewemail@test.test", "")
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}*/
