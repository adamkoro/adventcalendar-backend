package mariadb

import (
	"time"

	"gorm.io/gorm"
)

type Email struct {
	Key        uint      `gorm:"primaryKey"`
	Name       string    `gorm:"type:varchar(255);unique;not null" json:"name" binding:"required"`
	From       string    `gorm:"type:varchar(255);not null" json:"from" binding:"required"`
	To         string    `gorm:"type:varchar(5000);not null" json:"to" binding:"required"`
	Subject    string    `gorm:"type:varchar(255);not null" json:"subject" binding:"required"`
	Body       string    `gorm:"type:varchar(5000);not null" json:"body" binding:"required"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	ModifiedAt time.Time `gorm:"autoUpdateTime"`
}

type Repository struct {
	db *gorm.DB
}
