package handlers

import (
	"license-mnm/database"
	"license-mnm/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SDKGetSubscription returns current subscription for SDK
func SDKGetSubscription(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")

	var user models.User
	if err := database.DB.Where("api_key = ?", apiKey).Preload("Customer").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid API key"})
		return
	}

	if user.Customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	var subscription models.Subscription
	if err := database.DB.Where("customer_id = ? AND status = ?", user.Customer.ID, "active").Preload("Pack").First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No active subscription found"})
		return
	}

	isValid := subscription.ExpiresAt != nil && subscription.ExpiresAt.After(time.Now())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"subscription": gin.H{
			"id":         subscription.ID,
			"pack_name":  subscription.Pack.Name,
			"pack_sku":   subscription.Pack.SKU,
			"price":      subscription.Pack.Price,
			"status":     subscription.Status,
			"assigned_at": subscription.AssignedAt,
			"expires_at":  subscription.ExpiresAt,
			"is_valid":    isValid,
		},
	})
}

// SDKRequestSubscription creates a subscription request via SDK
func SDKRequestSubscription(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")

	var user models.User
	if err := database.DB.Where("api_key = ?", apiKey).Preload("Customer").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid API key"})
		return
	}

	if user.Customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	var req struct {
		PackSKU string `json:"pack_sku" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Check if customer has active subscription
	var activeSub models.Subscription
	if err := database.DB.Where("customer_id = ? AND status = ?", user.Customer.ID, "active").First(&activeSub).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Customer already has an active subscription"})
		return
	}

	// Get pack by SKU
	var pack models.SubscriptionPack
	if err := database.DB.Where("sku = ? AND deleted_at IS NULL", req.PackSKU).First(&pack).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Subscription pack not found"})
		return
	}

	subscription := models.Subscription{
		CustomerID:  user.Customer.ID,
		PackID:      pack.ID,
		Status:      "requested",
		RequestedAt: time.Now(),
	}

	if err := database.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create subscription request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Subscription request submitted successfully",
		"subscription": gin.H{
			"id":           subscription.ID,
			"status":       subscription.Status,
			"requested_at": subscription.RequestedAt,
		},
	})
}

// SDKDeactivateSubscription deactivates subscription via SDK
func SDKDeactivateSubscription(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")

	var user models.User
	if err := database.DB.Where("api_key = ?", apiKey).Preload("Customer").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid API key"})
		return
	}

	if user.Customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	var subscription models.Subscription
	if err := database.DB.Where("customer_id = ? AND status = ?", user.Customer.ID, "active").First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No active subscription found"})
		return
	}

	now := time.Now()
	subscription.Status = "inactive"
	subscription.DeactivatedAt = &now
	database.DB.Save(&subscription)

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Subscription deactivated successfully",
		"deactivated_at": now,
	})
}

// SDKGetSubscriptionHistory returns subscription history for SDK
func SDKGetSubscriptionHistory(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "desc")
	offset := (page - 1) * limit

	var user models.User
	if err := database.DB.Where("api_key = ?", apiKey).Preload("Customer").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid API key"})
		return
	}

	if user.Customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	var subscriptions []models.Subscription
	var total int64

	query := database.DB.Model(&models.Subscription{}).Where("customer_id = ?", user.Customer.ID)
	query.Count(&total)

	orderBy := "created_at DESC"
	if sort == "asc" {
		orderBy = "created_at ASC"
	}

	query.Order(orderBy).Offset(offset).Limit(limit).Preload("Pack").Find(&subscriptions)

	var history []map[string]interface{}
	for _, sub := range subscriptions {
		history = append(history, map[string]interface{}{
			"id":          sub.ID,
			"pack_name":   sub.Pack.Name,
			"status":      sub.Status,
			"assigned_at": sub.AssignedAt,
			"expires_at":  sub.ExpiresAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"history": history,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}







