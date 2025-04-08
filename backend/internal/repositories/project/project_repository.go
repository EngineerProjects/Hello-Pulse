package project

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hello-pulse.fr/internal/models/project"
	"hello-pulse.fr/internal/models/user"
)

// Repository handles database operations for projects
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new project repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new project
func (r *Repository) Create(project *project.Project) error {
	return r.db.Create(project).Error
}

// FindByID finds a project by ID
func (r *Repository) FindByID(id uuid.UUID) (*project.Project, error) {
	var project project.Project
	err := r.db.First(&project, "project_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// Update updates a project
func (r *Repository) Update(project *project.Project) error {
	return r.db.Save(project).Error
}

// Delete deletes a project
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&project.Project{}, "project_id = ?", id).Error
}

// FindByOrganization returns projects for a specific organization
func (r *Repository) FindByOrganization(orgID uuid.UUID) ([]project.Project, error) {
	var projects []project.Project
	err := r.db.Where("organization_id = ?", orgID).Find(&projects).Error
	return projects, err
}

// FindRootByOrganization returns root projects (with no parent) for a specific organization
func (r *Repository) FindRootByOrganization(orgID uuid.UUID) ([]project.Project, error) {
	var projects []project.Project
	err := r.db.Where("organization_id = ? AND parent_project_id IS NULL", orgID).Find(&projects).Error
	return projects, err
}

// FindByParent returns projects that have a specific parent
func (r *Repository) FindByParent(parentID uuid.UUID) ([]project.Project, error) {
	var projects []project.Project
	err := r.db.Where("parent_project_id = ?", parentID).Find(&projects).Error
	return projects, err
}

// AddParticipant adds a user as a participant to a project
func (r *Repository) AddParticipant(projectID, userID uuid.UUID) error {
	return r.db.Exec("INSERT INTO project_participants (project_id, user_id) VALUES (?, ?)", projectID, userID).Error
}

// RemoveParticipant removes a user as a participant from a project
func (r *Repository) RemoveParticipant(projectID, userID uuid.UUID) error {
	return r.db.Exec("DELETE FROM project_participants WHERE project_id = ? AND user_id = ?", projectID, userID).Error
}

// ClearParticipants removes all participants from a project
func (r *Repository) ClearParticipants(projectID uuid.UUID) error {
	return r.db.Exec("DELETE FROM project_participants WHERE project_id = ?", projectID).Error
}

// GetParticipants gets all participants of a project
func (r *Repository) GetParticipants(projectID uuid.UUID) ([]user.User, error) {
	var users []user.User
	err := r.db.Joins("JOIN project_participants ON users.user_id = project_participants.user_id").
		Where("project_participants.project_id = ?", projectID).
		Find(&users).Error
	return users, err
}