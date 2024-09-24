package common

import (
	"time"
	// "gorm.io/gorm"
)

type BaseModels struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt
}
