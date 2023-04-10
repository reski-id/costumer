package models

import (
	"time"
)

type Cart struct {
	ID         int        `json:"id"`
	UserID     int        `json:"userId" form:"userId" binding:"required"`
	Items      []CartItem `json:"items" gorm:"foreignKey:CartID"`
	TotalPrice float64    `json:"totalPrice"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type CartItem struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CartID    int       `json:"cartId" form:"cartId"`
	ProductID int       `json:"productId" form:"productId" binding:"required"`
	Quantity  uint32    `json:"quantity" form:"quantity" binding:"required"`
	Price     float64   `json:"price" form:"price" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
