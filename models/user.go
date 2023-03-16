package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"not null" json:"name"`
	Username  string `json:"username" form:"username" validate:"required"`
	Email     string `gorm:"not null;unique" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsAdmin   bool `gorm:"not null" json:"is_admin"`
}
