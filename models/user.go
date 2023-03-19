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
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	SKU         string  `gorm:"unique;not null"`
	Quantity    int     `gorm:"not null"`
	ImageURL    string  `gorm:"not null"`
}
