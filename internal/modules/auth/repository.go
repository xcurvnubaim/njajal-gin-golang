package auth

import (
	"github.com/google/uuid"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/e"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	RegisterUser(*UserModel) e.ApiError
	GetUserByEmail(string) (*UserModel, e.ApiError)
	GetUserByID(uuid.UUID) (*UserModel, e.ApiError)
	GetAllUser() ([]UserModel, e.ApiError)
	UpdateUser(*UserModel) e.ApiError
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{db}
}

func (r *authRepository) RegisterUser(user *UserModel) e.ApiError {
	result := r.db.Create(user)
	if result.Error != nil {
		return e.NewApiError(e.ERROR_REGISTER_REPOSITORY_FAILED, result.Error.Error())
	}

	return nil
}

func (r *authRepository) GetUserByEmail(email string) (*UserModel, e.ApiError) {
	user := &UserModel{}
	result := r.db.Where("email = ?", email).First(user)
	if result.Error != nil {
		return nil, e.NewApiError(e.ERROR_GET_USER_BY_EMAIL_REPOSITORY_FAILED, result.Error.Error())
	}

	return user, nil
}

func (r *authRepository) GetUserByID(id uuid.UUID) (*UserModel, e.ApiError) {
	user := &UserModel{}
	result := r.db.Where("id = ?", id).First(user)
	if result.Error != nil {
		return nil, e.NewApiError(e.ERROR_GET_USER_BY_ID_REPOSITORY_FAILED, result.Error.Error())
	}

	return user, nil
}

func (r *authRepository) GetAllUser() ([]UserModel, e.ApiError) {
	var users []UserModel
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, e.NewApiError(e.ERROR_GET_ALL_USER_REPOSITORY_FAILED, result.Error.Error())
	}

	return users, nil
}

func (r *authRepository) UpdateUser(user *UserModel) e.ApiError {
	result := r.db.Save(user)
	if result.Error != nil {
		return e.NewApiError(e.ERROR_UPDATE_USER_REPOSITORY_FAILED, result.Error.Error())
	}

	return nil
}