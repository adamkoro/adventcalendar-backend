package db

import (
	"strconv"

	"github.com/adamkoro/adventcalendar-backend/env"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
	sslmode    string
)

func init() {
	dbHost = env.GetDbHost()
	dbPort = env.GetDbPort()
	dbUser = env.GetDbUser()
	dbPassword = env.GetDbPassword()
	dbName = env.GetDbName()
	sslmode = env.GetDbSslMode()
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Connect() (*gorm.DB, error) {
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=" + sslmode
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func CreateUser(db *gorm.DB, username string, email string, password string) error {
	hashpass, err := hashPassword(password)
	if err != nil {
		return err
	}
	return db.Create(&User{Username: username, Email: email, Password: hashpass}).Error
}

func GetUser(db *gorm.DB, username string) (*User, error) {
	user := &User{}
	err := db.Where("username = ?", username).First(user).Error
	return user, err
}

func DeleteUser(db *gorm.DB, username string) error {
	return db.Where("username = ?", username).Delete(&User{}).Error
}

func UpdatePassword(db *gorm.DB, username string, password string) error {
	return db.Model(&User{}).Where("username = ?", username).Update("password", password).Error
}

func UpdateUsername(db *gorm.DB, username string, newUsername string) error {
	return db.Model(&User{}).Where("username = ?", username).Update("username", newUsername).Error
}
