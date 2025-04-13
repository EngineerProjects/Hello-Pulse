package project

import (
	"errors"
	"time"

	"github.com/google/uuid"
	
	"hello-pulse.fr/internal/models/project"
	projectrepo "hello-pulse.fr/internal/repositories/project"
)

// SummaryService handles project summary business logic
type SummaryService struct {
	summaryRepo *projectrepo.SummaryRepository
	projectRepo *projectrepo.Repository
}

// NewSummaryService creates a new summary service
func NewSummaryService(summaryRepo *projectrepo.SummaryRepository, projectRepo *projectrepo.Repository) *SummaryService {
	return &SummaryService{
		summaryRepo: summaryRepo,
		projectRepo: projectRepo,
	}
}

// CreateSummary creates a new project summary
func (s *SummaryService) CreateSummary(
	title string,
	content string,
	projectID uuid.UUID,
	createdBy uuid.UUID,
) (*project.Summary, error) {
	// Verify project exists
	proj, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	// Create summary
	summary := &project.Summary{
		Title:     title,
		Content:   content,
		ProjectID: proj.ProjectID,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.summaryRepo.Create(summary); err != nil {
		return nil, err
	}

	return summary, nil
}

// GetSummary retrieves a summary by ID
func (s *SummaryService) GetSummary(summaryID uuid.UUID) (*project.Summary, error) {
	return s.summaryRepo.FindByID(summaryID)
}

// GetProjectSummaries retrieves all summaries for a project
func (s *SummaryService) GetProjectSummaries(projectID uuid.UUID) ([]project.Summary, error) {
	return s.summaryRepo.FindByProject(projectID)
}

// UpdateSummary updates a summary's details
func (s *SummaryService) UpdateSummary(summaryID uuid.UUID, title string, content string, userID uuid.UUID) error {
	summary, err := s.summaryRepo.FindByID(summaryID)
	if err != nil {
		return errors.New("summary not found")
	}

	// Check if user is the creator of the summary
	if summary.CreatedBy != userID {
		return errors.New("only the creator can update this summary")
	}

	summary.Title = title
	summary.Content = content
	summary.UpdatedAt = time.Now()

	return s.summaryRepo.Update(summary)
}

// DeleteSummary deletes a project summary
func (s *SummaryService) DeleteSummary(summaryID uuid.UUID, userID uuid.UUID) error {
	summary, err := s.summaryRepo.FindByID(summaryID)
	if err != nil {
		return errors.New("summary not found")
	}

	// Check if user is the creator of the summary
	if summary.CreatedBy != userID {
		return errors.New("only the creator can delete this summary")
	}

	return s.summaryRepo.Delete(summaryID)
}

// DeleteProjectSummaries deletes all summaries for a project
func (s *SummaryService) DeleteProjectSummaries(projectID uuid.UUID) error {
	return s.summaryRepo.DeleteByProject(projectID)
}