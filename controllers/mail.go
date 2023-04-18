package controllers

import (
	"costumer/models"
	"costumer/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailController struct{}

func (controller MailController) SendSingleEmail(c *gin.Context) {
	email := c.Param("email")
	customer := getCustomerByEmail(email)
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	from := mail.NewEmail("Your Company Name", "your-email@company.com")
	to := mail.NewEmail(customer.Fullname, customer.Email) // use customer's Fullname field instead of Name
	subject := "New Product Launch"
	plainTextContent := "Check out our latest product and get 20% off your first purchase!"
	htmlContent := "<strong>Check out our latest product and get 20% off your first purchase!</strong>"
	sendEmail(from, to, subject, plainTextContent, htmlContent)

	c.JSON(http.StatusOK, gin.H{"message": "email sent successfully"})
}

func (controller MailController) SendBatchEmail(c *gin.Context) {
	var customers []models.User
	db, err := utils.Connect()
	if err != nil {
		panic("failed to connect database")
	}

	result := db.Find(&customers)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get customers from database"})
		return
	}

	from := mail.NewEmail("Your Company Name", "your-email@company.com")
	emails := []*mail.Email{}
	for _, customer := range customers {
		to := mail.NewEmail(customer.Fullname, customer.Email)
		emails = append(emails, to)
	}

	subject := "New Product Launch"
	plainTextContent := "Check out our latest product and get 20% off your first purchase!"
	htmlContent := "<strong>Check out our latest product and get 20% off your first purchase!</strong>"
	sendBatchEmail(from, emails, subject, plainTextContent, htmlContent)

	c.JSON(http.StatusOK, gin.H{"message": "emails sent successfully"})
}

func getCustomerByEmail(email string) *models.User {
	var customer models.User
	db, err := utils.Connect()
	if err != nil {
		panic("failed to connect database")
	}
	result := db.Where("email = ?", email).First(&customer)
	if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return &customer
}

func sendEmail(from *mail.Email, to *mail.Email, subject string, plainTextContent string, htmlContent string) {
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient("SENDGRID_API_KEY")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
	fmt.Println(response.Headers)
}

func sendBatchEmail(from *mail.Email, emails []*mail.Email, subject string, plainTextContent string, htmlContent string) {
	personalization := mail.NewPersonalization()
	for _, to := range emails {
		personalization.AddTos(to)
	}
	message := mail.NewSingleEmail(from, subject, &mail.Email{}, plainTextContent, htmlContent)
	message.AddPersonalizations(personalization)

	client := sendgrid.NewSendClient("SENDGRID_API_KEY")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
	fmt.Println(response.Headers)
}
