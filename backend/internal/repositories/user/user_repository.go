package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hello-pulse.fr/internal/models/user"
)

// Repository handles database operations for users
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new user
func (r *Repository) Create(user *user.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *Repository) FindByID(id uuid.UUID) (*user.User, error) {
	var user user.User
	err := r.db.First(&user, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *Repository) FindByEmail(email string) (*user.User, error) {
	var user user.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *Repository) Update(user *user.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&user.User{}, "user_id = ?", id).Error
}

// List returns a list of users
func (r *Repository) List(limit, offset int) ([]user.User, error) {
	var users []user.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// FindByOrganization returns users in a specific organization
func (r *Repository) FindByOrganization(orgID uuid.UUID) ([]user.User, error) {
	var users []user.User
	err := r.db.Where("organization_id = ?", orgID).Find(&users).Error
	return users, err
}