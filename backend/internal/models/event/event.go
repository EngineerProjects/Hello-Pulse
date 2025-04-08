package event

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Event represents a scheduled event in the system
type Event struct {
	EventID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title          string    `gorm:"not null" json:"title"`
	Date           time.Time `gorm:"not null" json:"date"`
	StartTime      time.Time `gorm:"not null" json:"startTime"`
	EndTime        time.Time `gorm:"not null" json:"endTime"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null" json:"organizationId"`
	CreatedByID    uuid.UUID `gorm:"type:uuid;not null" json:"createdById"`
	Importance     string    `gorm:"not null;default:'not important'" json:"importance"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// BeforeCreate is called by GORM before inserting a new record
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.EventID == uuid.Nil {
		e.EventID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the Event model
func (Event) TableName() string {
	return "events"
}