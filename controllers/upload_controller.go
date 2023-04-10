package controllers

import (
	"costumer/models"
	"costumer/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadController struct{}

// @Summary Create a new asset
// @Description Upload a new asset file and save its metadata to the database
// @Tags assets
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer {token}"
// @Param name formData string true "Asset name"
// @Param file formData file true "Asset file"
// @Success 201 {object} models.Upload
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /images [post]
func (controller UploadController) UploadAsset(c *gin.Context) {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can access"})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	name := c.PostForm("name")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Generate a unique filename for the uploaded file
	ext := filepath.Ext(file.Filename)
	filename := uuid.NewString() + ext

	// Save the uploaded file to the server
	err = c.SaveUploadedFile(file, filepath.Join("assets", filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	asset := models.Upload{Name: name, File: filename}
	if result := db.Create(&asset); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, asset)
}

// @Summary Delete an asset
// @Description Delete an asset file and its metadata from the database
// @Tags assets
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Asset ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /images/{id} [delete]
func (controller UploadController) DeleteAsset(c *gin.Context) {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can access"})
		return
	}

	db, err := utils.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	id := c.Param("id")

	asset := models.Upload{}
	if result := db.Where("id = ?", id).First(&asset); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Asset not found"})
		return
	}

	err = os.Remove(filepath.Join("assets", asset.File))
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		}
		return
	}

	if result := db.Delete(&asset); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "File Deleted Succesfully"})
}

func (controller UploadController) UploadAssetUsingS3(c *gin.Context) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}

	// Generate unique filename for S3 object
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext

	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_BUCKET := os.Getenv("AWS_BUCKET")
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region:      &AWS_REGION,
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Create S3 uploader
	uploader := s3manager.NewUploader(sess)

	// Open file
	fileObj, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	defer fileObj.Close()

	// Upload file to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: &AWS_BUCKET,
		Key:    aws.String(filename),
		Body:   fileObj,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"url": fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", AWS_BUCKET, AWS_REGION, filename),
	})
}

//delete

//upload to cgp

//delete from cgp
