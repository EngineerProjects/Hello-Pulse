package project

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project represents a project in the system
type Project struct {
	ProjectID       uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProjectName     string     `gorm:"not null" json:"name"`
	ProjectDesc     string     `json:"description"`
	OwnerID         uuid.UUID  `gorm:"type:uuid;not null" json:"ownerId"`
	OrganizationID  uuid.UUID  `gorm:"type:uuid;not null" json:"organizationId"`
	ParentProjectID *uuid.UUID `gorm:"type:uuid" json:"parentProjectId"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// BeforeCreate is called by GORM before inserting a new record
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ProjectID == uuid.Nil {
		p.ProjectID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the Project model
func (Project) TableName() string {
	return "projects"
}