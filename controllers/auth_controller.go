package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"costumer/models"
	"costumer/utils"
)

type AuthController struct{}

func (auth *AuthController) Login(c *gin.Context) {
	var loginData models.LoginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	var user models.User
	result := db.Where("username = ?", loginData.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (auth *AuthController) InsertLoginData(username, password string) error {
	db, err := utils.Connect()
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	loginData := models.LoginData{Username: username, Password: string(hash)}
	result := db.Create(&loginData)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
