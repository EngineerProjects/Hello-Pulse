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
	return &Repository{db: db}
}

// Create inserts a new file
func (r *Repository) Create(file *file.File) error {
	return r.db.Create(file).Error
}

// FindByID finds a file by ID
func (r *Repository) FindByID(id uuid.UUID) (*file.File, error) {
	var file file.File
	err := r.db.First(&file, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// Update updates a file
func (r *Repository) Update(file *file.File) error {
	return r.db.Save(file).Error
}

// Delete marks a file as deleted
func (r *Repository) SoftDelete(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&file.File{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": now,
	}).Error
}

// PermanentDelete permanently deletes a file
func (r *Repository) PermanentDelete(id uuid.UUID) error {
	return r.db.Delete(&file.File{}, "id = ?", id).Error
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

// FindExpiredDeleted returns files that were soft deleted before a given time
func (r *Repository) FindExpiredDeleted(threshold time.Time) ([]file.File, error) {
	var files []file.File
	err := r.db.Where("is_deleted = ? AND deleted_at <= ?", true, threshold).Find(&files).Error
	return files, err
}