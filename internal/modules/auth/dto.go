package auth

import "github.com/google/uuid"

type (
	LoginUserRequestDTO struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	LoginUserResponseDTO struct {
		Email string `json:"email"`
		Roles string `json:"roles"`
		Token string `json:"token"`
	}

	RegisterUserRequestDTO struct {
		Email           string `json:"email" binding:"min=5,required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	RegisterUserResponseDTO struct {
		Email string `json:"email"`
		Roles string `json:"roles"`
	}

	GetMeResponseDTO struct {
		Email string `json:"email"`
		Roles string `json:"roles"`
	}

	GetUser struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
	}

	GetAllUsersResponseDTO struct {
		Users []GetUser `json:"users"`
	}

	VerifyOTPRequestDTO struct {
		Email string `json:"email" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}

	VerifyOTPResponseDTO struct {
		Email      string `json:"email"`
		VerifiedAt string `json:"verified_at"`
	}
)
