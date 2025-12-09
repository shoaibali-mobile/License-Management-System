package handlers

import (
	"license-mnm/database"
	"license-mnm/models"
	"license-mnm/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminLogin handles admin login
func AdminLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ? AND role = ?", req.Email, "admin").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"token":      token,
		"email":      user.Email,
		"expires_in": 3600,
	})
}

// CustomerLogin handles customer login
func CustomerLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ? AND role = ?", req.Email, "customer").Preload("Customer").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	var name, phone string
	if user.Customer != nil {
		name = user.Customer.Name
		phone = user.Customer.Phone
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"token":      token,
		"name":       name,
		"phone":      phone,
		"expires_in": 3600,
	})
}

// CustomerSignup handles customer registration
func CustomerSignup(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Phone    string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password"})
		return
	}

	// Create user and customer in transaction
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		user := models.User{
			Email:        req.Email,
			PasswordHash: hashedPassword,
			Role:         "customer",
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		customer := models.Customer{
			UserID: user.ID,
			Name:   req.Name,
			Phone:  req.Phone,
		}
		if err := tx.Create(&customer).Error; err != nil {
			return err
		}

		// Generate token
		token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
		if err != nil {
			return err
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":    true,
			"message":    "Account created successfully",
			"token":      token,
			"name":       req.Name,
			"phone":      req.Phone,
			"expires_in": 3600,
		})

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create account"})
		return
	}
}

