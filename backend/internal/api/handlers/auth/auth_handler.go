package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hello-pulse.fr/internal/services/auth"
)

// Handler handles authentication API endpoints
type Handler struct {
	authService *auth.Service
}

// NewHandler creates a new authentication handler
func NewHandler(authService *auth.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

// RegisterRequest represents the register request payload
type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	user, token, err := h.authService.RegisterUser(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
		req.Phone,
		req.Address,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Set cookie
	c.SetCookie("token", token, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User registered successfully",
		"user": gin.H{
			"id":        user.UserID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid credentials",
		})
		return
	}

	// Set cookie
	c.SetCookie("token", token, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"user": gin.H{
			"id":        user.UserID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
	})
}

// Logout handles user logout
func (h *Handler) Logout(c *gin.Context) {
	token, err := c.Cookie("token")
	if err == nil {
		_ = h.authService.Logout(token)
	}

	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// Me returns the current authenticated user
func (h *Handler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}