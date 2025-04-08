package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hello-pulse.fr/internal/services/auth"
)

// AuthMiddleware creates a middleware for authentication
func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized",
			})
			c.Abort()
			return
		}

		// Validate session
		user, err := authService.ValidateSession(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid session",
			})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)

		c.Next()
	}
}