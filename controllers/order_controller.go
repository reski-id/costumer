package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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

	order.CustomerID = CustomerID
	order.OrderStatus = "Pending"

	result := db.Create(&order)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Create Order": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
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
	result := db.First(&order, c.Param("id"))
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
