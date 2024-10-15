package fileUpload

import (
	"crypto/sha256"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/e"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
)

type IUseCase interface {
	// GetFile(path string) ([]byte, error)
	UploadFile(*FileUploadRequestDTO) (*FileUploadResponseDTO, e.ApiError)
	DownloadFile(path *FileDownloadRequestDTO) (*FileDownloadResponseDTO, e.ApiError)
	GetAllFilesByUserId(userId string) (*GetAllFilesByUserIdResponseDTO, e.ApiError)
}

type useCase struct {
	repository  IRepository
	vaultClient *api.Client
}

func NewuseCase(vaultClient *api.Client, repository IRepository) *useCase {
	return &useCase{repository, vaultClient}
}

func (u *useCase) GetAllFilesByUserId(userId string) (*GetAllFilesByUserIdResponseDTO, e.ApiError) {
	files, err := u.repository.GetFileByUserId(userId)
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, "Internal Server Error")
	}
	var response GetAllFilesByUserIdResponseDTO
	for _, file := range files {
		response.Files = append(response.Files, FileUploadResponseDTO{
			FileId:         file.ID.String(),
			FileName:       file.FileName,
			EncryptionType: file.EncryptionType,
		})
	}

	return &response, nil
}

func (u *useCase) UploadFile(data *FileUploadRequestDTO) (*FileUploadResponseDTO, e.ApiError) {
	fileBytes, err := util.FileToBytes(data.File)
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, "Internal Server Error")
	}

	// Encrypt file
	key, err := util.GetUserKey(u.vaultClient, data.UserID, data.EncryptionType)
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, "Internal Server Error")
	}

	var encryptedFile []byte
	if data.EncryptionType == "aes" {
		encryptedFile, err = util.EncryptPlainTextAESCBC(fileBytes, key)
	} else if data.EncryptionType == "rc4" {
		encryptedFile, err = util.EncryptPlainTextRC4(fileBytes, key)
	} else if data.EncryptionType == "des" {
		encryptedFile, err = util.EncryptPlainTextDES(fileBytes, key)
	}

	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, "Internal Server Error")
	}

	// Save the encrypted file
	currentTime := time.Now().Format("2006010215")
	hash := sha256.New()
	hash.Write([]byte(data.File.Filename))
	filePath := fmt.Sprintf("uploads/files/%s_%x_%s", currentTime, hash.Sum(nil), filepath.Ext(data.File.Filename))

	if err := util.SaveFile(encryptedFile, filePath); err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, "Internal Server Error")
	}
	createFileUpload := NewFileUpload(uuid.New(), &data.File.Filename, &filePath, uuid.MustParse(data.UserID), data.EncryptionType)
	if err := u.repository.CreateFileUpload(createFileUpload); err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, "Internal Server Error")
	}
	fileId := createFileUpload.ID.String()
	return &FileUploadResponseDTO{
		FileId:         fileId,
		FileName:       createFileUpload.FileName,
		EncryptionType: createFileUpload.EncryptionType,
	}, nil

}

func (u *useCase) DownloadFile(data *FileDownloadRequestDTO) (*FileDownloadResponseDTO, e.ApiError) {
	file, err := u.repository.GetFileById(data.FileID)
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(404, "File not found")
	}
	// Decrypt file
	key, err := util.GetUserKey(u.vaultClient, data.UserID, file.EncryptionType)
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, "Internal Server Error")
	}

	fileBytes, err := util.ReadBytes(file.FilePath)
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, "Internal Server Error")
	}

	var decryptedFile []byte
	if file.EncryptionType == "aes" {
		decryptedFile, err = util.DecryptCipherTextAESCBC(fileBytes, key)
	} else if file.EncryptionType == "rc4" {
		decryptedFile, err = util.DecryptCipherTextRC4(fileBytes, key)
	} else if file.EncryptionType == "des" {
		decryptedFile, err = util.DecryptCipherTextDES(fileBytes, key)
	}

	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, "Internal Server Error")
	}

	return &FileDownloadResponseDTO{
		FileBytes: decryptedFile,
		FileName:  file.FileName,
	}, nil
}
