package fileUpload

import (
	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/common"
	// "gorm.io/gorm"
)

type (
	FileUploadModel struct {
		common.BaseModels
		FileName       string    `gorm:"column:file_name"`
		FilePath       string    `gorm:"column:file_path"`
		OwnerID        *uuid.UUID `gorm:"column:owner_id"`
		EncryptionType string    `gorm:"column:encryption_type"`
		KeyId          string    `gorm:"column:key_id"`
	}
)

func (FileUploadModel) TableName() string {
	return "file_uploads"
}

func NewFileUpload(id uuid.UUID, fileName, filePath *string, ownerID uuid.UUID, encryptionType, keyId string) *FileUploadModel {
	return &FileUploadModel{
		BaseModels: common.BaseModels{
			ID: id,
		},
		FileName:       *fileName,
		FilePath:       *filePath,
		OwnerID:        &ownerID,
		EncryptionType: encryptionType,
		KeyId:          keyId,
	}
}
