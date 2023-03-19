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
// @Accept json
// @Produce json
// @Param loginData body models.LoginData true "Login Data"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func (auth *AuthController) Login(c *gin.Context) {
	var loginData models.LoginData
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Token generation error"})
		return
	}

	response := models.TokenResponse{
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
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
	var registrationData models.User
	if err := c.ShouldBind(&registrationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database connection error"})
		return
	}

	var existingUser models.User
	result := db.Where("username = ?", registrationData.Username).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(registrationData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Password hashing error"})
		return
	}

	newUser := models.User{Username: registrationData.Username, Password: string(hash), Email: registrationData.Email, IsAdmin: registrationData.IsAdmin}
	result = db.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "User creation error"})
		return
	}

	token, err := utils.GenerateToken(newUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Token generation error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "username": newUser.Username, "email": newUser.Email, "is_admin": newUser.IsAdmin})
}
