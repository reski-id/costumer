package models

import (
	"time"
)

type User struct {
	Model
	Name      string `json:"name" form:"name" binding:"required"`
	Username  string `json:"username" form:"username" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
	Email     string `json:"email" form:"email" binding:"required"`
	Role      string `json:"role" form:"email" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TokenResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"is_admin"`
}
