package mariadb

import (
	"time"

	"gorm.io/gorm"
)

type Email struct {
	Key        uint      `gorm:"primaryKey"`
	Name       string    `gorm:"type:varchar(255);unique;not null"`
	From       string    `gorm:"type:varchar(255);not null"`
	To         string    `gorm:"type:varchar(5000);not null"`
	Subject    string    `gorm:"type:varchar(255);not null"`
	Body       string    `gorm:"type:varchar(5000);not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	ModifiedAt time.Time `gorm:"autoUpdateTime"`
}

type Repository struct {
	db *gorm.DB
}
