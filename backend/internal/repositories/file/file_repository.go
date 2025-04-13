// internal/repositories/file/file_repository.go
package file

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"hello-pulse.fr/internal/models/file"
)

// Repository handles database operations for files
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new file repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create inserts a new file
func (r *Repository) Create(file *file.File) error {
	return r.db.Create(file).Error
}

// FindByID finds a file by ID
func (r *Repository) FindByID(id uuid.UUID) (*file.File, error) {
	var f file.File
	err := r.db.First(&f, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// Update updates a file
func (r *Repository) Update(file *file.File) error {
	return r.db.Save(file).Error
}

// DeletePermanently permanently deletes a file
func (r *Repository) DeletePermanently(id uuid.UUID) error {
	return r.db.Delete(&file.File{}, "id = ?", id).Error
}

// SoftDelete marks a file as deleted
func (r *Repository) SoftDelete(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&file.File{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": now,
	}).Error
}

// Restore restores a soft-deleted file
func (r *Repository) Restore(id uuid.UUID) error {
	return r.db.Model(&file.File{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": false,
		"deleted_at": nil,
	}).Error
}

// FindByOrganization returns files for a specific organization
func (r *Repository) FindByOrganization(orgID uuid.UUID, includeDeleted bool) ([]file.File, error) {
	var files []file.File
	query := r.db.Where("organization_id = ?", orgID)
	
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	
	err := query.Find(&files).Error
	return files, err
}

// FindByUploader returns files uploaded by a specific user
func (r *Repository) FindByUploader(uploaderID uuid.UUID, includeDeleted bool) ([]file.File, error) {
	var files []file.File
	query := r.db.Where("uploader_id = ?", uploaderID)
	
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	
	err := query.Find(&files).Error
	return files, err
}

// FindByUploaderAndOrg returns files uploaded by a specific user in a specific organization
func (r *Repository) FindByUploaderAndOrg(uploaderID, orgID uuid.UUID, includeDeleted bool) ([]file.File, error) {
	var files []file.File
	query := r.db.Where("uploader_id = ? AND organization_id = ?", uploaderID, orgID)
	
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	
	err := query.Find(&files).Error
	return files, err
}

// FindExpiredDeleted returns files that were soft deleted before a given time
func (r *Repository) FindExpiredDeleted(threshold time.Time) ([]file.File, error) {
	var files []file.File
	err := r.db.Where("is_deleted = ? AND deleted_at <= ?", true, threshold).Find(&files).Error
	return files, err
}

// GetUserAccessibleFiles returns files that a user can access (owned or public within their org)
func (r *Repository) GetUserAccessibleFiles(userID, orgID uuid.UUID) ([]file.File, error) {
	var files []file.File
	err := r.db.Where(
		"(organization_id = ? AND is_deleted = ? AND (uploader_id = ? OR is_public = ?))",
		orgID, false, userID, true).
		Find(&files).Error
	return files, err
}

// CountFiles returns the count of files based on the provided conditions
func (r *Repository) CountFiles(conditions map[string]interface{}) (int64, error) {
	var count int64
	err := r.db.Model(&file.File{}).Where(conditions).Count(&count).Error
	return count, err
}

// CountFilesByOrganization returns the count of files in an organization
func (r *Repository) CountFilesByOrganization(orgID uuid.UUID, includeDeleted bool) (int64, error) {
	conditions := map[string]interface{}{
		"organization_id": orgID,
	}
	
	if !includeDeleted {
		conditions["is_deleted"] = false
	}
	
	return r.CountFiles(conditions)
}

// CountFilesByUploader returns the count of files uploaded by a user
func (r *Repository) CountFilesByUploader(uploaderID uuid.UUID, includeDeleted bool) (int64, error) {
	conditions := map[string]interface{}{
		"uploader_id": uploaderID,
	}
	
	if !includeDeleted {
		conditions["is_deleted"] = false
	}
	
	return r.CountFiles(conditions)
}

// GetTotalFileSizeByOrganization returns the total size of files in an organization
func (r *Repository) GetTotalFileSizeByOrganization(orgID uuid.UUID, includeDeleted bool) (int64, error) {
	var totalSize int64
	
	query := r.db.Model(&file.File{}).
		Select("SUM(size) as total_size").
		Where("organization_id = ?", orgID)
	
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	
	err := query.Scan(&totalSize).Error
	return totalSize, err
}

// GetFilesByContentType returns files of a specific content type in an organization
func (r *Repository) GetFilesByContentType(orgID uuid.UUID, contentType string, includeDeleted bool) ([]file.File, error) {
	var files []file.File
	
	query := r.db.Where("organization_id = ? AND content_type LIKE ?", orgID, contentType+"%")
	
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	
	err := query.Find(&files).Error
	return files, err
}