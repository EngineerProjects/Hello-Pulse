package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/user"
	"hello-pulse.fr/pkg/security"
)

type contextKey string
const (
	// ContextKeyOrgID is the key for organization ID in the context
	ContextKeyOrgID contextKey = "organizationID"
	
	// ContextKeyUserID is the key for user ID in the context
	ContextKeyUserID contextKey = "userID"
)

// SecurityMiddleware adds security context to the request
func SecurityMiddleware(securityService *security.AuthorizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user from context (previously set by AuthMiddleware)
		currentUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized",
			})
			c.Abort()
			return
		}
		
		user := currentUser.(*user.User)
		
		// Add user ID to context
		ctx := context.WithValue(c.Request.Context(), ContextKeyUserID, user.UserID)
		
		// Add organization ID to context if user belongs to an organization
		if user.OrganizationID != nil {
			ctx = context.WithValue(ctx, ContextKeyOrgID, *user.OrganizationID)
		}
		
		// Update the request context
		c.Request = c.Request.WithContext(ctx)
		
		c.Next()
	}
}

// AdminRequiredMiddleware ensures that the user is an admin
func AdminRequiredMiddleware(securityService *security.AuthorizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user from context
		currentUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized",
			})
			c.Abort()
			return
		}
		
		user := currentUser.(*user.User)
		
		// Check if user is an admin
		if user.Role != security.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Admin privileges required",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// OrganizationRequiredMiddleware ensures the user belongs to an organization
func OrganizationRequiredMiddleware(securityService *security.AuthorizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user from context
		currentUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized",
			})
			c.Abort()
			return
		}
		
		user := currentUser.(*user.User)
		
		// Check if user belongs to an organization
		if user.OrganizationID == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "User does not belong to an organization",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// ResourceOwnerMiddleware ensures the user is the owner of the resource
func ResourceOwnerMiddleware(securityService *security.AuthorizationService, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get resource ID from URL parameter
		resourceID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid resource ID",
			})
			c.Abort()
			return
		}
		
		// Get current user from context
		currentUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized",
			})
			c.Abort()
			return
		}
		
		user := currentUser.(*user.User)
		
		var canAccess bool
		
		// Check if user can modify the resource
		switch resourceType {
		case "file":
			canAccess, err = securityService.CanModifyFile(c.Request.Context(), user.UserID, resourceID)
		case "project":
			canAccess, err = securityService.CanModifyProject(c.Request.Context(), user.UserID, resourceID)
		case "event":
			canAccess, err = securityService.CanModifyEvent(c.Request.Context(), user.UserID, resourceID)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Unknown resource type",
			})
			c.Abort()
			return
		}
		
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Resource not found",
			})
			c.Abort()
			return
		}
		
		if !canAccess {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "You do not have permission to modify this resource",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// Helper functions to extract values from context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value := ctx.Value(ContextKeyUserID)
	if value == nil {
		return uuid.Nil, false
	}
	
	userID, ok := value.(uuid.UUID)
	return userID, ok
}

func GetOrgIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value := ctx.Value(ContextKeyOrgID)
	if value == nil {
		return uuid.Nil, false
	}
	
	orgID, ok := value.(uuid.UUID)
	return orgID, ok
}