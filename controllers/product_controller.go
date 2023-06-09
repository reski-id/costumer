package controllers

import (
	"costumer/models"
	"costumer/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductController struct{}

// @Summary Create a new product
// @Description Create a new product with the specified details
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer {token}"
// @Param product body models.Product true "Product details"
// @Success 201 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products [post]
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

// // GetProducts godoc
// // @Summary Get a list of products
// // @Description Get a list of products with pagination support
// // @Tags products
// // @Accept json
// // @Produce json
// // @Param page query int false "Page number (default 1)"
// // @Param limit query int false "Number of items per page (default 10)"
// // @Success 200 {object} models.ProductsResponse
// // @Failure 400 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Router /products [get]
// func (controller ProductController) GetProducts(c *gin.Context) {
// 	// all user can access
// 	db, err := utils.Connect()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
// 		return
// 	}
// 	var products []models.Product
// 	page, err := strconv.Atoi(c.Query("page"))
// 	if err != nil {
// 		page = 1
// 	}
// 	limit, err := strconv.Atoi(c.Query("limit"))
// 	if err != nil {
// 		limit = 10
// 	}
// 	offset := (page - 1) * limit
// 	if result := db.Offset(offset).Limit(limit).Find(&products); result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
// 		return
// 	}
// 	var total int64
// 	if result := db.Model(&models.Product{}).Count(&total); result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
// 		return
// 	}
// 	response := models.ProductsResponse{
// 		Data:  products,
// 		Count: total,
// 	}
// 	c.JSON(http.StatusOK, response)
// }

// // @Summary Get a product by ID
// // @Description Retrieve a product by ID
// // @Tags products
// // @Accept json
// // @Produce json
// // @Security ApiKeyAuth
// // @Param Authorization header string true "Bearer {token}"
// // @Param id path int true "Product ID"
// // @Success 200 {object} models.Product
// // @Failure 404 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// // @Router /products/{id} [get]
// func (controller ProductController) GetProduct(c *gin.Context) {
// 	// all user can access
// 	db, err := utils.Connect()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
// 		return
// 	}
// 	var product models.Product
// 	if result := db.First(&product, c.Param("id")); result.Error != nil {
// 		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, product)
// }

// UpdateProduct updates an existing product
// Only admin can update a product
// @Summary Update a product
// @Description Update a product with the specified ID
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product details"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id} [put]
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

// DeleteProduct deletes an existing product
// Only admin can delete a product
// @Summary Delete a product
// @Description Delete a product with the specified ID
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Product ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id} [delete]
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

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Product deleted successfully"})
}

// SearchProduct searches for products with a matching name
// @Summary Search products
// @Description Search products with a matching name
// @Tags products
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} models.Product
// @Failure 500 {object} models.ErrorResponse
// @Router /products/search [get]
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

// GetProducts godoc
// @Summary Get a list of products
// @Description Get a list of products with pagination support
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Number of items per page (default 10)"
// @Success 200 {object} models.ProductsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products [get]
func (controller ProductController) GetProducts(c *gin.Context) {
	// all user can access
	var response models.ProductsResponse

	// Check if the data exists in Redis
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	cacheKey := fmt.Sprintf("products:%s:%s", page, limit)

	// Initialize Redis client
	redisClient, err := utils.ConnectRedis(utils.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	cachedData, err := redisClient.Get(cacheKey).Result()
	if err == nil {
		// Data exists in Redis
		if err = json.Unmarshal([]byte(cachedData), &response); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// Data not found in Redis, query the database and cache the result
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	var products []models.Product

	pageInt, err := utils.StringToInt(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid page number"})
		return
	}

	limitInt, err := utils.StringToInt(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid limit"})
		return
	}

	offset := (pageInt - 1) * limitInt

	if result := db.Offset(offset).Limit(limitInt).Find(&products); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	var total int64
	if result := db.Model(&models.Product{}).Count(&total); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	response = models.ProductsResponse{
		Data:  products,
		Count: total,
	}

	// Cache the data
	jsonData, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	if _, err = redisClient.Set(cacheKey, jsonData, time.Minute*5).Result(); err != nil {
		log.Printf("Failed to cache data: %s\n", err.Error())
	}

	c.JSON(http.StatusOK, response)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Retrieve a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id} [get]
func (controller ProductController) GetProduct(c *gin.Context) {
	// Check if the data exists in Redis
	cacheKey := fmt.Sprintf("product:%s", c.Param("id"))

	// Initialize Redis client
	redisConfig := utils.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	redisClient, err := utils.ConnectRedis(redisConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	cachedData, err := redisClient.Get(cacheKey).Result()
	if err == nil {
		// Data exists in Redis
		var product models.Product
		if err = json.Unmarshal([]byte(cachedData), &product); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)
		return
	}

	// Data not found in Redis, query the database and cache the result
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

	// Cache the data
	jsonData, err := json.Marshal(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	if _, err = redisClient.Set(cacheKey, jsonData, time.Minute*5).Result(); err != nil {
		log.Printf("Failed to cache data: %s\n", err.Error())
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, product)
}
