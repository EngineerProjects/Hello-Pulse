package organization

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization represents an organization in the system
type Organization struct {
	OrganizationID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrganizationName string    `gorm:"not null;unique" json:"name"`
	OwnerID          uuid.UUID `gorm:"type:uuid;not null" json:"ownerId"`
	OpenAIAPIKey     string    `gorm:"type:text" json:"-"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// BeforeCreate is called by GORM before inserting a new record
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	if o.OrganizationID == uuid.Nil {
		o.OrganizationID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the Organization model
func (Organization) TableName() string {
	return "organizations"
}