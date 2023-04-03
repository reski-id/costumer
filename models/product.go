package models

type Product struct {
	Model
	Name        string  `json:"name" form:"name" validate:"required"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price" validate:"required"`
	SKU         string  `json:"sku" form:"sku"`
	Quantity    int     `json:"qty" form:"qty" validate:"required"`
	ImageURL    string  `json:"file" form:"file"`
}

// ProductsResponse is the response object for GetProducts function
type ProductsResponse struct {
	Data  []Product `json:"data"`
	Count int64     `json:"count"`
}
