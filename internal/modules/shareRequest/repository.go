package shareRequest

import (
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/fileUpload"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/profile"
	"gorm.io/gorm"
)

type IRepository interface {
	CreateShareRequest(data *ShareRequestModel) error
	GetShareRequestById(id string) (*ShareRequestModel, error)
	GetShareRequestToUserId(userId string) ([]GetShareRequestToDomain, error)
	GetShareRequestByUserId(userId string) ([]GetShareRequestByDomain, error)
	UpdateShareRequestStatus(data *ShareRequestModel) error
	BeginTransaction() *gorm.DB
	UpdateShareRequestUserDataWithTransaction(tx *gorm.DB, data *ShareRequestModel) error
	BatchInsertFileUploadWithTransaction(tx *gorm.DB, data []fileUpload.FileUploadModel) error
	BatchInsertShareRequestFileWithTransaction(tx *gorm.DB, data []ShareRequestFileModel) error
	GetEmailFromUserId(userId string) (*profile.ProfileModel, error)
	GetShareRequestDetailsById(id string) (*GetShareRequestDetailsByIdDomain, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *repository) CreateShareRequest(data *ShareRequestModel) error {
	result := r.db.Create(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) GetShareRequestById(id string) (*ShareRequestModel, error) {
	var shareRequest ShareRequestModel
	result := r.db.Where("id = ?", id).First(&shareRequest)
	if result.Error != nil {
		return nil, result.Error
	}

	return &shareRequest, nil
}

func (r *repository) GetShareRequestToUserId(userId string) ([]GetShareRequestToDomain, error) {
	var shareRequests []GetShareRequestToDomain
	result := r.db.Table("share_requests").
		Select("users.username as request_by_name, share_requests.rsa_public_key, share_requests.status, share_requests.id, share_requests.created_at").
		Joins("join users on share_requests.request_by = users.id").
		Where("request_to = ?", userId).Find(&shareRequests)
	if result.Error != nil {
		return nil, result.Error
	}

	return shareRequests, nil
}

func (r *repository) GetShareRequestByUserId(userId string) ([]GetShareRequestByDomain, error) {
	var shareRequests []GetShareRequestByDomain
	result := r.db.Table("share_requests").
		Select("users.username as request_to_name, share_requests.rsa_public_key, share_requests.status, share_requests.id, share_requests.created_at").
		Joins("join users on share_requests.request_to = users.id").
		Where("request_by = ?", userId).Find(&shareRequests)
	if result.Error != nil {
		return nil, result.Error
	}

	return shareRequests, nil
}

func (r *repository) UpdateShareRequestStatus(data *ShareRequestModel) error {
	result := r.db.Model(data).Update("status", data.Status)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) UpdateShareRequestUserDataWithTransaction(tx *gorm.DB, data *ShareRequestModel) error {
	result := tx.Model(data).Update("user_profile_json", data.UserProfileJson)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) BatchInsertShareRequestFileWithTransaction(tx *gorm.DB, data []ShareRequestFileModel) error {
	result := tx.Create(&data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) BatchInsertFileUploadWithTransaction(tx *gorm.DB, data []fileUpload.FileUploadModel) error {
	result := tx.Create(&data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) GetEmailFromUserId(userId string) (*profile.ProfileModel, error) {
	var profile profile.ProfileModel
	result := r.db.Where("user_id = ?", userId).Order("created_at ASC").Limit(1).Find(&profile)
	if result.Error != nil {
		return nil, result.Error
	}

	return &profile, nil
}

// GetShareRequestDetailsById retrieves ShareRequest details with associated files
func (r *repository) GetShareRequestDetailsById(id string) (*GetShareRequestDetailsByIdDomain, error) {
	var files []GetShareRequestFilesDomain

	result := r.db.Table("share_request_files").
		Select("share_request_files.share_request_id, file_uploads.file_name, file_uploads.encryption_type, share_request_files.file_id").
		Joins("join file_uploads on share_request_files.file_id = file_uploads.id").
		Where("share_request_files.share_request_id = ?", id).
		Find(&files)

	if result.Error != nil {
		return nil, result.Error
	}

	var shareRequestDetails GetShareRequestDetailsByIdDomain
	result = r.db.Table("share_requests").
		Select("share_requests.id, share_requests.user_profile_json, share_requests.status").
		Where("share_requests.id = ?", id).
		First(&shareRequestDetails)

	if result.Error != nil {
		return nil, result.Error
	}

	shareRequestDetails.Files = files

	return &shareRequestDetails, nil
}
