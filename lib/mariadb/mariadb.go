package mariadb

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, ctx *context.Context) *Repository {
	return &Repository{
		Db:  db,
		Ctx: ctx,
	}
}

func (r *Repository) Connect(username, password, host, database string, port int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func (r *Repository) Migrate() error {
	return r.Db.WithContext(*r.Ctx).AutoMigrate(&Email{})
}

func (r *Repository) Close() error {
	db, err := r.Db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (r *Repository) CreateEmail(email *Email) error {
	return r.Db.WithContext(*r.Ctx).Create(email).Error
}

func (r *Repository) GetAllEmails() ([]Email, error) {
	var emails []Email
	err := r.Db.WithContext(*r.Ctx).Find(&emails).Error
	return emails, err
}

func (r *Repository) DeleteEmailByName(name string) error {
	return r.Db.WithContext(*r.Ctx).Where("name = ?", name).Delete(&Email{}).Error
}

func (r *Repository) GetEmailByName(name string) (*Email, error) {
	email := &Email{}
	err := r.Db.WithContext(*r.Ctx).Where("name = ?", name).First(email).Error
	return email, err
}

func (r *Repository) UpdateEmail(key uint, email *Email) error {
	return r.Db.WithContext(*r.Ctx).Model(&Email{}).Where("key = ?", key).Updates(email).Error
}

func (r *Repository) Ping() error {
	db, err := r.Db.WithContext(*r.Ctx).DB()
	if err != nil {
		return err
	}
	return db.Ping()
}
