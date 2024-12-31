package auth

import "github.com/google/uuid"

type (
	LoginUserRequestDTO struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	LoginUserResponseDTO struct {
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Roles       string `json:"roles"`
		Token       string `json:"token"`
	}

	RegisterUserRequestDTO struct {
		Username        string `json:"username" binding:"min=5,required"`
		PhoneNumber     string `json:"phone_number" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	RegisterUserResponseDTO struct {
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
	}

	GetMeResponseDTO struct {
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Roles       string `json:"roles"`
	}

	GetUser struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
	}

	GetAllUsersResponseDTO struct {
		Users []GetUser `json:"users"`
	}

	VerifyOTPRequestDTO struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		OTP         string `json:"otp" binding:"required"`
	}
)
