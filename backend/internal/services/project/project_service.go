package project

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"hello-pulse.fr/internal/models/project"
	"hello-pulse.fr/internal/models/user"
	projectrepo "hello-pulse.fr/internal/repositories/project"
	userrepo "hello-pulse.fr/internal/repositories/user"
)

// Service handles project business logic
type Service struct {
	projectRepo *projectrepo.Repository
	userRepo    *userrepo.Repository
	summaryRepo *projectrepo.SummaryRepository
}

// NewService creates a new project service
func NewService(projectRepo *projectrepo.Repository, userRepo *userrepo.Repository, summaryRepo *projectrepo.SummaryRepository) *Service {
	return &Service{
		projectRepo: projectRepo,
		userRepo:    userRepo,
		summaryRepo: summaryRepo,
	}
}

// CreateProject creates a new project
func (s *Service) CreateProject(name, description string, ownerID, orgID uuid.UUID, parentProjectID *uuid.UUID) (*project.Project, error) {
	// Validate owner and organization
	owner, err := s.userRepo.FindByID(ownerID)
	if err != nil {
		return nil, errors.New("owner not found")
	}

	if owner.OrganizationID == nil || *owner.OrganizationID != orgID {
		return nil, errors.New("owner does not belong to the specified organization")
	}

	// Create project
	project := &project.Project{
		ProjectName:     name,
		ProjectDesc:     description,
		OwnerID:         ownerID,
		OrganizationID:  orgID,
		ParentProjectID: parentProjectID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	// Add owner as a participant
	if err := s.projectRepo.AddParticipant(project.ProjectID, owner.UserID); err != nil {
		return nil, err
	}

	return project, nil
}

// GetProject retrieves a project by ID
func (s *Service) GetProject(projectID uuid.UUID) (*project.Project, error) {
	return s.projectRepo.FindByID(projectID)
}

// GetOrganizationProjects retrieves all projects for an organization
func (s *Service) GetOrganizationProjects(orgID uuid.UUID) ([]project.Project, error) {
	return s.projectRepo.FindByOrganization(orgID)
}

// GetRootProjects retrieves all root projects (with no parent) for an organization
func (s *Service) GetRootProjects(orgID uuid.UUID) ([]project.Project, error) {
	return s.projectRepo.FindRootByOrganization(orgID)
}

// GetChildProjects retrieves all child projects for a given parent project
func (s *Service) GetChildProjects(parentID uuid.UUID) ([]project.Project, error) {
	return s.projectRepo.FindByParent(parentID)
}

// UpdateProject updates a project's details
func (s *Service) UpdateProject(projectID uuid.UUID, name, description string) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return err
	}

	project.ProjectName = name
	project.ProjectDesc = description
	project.UpdatedAt = time.Now()

	return s.projectRepo.Update(project)
}

// DeleteProject deletes a project and all its children
func (s *Service) DeleteProject(projectID uuid.UUID) error {
	// Get all child projects
	children, err := s.projectRepo.FindByParent(projectID)
	if err != nil {
		return err
	}

	// Recursively delete children
	for _, child := range children {
		if err := s.DeleteProject(child.ProjectID); err != nil {
			return err
		}
	}

	// Delete all project summaries
	if err := s.summaryRepo.DeleteByProject(projectID); err != nil {
		return err
	}

	// Clear project participants
	if err := s.projectRepo.ClearParticipants(projectID); err != nil {
		return err
	}

	// Delete the project
	return s.projectRepo.Delete(projectID)
}

// AddParticipant adds a user to a project
func (s *Service) AddParticipant(projectID, userID uuid.UUID) error {
	// Verify project exists
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify user and project belong to the same organization
	if user.OrganizationID == nil || *user.OrganizationID != project.OrganizationID {
		return errors.New("user does not belong to the project's organization")
	}

	// Add participant
	return s.projectRepo.AddParticipant(projectID, userID)
}

// RemoveParticipant removes a user from a project
func (s *Service) RemoveParticipant(projectID, userID uuid.UUID) error {
	// Cannot remove the owner
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	if project.OwnerID == userID {
		return errors.New("cannot remove project owner from participants")
	}

	return s.projectRepo.RemoveParticipant(projectID, userID)
}

// GetProjectParticipants retrieves all participants of a project
func (s *Service) GetProjectParticipants(projectID uuid.UUID) ([]user.User, error) {
	return s.projectRepo.GetParticipants(projectID)
}