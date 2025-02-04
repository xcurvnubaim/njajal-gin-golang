package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/modules/common"
)

type (
	UserModel struct {
		common.BaseModels
		Email        string     `gorm:"unique;not null"`
		Password     string     `gorm:"not null"`
		Role         string     `gorm:"default:'user'"`
		Otp          string     `gorm:"not null"`
		OtpExpiredAt time.Time  `gorm:"not null"`
		VerifiedAt   *time.Time `gorm:"default:null"`
	}

	PayloadToken struct {
		ID   uuid.UUID
		Role string
	}
)

func (UserModel) TableName() string {
	return "users"
}

func NewUser(email, password, otp string, otpExpiredAt time.Time) *UserModel {
	return &UserModel{
		BaseModels:   common.NewBaseModels(),
		Email:        email,
		Password:     password,
		Otp:          otp,
		OtpExpiredAt: otpExpiredAt,
	}
}
