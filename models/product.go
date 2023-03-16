package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	SKU         string  `gorm:"unique;not null"`
	Quantity    int     `gorm:"not null"`
	ImageURL    string  `gorm:"not null"`
}
