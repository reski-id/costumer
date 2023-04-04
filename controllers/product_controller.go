package controllers

import (
	"costumer/models"
	"costumer/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct{}

func (controller ProductController) CreateProduct(c *gin.Context) {
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

	var product models.Product
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	if result := db.Create(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (controller ProductController) GetProducts(c *gin.Context) {
	// all user can access
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

func (controller ProductController) GetProduct(c *gin.Context) {
	// all user can access
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

func (controller ProductController) UpdateProduct(c *gin.Context) {
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

	var product models.Product
	if result := db.First(&product, c.Param("id")); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found"})
		return
	}

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

func (controller ProductController) DeleteProduct(c *gin.Context) {
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

func (controller ProductController) SearchProduct(c *gin.Context) {

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var Product []models.Product
	query := "%" + c.Query("query") + "%"

	result := db.Where("name LIKE ?", query).Find(&Product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, Product)
}
