package seeder

import (
	"costumer/models"
	"costumer/utils"
)

func CreateMigration() {
	db, err := utils.Connect()
	if err != nil {
		panic("Failed to connect to database!")
	}

	// Auto migrate all entities
	db.AutoMigrate(&models.User{}, &models.Order{}, &models.Product{}, &models.Upload{}, &models.Cart{}, &models.CartItem{})
}
