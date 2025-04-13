package project

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hello-pulse.fr/internal/models/project"
)

// SummaryRepository handles database operations for project summaries
type SummaryRepository struct {
	db *gorm.DB
}

// NewSummaryRepository creates a new summary repository
func NewSummaryRepository(db *gorm.DB) *SummaryRepository {
	return &SummaryRepository{db: db}
}

// Create inserts a new project summary
func (r *SummaryRepository) Create(summary *project.Summary) error {
	return r.db.Create(summary).Error
}

// FindByID finds a summary by ID
func (r *SummaryRepository) FindByID(id uuid.UUID) (*project.Summary, error) {
	var summary project.Summary
	err := r.db.First(&summary, "summary_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// Update updates a summary
func (r *SummaryRepository) Update(summary *project.Summary) error {
	return r.db.Save(summary).Error
}

// Delete deletes a summary
func (r *SummaryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&project.Summary{}, "summary_id = ?", id).Error
}

// FindByProject returns all summaries for a specific project
func (r *SummaryRepository) FindByProject(projectID uuid.UUID) ([]project.Summary, error) {
	var summaries []project.Summary
	err := r.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&summaries).Error
	return summaries, err
}

// DeleteByProject deletes all summaries for a project
func (r *SummaryRepository) DeleteByProject(projectID uuid.UUID) error {
	return r.db.Delete(&project.Summary{}, "project_id = ?", projectID).Error
}