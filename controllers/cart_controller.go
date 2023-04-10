package controllers

import (
	"costumer/models"
	"costumer/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartController struct{}

// @Summary Add item to cart
// @Description Add an item to the customer's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer {token}"
// @Param item body models.CartItem true "Cart item object"
// @Success 200 {object} models.Cart
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /cart [post]
func (controller CartController) AddToCart(c *gin.Context) {
	CustomerID, _, err := utils.ExtractData(c)
	if CustomerID == -1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Please login first"})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Unable to connect to database"})
		return
	}

	var item models.CartItem
	err = c.ShouldBind(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	cart := models.Cart{
		UserID: CustomerID,
	}

	result := db.Where(&cart).First(&cart)
	if result.Error != nil {
		db.Create(&cart)
	}

	item.ID = cart.ID

	result = db.Create(&item)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Unable to add item to cart"})
		return
	}
	c.JSON(http.StatusOK, cart)

	//problem
	// -duplicate entry for cart_item

	// next
	// get price from product * quantity save it on cart

}
