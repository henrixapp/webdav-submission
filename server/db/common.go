package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//BaseObject is the minimum Object for GORM, containing timestamps and an UUID
type BaseObject struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *BaseObject) BeforeCreate(scope *gorm.DB) error {
	uuid := uuid.New()
	base.ID = uuid
	return nil
}
