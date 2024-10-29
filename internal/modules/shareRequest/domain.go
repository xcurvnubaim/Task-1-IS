package shareRequest

import (
	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/common"
)

type (
	ShareRequestModel struct {
		common.BaseModels
		RequestBy       uuid.UUID `gorm:"column:request_by"`
		RequestTo       uuid.UUID `gorm:"column:request_to"`
		RSAPublicKey    string    `gorm:"column:rsa_public_key"`
		Status          string    `gorm:"column:status"`
		UserProfileJson string    `gorm:"column:user_profile_json"`
	}

	ShareRequestFileModel struct {
		ShareRequestId uuid.UUID `gorm:"column:share_request_id"`
		FileId         uuid.UUID `gorm:"column:file_id"`
	}

	GetShareRequestToDomain struct {
		ID            string `gorm:"column:id"`
		RequestByName string `gorm:"request_by_name"`
		RSAPublicKey  string `gorm:"column:rsa_public_key"`
		Status        string `gorm:"column:status"`
		CreatedAt     string `gorm:"column:created_at"`
	}

	GetShareRequestByDomain struct {
		ID            string `gorm:"column:id"`
		RequestToName string `gorm:"request_by_name"`
		RSAPublicKey  string `gorm:"column:rsa_public_key"`
		Status        string `gorm:"column:status"`
		CreatedAt     string `gorm:"column:created_at"`
	}

	GetShareRequestFilesDomain struct {
		ShareRequestId string `gorm:"column:share_request_id"`
		FileId         string `gorm:"column:file_id"`
		FileName       string `gorm:"column:file_name"`
		EncryptionType string `gorm:"column:encryption_type"`
	}

	GetShareRequestDetailsByIdDomain struct {
		ID              string                       `gorm:"column:id"`
		UserProfileJson string                       `gorm:"column:user_profile_json"`
		Status          string                       `gorm:"column:status"`
		Files           []GetShareRequestFilesDomain `gorm:"foreignKey:ShareRequestId;references:ID"`
	}
)

func (ShareRequestModel) TableName() string {
	return "share_requests"
}

func (ShareRequestFileModel) TableName() string {
	return "share_request_files"
}

func NewShareRequest(id uuid.UUID, request_by, request_to *uuid.UUID, rsa_pub_key, status, user_profile_json string) *ShareRequestModel {
	return &ShareRequestModel{
		BaseModels: common.BaseModels{
			ID: id,
		},
		RequestBy:       *request_by,
		RequestTo:       *request_to,
		RSAPublicKey:    rsa_pub_key,
		Status:          status,
		UserProfileJson: user_profile_json,
	}
}

func NewShareRequestFile(share_request_id, file_id *uuid.UUID) *ShareRequestFileModel {
	return &ShareRequestFileModel{
		ShareRequestId: *share_request_id,
		FileId:         *file_id,
	}
}
