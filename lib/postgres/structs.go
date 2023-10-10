package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Key        uint      `gorm:"primaryKey"`
	Username   string    `gorm:"type:varchar(32);unique;not null" validate:"required,min=1,max=32,alphanum"`
	Email      string    `gorm:"type:varchar(255);unique;not null" validate:"required,min=1,max=255,email"`
	Password   string    `gorm:"type:varchar(255);not null" validate:"required,min=1,max=255"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	ModifiedAt time.Time `gorm:"autoUpdateTime"`
}
type CreateUserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32,alphanum"`
	Email    string `json:"email" binding:"required" validate:"required,min=1,max=255,email"`
	Password string `json:"password" binding:"required" validate:"required,min=1,max=255"`
}
type DeleteUserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32,alphanum"`
}
type UpdateUserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32,alphanum"`
	Email    string `json:"email" validate:"omitempty,min=1,max=255,email"`
	Password string `json:"password" validate:"omitempty"`
}
type UserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32,alphanum"`
}
type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=1,max=32,alphanum"`
	Password string `json:"password" binding:"required" validate:"required,min=1,max=255"`
}
type Repository struct {
	Db  *gorm.DB
	Ctx *context.Context
}
