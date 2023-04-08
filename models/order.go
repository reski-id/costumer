package models

import "github.com/google/uuid"

type Order struct {
	ID          uuid.UUID `json:"id" form:"id"`
	CustomerID  int       `json:"customerId" form:"customerId" validate:"required"`
	ProductID   int       `json:"productId" form:"productId" validate:"required"`
	Quantity    int       `json:"quantity" form:"quantity" validate:"required"`
	OrderStatus string    `json:"orderStatus" form:"orderStatus"`
}

type OrderPagination struct {
	Orders       []Order `json:"orders" form:"orders"`
	TotalRecords int64   `json:"totalRecords" form:"totalRecords"`
}
