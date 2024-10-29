package shareRequest

import (
	"gorm.io/gorm"
)

type IRepository interface {
	CreateShareRequest(data *ShareRequestModel) error
	GetShareRequestById(id string) (*ShareRequestModel, error)
	GetShareRequestToUserId(userId string) ([]GetShareRequestToDomain, error)
	GetShareRequestByUserId(userId string) ([]GetShareRequestByDomain, error)
	UpdateShareRequestStatus(data *ShareRequestModel) error
	UpdateShareRequestUserData(data *ShareRequestModel) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
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

func (r *repository) UpdateShareRequestUserData(data *ShareRequestModel) error {
	result := r.db.Model(data).Update("user_profile_json", data.UserProfileJson)
	if result.Error != nil {
		return result.Error
	}

	return nil
}