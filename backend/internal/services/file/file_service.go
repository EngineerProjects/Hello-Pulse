package file

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/file"
	filerepo "hello-pulse.fr/internal/repositories/file"
)

// Service handles file operations
type Service struct {
	fileRepo *filerepo.Repository
}

// NewService creates a new file service
func NewService(fileRepo *filerepo.Repository) *Service {
	return &Service{
		fileRepo: fileRepo,
	}
}

// UploadFile uploads a file to storage and creates a database record
func (s *Service) UploadFile(
	fileName string,
	bucketName string,
	objectName string,
	contentType string,
	uploaderID uuid.UUID,
	organizationID uuid.UUID,
	isPublic bool,
) (*file.File, error) {
	// Create a new file record
	newFile := &file.File{
		FileName:       fileName,
		BucketName:     bucketName,
		ObjectName:     objectName,
		ContentType:    contentType,
		UploadedAt:     time.Now(),
		UploaderID:     uploaderID,
		OrganizationID: organizationID,
		IsPublic:       isPublic,
	}

	// Save file record to database
	if err := s.fileRepo.Create(newFile); err != nil {
		return nil, err
	}

	return newFile, nil
}

// GetFile retrieves a file by ID
func (s *Service) GetFile(fileID uuid.UUID) (*file.File, error) {
	return s.fileRepo.FindByID(fileID)
}

// SoftDeleteFile marks a file as deleted
func (s *Service) SoftDeleteFile(fileID, userID uuid.UUID) error {
	// Get file
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return errors.New("file not found")
	}

	// Check if user is the uploader
	if file.UploaderID != userID {
		return errors.New("you don't have permission to delete this file")
	}

	// Soft delete the file
	return s.fileRepo.SoftDelete(fileID)
}

// GetOrganizationFiles retrieves all files for an organization
func (s *Service) GetOrganizationFiles(orgID uuid.UUID, includeDeleted bool) ([]file.File, error) {
	return s.fileRepo.FindByOrganization(orgID, includeDeleted)
}

// GetUserFiles retrieves all files uploaded by a user
func (s *Service) GetUserFiles(userID uuid.UUID, includeDeleted bool) ([]file.File, error) {
	return s.fileRepo.FindByUploader(userID, includeDeleted)
}

// RestoreFile restores a soft-deleted file
func (s *Service) RestoreFile(fileID, userID uuid.UUID) error {
	// Get file
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return errors.New("file not found")
	}

	// Check if user is the uploader
	if file.UploaderID != userID {
		return errors.New("you don't have permission to restore this file")
	}

	// Check if file is deleted
	if !file.IsDeleted {
		return errors.New("file is not deleted")
	}

	// Restore file
	file.IsDeleted = false
	file.DeletedAt = nil
	return s.fileRepo.Update(file)
}