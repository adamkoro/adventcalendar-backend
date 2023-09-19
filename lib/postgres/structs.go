package postgres

import (
	"context"
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

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
}
type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Repository struct {
	Db  *gorm.DB
	Ctx *context.Context
}
