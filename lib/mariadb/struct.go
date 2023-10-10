package mariadb

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Email struct {
	Key        uint      `gorm:"primaryKey" json:"key"`
	Name       string    `gorm:"type:varchar(255);unique;not null" json:"name" binding:"required" validate:"required,min=1,max=255,alpha"`
	From       string    `gorm:"type:varchar(255);not null" json:"from" binding:"required" validate:"required,min=1,max=255,email"`
	To         string    `gorm:"type:varchar(5000);not null" json:"to" binding:"required" validate:"required,min=1,max=5000"`
	Subject    string    `gorm:"type:varchar(255);not null" json:"subject" binding:"required" validate:"required,min=1,max=255"`
	Body       string    `gorm:"type:varchar(5000);not null" json:"body" binding:"required" validate:"required,min=1,max=5000"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	ModifiedAt time.Time `gorm:"autoUpdateTime" json:"modified_at"`
}
type UpdateEmailRequest struct {
	Key     uint   `json:"key" binding:"required"`
	Name    string `json:"name" validate:"omitempty,min=1,max=255,alpha"`
	From    string `json:"from" validate:"omitempty,min=1,max=255,email"`
	To      string `json:"to" validate:"omitempty,min=1,max=5000"`
	Subject string `json:"subject" validate:"omitempty,min=1,max=255"`
	Body    string `json:"body" validate:"omitempty,min=1,max=5000"`
}
type DeleteEmailRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=255,alpha"`
}
type MQMessage struct {
	EmailTo string `json:"emailto" binding:"required" validate:"required,min=1,max=5000"`
	Subject string `json:"subject" binding:"required" validate:"required,min=1,max=255"`
	Message string `json:"message" binding:"required" validate:"required,min=1,max=5000"`
}
type EmailRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=255,alpha"`
}
type Repository struct {
	Db  *gorm.DB
	Ctx *context.Context
}
