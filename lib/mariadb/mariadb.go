package mariadb

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Connect(username, password, host, database string, port int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

func (r *Repository) CreateEmail(email *Email) error {
	return r.db.Create(email).Error
}

func (r *Repository) GetAllEmails() ([]Email, error) {
	var emails []Email
	err := r.db.Find(&emails).Error
	return emails, err
}

func (r *Repository) DeleteEmailByName(name string) error {
	return r.db.Where("name = ?", name).Delete(&Email{}).Error
}

func (r *Repository) GetEmailByName(name string) (*Email, error) {
	email := &Email{}
	err := r.db.Where("name = ?", name).First(email).Error
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
