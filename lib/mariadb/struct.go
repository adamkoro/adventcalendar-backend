package mariadb

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Email struct {
	Key        uint      `gorm:"primaryKey" json:"key"`
	Name       string    `gorm:"type:varchar(255);unique;not null" json:"name" binding:"required" validate:"required,min=1,max=255"`
	From       string    `gorm:"type:varchar(255);not null" json:"from" binding:"required" validate:"required,min=1,max=255"`
	To         string    `gorm:"type:varchar(5000);not null" json:"to" binding:"required" validate:"required,min=1,max=5000"`
	Subject    string    `gorm:"type:varchar(255);not null" json:"subject" binding:"required" validate:"required,min=1,max=255"`
	Body       string    `gorm:"type:varchar(5000);not null" json:"body" binding:"required" validate:"required,min=1,max=5000"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	ModifiedAt time.Time `gorm:"autoUpdateTime" json:"modified_at"`
}

type UpdateEmail struct {
	Key     uint   `json:"key"`
	Name    string `json:"name"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type MQMessage struct {
	EmailTo string `json:"emailto" binding:"required" validate:"required,min=1,max=5000"`
	Subject string `json:"subject" binding:"required" validate:"required,min=1,max=255"`
	Message string `json:"message" binding:"required" validate:"required,min=1,max=5000"`
}

type EmailRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=255"`
}

type Repository struct {
	Db  *gorm.DB
	Ctx *context.Context
}
