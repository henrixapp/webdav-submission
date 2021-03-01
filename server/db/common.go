package db

import (
	"time"

	"github.com/google/uuid"
)

//BaseObject is the minimum Object for GORM, containing timestamps and an UUID
type BaseObject struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
