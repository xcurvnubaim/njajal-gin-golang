package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// "gorm.io/gorm"
)

type BaseModels struct {
	ID        uuid.UUID     `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP, onUpdate:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt 
}
