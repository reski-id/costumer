package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"costumer/models"
	"costumer/utils"
)

type AuthController struct{}

// Login godoc
// @Summary Login to the system
// @Description Login to the system with username and password
// @Tags Auth
// @Accept json, multipart/form-data
// @Produce json
// @Param loginData body models.LoginData true "Login Data"
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func (auth *AuthController) Login(c *gin.Context) {
	var loginData models.LoginData
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database connection error"})
		return
	}

	var user models.User
	result := db.Where("username = ?", loginData.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid username or password"})
		return
	}

	token, err := utils.GenerateToken(int(user.ID), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Token generation error"})
		return
	}

	response := models.TokenResponse{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Token:    token,
	}

	c.JSON(http.StatusOK, response)
}

// Register godoc
// @Summary Register to the system
// @Description Register to the system with username, password, email, and isAdmin flag
// @Tags Auth
// @Accept json
// @Produce json
// @Param registrationData body models.User true "Registration Data"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func (auth *AuthController) Register(c *gin.Context) {
	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	var customer models.User
	err = c.ShouldBind(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	//cek
	var existingUser models.User
	checkUsername := db.Where("username = ?", customer.Username).First(&existingUser)
	if checkUsername.Error == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Username already exists"})
		return
	}
	checkEmail := db.Where("email = ?", customer.Email).First(&existingUser)
	if checkEmail.Error == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Email already exists"})
		return
	}

	customer.Role = "user"

	hash, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Password hashing error"})
		return
	}

	newCustomer := models.User{
		Fullname:    customer.Fullname,
		Username:    customer.Username,
		Password:    string(hash),
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Address:     customer.Address,
		Role:        customer.Role}

	result := db.Create(&newCustomer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	response := models.CreateUserResponse{
		Fullname:    newCustomer.Fullname,
		Username:    newCustomer.Username,
		Password:    customer.Password,
		Email:       newCustomer.Email,
		Role:        newCustomer.Role,
		PhoneNumber: newCustomer.PhoneNumber,
		Address:     newCustomer.Address,
	}

	c.JSON(http.StatusOK, response)
}
