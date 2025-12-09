package main

import (
	"license-mnm/database"
	"license-mnm/handlers"
	"license-mnm/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Create Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-API-Key"}
	r.Use(cors.New(config))

	// Public authentication endpoints (no auth required)
	api := r.Group("/api")
	{
		api.POST("/admin/login", handlers.AdminLogin)
		api.POST("/customer/login", handlers.CustomerLogin)
		api.POST("/customer/signup", handlers.CustomerSignup)
	}

	// Protected admin endpoints (JWT + Admin role required)
	adminV1 := r.Group("/api/v1/admin")
	adminV1.Use(middleware.AuthMiddleware())
	adminV1.Use(middleware.AdminOnly())
	{
		adminV1.GET("/dashboard", handlers.GetDashboard)
		adminV1.GET("/customers", handlers.ListCustomers)
		adminV1.POST("/customers", handlers.CreateCustomer)
		adminV1.GET("/customers/:customer_id", handlers.GetCustomer)
		adminV1.PUT("/customers/:customer_id", handlers.UpdateCustomer)
		adminV1.DELETE("/customers/:customer_id", handlers.DeleteCustomer)
		adminV1.GET("/subscription-packs", handlers.ListSubscriptionPacks)
		adminV1.POST("/subscription-packs", handlers.CreateSubscriptionPack)
		adminV1.PUT("/subscription-packs/:pack_id", handlers.UpdateSubscriptionPack)
		adminV1.DELETE("/subscription-packs/:pack_id", handlers.DeleteSubscriptionPack)
		adminV1.GET("/subscriptions", handlers.ListSubscriptions)
		adminV1.POST("/subscriptions/:subscription_id/approve", handlers.ApproveSubscription)
		adminV1.POST("/customers/:customer_id/assign-subscription", handlers.AssignSubscription)
		adminV1.DELETE("/customers/:customer_id/subscription/:subscription_id", handlers.UnassignSubscription)
	}

	// Protected customer endpoints (JWT + Customer role required)
	customerV1 := r.Group("/api/v1/customer")
	customerV1.Use(middleware.AuthMiddleware())
	customerV1.Use(middleware.CustomerOnly())
	{
		customerV1.GET("/subscription", handlers.GetCustomerSubscription)
		customerV1.POST("/subscription", handlers.RequestSubscription)
		customerV1.DELETE("/subscription", handlers.DeactivateSubscription)
		customerV1.GET("/subscription-history", handlers.GetSubscriptionHistory)
	}

	// SDK authentication (no auth required)
	sdk := r.Group("/sdk")
	{
		sdk.POST("/auth/login", handlers.SDKLogin)
	}

	// SDK protected endpoints (API Key required)
	sdkV1 := r.Group("/sdk/v1")
	sdkV1.Use(middleware.APIKeyAuth())
	{
		sdkV1.GET("/subscription", handlers.SDKGetSubscription)
		sdkV1.POST("/subscription", handlers.SDKRequestSubscription)
		sdkV1.DELETE("/subscription", handlers.SDKDeactivateSubscription)
		sdkV1.GET("/subscription-history", handlers.SDKGetSubscriptionHistory)
	}

	// Start server
	r.Run("0.0.0.0:8080")
}





