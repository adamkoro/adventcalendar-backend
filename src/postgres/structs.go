package db

import "time"

type User struct {
	Key        uint      `gorm:"primaryKey"`
	Username   string    `gorm:"unique, not null"`
	Email      string    `gorm:"unique, not null"`
	Password   string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	ModifiedAt time.Time `gorm:"autoUpdateTime"`
}
