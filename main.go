package main

import (
	"costumer/controllers"

	"github.com/gin-gonic/gin"

	docs "costumer/docs"
	seed "costumer/seeder"

	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           Swagger Costumer APP
// @version         2.0
// @description     This is a swagger documentation for Costumer APP.
// @BasePath        /api/v1
// @host            localhost:8080
// @schemes         http https
// @SecurityDefinition  jwt
// @Security        jwt
func main() {

	//migrate and seeder
	seed.CreateMigration()
	seed.SeedUsers()
	seed.SeedProducts()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	docs.SwaggerInfo.BasePath = "/api/v1"

	customerController := controllers.CustomerController{}
	orderController := controllers.OrderController{}
	authController := controllers.AuthController{}
	productController := controllers.ProductController{}

	v1 := router.Group("/api/v1")

	v1.POST("/login", authController.Login)
	v1.POST("/register", authController.Register)

	v1.POST("/customers", customerController.CreateCustomer)
	v1.PUT("/customers/:id", customerController.UpdateCustomer)
	v1.DELETE("/customers/:id", customerController.DeleteCustomer)
	v1.GET("/customers", customerController.GetCustomers)
	v1.GET("/customers/:id", customerController.GetCustomer)
	v1.GET("/customers/search", customerController.SearchCustomers)

	v1.POST("/products", productController.CreateProduct)
	v1.PUT("/products/:id", productController.UpdateProduct)
	v1.DELETE("/products/:id", productController.DeleteProduct)
	v1.GET("/products", productController.GetProducts)
	v1.GET("/products/:id", productController.GetProduct)
	v1.GET("/products/search", productController.SearchProduct)

	v1.POST("/orders", orderController.CreateOrder)
	v1.POST("/multiorder", orderController.CreateOrderMulti)
	v1.PUT("/orders/:id", orderController.UpdateOrder)
	v1.DELETE("/orders/:id", orderController.DeleteOrder)
	v1.GET("/orders", orderController.GetOrders)
	v1.GET("/myorder", orderController.GetMyOrders)
	v1.PUT("/myorder/:id", orderController.UpdateMyOrder)
	v1.GET("/orders/:id", orderController.GetOrder)
	v1.GET("/orders/search", orderController.SearchOrders)

	router.Run(":8080")
}
