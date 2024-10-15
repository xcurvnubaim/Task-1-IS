package fileUpload

import "mime/multipart"

type (
	FileUploadRequestDTO struct {
		UserID         string
		EncryptionType string                `form:"encryption_type" binding:"required"`
		File           *multipart.FileHeader `form:"file" binding:"required"`
	}

	FileUploadResponseDTO struct {
		FileId string `json:"file_id"`
		// FilePath       string `json:"file_path"`
		FileName       string `json:"file_name"`
		EncryptionType string `json:"encryption_type"`
	}

	FileDownloadRequestDTO struct {
		UserID string
		FileID string
	}

	FileDownloadResponseDTO struct {
		FileBytes []byte
		FileName  string
	}

	GetAllFilesByUserIdResponseDTO struct {
		Files []FileUploadResponseDTO `json:"files"`
	}
)
