package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name        string `json:"name" form:"name" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" validate:"required"`
	Address     string `json:"address" form:"address"`
}