package organization

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hello-pulse.fr/internal/models/organization"
)

// Repository handles database operations for organizations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new organization repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new organization
func (r *Repository) Create(org *organization.Organization) error {
	return r.db.Create(org).Error
}

// FindByID finds an organization by ID
func (r *Repository) FindByID(id uuid.UUID) (*organization.Organization, error) {
	var org organization.Organization
	err := r.db.First(&org, "organization_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// FindByName finds an organization by name
func (r *Repository) FindByName(name string) (*organization.Organization, error) {
	var org organization.Organization
	err := r.db.First(&org, "organization_name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// Update updates an organization
func (r *Repository) Update(org *organization.Organization) error {
	return r.db.Save(org).Error
}

// Delete deletes an organization
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&organization.Organization{}, "organization_id = ?", id).Error
}

// List returns a list of organizations
func (r *Repository) List(limit, offset int) ([]organization.Organization, error) {
	var orgs []organization.Organization
	err := r.db.Limit(limit).Offset(offset).Find(&orgs).Error
	return orgs, err
}

// FindByOwner returns organizations owned by a specific user
func (r *Repository) FindByOwner(ownerID uuid.UUID) ([]organization.Organization, error) {
	var orgs []organization.Organization
	err := r.db.Where("owner_id = ?", ownerID).Find(&orgs).Error
	return orgs, err
}