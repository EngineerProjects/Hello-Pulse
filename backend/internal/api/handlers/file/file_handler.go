package file

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/user"
	"hello-pulse.fr/internal/services/file"
)

// Handler handles file API endpoints
type Handler struct {
	fileService *file.Service
}

// NewHandler creates a new file handler
func NewHandler(fileService *file.Service) *Handler {
	return &Handler{
		fileService: fileService,
	}
}

// UploadFile handles file upload
func (h *Handler) UploadFile(c *gin.Context) {
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

	// Parse isPublic value from form
	isPublic := c.PostForm("isPublic") == "true"

	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No file uploaded",
		})
		return
	}

	// Check file size - 100 MB limit
	if file.Size > 100*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "File size exceeds the 100 MB limit",
		})
		return
	}

	// Upload file
	uploadedFile, err := h.fileService.UploadFile(
		c.Request.Context(),
		file,
		user.UserID,
		*user.OrganizationID,
		isPublic,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upload file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"file": gin.H{
			"id":          uploadedFile.ID,
			"fileName":    uploadedFile.FileName,
			"contentType": uploadedFile.ContentType,
			"size":        uploadedFile.Size,
			"uploadedAt":  uploadedFile.UploadedAt,
			"isPublic":    uploadedFile.IsPublic,
		},
	})
}

// GetUserFiles handles retrieving files uploaded by the current user
func (h *Handler) GetUserFiles(c *gin.Context) {
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

	// Parse includeDeleted query parameter
	includeDeleted := c.Query("includeDeleted") == "true"

	// Get files
	files, err := h.fileService.GetUserFiles(user.UserID, includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve files",
		})
		return
	}

	// Format response
	var formattedFiles []gin.H
	for _, file := range files {
		formattedFiles = append(formattedFiles, gin.H{
			"id":          file.ID,
			"fileName":    file.FileName,
			"contentType": file.ContentType,
			"size":        file.Size,
			"uploadedAt":  file.UploadedAt,
			"isPublic":    file.IsPublic,
			"isDeleted":   file.IsDeleted,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"files":   formattedFiles,
	})
}

// GetOrganizationFiles handles retrieving files for the user's organization
func (h *Handler) GetOrganizationFiles(c *gin.Context) {
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

	// Parse includeDeleted query parameter
	includeDeleted := c.Query("includeDeleted") == "true"

	// Get files
	files, err := h.fileService.GetOrganizationFiles(*user.OrganizationID, includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve files",
		})
		return
	}

	// Format response
	var formattedFiles []gin.H
	for _, file := range files {
		formattedFiles = append(formattedFiles, gin.H{
			"id":          file.ID,
			"fileName":    file.FileName,
			"contentType": file.ContentType,
			"size":        file.Size,
			"uploadedAt":  file.UploadedAt,
			"isPublic":    file.IsPublic,
			"isDeleted":   file.IsDeleted,
			"uploaderId":  file.UploaderID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"files":   formattedFiles,
	})
}

// GetFileURL handles generating a presigned URL for a file
func (h *Handler) GetFileURL(c *gin.Context) {
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

	// Parse file ID from URL parameter
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	// Get file URL
	url, err := h.fileService.GetFileURL(c.Request.Context(), fileID, user.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"url":     url,
	})
}

// SoftDeleteFile handles soft-deleting a file
func (h *Handler) SoftDeleteFile(c *gin.Context) {
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

	// Parse file ID from URL parameter
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	// Soft delete the file
	if err := h.fileService.SoftDeleteFile(c.Request.Context(), fileID, user.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File deleted successfully",
	})
}

// RestoreFile handles restoring a soft-deleted file
func (h *Handler) RestoreFile(c *gin.Context) {
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

	// Parse file ID from URL parameter
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	// Restore the file
	if err := h.fileService.RestoreFile(c.Request.Context(), fileID, user.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File restored successfully",
	})
}

// GetFileTypes returns a list of supported file types and their extensions
func (h *Handler) GetFileTypes(c *gin.Context) {
	fileTypes := h.fileService.GetSupportedFileTypes()

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"fileTypes": fileTypes,
	})
}

// GetFileByID retrieves file details by ID
func (h *Handler) GetFileByID(c *gin.Context) {
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

	// Parse file ID from URL parameter
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	// Get file
	file, err := h.fileService.GetFile(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "File not found",
		})
		return
	}

	// Check if user has access to the file
	if file.UploaderID != user.UserID && !file.IsPublic && (user.OrganizationID == nil || file.OrganizationID != *user.OrganizationID) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Access denied",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"file": gin.H{
			"id":          file.ID,
			"fileName":    file.FileName,
			"contentType": file.ContentType,
			"size":        file.Size,
			"uploadedAt":  file.UploadedAt,
			"isPublic":    file.IsPublic,
			"isDeleted":   file.IsDeleted,
			"uploaderId":  file.UploaderID,
		},
	})
}

// BatchDeleteFiles handles batch deletion of files
func (h *Handler) BatchDeleteFiles(c *gin.Context) {
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

	// Parse request body
	var req struct {
		FileIDs []string `json:"fileIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request",
		})
		return
	}

	// Process each file
	var failed []string
	for _, idStr := range req.FileIDs {
		fileID, err := uuid.Parse(idStr)
		if err != nil {
			failed = append(failed, idStr)
			continue
		}

		if err := h.fileService.SoftDeleteFile(c.Request.Context(), fileID, user.UserID); err != nil {
			failed = append(failed, idStr)
		}
	}

	if len(failed) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"success":     false,
			"message":     "Some files could not be deleted",
			"failedFiles": failed,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All files deleted successfully",
	})
}

// RunCleanup handles manual cleanup of expired deleted files
func (h *Handler) RunCleanup(c *gin.Context) {
	// This endpoint should be admin-only
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}
	user := currentUser.(*user.User)

	// Check if user is an admin
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only administrators can run cleanup",
		})
		return
	}

	// Parse days parameter with default of 30 days
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	// Calculate threshold date
	threshold := time.Now().AddDate(0, 0, -days)

	// Run cleanup
	if err := h.fileService.CleanupExpiredFiles(c.Request.Context(), threshold); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Cleanup failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cleanup completed successfully",
	})
}