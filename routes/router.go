package routes

import (
	"costumer/controllers"
	"costumer/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Unauthenticated routes
	authController := controllers.AuthController{}
	router.POST("/login", authController.Login)

	// Authenticated routes
	api := router.Group("/api")
	api.Use(utils.AdminAuthMiddleware())
	{
		// Customer routes
		customerController := controllers.CustomerController{}
		api.GET("/customers", customerController.GetCustomers)
		api.GET("/customers/:id", customerController.GetCustomer)
		api.POST("/customers", customerController.CreateCustomer)
		api.PUT("/customers/:id", customerController.UpdateCustomer)
		api.DELETE("/customers/:id", customerController.DeleteCustomer)
		api.GET("/customers/search", customerController.SearchCustomers)

		// Order routes
		orderController := controllers.OrderController{}
		api.GET("/orders", orderController.GetOrders)
		api.GET("/orders/:id", orderController.GetOrder)
		api.POST("/orders", orderController.CreateOrder)
		api.PUT("/orders/:id", orderController.UpdateOrder)
		api.DELETE("/orders/:id", orderController.DeleteOrder)
		api.GET("/orders/search", orderController.SearchOrders)
	}

	return router
}
