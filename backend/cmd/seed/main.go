package main

import (
	"license-mnm/database"
	"license-mnm/models"
	"license-mnm/utils"
	"log"
)

// Run this once to create an admin user
// Usage: go run cmd/seed/main.go
func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Check if admin already exists
	var existingAdmin models.User
	if err := database.DB.Where("email = ? AND role = ?", "admin@example.com", "admin").First(&existingAdmin).Error; err == nil {
		log.Println("Admin user already exists")
		return
	}

	// Create admin user
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	admin := models.User{
		Email:        "admin@example.com",
		PasswordHash: hashedPassword,
		Role:         "admin",
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		log.Fatal("Failed to create admin:", err)
	}

	log.Println("Admin user created successfully!")
	log.Println("Email: admin@example.com")
	log.Println("Password: admin123")
}





