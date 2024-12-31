package auth

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/configs"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/e"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/mail"
	"golang.org/x/crypto/bcrypt"
)

type IAuthUseCase interface {
	RegisterUser(*RegisterUserRequestDTO) (*RegisterUserResponseDTO, e.ApiError)
	LoginUser(*LoginUserRequestDTO) (*LoginUserResponseDTO, e.ApiError)
	GetMe(uuid.UUID) (*GetMeResponseDTO, e.ApiError)
	HashPassword(string) (string, error)
	VerifyPassword(string, string) bool
	generateOTPCode() string
	GenerateToken(PayloadToken) (string, error)
	GetAllUser() (*GetAllUsersResponseDTO, e.ApiError)
	VerifyOTPcode(*UserModel, string) error
	VerifyUser(*VerifyOTPRequestDTO) (*VerifyOTPResponseDTO, e.ApiError)
}

type authUseCase struct {
	authRepository IAuthRepository
}

func NewAuthUseCase(authRepository IAuthRepository) *authUseCase {
	return &authUseCase{
		authRepository,
	}
}

func (uc *authUseCase) RegisterUser(data *RegisterUserRequestDTO) (*RegisterUserResponseDTO, e.ApiError) {

	// Check email already registered
	userCheck, _ := uc.authRepository.GetUserByEmail(data.Email)

	if userCheck != nil {
		return nil, e.NewApiError(400, "Email already registered")
	}

	hashedPassword, err := uc.HashPassword(data.Password)
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_BCRYPT_HASH_FAILED))
	}

	user := NewUser(data.Email, hashedPassword, uc.generateOTPCode(), time.Now().Add(time.Minute*15))
	
	if err := uc.authRepository.RegisterUser(user); err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}

	if err := sendOTPcode(user.Otp, data.Email); err != nil {
		log.Println(err.Error())
	}

	return &RegisterUserResponseDTO{
		Email: data.Email,
		Roles: user.Role,
	}, nil
}

func (uc *authUseCase) generateOTPCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (uc *authUseCase) LoginUser(data *LoginUserRequestDTO) (*LoginUserResponseDTO, e.ApiError) {
	user, err := uc.authRepository.GetUserByEmail(data.Email)
	if err != nil {
		return nil, e.NewApiError(400, "User not found")
	}

	// Check if user is verified
	if user.VerifiedAt == nil {
		return nil, e.NewApiError(400, "User is not verified")
	}

	if !uc.VerifyPassword(user.Password, data.Password) {
		return nil, e.NewApiError(400, "Password is incorrect")
	}

	payloadToken := PayloadToken{
		ID:   user.ID,
		Role: user.Role,
	}

	token, errToken := uc.GenerateToken(payloadToken)
	if errToken != nil {
		log.Println(errToken.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_GENERATE_TOKEN_FAILED))
	}

	return &LoginUserResponseDTO{
		Email: user.Email,
		Roles: user.Role,
		Token: token,
	}, nil
}

func (uc *authUseCase) GetMe(userID uuid.UUID) (*GetMeResponseDTO, e.ApiError) {
	user, err := uc.authRepository.GetUserByID(userID)
	if err != nil {
		return &GetMeResponseDTO{}, e.NewApiError(404, "User not found")
	}
	return &GetMeResponseDTO{
		Email: user.Email,
		Roles: user.Role,
	}, nil
}

func (uc *authUseCase) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (uc *authUseCase) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (uc *authUseCase) GenerateToken(payloadToken PayloadToken) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = payloadToken.ID
	claims["role"] = payloadToken.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Menggunakan secret key dari variabel lingkungan
	secretKey := configs.Config.JWT_SECRET
	if secretKey == "" {
		// Handle kasus dimana secret key tidak ditemukan
		return "", errors.New("secret key for JWT is not set")
	}
	//create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (uc *authUseCase) GetAllUser() (*GetAllUsersResponseDTO, e.ApiError) {
	users, err := uc.authRepository.GetAllUser()
	if err != nil {
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}

	var response []GetUser
	for _, user := range users {
		response = append(response, GetUser{
			ID:    user.ID,
			Email: user.Email,
		})
	}

	return &GetAllUsersResponseDTO{
		Users: response,
	}, nil
}

func (uc *authUseCase) VerifyOTPcode(user *UserModel, code string) error {
	// Check if OTP code is valid
	if user.Otp == "" || code != user.Otp {
		return errors.New("invalid OTP code")
	}

	// Check if OTP code is expired
	if time.Now().After(user.OtpExpiredAt) {
		return errors.New("OTP code is expired")
	}

	// Reset OTP code
	user.Otp = ""
	user.OtpExpiredAt = time.Now()

	return nil
}

func sendOTPcode(otp, email string) error {
	// generate 6 digit random number
	// otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	bodyEmail := templateSendEmail(otp)
	err := mail.SendEmail(email, "Your OTP Code", bodyEmail)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (uc *authUseCase) VerifyUser(data *VerifyOTPRequestDTO) (*VerifyOTPResponseDTO, e.ApiError) {
	user, err := uc.authRepository.GetUserByEmail(data.Email)
	if err != nil {
		return nil, e.NewApiError(400, "User not found")
	}
	
	errApi := uc.VerifyOTPcode(user, data.OTP)
	if errApi != nil {
		return nil, e.NewApiError(400, errApi.Error())
	}

	now := time.Now()
	user.VerifiedAt = &now

	if err := uc.authRepository.UpdateUser(user); err != nil {
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}

	return &VerifyOTPResponseDTO{
		Email:      user.Email,
		VerifiedAt: user.VerifiedAt.Format("2006-01-02 15:04:05"),
	}, nil
}