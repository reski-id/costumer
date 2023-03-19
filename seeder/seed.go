package seeder

import (
	"log"

	"costumer/models"
	"costumer/utils"
)

func SeedUsers() {
	db, err := utils.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %s", err.Error())
	}

	// check if any user already exists in the database
	var user models.User
	if db.First(&user).Error == nil {
		log.Println("users already seeded")
		return
	}

	// migrate the user table
	db.AutoMigrate(&models.User{})

	// create some users
	users := []models.User{
		{Username: "john_doe", Password: "password1", IsAdmin: false, Email: "jhon@gmail.com", Name: "Jhon"},
		{Username: "jane_doe", Password: "password2", IsAdmin: true, Email: "adm@gmail.com", Name: "Jhane"},
		{Username: "bob_smith", Password: "password3", IsAdmin: false, Email: "bob@gmail.com", Name: "Bob"},
	}

	for i := range users {
		users[i].ID = uint(i) + 1
		err = db.Create(&users[i]).Error
		if err != nil {
			log.Fatalf("failed to seed users: %s", err.Error())
		}
	}

	log.Println("users seeded")
}

func SeedCustomers() {
	db, err := utils.Connect()
	if err != nil {
		panic("Failed to connect to database!")
	}

	// check if any customer already exists in the database
	var customer models.Customer
	if db.First(&customer).Error == nil {
		log.Println("customers already seeded")
		return
	}

	// create some sample customers
	customers := []models.Customer{
		{Name: "John Doe", Email: "john.doe@example.com", PhoneNumber: "1234567890", Address: "123 Main St"},
		{Name: "Jane Smith", Email: "jane.smith@example.com", PhoneNumber: "0987654321", Address: "456 Elm St"},
	}

	// insert customers into the database
	for _, c := range customers {
		result := db.Create(&c)
		if result.Error != nil {
			panic("Failed to insert customer!")
		}
	}
}

func SeedProducts() {
	db, err := utils.Connect()
	if err != nil {
		panic("Failed to connect to database!")
	}

	// check if any customer already exists in the database
	var product models.Product
	if db.First(&product).Error == nil {
		log.Println("products already seeded")
		return
	}

	// create some sample products
	products := []models.Product{
		{Name: "Dell XPS 13", Description: "13-inch laptop with 11th Gen Intel Core processor", Price: 1199.99, SKU: "LAPTOP-001", Quantity: 10, ImageURL: "https://example.com/laptop.jpg"},
		{Name: "iPhone 13", Description: "6.1-inch smartphone with A15 Bionic chip and 5G connectivity", Price: 799.99, SKU: "PHONE-001", Quantity: 5, ImageURL: "https://example.com/phone.jpg"},
		{Name: "Apple Watch Series 7", Description: "Smartwatch with always-on Retina display and ECG app", Price: 399.99, SKU: "WATCH-001", Quantity: 3, ImageURL: "https://example.com/watch.jpg"},
		{Name: "AirPods Pro", Description: "Wireless earbuds with active noise cancellation", Price: 249.99, SKU: "EARBUDS-001", Quantity: 7, ImageURL: "https://example.com/earbuds.jpg"},
		{Name: "Fitbit Charge 5", Description: "Advanced fitness tracker with EDA and ECG sensors", Price: 179.99, SKU: "FITNESS-001", Quantity: 2, ImageURL: "https://example.com/fitness.jpg"},
	}

	// insert products into the database
	for _, p := range products {
		result := db.Create(&p)
		if result.Error != nil {
			panic("Failed to insert product!")
		}
	}
}
