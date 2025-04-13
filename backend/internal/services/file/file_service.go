package file

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/file"
	fileRepo "hello-pulse.fr/internal/repositories/file"
	"hello-pulse.fr/pkg/storage"
)

// Service handles file operations
type Service struct {
	fileRepo    *fileRepo.Repository
	minioClient *storage.MinioClient
	bucketName  string
}

// NewService creates a new file service
func NewService(repo *fileRepo.Repository, minioClient *storage.MinioClient, bucketName string) *Service {
	return &Service{
		fileRepo:    repo,
		minioClient: minioClient,
		bucketName:  bucketName,
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
	// Open the file
	f, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer f.Close()

	// Determine the folder based on file type
	folder := storage.GetFolderForFileType(fileHeader.Filename)

	// Generate a unique object name
	objectName := storage.GenerateObjectName(organizationID, folder, fileHeader.Filename)

	// Get content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload file to MinIO
	if _, err := s.minioClient.UploadFile(
		ctx,
		s.bucketName,
		objectName,
		f,
		fileHeader.Size,
		contentType,
	); err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Create database record
	fileRecord := &file.File{
		FileName:       fileHeader.Filename,
		BucketName:     s.bucketName,
		ObjectName:     objectName,
		ContentType:    contentType,
		Size:           fileHeader.Size,
		UploadedAt:     time.Now(),
		UploaderID:     uploaderID,
		OrganizationID: organizationID,
		IsPublic:       isPublic,
	}

	if err := s.fileRepo.Create(fileRecord); err != nil {
		// Try to delete the file from storage if database record creation fails
		_ = s.minioClient.DeleteFile(ctx, s.bucketName, objectName)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	return fileRecord, nil
}

// GetFileURL generates a presigned URL for a file
func (s *Service) GetFileURL(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (string, error) {
	// Get file record
	fileRecord, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return "", fmt.Errorf("file not found: %w", err)
	}

	// Check if user has access to the file
	if fileRecord.UploaderID != userID && !fileRecord.IsPublic {
		return "", errors.New("access denied: you don't have permission to access this file")
	}

	// Check if file is deleted
	if fileRecord.IsDeleted {
		return "", errors.New("file is deleted")
	}

	// Generate presigned URL with 1 hour expiration
	url, err := s.minioClient.GetFileURL(
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
	// Get file record
	fileRecord, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check if user has permission to delete the file
	if fileRecord.UploaderID != userID {
		return errors.New("access denied: only the uploader can delete this file")
	}

	// Check if file is already deleted
	if fileRecord.IsDeleted {
		return errors.New("file is already deleted")
	}

	// Soft delete the file
	if err := s.fileRepo.SoftDelete(fileID); err != nil {
		return fmt.Errorf("failed to mark file as deleted: %w", err)
	}

	return nil
}

// RestoreFile restores a soft-deleted file
func (s *Service) RestoreFile(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) error {
	// Get file record
	fileRecord, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check if user has permission to restore the file
	if fileRecord.UploaderID != userID {
		return errors.New("access denied: only the uploader can restore this file")
	}

	// Check if file is actually deleted
	if !fileRecord.IsDeleted {
		return errors.New("file is not deleted")
	}

	// Restore the file
	if err := s.fileRepo.Restore(fileID); err != nil {
		return fmt.Errorf("failed to restore file: %w", err)
	}

	return nil
}

// DeleteFilePermanently permanently deletes a file
func (s *Service) DeleteFilePermanently(ctx context.Context, fileID uuid.UUID) error {
	// Get file record
	fileRecord, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Delete the file from storage
	if err := s.minioClient.DeleteFile(ctx, fileRecord.BucketName, fileRecord.ObjectName); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete the file record from database
	if err := s.fileRepo.DeletePermanently(fileID); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

// GetUserFiles gets all files uploaded by a user
func (s *Service) GetUserFiles(userID uuid.UUID, includeDeleted bool) ([]file.File, error) {
	return s.fileRepo.FindByUploader(userID, includeDeleted)
}

// GetOrganizationFiles gets all files for an organization
func (s *Service) GetOrganizationFiles(orgID uuid.UUID, includeDeleted bool) ([]file.File, error) {
	return s.fileRepo.FindByOrganization(orgID, includeDeleted)
}

// GetFile gets a file by ID
func (s *Service) GetFile(fileID uuid.UUID) (*file.File, error) {
	return s.fileRepo.FindByID(fileID)
}

// UpdateFileVisibility updates a file's public/private status
func (s *Service) UpdateFileVisibility(ctx context.Context, fileID uuid.UUID, userID uuid.UUID, isPublic bool) error {
	// Get file record
	fileRecord, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check if user has permission to update the file
	if fileRecord.UploaderID != userID {
		return errors.New("access denied: only the uploader can update this file")
	}

	// Update the file visibility
	fileRecord.IsPublic = isPublic
	if err := s.fileRepo.Update(fileRecord); err != nil {
		return fmt.Errorf("failed to update file visibility: %w", err)
	}

	return nil
}

// CleanupExpiredFiles permanently deletes files that were soft-deleted before a threshold
func (s *Service) CleanupExpiredFiles(ctx context.Context, threshold time.Time) error {
	// Find files to delete
	files, err := s.fileRepo.FindExpiredDeleted(threshold)
	if err != nil {
		return fmt.Errorf("failed to find expired files: %w", err)
	}

	// Delete each file
	for _, f := range files {
		// Delete from storage
		if err := s.minioClient.DeleteFile(ctx, f.BucketName, f.ObjectName); err != nil {
			fmt.Printf("Warning: failed to delete file %s from storage: %v\n", f.ObjectName, err)
			continue
		}

		// Delete from database
		if err := s.fileRepo.DeletePermanently(f.ID); err != nil {
			fmt.Printf("Warning: failed to delete file %s from database: %v\n", f.ID, err)
		}
	}

	return nil
}

// GetSupportedFileTypes returns a map of supported file types
func (s *Service) GetSupportedFileTypes() map[string][]string {
	return file.GetSupportedFileTypes()
}

// GetUserAccessibleFiles returns files that a user can access (owned or public within their org)
func (s *Service) GetUserAccessibleFiles(userID, orgID uuid.UUID) ([]file.File, error) {
	return s.fileRepo.GetUserAccessibleFiles(userID, orgID)
}