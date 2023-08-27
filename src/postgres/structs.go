package db

import "time"

type User struct {
	Key        uint   `gorm:"primaryKey"`
	Username   string `gorm:"unique"`
	Email      string `gorm:"unique"`
	Password   string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	ModifiedAt time.Time `gorm:"autoUpdateTime"`
}
