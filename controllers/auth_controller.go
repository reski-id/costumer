package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"

	"costumer/models"
	"costumer/trace"
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

	_, span := trace.NewSpan(c, "AuthController.Login", nil)
	defer span.End()
	span.SetAttributes(
		attribute.String("Method", c.Request.Method),
		attribute.String("URL", c.Request.RequestURI),
		attribute.String("Proto", c.Request.Proto),
		attribute.String("RemoteAddr", c.Request.RemoteAddr),
		attribute.String("Host", c.Request.Host))
	if err := c.ShouldBind(&loginData); err != nil {
		trace.AddSpanError(span, err)
		trace.FailSpan(span, err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		trace.AddSpanError(span, err)
		trace.FailSpan(span, err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database connection error"})
		return
	}

	var user models.User
	result := db.Where("username = ?", loginData.Username).First(&user)
	if result.Error != nil {
		_, span := trace.NewSpan(c, "AuthController.Login:FindUser", nil)
		defer span.End()
		trace.AddSpanError(span, result.Error)
		trace.FailSpan(span, result.Error.Error())
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		_, span := trace.NewSpan(c, "AuthController.CompareHashAndPassword", nil)
		defer span.End()
		trace.AddSpanError(span, err)
		trace.FailSpan(span, err.Error())
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid username or password"})
		return
	}

	token, err := utils.GenerateToken(int(user.ID), user.Role)
	if err != nil {
		_, span := trace.NewSpan(c, "AuthController.GenerateToken", nil)
		defer span.End()
		trace.AddSpanError(span, err)
		trace.FailSpan(span, err.Error())
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
