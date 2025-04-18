package organization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/user"
	"hello-pulse.fr/internal/services/organization"
	"hello-pulse.fr/pkg/security"
)

// Handler handles organization API endpoints
type Handler struct {
	orgService *organization.Service
	securityService *security.AuthorizationService 
}

// NewHandler creates a new organization handler
func NewHandler(orgService *organization.Service, securityService *security.AuthorizationService) *Handler {
	return &Handler{
		orgService: orgService,
		securityService: securityService,
	}
}

// CreateOrganizationRequest represents the create organization request payload
type CreateOrganizationRequest struct {
	OrganizationName string `json:"name" binding:"required"`
}

// JoinOrganizationRequest represents the join organization request payload
type JoinOrganizationRequest struct {
	Code string `json:"code" binding:"required"`
}

// CreateInviteCodeRequest represents the create invite code request payload
type CreateInviteCodeRequest struct {
	ExpirationTimeMs int64 `json:"expirationTimeMs" binding:"required"`
}

// CreateOrganization handles organization creation
func (h *Handler) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user already belongs to an organization
	if user.OrganizationID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User already belongs to an organization",
		})
		return
	}

	org, err := h.orgService.CreateOrganization(req.OrganizationName, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Organization created successfully",
		"organization": gin.H{
			"id":   org.OrganizationID,
			"name": org.OrganizationName,
		},
	})
}

// JoinOrganization handles joining an organization with an invite code
func (h *Handler) JoinOrganization(c *gin.Context) {
	var req JoinOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user already belongs to an organization
	if user.OrganizationID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User already belongs to an organization",
		})
		return
	}

	if err := h.orgService.JoinOrganization(user.UserID, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Joined organization successfully",
	})
}

// CreateInviteCode handles creating an invite code for an organization
func (h *Handler) CreateInviteCode(c *gin.Context) {
	var req CreateInviteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user belongs to an organization
	if user.OrganizationID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User does not belong to an organization",
		})
		return
	}

	// Check if user can create invite codes
	canCreate, err := h.securityService.CanCreateInviteCode(c.Request.Context(), user.UserID, *user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to check permissions",
		})
		return
	}
	
	if !canCreate {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only administrators can create invite codes",
		})
		return
	}

	inviteCode, err := h.orgService.CreateInviteCode(*user.OrganizationID, req.ExpirationTimeMs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    inviteCode.Value,
	})
}

// GetInviteCodes handles retrieving all invite codes for an organization
func (h *Handler) GetInviteCodes(c *gin.Context) {
	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user belongs to an organization
	if user.OrganizationID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User does not belong to an organization",
		})
		return
	}

	// Check if user is an admin
	isAdmin, err := h.securityService.IsUserAdmin(c.Request.Context(), user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to check user permissions",
		})
		return
	}
	
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only administrators can view invite codes",
		})
		return
	}

	inviteCodes, err := h.orgService.GetInviteCodes(*user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve invite codes",
		})
		return
	}

	type codeResponse struct {
		ID               string `json:"id"`
		Code             string `json:"code"`
		ExpirationTimeMs int64  `json:"expirationTimeMs"`
	}

	var codes []codeResponse
	for _, code := range inviteCodes {
		codes = append(codes, codeResponse{
			ID:               code.InviteCodeID.String(),
			Code:             code.Value,
			ExpirationTimeMs: code.ExpirationTime.UnixMilli(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"codes":   codes,
	})
}

// DeleteInviteCode handles deleting an invite code for an organization
func (h *Handler) DeleteInviteCode(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user belongs to an organization
	if user.OrganizationID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User does not belong to an organization",
		})
		return
	}

	// Check if user is an admin
	isAdmin, err := h.securityService.IsUserAdmin(c.Request.Context(), user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to check user permissions",
		})
		return
	}
	
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only administrators can delete invite codes",
		})
		return
	}

	// Parse the invite code ID
	inviteCodeID, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid invite code ID",
		})
		return
	}

	// Delete the invite code
	if err := h.orgService.DeleteInviteCode(inviteCodeID, *user.OrganizationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Invite code deleted successfully",
	})
}