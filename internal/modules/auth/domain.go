package auth

import (
	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/common"
)

type (
	RegisterUserDomain struct {
		Id       uuid.UUID
		Email    string
		Password string
	}

	UserModel struct {
		common.BaseModels
		Email    string `gorm:"unique;not null"`
		Password string `gorm:"not null"`
		Role     string `gorm:"default:'user'"`
	}

	PayloadToken struct {
		ID   uuid.UUID
		Role string
	}
)

func (UserModel) TableName() string {
	return "users"
}
