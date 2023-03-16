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

	// migrate the user table
	db.AutoMigrate(&models.User{})

	// create some users
	users := []models.User{
		{Username: "john_doe", Password: "password1", IsAdmin: false},
		{Username: "jane_doe", Password: "password2", IsAdmin: true},
		{Username: "bob_smith", Password: "password3", IsAdmin: false},
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
