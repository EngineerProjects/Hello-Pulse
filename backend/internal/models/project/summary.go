package project

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Summary represents a project summary
type Summary struct {
	SummaryID  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProjectID  uuid.UUID `gorm:"type:uuid;not null;index" json:"projectId"`
	Title      string    `gorm:"not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedBy  uuid.UUID `gorm:"type:uuid;not null" json:"createdBy"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Project    *Project  `gorm:"-"` // This is just for association, not stored in DB
}

// BeforeCreate is called by GORM before inserting a new record
func (s *Summary) BeforeCreate(tx *gorm.DB) error {
	if s.SummaryID == uuid.Nil {
		s.SummaryID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the Summary model
func (Summary) TableName() string {
	return "project_summaries"
}