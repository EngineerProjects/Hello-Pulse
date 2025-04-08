package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hello-pulse.fr/internal/models/auth"
)

// Repository handles database operations for authentication sessions
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new session
func (r *Repository) Create(session *auth.Session) error {
	return r.db.Create(session).Error
}

// FindByID finds a session by ID
func (r *Repository) FindByID(id uuid.UUID) (*auth.Session, error) {
	var session auth.Session
	err := r.db.First(&session, "session_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindByToken finds a session by token
func (r *Repository) FindByToken(token string) (*auth.Session, error) {
	var session auth.Session
	err := r.db.First(&session, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindByUserID finds all sessions for a user
func (r *Repository) FindByUserID(userID uuid.UUID) ([]auth.Session, error) {
	var sessions []auth.Session
	err := r.db.Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

// Update updates a session
func (r *Repository) Update(session *auth.Session) error {
	return r.db.Save(session).Error
}

// Delete deletes a session
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&auth.Session{}, "session_id = ?", id).Error
}

// DeleteByToken deletes a session by token
func (r *Repository) DeleteByToken(token string) error {
	return r.db.Delete(&auth.Session{}, "token = ?", token).Error
}

// DeleteByUserID deletes all sessions for a user
func (r *Repository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Delete(&auth.Session{}, "user_id = ?", userID).Error
}