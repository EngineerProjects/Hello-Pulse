package invite

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InviteCode represents an invitation code for joining an organization
type InviteCode struct {
	InviteCodeID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Value          string    `gorm:"not null;unique" json:"value"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null" json:"organizationId"`
	ExpirationTime time.Time `gorm:"not null" json:"expirationTime"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// BeforeCreate is called by GORM before inserting a new record
func (i *InviteCode) BeforeCreate(tx *gorm.DB) error {
	if i.InviteCodeID == uuid.Nil {
		i.InviteCodeID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the InviteCode model
func (InviteCode) TableName() string {
	return "invite_codes"
}