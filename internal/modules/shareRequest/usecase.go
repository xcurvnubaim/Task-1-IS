package shareRequest

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/common"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/fileUpload"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/profile"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/e"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
)

type IUseCase interface {
	CreateShareRequest(data *CreateShareRequestDTO) (*CreateShareResponseDTO, e.ApiError)
	GetAllShareRequestByUser(data *GetAllShareRequestDTO) (*GetAllShareRequestResponseDTO, e.ApiError)
	GetAllShareRequestToUser(data *GetAllShareRequestDTO) (*GetAllShareRequestResponseDTO, e.ApiError)
	UpdateShareRequestStatus(data *UpdateShareRequestDTO) (*UpdateShareRequestResponseDTO, e.ApiError)
}

type useCase struct {
	repository        IRepository
	vaultClient       *api.Client
	profileUseCase    profile.IProfileUseCase
	fileUploadUseCase fileUpload.IUseCase
}

func NewuseCase(vaultClient *api.Client, repository IRepository, profileUseCase profile.IProfileUseCase, fileUploadUseCase fileUpload.IUseCase) *useCase {
	return &useCase{
		repository,
		vaultClient,
		profileUseCase,
		fileUploadUseCase,
	}
}

func (u *useCase) CreateShareRequest(data *CreateShareRequestDTO) (*CreateShareResponseDTO, e.ApiError) {
	requestTo, err := uuid.Parse(data.RequestTo)
	if err != nil {
		log.Println("Error parsing request to UUID")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_PARSING_UUID_USER))
	}

	rsaKeyPair, err := util.GenerateRSAKeyPair()
	if err != nil {
		log.Println("Error generating RSA key pair")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_GENERATE_RSA_KEYPAIR))
	}

	shareRequest := NewShareRequest(uuid.New(), &data.UserId, &requestTo, rsaKeyPair.PublicKey, RequestStatus.Pending, "")
	err = u.repository.CreateShareRequest(shareRequest)
	if err != nil {
		log.Println("Error creating share request")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_CREATE_SHARE_REQUEST))
	}

	return &CreateShareResponseDTO{
		ID:        shareRequest.ID.String(),
		CreatedAt: shareRequest.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *useCase) GetAllShareRequestByUser(data *GetAllShareRequestDTO) (*GetAllShareRequestResponseDTO, e.ApiError) {
	shareRequests, err := u.repository.GetShareRequestByUserId(data.UserId.String())
	if err != nil {
		log.Println("Error getting all share requests")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_GET_ALL_SHARE_REQUEST))
	}

	var response []GetShareRequestResponseDTO
	for _, shareRequest := range shareRequests {
		response = append(response, GetShareRequestResponseDTO{
			ID:            shareRequest.ID,
			RequestByName: shareRequest.RequestToName,
			Status:        shareRequest.Status,
			RSAPublicKey:  shareRequest.RSAPublicKey,
			CreatedAt:     shareRequest.CreatedAt,
		})
	}

	return &GetAllShareRequestResponseDTO{
		Request: response,
	}, nil
}

func (u *useCase) GetAllShareRequestToUser(data *GetAllShareRequestDTO) (*GetAllShareRequestResponseDTO, e.ApiError) {
	shareRequests, err := u.repository.GetShareRequestToUserId(data.UserId.String())
	if err != nil {
		log.Println("Error getting all share requests")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_GET_ALL_SHARE_REQUEST))
	}

	var response []GetShareRequestResponseDTO
	for _, shareRequest := range shareRequests {
		response = append(response, GetShareRequestResponseDTO{
			ID:            shareRequest.ID,
			RequestByName: shareRequest.RequestByName,
			Status:        shareRequest.Status,
			RSAPublicKey:  shareRequest.RSAPublicKey,
			CreatedAt:     shareRequest.CreatedAt,
		})
	}

	return &GetAllShareRequestResponseDTO{
		Request: response,
	}, nil
}

func (u *useCase) UpdateShareRequestStatus(data *UpdateShareRequestDTO) (*UpdateShareRequestResponseDTO, e.ApiError) {
	shareRequest, err := u.repository.GetShareRequestById(data.ID)
	if err != nil {
		log.Println("Error getting share request")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_GET_SHARE_REQUEST))
	}

	if shareRequest.RequestTo.String() != data.UserId.String() {
		log.Println("Error user not authorized to update share request")
		return nil, e.NewApiError(403, fmt.Sprintf("Forbidden: %d", e.ERROR_USER_NOT_AUTHORIZED))
	}

	shareRequest.Status = data.Status
	err = u.repository.UpdateShareRequestStatus(shareRequest)
	if err != nil {
		log.Println("Error updating share request status")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_UPDATE_SHARE_REQUEST_STATUS))
	}

	if data.Status == RequestStatus.Accepted {
		u.createDetailShareRequest(shareRequest.RequestTo, shareRequest.ID, shareRequest.RSAPublicKey)
	}

	return &UpdateShareRequestResponseDTO{
		ID:     shareRequest.ID.String(),
		Status: shareRequest.Status,
	}, nil
}

func (u *useCase) createDetailShareRequest(userId, requestId uuid.UUID, pub_rsa_key string) error {
	aes_key, err := util.GenerateAESKey(); if err != nil {
		log.Println("Error generating AES key")
		return err
	} 

	profile, err := u.profileUseCase.GetProfile(userId)
	if err != nil {
		log.Println("Error getting profile")
		return errors.New(err.Error())
	}
	jsonProfile, errs := json.Marshal(profile)
	if errs != nil {
		log.Println("Error marshalling profile")
		return err
	}
	encryptedJsonProfile, errs := u.encryptData(jsonProfile, aes_key)
	if errs != nil {
		log.Println("Error encrypting profile")
		return err
	}
	errs = u.repository.UpdateShareRequestUserData(&ShareRequestModel{
		BaseModels: common.BaseModels{
			ID: requestId,
		},
		UserProfileJson: encryptedJsonProfile,
	})
	if errs != nil {
		log.Println("Error updating share request user data")
		return err
	}

	dataFile, err := u.fileUploadUseCase.GetAllFilesByUserId(userId.String())
	if err != nil {
		log.Println("Error getting all files by user id")
		return errors.New(err.Error())
	}

	for _, file := range dataFile.Files {
		fileRes, err := u.fileUploadUseCase.DownloadFile(&fileUpload.FileDownloadRequestDTO{
			FileID: file.FileId,
			UserID: userId.String(),
		})

		if err != nil {
			log.Println("Error downloading file")
			return errors.New(err.Error())
		}

		_, errs := util.EncryptPlainTextAESCBC(fileRes.FileBytes, aes_key)

		if err != nil {
			log.Println("Error encrypting file")
			return errors.New(errs.Error())
		}
	}

	return nil
}

func (u *useCase) encryptData(data []byte, key string) (string, error) {
	encryptedData, err := util.EncryptPlainTextAESCBC(data, key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func (u *useCase) decryptData(data string, key string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode data: %w", err)
	}

	decryptedData, err := util.DecryptCipherTextAESCBC(decodedData, key)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %w", err)
	}

	return string(decryptedData), nil
}