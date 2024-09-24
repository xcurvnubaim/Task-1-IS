package common

import (
	"time"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type BaseModels struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP, onUpdate:CURRENT_TIMESTAMP"`
	// DeletedAt gorm.DeletedAt
}
