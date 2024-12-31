package common

import (
	"gorm.io/gorm"
)

type IRepository interface {
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}
