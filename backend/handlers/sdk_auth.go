package handlers

import (
	"license-mnm/database"
	"license-mnm/models"
	"license-mnm/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SDKLogin handles SDK authentication and generates API key
func SDKLogin(c *gin.Context) {
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

	// Generate or retrieve API key
	var apiKey string
	if user.APIKey == "" {
		newKey, err := utils.GenerateAPIKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate API key"})
			return
		}
		apiKey = newKey
		user.APIKey = apiKey
		database.DB.Save(&user)
	} else {
		apiKey = user.APIKey
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
		"api_key":    apiKey,
		"token":      token,
		"name":       name,
		"phone":      phone,
		"expires_in": 3600,
	})
}





