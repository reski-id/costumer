package controllers

import (
	"costumer/models"
	"costumer/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerController struct{}

// @Summary Get a list of customers
// @Description Get a list of customers, paginated by `page` and `limit` query parameters
// @Tags Customers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param page query integer false "Page number (default 1)"
// @Param limit query integer false "Number of customers per page (default 10)"
// @Success 200 {array} models.Customer
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /customers [get]
func (controller CustomerController) GetCustomers(c *gin.Context) {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can Access"})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database connection error"})
		return
	}

	var customers []models.Customer
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset := (page - 1) * limit

	result := db.Offset(offset).Limit(limit).Find(&customers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// @Summary Get a customer
// @Description Get a single customer by ID
// @Tags Customers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Customer ID"
// @Success 200 {object} models.Customer
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /customers/{id} [get]
func (controller CustomerController) GetCustomer(c *gin.Context) {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can Access"})
		return
	}
	fmt.Println(role)
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	var customer models.Customer
	result := db.First(&customer, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary Create a customer
// @Description Create a new customer
// @Tags Customers
// @Accept json
// @Produce json
// @Param customer body models.Customer true "Customer object"
// @Success 200 {object} models.Customer
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /customers [post]
func (controller CustomerController) CreateCustomer(c *gin.Context) {
	// all user can access
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var customer models.Customer
	err = c.ShouldBind(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	result := db.Create(&customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// UpdateCustomer godoc
// @Summary Update a customer by ID
// @Description Update a customer by ID
// @Tags Customers
// @Param id path int true "Customer ID"
// @Accept json
// @Produce json
// @Param customer body models.Customer true "Customer data"
// @Success 200 {object} models.Customer
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /customers/{id} [put]
func (controller CustomerController) UpdateCustomer(c *gin.Context) {
	// all user can access
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var customer models.Customer
	result := db.First(&customer, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Customer not found"})
		return
	}

	err = c.ShouldBind(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	result = db.Save(&customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer godoc
// @Summary Delete a customer by ID
// @Description Delete a customer by ID
// @Tags Customers
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Customer ID"
// @Produce json
// @Success 200 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /customers/{id} [delete]
func (controller CustomerController) DeleteCustomer(c *gin.Context) {
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

	var customer models.Customer
	result := db.First(&customer, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Customer not found"})
		return
	}

	result = db.Delete(&customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

// SearchCustomers godoc
// @Summary Search customers by name
// @Description Search customers by name
// @Tags Customers
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param query query string true "Search query"
// @Success 200 {object} []models.Customer
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /customers/search [get]
func (controller CustomerController) SearchCustomers(c *gin.Context) {
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

	var customers []models.Customer
	query := "%" + c.Query("query") + "%"

	result := db.Where("name LIKE ?", query).Find(&customers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}
