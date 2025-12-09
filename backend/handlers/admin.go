package handlers

import (
	"license-mnm/database"
	"license-mnm/models"
	"license-mnm/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetDashboard returns admin dashboard data
func GetDashboard(c *gin.Context) {
	var totalCustomers int64
	var activeSubscriptions int64
	var pendingRequests int64
	var totalRevenue float64

	database.DB.Model(&models.Customer{}).Where("deleted_at IS NULL").Count(&totalCustomers)
	database.DB.Model(&models.Subscription{}).Where("status = ?", "active").Count(&activeSubscriptions)
	database.DB.Model(&models.Subscription{}).Where("status = ?", "requested").Count(&pendingRequests)

	// Calculate total revenue from active subscriptions
	var subscriptions []models.Subscription
	database.DB.Where("status = ?", "active").Preload("Pack").Find(&subscriptions)
	for _, sub := range subscriptions {
		totalRevenue += sub.Pack.Price
	}

	// Get recent activities (last 10 subscriptions)
	var recentActivities []map[string]interface{}
	var recentSubs []models.Subscription
	database.DB.Order("created_at DESC").Limit(10).Preload("Customer").Preload("Pack").Find(&recentSubs)

	for _, sub := range recentSubs {
		activity := map[string]interface{}{
			"type":      "subscription_" + sub.Status,
			"customer":  sub.Customer.Name,
			"pack":      sub.Pack.Name,
			"timestamp": sub.CreatedAt.Format(time.RFC3339),
		}
		recentActivities = append(recentActivities, activity)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_customers":      totalCustomers,
			"active_subscriptions":  activeSubscriptions,
			"pending_requests":     pendingRequests,
			"total_revenue":        totalRevenue,
			"recent_activities":    recentActivities,
		},
	})
}

// ListCustomers returns paginated list of customers
func ListCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	offset := (page - 1) * limit

	var customers []models.Customer
	var total int64

	query := database.DB.Model(&models.Customer{}).Where("deleted_at IS NULL")
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	query.Offset(offset).Limit(limit).Preload("User").Find(&customers)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"customers": customers,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CreateCustomer creates a new customer
func CreateCustomer(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Check if email exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email already registered"})
		return
	}

	// Create user with default password
	hashedPassword, _ := utils.HashPassword("password123") // Default password
	user := models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "customer",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user"})
		return
	}

	customer := models.Customer{
		UserID: user.ID,
		Name:   req.Name,
		Phone:  req.Phone,
	}

	if err := database.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"customer": customer,
	})
}

// GetCustomer returns customer details
func GetCustomer(c *gin.Context) {
	id := c.Param("customer_id")

	var customer models.Customer
	if err := database.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("User").Preload("Subscriptions.Pack").First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"customer": customer,
	})
}

// UpdateCustomer updates customer information
func UpdateCustomer(c *gin.Context) {
	id := c.Param("customer_id")

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var customer models.Customer
	if err := database.DB.Where("id = ? AND deleted_at IS NULL", id).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}

	database.DB.Save(&customer)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"customer": customer,
	})
}

// DeleteCustomer soft deletes a customer
func DeleteCustomer(c *gin.Context) {
	id := c.Param("customer_id")

	now := time.Now()
	if err := database.DB.Model(&models.Customer{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Customer deleted successfully",
	})
}

// ListSubscriptionPacks returns paginated list of subscription packs
func ListSubscriptionPacks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var packs []models.SubscriptionPack
	var total int64

	database.DB.Model(&models.SubscriptionPack{}).Where("deleted_at IS NULL").Count(&total)
	database.DB.Where("deleted_at IS NULL").Offset(offset).Limit(limit).Find(&packs)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"packs":   packs,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CreateSubscriptionPack creates a new subscription pack
func CreateSubscriptionPack(c *gin.Context) {
	var req struct {
		Name          string  `json:"name" binding:"required"`
		Description   string  `json:"description"`
		SKU           string  `json:"sku" binding:"required"`
		Price         float64 `json:"price" binding:"required"`
		ValidityMonths int    `json:"validity_months" binding:"required,min=1,max=12"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	pack := models.SubscriptionPack{
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		Price:         req.Price,
		ValidityMonths: req.ValidityMonths,
	}

	if err := database.DB.Create(&pack).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "SKU already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"pack":    pack,
	})
}

// UpdateSubscriptionPack updates a subscription pack
func UpdateSubscriptionPack(c *gin.Context) {
	id := c.Param("pack_id")

	var req struct {
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		SKU           string  `json:"sku"`
		Price         float64 `json:"price"`
		ValidityMonths int    `json:"validity_months"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var pack models.SubscriptionPack
	if err := database.DB.Where("id = ? AND deleted_at IS NULL", id).First(&pack).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Subscription pack not found"})
		return
	}

	if req.Name != "" {
		pack.Name = req.Name
	}
	if req.Description != "" {
		pack.Description = req.Description
	}
	if req.SKU != "" {
		pack.SKU = req.SKU
	}
	if req.Price > 0 {
		pack.Price = req.Price
	}
	if req.ValidityMonths > 0 {
		pack.ValidityMonths = req.ValidityMonths
	}

	database.DB.Save(&pack)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"pack":    pack,
	})
}

// DeleteSubscriptionPack soft deletes a subscription pack
func DeleteSubscriptionPack(c *gin.Context) {
	id := c.Param("pack_id")

	now := time.Now()
	if err := database.DB.Model(&models.SubscriptionPack{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Subscription pack not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Subscription pack deleted successfully",
	})
}

// ListSubscriptions returns all subscriptions
func ListSubscriptions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	offset := (page - 1) * limit

	var subscriptions []models.Subscription
	var total int64

	query := database.DB.Model(&models.Subscription{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset(offset).Limit(limit).Preload("Customer").Preload("Pack").Find(&subscriptions)

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"subscriptions": subscriptions,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ApproveSubscription approves a subscription request
func ApproveSubscription(c *gin.Context) {
	id := c.Param("subscription_id")

	var subscription models.Subscription
	if err := database.DB.Where("id = ?", id).First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Subscription not found"})
		return
	}

	if subscription.Status != "requested" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Subscription is not in requested status"})
		return
	}

	now := time.Now()
	subscription.Status = "approved"
	subscription.ApprovedAt = &now
	database.DB.Save(&subscription)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Subscription approved successfully",
	})
}

// AssignSubscription assigns a subscription pack to a customer
func AssignSubscription(c *gin.Context) {
	customerID := c.Param("customer_id")

	var req struct {
		PackID int `json:"pack_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Check if customer has active subscription
	var activeSubCount int64
	if err := database.DB.Model(&models.Subscription{}).Where("customer_id = ? AND status = ?", customerID, "active").Count(&activeSubCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to check existing subscription"})
		return
	}
	
	if activeSubCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Customer already has an active subscription"})
		return
	}

	// Get pack
	var pack models.SubscriptionPack
	if err := database.DB.Where("id = ? AND deleted_at IS NULL", req.PackID).First(&pack).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Subscription pack not found"})
		return
	}

	now := time.Now()
	expiresAt := now.AddDate(0, pack.ValidityMonths, 0)

	subscription := models.Subscription{
		CustomerID:  uint(parseInt(customerID)),
		PackID:      uint(req.PackID),
		Status:      "active",
		AssignedAt:  &now,
		ExpiresAt:   &expiresAt,
	}

	if err := database.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to assign subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Subscription assigned successfully",
	})
}

// UnassignSubscription removes a subscription assignment
func UnassignSubscription(c *gin.Context) {
	customerID := c.Param("customer_id")
	subscriptionID := c.Param("subscription_id")

	if err := database.DB.Where("id = ? AND customer_id = ?", subscriptionID, customerID).Delete(&models.Subscription{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Subscription unassigned successfully",
	})
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

