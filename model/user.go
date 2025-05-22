package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id" gorm:"column:id;primary_key"`
	Email     string         `json:"email" gorm:"column:email;unique;not null"`
	Password  string         `json:"-" gorm:"column:password;not null"` // "-" means this field won't be included in JSON
	FirstName string         `json:"firstName" gorm:"column:first_name"`
	LastName  string         `json:"lastName" gorm:"column:last_name"`
	Role      string         `json:"role" gorm:"column:role;default:'user'"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
