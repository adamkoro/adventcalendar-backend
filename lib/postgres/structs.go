package postgres

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Key        uint      `gorm:"primaryKey"`
	Username   string    `gorm:"type:varchar(32);unique;not null"`
	Email      string    `gorm:"type:varchar(255);unique;not null"`
	Password   string    `gorm:"type:varchar(255);not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	ModifiedAt time.Time `gorm:"autoUpdateTime"`
}

type Repository struct {
	db *gorm.DB
}
