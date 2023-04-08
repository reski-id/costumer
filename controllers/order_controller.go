package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"costumer/models"
	"costumer/utils"
)

type OrderController struct{}

// @Summary Get orders
// @Description Retrieve a list of orders
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items to retrieve per page (default: 10)"
// @Success 200 {array} models.Order
// @Failure 500 {object} models.ErrorResponse
// @Router /orders [get]
func (controller OrderController) GetOrders(c *gin.Context) {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can Access"})
		return
	}
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var orders []models.Order
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset := (page - 1) * limit

	result := db.Offset(offset).Limit(limit).Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// @Summary Get order
// @Description Retrieve an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/{id} [get]
func (controller OrderController) GetOrder(c *gin.Context) {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can access"})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var order models.Order
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid order ID"})
		return
	}

	result := db.Where("id = ?", orderID).First(&order)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Order object"
// @Success 200 {object} models.Order
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders [post]
func (controller OrderController) CreateOrder(c *gin.Context) {
	CustomerID, _, err := utils.ExtractData(c)
	if CustomerID == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"Silahkan Login Terlebih dahulu": err.Error()})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Connect to Databases": err.Error()})
		return
	}

	var order models.Order
	err = c.ShouldBind(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error Error Bind": err.Error()})
		return
	}

	order.ID = uuid.New() //2c5d88f7-a8b4-4dda-a11a-dce64cddb9e7

	order.CustomerID = CustomerID
	order.OrderStatus = "Pending"

	result := db.Create(&order)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Create Order": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateOrderMulti godoc
// @Summary Create multiple orders
// @Description Create multiple orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param customerId path int true "Customer ID"
// @Param productId formData []int true "Product IDs"
// @Param quantity formData []int true "Quantities"
// @Success 200 {object} []models.Order
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /ordermulti [post]
func (controller OrderController) CreateOrderMulti(c *gin.Context) {
	CustomerID, _, err := utils.ExtractData(c)
	if CustomerID == -1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Silahkan Login Terlebih dahulu"})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "error Connect to Databases"})
		return
	}

	var orderRequest struct {
		ProductIDs []int `json:"productId" form:"productId" binding:"required"`
		Quantities []int `json:"quantity" form:"quantity" binding:"required"`
	}

	if err := c.ShouldBind(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "error Bind"})
		return
	}

	if len(orderRequest.ProductIDs) != len(orderRequest.Quantities) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "productId and quantities length must match"})
		return
	}

	var orders []models.Order
	for i := 0; i < len(orderRequest.ProductIDs); i++ {
		order := models.Order{
			ID:          uuid.New(),
			CustomerID:  CustomerID,
			ProductID:   orderRequest.ProductIDs[i],
			Quantity:    orderRequest.Quantities[i],
			OrderStatus: "Pending",
		}
		orders = append(orders, order)
	}

	result := db.Create(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "error Create Order"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// @Summary Update order
// @Description Update an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param order body models.Order true "Order object"
// @Success 200 {object} models.Order
// @Failure 404 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/{id} [put]
func (controller OrderController) UpdateOrder(c *gin.Context) {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can Access"})
		return
	}
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var order models.Order
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid order ID"})
		return
	}

	result := db.Where("id = ?", orderID).First(&order)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Order not found"})
		return
	}

	err = c.ShouldBind(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	result = db.Save(&order)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Delete an order
// @Description Delete an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/{id} [delete]
func (controller OrderController) DeleteOrder(c *gin.Context) {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can Access"})
		return
	}
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var order models.Order
	result := db.First(&order, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Order not found"})
		return
	}

	result = db.Delete(&order)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}

// @Summary Search orders
// @Description Search for orders by customer ID
// @Tags orders
// @Accept json
// @Produce json
// @Param query query string true "Customer ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items to retrieve per page (default: 10)"
// @Success 200 {array} models.Order
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/search [get]
func (controller OrderController) SearchOrders(c *gin.Context) {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can Access"})
		return
	}
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var orders []models.Order
	query := c.Query("query")

	result := db.Where("customer_id", query).Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// @Summary Get my orders
// @Description Retrieve a list of orders placed by the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} models.Order
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /myorder [get]
func (controller OrderController) GetMyOrders(c *gin.Context) {
	userID, _, err := utils.ExtractData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var orders []models.Order
	result := db.Where("customer_id = ?", userID).Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

//update my order
