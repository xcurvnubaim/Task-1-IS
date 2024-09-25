package profile

import (
	"mime/multipart"

	"github.com/google/uuid"
)

// CreateProfileRequest represents the request payload for creating a profile
type (
	CreateProfileRequestDTO struct {
		Id                 uuid.UUID
		Fullname           string                `form:"fullname" binding:"required"`        // User's full name
		ProfilePicture     *multipart.FileHeader `form:"profile_picture" binding:"required"` // Uploaded profile picture
		ProfilePicturePath string
	}

	CreateProfileResponseDTO struct {
		Fullname       string `json:"fullname"`        // User's full name
		ProfilePicture string `json:"profile_picture"` // Profile picture URL
	}

	GetProfileResponseDTO struct {
		Email          string `json:"email"`           // User's email
		Roles          string `json:"roles"`           // User's roles
		Fullname       string `json:"fullname"`        // User's full name
		ProfilePicture string `json:"profile_picture"` // Profile picture URL
	}
)
