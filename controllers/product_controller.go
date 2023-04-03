package controllers

import (
	"costumer/models"
	"costumer/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct{}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param product body models.Product true "Product information"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /product [post]
func (controller ProductController) CreateProduct(c *gin.Context) {
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var product models.Product
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Check if the user is an admin
	// claims, exists := c.Get("claims")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
	// 	return
	// }
	// isAdmin := claims.(utils.Claims).IsAdmin
	// if !isAdmin {
	// 	c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
	// 	return
	// }

	if result := db.Create(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProducts retrieves all products with pagination support
// @Summary Get a list of products
// @Description Get a list of products with pagination support
// @Tags Products
// @ID get-products
// @Produce  json
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Number of products per page, default is 10"
// @Success 200 {object} ProductsResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products [get]
func (controller ProductController) GetProducts(c *gin.Context) {
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var products []models.Product

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	offset := (page - 1) * limit

	if result := db.Offset(offset).Limit(limit).Find(&products); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	var total int64
	if result := db.Model(&models.Product{}).Count(&total); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	response := models.ProductsResponse{
		Data:  products,
		Count: total,
	}

	c.JSON(http.StatusOK, response)
}

// GetProduct retrieves a product by ID
// @Summary Get a product by ID
// @Description Retrieve a product using its unique identifier
// @ID get-product-by-id
// @Produce json
// @Tags Products
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id} [get]
func (controller ProductController) GetProduct(c *gin.Context) {
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var product models.Product
	if result := db.First(&product, c.Param("id")); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates a product by ID
// @Summary Update a product
// @Description Update a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product object"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id} [put]
func (controller ProductController) UpdateProduct(c *gin.Context) {
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var product models.Product
	if result := db.First(&product, c.Param("id")); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found"})
		return
	}

	// // Check if the user is an admin
	// claims, exists := c.Get("claims")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
	// 	return
	// }
	// isAdmin := claims.(utils.Claims).IsAdmin
	// if !isAdmin {
	// 	c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
	// 	return
	// }

	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if result := db.Save(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product by ID
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id} [delete]
func (controller ProductController) DeleteProduct(c *gin.Context) {
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// // Check if the user is an admin
	// claims, exists := c.Get("claims")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
	// 	return
	// }
	// isAdmin := claims.(utils.Claims).IsAdmin
	// if !isAdmin {
	// 	c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
	// 	return
	// }

	// Get the product ID from the URL parameters
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	// Find the product with the specified ID
	var product models.Product
	result := db.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found"})
		return
	}

	// Delete the product
	result = db.Delete(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ErrorResponse{Error: "Product deleted successfully"})
}
