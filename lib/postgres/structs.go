package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Key        uint      `gorm:"primaryKey"`
	Username   string    `gorm:"type:varchar(32);unique;not null" validate:"required,min=1,max=32"`
	Email      string    `gorm:"type:varchar(255);unique;not null" validate:"required,min=1,max=255"`
	Password   string    `gorm:"type:varchar(255);not null" validate:"required,min=1,max=255"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	ModifiedAt time.Time `gorm:"autoUpdateTime"`
}
type CreateUserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32"`
	Email    string `json:"email" binding:"required" validate:"required,min=1,max=255"`
	Password string `json:"password" binding:"required" validate:"required,min=1,max=255"`
}
type DeleteUserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32"`
}
type UpdateUserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32"`
}
type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32"`
	Password string `json:"password" binding:"required" validate:"required,min=1,max=255"`
}
type Repository struct {
	Db  *gorm.DB
	Ctx *context.Context
}
