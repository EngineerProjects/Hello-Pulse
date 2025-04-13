package file

import (
	"fmt"
	"net/http"
	
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

	// Generate a unique object name
	objectName := fmt.Sprintf("%s/%s", user.OrganizationID.String(), uuid.New().String())
	
	// Determine bucket name based on file extension
	bucketName := "hello-pulse"
	
	// Get content type from file
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	
	// Upload file record to database
	fileRecord, err := h.fileService.UploadFile(
		file.Filename,
		bucketName,
		objectName,
		contentType,
		user.UserID,
		*user.OrganizationID,
		isPublic,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to save file metadata",
		})
		return
	}

	// In a real implementation, you would upload the file to storage here
	// For now, we'll just return the file metadata
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"file": gin.H{
			"id":          fileRecord.ID,
			"fileName":    fileRecord.FileName,
			"contentType": fileRecord.ContentType,
			"uploadedAt":  fileRecord.UploadedAt,
			"isPublic":    fileRecord.IsPublic,
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

	// Get file ID from URL parameter
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	// Soft delete file
	if err := h.fileService.SoftDeleteFile(fileID, user.UserID); err != nil {
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

	// Get file ID from URL parameter
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	// Restore file
	if err := h.fileService.RestoreFile(fileID, user.UserID); err != nil {
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
	// Define supported file types
	fileTypes := map[string][]string{
		"documents": {".pdf", ".doc", ".docx", ".txt", ".rtf", ".odt", ".md"},
		"images":    {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"},
		"audio":     {".mp3", ".wav", ".ogg", ".flac", ".m4a"},
		"video":     {".mp4", ".mov", ".avi", ".mkv", ".webm"},
		"archives":  {".zip", ".rar", ".7z", ".tar", ".gz"},
		"spreadsheets": {".xls", ".xlsx", ".csv", ".ods"},
		"presentations": {".ppt", ".pptx", ".odp"},
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"fileTypes": fileTypes,
	})
}

// GetFileByID retrieves a file by ID
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

	// Get file ID from URL parameter
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
	if file.OrganizationID != *user.OrganizationID && !file.IsPublic && file.UploaderID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "You don't have permission to access this file",
		})
		return
	}

	// Format response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"file": gin.H{
			"id":          file.ID,
			"fileName":    file.FileName,
			"contentType": file.ContentType,
			"uploadedAt":  file.UploadedAt,
			"isPublic":    file.IsPublic,
			"isDeleted":   file.IsDeleted,
			"uploaderId":  file.UploaderID,
			"bucketName":  file.BucketName,
			"objectName":  file.ObjectName,
		},
	})
}