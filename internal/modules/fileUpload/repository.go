package fileUpload

import (
	"gorm.io/gorm"
)

type IRepository interface {
	CreateFileUpload(*FileUploadModel) error
	GetFileById(id string) (*FileUploadModel, error)
	GetFileByUserId(userId string) ([]FileUploadModel, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateFileUpload(data *FileUploadModel) error {
	result := r.db.Create(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) GetFileById(id string) (*FileUploadModel, error) {
	var file FileUploadModel
	result := r.db.Where("id = ?", id).First(&file)
	if result.Error != nil {
		return nil, result.Error
	}

	return &file, nil
}

func (r *repository) GetFileByUserId(userId string) ([]FileUploadModel, error) {
	var files []FileUploadModel
	result := r.db.Where("owner_id = ?", userId).Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}

	return files, nil
}