package mariadb

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Connect(username, password, host, database string, port int) *sql.DB {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Println(err)
	}
	return db
}

func (r *Repository) Migrate() error {
	return r.db.AutoMigrate(&Email{})
}

func (r *Repository) Close() error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (r *Repository) Create(email *Email) error {
	return r.db.Create(email).Error
}

func (r *Repository) FindAll() ([]Email, error) {
	var emails []Email
	err := r.db.Find(&emails).Error
	return emails, err
}

func (r *Repository) DeleteEmail(key uint) error {
	return r.db.Delete(&Email{}, key).Error
}

func (r *Repository) FindEmail(key uint) (Email, error) {
	var email Email
	err := r.db.First(&email, key).Error
	return email, err
}

func (r *Repository) UpdateEmail(key uint, email *Email) error {
	return r.db.Model(&Email{}).Where("key = ?", key).Updates(email).Error
}

func (r *Repository) Ping() error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}
