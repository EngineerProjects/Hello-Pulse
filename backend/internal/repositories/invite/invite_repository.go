package invite

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hello-pulse.fr/internal/models/invite"
)

// Repository handles database operations for invite codes
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new invite repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new invite code
func (r *Repository) Create(invite *invite.InviteCode) error {
	return r.db.Create(invite).Error
}

// FindByID finds an invite code by ID
func (r *Repository) FindByID(id uuid.UUID) (*invite.InviteCode, error) {
	var invite invite.InviteCode
	err := r.db.First(&invite, "invite_code_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

// FindByCode finds an invite code by value
func (r *Repository) FindByCode(code string) (*invite.InviteCode, error) {
	var invite invite.InviteCode
	err := r.db.First(&invite, "value = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

// FindByOrganization finds all invite codes for an organization
func (r *Repository) FindByOrganization(orgID uuid.UUID) ([]invite.InviteCode, error) {
	var invites []invite.InviteCode
	err := r.db.Where("organization_id = ?", orgID).Find(&invites).Error
	return invites, err
}

// Update updates an invite code
func (r *Repository) Update(invite *invite.InviteCode) error {
	return r.db.Save(invite).Error
}

// Delete deletes an invite code
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&invite.InviteCode{}, "invite_code_id = ?", id).Error
}