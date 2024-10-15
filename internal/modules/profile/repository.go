package profile

import (
	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/e"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IProfileRepository interface {
	CreateProfile(data *ProfileModel) e.ApiError
	GetProfileById(id uuid.UUID) (*GetProfileDomain, e.ApiError)
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *profileRepository {
	return &profileRepository{db}
}

func (r *profileRepository) CreateProfile(profile *ProfileModel) e.ApiError {
	result := r.db.Clauses(clause.Returning{}).Create(profile)
	if result.Error != nil {
		return e.NewApiError(e.ERROR_CREATE_PROFILE_REPOSITORY_FAILED, result.Error.Error())
	}

	return nil
}

func (r *profileRepository) GetProfileById(id uuid.UUID) (*GetProfileDomain, e.ApiError) {
	var profile GetProfileDomain

	// Execute the query and check for errors
	if err := r.db.Table("users").
		Select("profiles.email as email, users.role as roles, profiles.fullname as fullname, profiles.profile_picture as profile_picture, profiles.phone as phone, profiles.address as address, profiles.nik as nik").
		Joins("join profiles on users.id = profiles.user_id").
		Where("users.id = ?", id). // Ensure you're using "users.id" for the WHERE clause
		Scan(&profile).Error; err != nil {

		// Handle other potential errors
		return nil, e.NewApiError(e.ERROR_GET_PROFILE_BY_ID_REPOSITORY, err.Error())
	}

	return &profile, nil
}
