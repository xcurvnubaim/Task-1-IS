package shareRequest

import "github.com/google/uuid"

type (
	CreateShareRequestDTO struct {
		UserId uuid.UUID 
		RequestTo string `json:"request_to" binding:"required"`
	}

	CreateShareResponseDTO struct {
		ID string `json:"id"`
		CreatedAt string `json:"created_at"`
	}

	UpdateShareRequestDTO struct {
		UserId uuid.UUID
		ID string `json:"id" binding:"required"`
		Status string `json:"status" binding:"required,oneof=accepted rejected"`
	}

	UpdateShareRequestResponseDTO struct {
		ID string `json:"id"`
		Status string `json:"status"`
	}

	GetAllShareRequestDTO struct {
		UserId uuid.UUID
	}	

	GetShareRequestResponseDTO struct {
		ID string `json:"id"`
		RequestByName string `json:"request_by_name"`
		Status string `json:"status"`
		RSAPublicKey string `json:"rsa_public_key"`
		CreatedAt string `json:"created_at"`
	}	

	GetAllShareRequestResponseDTO struct {
		Request []GetShareRequestResponseDTO `json:"request"`
	}

	UpdateShareRequestStatusDTO struct {
		RequestID string `json:"request_id" binding:"required"`
		Status string `json:"status" binding:"required"`
	}
)
