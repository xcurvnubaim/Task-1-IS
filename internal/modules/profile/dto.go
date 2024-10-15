package profile

import (
	"mime/multipart"

	"github.com/google/uuid"
)

// CreateProfileRequest represents the request payload for creating a profile
type (
	CreateProfileRequestDTO struct {
		Id                 uuid.UUID
		Fullname           string               `form:"fullname"`        // User's full name
		Email              string               `form:"email"`           // User's email
		Phone              string               `form:"phone"`           // User's phone number
		Address            string               `form:"address"`         // User's address
		Nik                string               `form:"nik"`             // User's National Identity Number
		ProfilePicture     *multipart.FileHeader `form:"profile_picture"` // Uploaded profile picture
		ProfilePictureByte []byte
		ProfilePicturePath string
	}

	CreateProfileResponseDTO struct {
		Fullname       *string `json:"fullname"`        // User's full name
		Email          *string `json:"email"`           // User's email
		Phone          *string `json:"phone"`           // User's phone number
		Address        *string `json:"address"`         // User's address
		Nik            *string `json:"nik"`             // User's National Identity Number
		ProfilePicture *string  // Profile picture URL
	}

	GetProfileResponseDTO struct {
		Email          *string `json:"email"`           // User's email
		Roles          *string `json:"roles"`           // User's roles
		Fullname       *string `json:"fullname"`        // User's full name
		Phone          *string `json:"phone"`           // User's phone number
		Address        *string `json:"address"`         // User's address
		Nik            *string `json:"nik"`             // User's National Identity Number
		ProfilePicture *string  // Profile picture URL
	}
)
