package middleware

import (
	"net/http"
	"strings"

	"license-mnm/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminOnly middleware ensures only admin users can access
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CustomerOnly middleware ensures only customer users can access
func CustomerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "customer" {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Customer access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// APIKeyAuth middleware validates API key for SDK endpoints
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "X-API-Key header required"})
			c.Abort()
			return
		}

		// Validate API key in database
		// This will be implemented in the handler
		c.Set("api_key", apiKey)
		c.Next()
	}
}







