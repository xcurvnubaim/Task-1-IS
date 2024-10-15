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
		Phone          string `gorm:"column:phone"`
		Address        string `gorm:"column:address"`
		Nik            string `gorm:"column:nik"`
		ProfilePicture string `gorm:"column:profile_picture"`
	}

	ProfileModel struct {
		common.BaseModels
		UserId         uuid.UUID `gorm:"not null, column:user_id"`
		Fullname       string
		Email          string
		Phone          string
		Address        string
		Nik            string
		ProfilePicture string
	}
)

func (ProfileModel) TableName() string {
	return "profiles"
}

func NewProfile(userId uuid.UUID, fullname, email, phone, address, nik, profilePicture *string) *ProfileModel {
	return &ProfileModel{
		UserId:         userId,
		Fullname:       *fullname,
		Email:          *email,
		Phone:          *phone,
		Address:        *address,
		Nik:            *nik,
		ProfilePicture: *profilePicture,
	}
}
