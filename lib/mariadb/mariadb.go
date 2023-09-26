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

func (r *Repository) DeleteEmailByName(email *DeleteEmailRequest) error {
	return r.Db.WithContext(*r.Ctx).Where("name = ?", email.Name).Delete(&Email{}).Error
}

func (r *Repository) GetEmailByName(email *EmailRequest) (*Email, error) {
	dbEmail := &Email{}
	err := r.Db.WithContext(*r.Ctx).Where("name = ?", email.Name).First(dbEmail).Error
	return dbEmail, err
}

func (r *Repository) UpdateEmail(email *UpdateEmailRequest) error {
	return r.Db.WithContext(*r.Ctx).Model(&Email{}).Where("`key` = ?", email.Key).Updates(Email{
		Name:    email.Name,
		From:    email.From,
		To:      email.To,
		Subject: email.Subject,
		Body:    email.Body,
	}).Error
}

func (r *Repository) Ping() error {
	db, err := r.Db.WithContext(*r.Ctx).DB()
	if err != nil {
		return err
	}
	return db.Ping()
}
