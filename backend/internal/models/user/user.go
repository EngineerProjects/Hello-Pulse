package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	UserID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `gorm:"unique" json:"email"`
	PasswordHash   string    `json:"-"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	OrganizationID *uuid.UUID `json:"organizationId"`
	Role           string    `json:"role"`
	LastActive     time.Time `json:"lastActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// BeforeCreate is called by GORM before inserting a new record
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}