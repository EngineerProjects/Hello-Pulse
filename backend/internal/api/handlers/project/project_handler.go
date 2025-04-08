package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/user"
	"hello-pulse.fr/internal/services/project"
)

// Handler handles project API endpoints
type Handler struct {
	projectService *project.Service
}

// NewHandler creates a new project handler
func NewHandler(projectService *project.Service) *Handler {
	return &Handler{
		projectService: projectService,
	}
}

// CreateProjectRequest represents the create project request payload
type CreateProjectRequest struct {
	ProjectName     string  `json:"projectName" binding:"required"`
	ProjectDesc     string  `json:"projectDesc"`
	ParentProjectID *string `json:"parentProjectId,omitempty"`
}

// UpdateProjectRequest represents the update project request payload
type UpdateProjectRequest struct {
	ProjectName string `json:"projectName" binding:"required"`
	ProjectDesc string `json:"projectDesc"`
}

// AddParticipantRequest represents the add participant request payload
type AddParticipantRequest struct {
	ProjectID string `json:"projectId" binding:"required"`
	UserID    string `json:"userId" binding:"required"`
}

// CreateProject handles project creation
func (h *Handler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
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

	if user.OrganizationID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User does not belong to an organization",
		})
		return
	}

	var parentProjectID *uuid.UUID
	if req.ParentProjectID != nil && *req.ParentProjectID != "" {
		parsed, err := uuid.Parse(*req.ParentProjectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid parent project ID",
			})
			return
		}
		parentProjectID = &parsed
	}

	project, err := h.projectService.CreateProject(
		req.ProjectName,
		req.ProjectDesc,
		user.UserID,
		*user.OrganizationID,
		parentProjectID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project created successfully",
		"project": project,
	})
}

// GetProjects handles retrieving all projects for the user's organization
func (h *Handler) GetProjects(c *gin.Context) {
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

	if user.OrganizationID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User does not belong to an organization",
		})
		return
	}

	// Get only root projects (no parent)
	projects, err := h.projectService.GetRootProjects(*user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve projects",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"projects": projects,
	})
}

// GetProject handles retrieving a single project by ID
func (h *Handler) GetProject(c *gin.Context) {
	id := c.Param("id")
	projectID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID",
		})
		return
	}

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

	// Get child projects
	childProjects, err := h.projectService.GetChildProjects(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve child projects",
		})
		return
	}

	// Get participants
	participants, err := h.projectService.GetProjectParticipants(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve project participants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"project":       project,
		"childProjects": childProjects,
		"participants":  participants,
	})
}

// UpdateProject handles updating a project
func (h *Handler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	projectID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID",
		})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	// Verify project exists and user has access
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

	// Check if user is the owner of the project
	if project.OwnerID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only the project owner can update the project",
		})
		return
	}

	if err := h.projectService.UpdateProject(projectID, req.ProjectName, req.ProjectDesc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update project",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project updated successfully",
	})
}

// DeleteProject handles deleting a project
func (h *Handler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	projectID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID",
		})
		return
	}

	// Verify project exists and user has access
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

	// Check if user is the owner of the project
	if project.OwnerID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only the project owner can delete the project",
		})
		return
	}

	if err := h.projectService.DeleteProject(projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete project",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project deleted successfully",
	})
}

// AddParticipant handles adding a user to a project
func (h *Handler) AddParticipant(c *gin.Context) {
	var req AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request payload",
		})
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID",
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	// Verify project exists and current user has access
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

	// Check if user is the owner of the project
	if project.OwnerID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only the project owner can add participants",
		})
		return
	}

	if err := h.projectService.AddParticipant(projectID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Participant added successfully",
	})
}