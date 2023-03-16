package main

import (
	"github.com/gin-gonic/gin"

	"costumer/controllers"
	seed "costumer/seeder"
	"costumer/utils"
)

func main() {

	//migrate and seeder
	seed.CreateMigration()
	seed.SeedUsers()
	seed.SeedCustomers()
	seed.SeedProducts()

	router := gin.Default()

	customerController := controllers.CustomerController{}
	orderController := controllers.OrderController{}
	authController := controllers.AuthController{}
	productController := controllers.ProductController{}

	v1 := router.Group("/api/v1")

	v1.POST("/login", authController.Login)
	v1.Use(utils.AdminAuthMiddleware())

	v1.GET("/customers", customerController.GetCustomers)
	v1.GET("/customers/:id", customerController.GetCustomer)
	v1.POST("/customers", customerController.CreateCustomer)
	v1.PUT("/customers/:id", customerController.UpdateCustomer)
	v1.DELETE("/customers/:id", customerController.DeleteCustomer)
	v1.GET("/orders", orderController.GetOrders)
	v1.GET("/orders/:id", orderController.GetOrder)
	v1.POST("/orders", orderController.CreateOrder)
	v1.PUT("/orders/:id", orderController.UpdateOrder)
	v1.DELETE("/orders/:id", orderController.DeleteOrder)

	v1.GET("/products", productController.GetProducts)
	v1.GET("/products/:id", productController.GetProduct)
	v1.POST("/products", productController.CreateProduct)
	v1.PUT("/products/:id", productController.UpdateProduct)
	v1.DELETE("/products/:id", productController.DeleteProduct)

	router.Run(":8080")
}
