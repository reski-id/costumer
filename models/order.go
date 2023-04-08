package models

type Order struct {
	ID          string `json:"id" form:"id"`
	CustomerID  int    `json:"customerId" form:"customerId" validate:"required"`
	ProductID   int    `json:"productId" form:"productId" validate:"required"`
	Quantity    int    `json:"quantity" form:"quantity" validate:"required"`
	OrderStatus string `json:"orderStatus" form:"orderStatus"`
}

type OrderPagination struct {
	Orders       []Order `json:"orders" form:"orders"`
	TotalRecords int64   `json:"totalRecords" form:"totalRecords"`
}
