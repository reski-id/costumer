package controllers

import (
	"costumer/models"
	"costumer/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "File Required"})
		return
	}

	// Generate a unique filename for the uploaded file
	ext := filepath.Ext(file.Filename)
	filename := uuid.NewString() + ext

	// Open the file using the Open method of the FileHeader type
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	defer f.Close()

	// Create a new S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Replace with your preferred region
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Upload the file to the S3 bucket
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String("myawsgin"), // Replace with your bucket name
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Return the uploaded file metadata
	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"url":      "https://myawsgin.s3.us-east-1.amazonaws.com/" + filename,
	})
}

func (controller UploadController) DeleteAssetsInS3(c *gin.Context) {
	// Get the filename to be deleted from the request URL parameter
	filename := c.Param("filename")

	// Create a new S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Replace with your preferred region
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Delete the file from the S3 bucket
	_, err = s3.New(sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("aws-bucket-alterra"), // Replace with your bucket name
		Key:    aws.String(filename),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, models.MessageResponse{Message: "File Deleted Succesfully"})
}

//upload to cgp

//delete from cgp
