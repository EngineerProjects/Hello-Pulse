package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/user"
	"hello-pulse.fr/internal/services/project"
)

// SummaryHandler handles project summary API endpoints
type SummaryHandler struct {
	summaryService *project.SummaryService
	projectService *project.Service
}

// NewSummaryHandler creates a new summary handler
func NewSummaryHandler(summaryService *project.SummaryService, projectService *project.Service) *SummaryHandler {
	return &SummaryHandler{
		summaryService: summaryService,
		projectService: projectService,
	}
}

// CreateSummaryRequest represents the create summary request payload
type CreateSummaryRequest struct {
	ProjectID string `json:"projectId" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// UpdateSummaryRequest represents the update summary request payload
type UpdateSummaryRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CreateSummary handles creating a new project summary
func (h *SummaryHandler) CreateSummary(c *gin.Context) {
	var req CreateSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	// Parse project ID
	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID",
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

	// Create summary
	summary, err := h.summaryService.CreateSummary(
		req.Title,
		req.Content,
		projectID,
		user.UserID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Summary created successfully",
		"summary": gin.H{
			"id":        summary.SummaryID,
			"title":     summary.Title,
			"content":   summary.Content,
			"projectId": summary.ProjectID,
			"createdBy": summary.CreatedBy,
			"createdAt": summary.CreatedAt,
			"updatedAt": summary.UpdatedAt,
		},
	})
}

// GetProjectSummaries handles retrieving all summaries for a project
func (h *SummaryHandler) GetProjectSummaries(c *gin.Context) {
	// Get project ID from URL parameter
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID",
		})
		return
	}

	// Get project details to verify user has access
	project, err := h.projectService.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Project not found",
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

	// Check if user has access to this project (belongs to same organization)
	if user.OrganizationID == nil || *user.OrganizationID != project.OrganizationID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Access denied",
		})
		return
	}

	// Get summaries
	summaries, err := h.summaryService.GetProjectSummaries(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve summaries",
		})
		return
	}

	// Format summaries for response
	var formattedSummaries []gin.H
	for _, summary := range summaries {
		formattedSummaries = append(formattedSummaries, gin.H{
			"id":        summary.SummaryID,
			"title":     summary.Title,
			"content":   summary.Content,
			"projectId": summary.ProjectID,
			"createdBy": summary.CreatedBy,
			"createdAt": summary.CreatedAt,
			"updatedAt": summary.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"summaries": formattedSummaries,
	})
}

// GetSummary handles retrieving a single summary by ID
func (h *SummaryHandler) GetSummary(c *gin.Context) {
	// Get summary ID from URL parameter
	summaryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid summary ID",
		})
		return
	}

	// Get summary
	summary, err := h.summaryService.GetSummary(summaryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Summary not found",
		})
		return
	}

	// Get project details to verify user has access
	project, err := h.projectService.GetProject(summary.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve project information",
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

	// Check if user has access to this project (belongs to same organization)
	if user.OrganizationID == nil || *user.OrganizationID != project.OrganizationID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Access denied",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"summary": gin.H{
			"id":        summary.SummaryID,
			"title":     summary.Title,
			"content":   summary.Content,
			"projectId": summary.ProjectID,
			"createdBy": summary.CreatedBy,
			"createdAt": summary.CreatedAt,
			"updatedAt": summary.UpdatedAt,
		},
	})
}

// UpdateSummary handles updating a summary
func (h *SummaryHandler) UpdateSummary(c *gin.Context) {
	// Get summary ID from URL parameter
	summaryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid summary ID",
		})
		return
	}

	var req UpdateSummaryRequest
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

	// Update summary
	if err := h.summaryService.UpdateSummary(summaryID, req.Title, req.Content, user.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Summary updated successfully",
	})
}

// DeleteSummary handles deleting a summary
func (h *SummaryHandler) DeleteSummary(c *gin.Context) {
	// Get summary ID from URL parameter
	summaryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid summary ID",
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

	// Delete summary
	if err := h.summaryService.DeleteSummary(summaryID, user.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Summary deleted successfully",
	})
}