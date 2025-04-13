// internal/services/file/file_service.go
package file

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/file"
	fileRepo "hello-pulse.fr/internal/repositories/file"
	"hello-pulse.fr/pkg/security"
	"hello-pulse.fr/pkg/storage"
)

// Service handles file operations
type Service struct {
	repository      *fileRepo.Repository
	storageProvider storage.Provider
	defaultBucket   string
	securityService *security.AuthorizationService
}

// NewService creates a new file service
func NewService(
	repo *fileRepo.Repository, 
	storageProvider storage.Provider, 
	defaultBucket string,
	securityService *security.AuthorizationService,
) *Service {
	return &Service{
		repository:      repo,
		storageProvider: storageProvider,
		defaultBucket:   defaultBucket,
		securityService: securityService,
	}
}

// UploadFile uploads a file and creates a database record
func (s *Service) UploadFile(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	uploaderID uuid.UUID,
	organizationID uuid.UUID,
	isPublic bool,
) (*file.File, error) {
	// Verify the user belongs to the organization
	if err := s.securityService.ValidateUserAccess(ctx, uploaderID, organizationID); err != nil {
		return nil, fmt.Errorf("unauthorized upload attempt: %w", err)
	}

	// Open the file
	f, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer f.Close()

	// Determine the file category
	category := storage.GetFileCategory(fileHeader.Filename)

	// Generate a storage path
	objectName := storage.GenerateObjectName(organizationID, category, fileHeader.Filename)

	// Get content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload file to storage provider
	_, err = s.storageProvider.UploadFile(
		ctx,
		s.defaultBucket,
		objectName,
		f,
		fileHeader.Size,
		contentType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Create database record
	fileRecord := &file.File{
		ID:             uuid.New(),
		FileName:       fileHeader.Filename,
		BucketName:     s.defaultBucket,
		ObjectName:     objectName,
		ContentType:    contentType,
		Size:           fileHeader.Size,
		UploadedAt:     time.Now(),
		UploaderID:     uploaderID,
		OrganizationID: organizationID,
		IsPublic:       isPublic,
	}

	if err := s.repository.Create(fileRecord); err != nil {
		// Try to delete the file from storage if database record creation fails
		_ = s.storageProvider.DeleteFile(ctx, s.defaultBucket, objectName)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	return fileRecord, nil
}

// GetFileURL generates a presigned URL for a file
func (s *Service) GetFileURL(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (string, error) {
	// Get file record
	fileRecord, err := s.repository.FindByID(fileID)
	if err != nil {
		return "", fmt.Errorf("file not found: %w", err)
	}

	// Check if user has access to the file
	canAccess, err := s.securityService.CanAccessFile(ctx, userID, fileID)
	if err != nil {
		return "", fmt.Errorf("error checking file access: %w", err)
	}
	
	if !canAccess {
		return "", security.ErrAccessDenied
	}

	// Generate presigned URL with 1 hour expiration
	url, err := s.storageProvider.GetFileURL(
		ctx,
		fileRecord.BucketName,
		fileRecord.ObjectName,
		time.Hour,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate file URL: %w", err)
	}

	return url, nil
}

// SoftDeleteFile marks a file as deleted
func (s *Service) SoftDeleteFile(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) error {
	// Check if user can modify the file
	canModify, err := s.securityService.CanModifyFile(ctx, userID, fileID)
	if err != nil {
		return fmt.Errorf("error checking file modification permissions: %w", err)
	}
	
	if !canModify {
		return security.ErrAccessDenied
	}

	// Get file record to check if it's already deleted
	fileRecord, err := s.repository.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check if file is already deleted
	if fileRecord.IsDeleted {
		return fmt.Errorf("file is already deleted")
	}

	// Soft delete the file
	if err := s.repository.SoftDelete(fileID); err != nil {
		return fmt.Errorf("failed to mark file as deleted: %w", err)
	}

	return nil
}

// RestoreFile restores a soft-deleted file
func (s *Service) RestoreFile(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) error {
	// Check if user can modify the file
	canModify, err := s.securityService.CanModifyFile(ctx, userID, fileID)
	if err != nil {
		return fmt.Errorf("error checking file modification permissions: %w", err)
	}
	
	if !canModify {
		return security.ErrAccessDenied
	}

	// Get file record to check if it's actually deleted
	fileRecord, err := s.repository.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check if file is actually deleted
	if !fileRecord.IsDeleted {
		return fmt.Errorf("file is not deleted")
	}

	// Restore the file
	if err := s.repository.Restore(fileID); err != nil {
		return fmt.Errorf("failed to restore file: %w", err)
	}

	return nil
}

// DeleteFilePermanently permanently deletes a file
func (s *Service) DeleteFilePermanently(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) error {
	// Check if user can modify the file
	canModify, err := s.securityService.CanModifyFile(ctx, userID, fileID)
	if err != nil {
		return fmt.Errorf("error checking file modification permissions: %w", err)
	}
	
	if !canModify {
		return security.ErrAccessDenied
	}

	// Get file record
	fileRecord, err := s.repository.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Delete the file from storage
	if err := s.storageProvider.DeleteFile(ctx, fileRecord.BucketName, fileRecord.ObjectName); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete the file record from database
	if err := s.repository.DeletePermanently(fileID); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

// GetUserFiles gets all files uploaded by a user
func (s *Service) GetUserFiles(ctx context.Context, userID uuid.UUID, includeDeleted bool) ([]file.File, error) {
	// Get user's organization ID
	orgID, err := s.securityService.GetUserOrganizationID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user's organization: %w", err)
	}

	// Only return files from the user's organization
	return s.repository.FindByUploaderAndOrg(userID, orgID, includeDeleted)
}

// GetOrganizationFiles gets all files for an organization
func (s *Service) GetOrganizationFiles(ctx context.Context, userID, orgID uuid.UUID, includeDeleted bool) ([]file.File, error) {
	// Verify the user belongs to the organization
	if err := s.securityService.ValidateUserAccess(ctx, userID, orgID); err != nil {
		return nil, fmt.Errorf("unauthorized access attempt: %w", err)
	}

	return s.repository.FindByOrganization(orgID, includeDeleted)
}

// GetFile gets a file by ID
func (s *Service) GetFile(ctx context.Context, fileID, userID uuid.UUID) (*file.File, error) {
	// Get file record
	fileRecord, err := s.repository.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Check if user has access to the file
	canAccess, err := s.securityService.CanAccessFile(ctx, userID, fileID)
	if err != nil {
		return nil, fmt.Errorf("error checking file access: %w", err)
	}
	
	if !canAccess {
		return nil, security.ErrAccessDenied
	}

	return fileRecord, nil
}

// UpdateFileVisibility updates a file's public/private status
func (s *Service) UpdateFileVisibility(ctx context.Context, fileID uuid.UUID, userID uuid.UUID, isPublic bool) error {
	// Check if user can modify the file
	canModify, err := s.securityService.CanModifyFile(ctx, userID, fileID)
	if err != nil {
		return fmt.Errorf("error checking file modification permissions: %w", err)
	}
	
	if !canModify {
		return security.ErrAccessDenied
	}

	// Get file record
	fileRecord, err := s.repository.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Update the file visibility
	fileRecord.IsPublic = isPublic

	// Update the file record
	if err := s.repository.Update(fileRecord); err != nil {
		return fmt.Errorf("failed to update file record: %w", err)
	}

	return nil
}

// CleanupExpiredFiles permanently deletes files that were soft-deleted before a threshold
func (s *Service) CleanupExpiredFiles(ctx context.Context, threshold time.Time) error {
	// Find files to delete
	filesToDelete, err := s.repository.FindExpiredDeleted(threshold)
	if err != nil {
		return fmt.Errorf("failed to find expired files: %w", err)
	}

	var errors []string

	// Delete each file
	for _, fileRecord := range filesToDelete {
		// Delete from storage
		if err := s.storageProvider.DeleteFile(ctx, fileRecord.BucketName, fileRecord.ObjectName); err != nil {
			errors = append(errors, fmt.Sprintf("failed to delete file %s from storage: %v", fileRecord.ObjectName, err))
			continue
		}

		// Delete from database
		if err := s.repository.DeletePermanently(fileRecord.ID); err != nil {
			errors = append(errors, fmt.Sprintf("failed to delete file %s from database: %v", fileRecord.ID, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during cleanup: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetSupportedFileTypes returns a map of supported file types
func (s *Service) GetSupportedFileTypes() map[string][]string {
	return storage.GetSupportedFileTypes()
}

// DownloadFile downloads a file from storage
func (s *Service) DownloadFile(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (io.ReadCloser, string, error) {
	// Get file record
	fileRecord, err := s.repository.FindByID(fileID)
	if err != nil {
		return nil, "", fmt.Errorf("file not found: %w", err)
	}

	// Check if user has access to the file
	canAccess, err := s.securityService.CanAccessFile(ctx, userID, fileID)
	if err != nil {
		return nil, "", fmt.Errorf("error checking file access: %w", err)
	}
	
	if !canAccess {
		return nil, "", security.ErrAccessDenied
	}

	// Check if file is deleted
	if fileRecord.IsDeleted {
		return nil, "", fmt.Errorf("file is deleted")
	}

	// Download file from storage
	reader, err := s.storageProvider.DownloadFile(
		ctx,
		fileRecord.BucketName,
		fileRecord.ObjectName,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download file: %w", err)
	}

	return reader, fileRecord.ContentType, nil
}

// BatchSoftDeleteFiles marks multiple files as deleted
func (s *Service) BatchSoftDeleteFiles(ctx context.Context, fileIDs []uuid.UUID, userID uuid.UUID) ([]uuid.UUID, error) {
	var failedFiles []uuid.UUID

	for _, fileID := range fileIDs {
		if err := s.SoftDeleteFile(ctx, fileID, userID); err != nil {
			failedFiles = append(failedFiles, fileID)
		}
	}

	return failedFiles, nil
}