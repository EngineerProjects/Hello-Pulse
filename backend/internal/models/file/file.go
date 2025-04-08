package file

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// File represents a file stored in the system
type File struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FileName       string     `gorm:"type:varchar(255);not null" json:"fileName"`
	BucketName     string     `gorm:"type:varchar(255);not null" json:"bucketName"`
	ObjectName     string     `gorm:"type:varchar(255);not null" json:"objectName"`
	ContentType    string     `gorm:"type:varchar(50)" json:"contentType"`
	UploadedAt     time.Time  `gorm:"autoCreateTime" json:"uploadedAt"`
	IsDeleted      bool       `gorm:"default:false" json:"isDeleted"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty"`
	UploaderID     uuid.UUID  `gorm:"type:uuid;not null" json:"uploaderId"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null" json:"organizationId"`
	IsPublic       bool       `gorm:"default:false" json:"isPublic"`
}

// BeforeCreate is called by GORM before inserting a new record
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the File model
func (File) TableName() string {
	return "files"
}