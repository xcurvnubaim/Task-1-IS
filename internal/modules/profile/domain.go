package profile

import (
	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/common"
	// "gorm.io/gorm"
)

type (
	GetProfileDomain struct {
		Email          string `gorm:"column:email"`
		Roles          string `gorm:"column:roles"`
		FullName       string `gorm:"column:fullname"`
		ProfilePicture string `gorm:"column:profile_picture"`
	}

	ProfileModel struct {
		common.BaseModels
		UserId         uuid.UUID `gorm:"not null, column:user_id"`
		Fullname       string    `gorm:"not null"`
		ProfilePicture string    `gorm:"not null"`
	}
)

func (ProfileModel) TableName() string {
	return "profiles"
}

func NewProfile(userId uuid.UUID, fullname, profilePicture string) *ProfileModel {
	return &ProfileModel{
		UserId:         userId,
		Fullname:       fullname,
		ProfilePicture: profilePicture,
	}
}