package auth

import (
	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/common"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/e"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	RegisterUser(data *RegisterUserDomain) e.ApiError
	GetUserByEmail(email string) (*UserModel, e.ApiError)
	GetUserByID(id uuid.UUID) (*UserModel, e.ApiError)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{db}
}

func (r *authRepository) RegisterUser(data *RegisterUserDomain) e.ApiError {
	user := &UserModel{
		BaseModels: common.BaseModels{
            ID:        data.Id, // or however you assign the ID
        },
		Email:    data.Email,
		Password: data.Password,
	}

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