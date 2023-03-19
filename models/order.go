package models

import (
	"github.com/google/uuid"
)

type Order struct {
	Model
	CustomerID uuid.UUID `json:"customerId" form:"customerId" validate:"required" gorm:"type:char(36);primary_key"`
	Customer   Customer  `json:"customer" gorm:"foreignKey:CustomerID"  form:"customer"  validate:"required"`
	Name       string    `json:"name" form:"name" validate:"required"`
	Quantity   int       `json:"quantity" form:"quantity" validate:"required"`
}

type OrderPagination struct {
	Orders       []Order `json:"orders" form:"orders"`
	TotalRecords int64   `json:"totalRecords" form:"totalRecords"`
}
