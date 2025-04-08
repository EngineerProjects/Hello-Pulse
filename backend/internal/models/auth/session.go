package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a user session
type Session struct {
	SessionID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	Token     string    `gorm:"unique;not null" json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// BeforeCreate is called by GORM before inserting a new record
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.SessionID == uuid.Nil {
		s.SessionID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the Session model
func (Session) TableName() string {
	return "sessions"
}