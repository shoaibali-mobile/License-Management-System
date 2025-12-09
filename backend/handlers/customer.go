package handlers

import (
	"license-mnm/database"
	"license-mnm/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetCustomerSubscription returns customer's current subscription
func GetCustomerSubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var customer models.Customer
	if err := database.DB.Where("user_id = ?", userID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	var subscription models.Subscription
	if err := database.DB.Where("customer_id = ? AND status = ?", customer.ID, "active").Preload("Pack").First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No active subscription found"})
		return
	}

	isValid := subscription.ExpiresAt != nil && subscription.ExpiresAt.After(time.Now())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"subscription": gin.H{
			"id":      subscription.ID,
			"pack": gin.H{
				"name":            subscription.Pack.Name,
				"sku":             subscription.Pack.SKU,
				"price":           subscription.Pack.Price,
				"validity_months": subscription.Pack.ValidityMonths,
			},
			"status":       subscription.Status,
			"assigned_at":  subscription.AssignedAt,
			"expires_at":   subscription.ExpiresAt,
			"is_valid":     isValid,
		},
	})
}

// RequestSubscription creates a subscription request
func RequestSubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		SKU string `json:"sku" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var customer models.Customer
	if err := database.DB.Where("user_id = ?", userID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	// Check if customer has active subscription
	var activeSub models.Subscription
	if err := database.DB.Where("customer_id = ? AND status = ?", customer.ID, "active").First(&activeSub).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Customer already has an active subscription"})
		return
	}

	// Get pack by SKU
	var pack models.SubscriptionPack
	if err := database.DB.Where("sku = ? AND deleted_at IS NULL", req.SKU).First(&pack).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Subscription pack not found"})
		return
	}

	subscription := models.Subscription{
		CustomerID:  customer.ID,
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

// DeactivateSubscription deactivates customer's active subscription
func DeactivateSubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var customer models.Customer
	if err := database.DB.Where("user_id = ?", userID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	var subscription models.Subscription
	if err := database.DB.Where("customer_id = ? AND status = ?", customer.ID, "active").First(&subscription).Error; err != nil {
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

// GetSubscriptionHistory returns customer's subscription history
func GetSubscriptionHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "desc")
	offset := (page - 1) * limit

	var customer models.Customer
	if err := database.DB.Where("user_id = ?", userID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	var subscriptions []models.Subscription
	var total int64

	query := database.DB.Model(&models.Subscription{}).Where("customer_id = ?", customer.ID)
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





