package shareRequest

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/common"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/fileUpload"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/profile"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/e"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/email"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
)

type IUseCase interface {
	CreateShareRequest(data *CreateShareRequestDTO) (*CreateShareResponseDTO, e.ApiError)
	GetAllShareRequestByUser(data *GetAllShareRequestDTO) (*GetAllShareRequestResponseDTO, e.ApiError)
	GetAllShareRequestToUser(data *GetAllShareRequestDTO) (*GetAllShareRequestResponseDTO, e.ApiError)
	UpdateShareRequestStatus(data *UpdateShareRequestDTO) (*UpdateShareRequestResponseDTO, e.ApiError)
	GetShareRequestDetailsById(data *GetShareRequestDetailsByIdRequestDTO) (*GetShareRequestDetailsByIdResponseDTO, e.ApiError)
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
	// check email
	profile, err := u.repository.GetEmailFromUserId(data.UserId.String())
	if err != nil || profile.Email == "" {
		log.Println("Error getting email from user id")
		return nil, e.NewApiError(400, "Email not found")
	}

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

	requestId := uuid.New()

	if err := util.StoreRequestShareKey(u.vaultClient, requestId.String(), rsaKeyPair.PrivateKey, "rsa"); err != nil {
		log.Println("Error storing RSA private key")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_STORE_RSA_PRIVATE_KEY))
	}

	shareRequest := NewShareRequest(requestId, &data.UserId, &requestTo, rsaKeyPair.PublicKey, RequestStatus.Pending, "")
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
			RequestToName: shareRequest.RequestToName,
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

func (u *useCase) GetShareRequestDetailsById(data *GetShareRequestDetailsByIdRequestDTO) (*GetShareRequestDetailsByIdResponseDTO, e.ApiError) {
	// Get detail
	details, err := u.repository.GetShareRequestDetailsById(data.Id)
	if err != nil {
		log.Println("Error getting share request files")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_GET_SHARE_REQUEST_DETAILS))
	}

	// Get RSA private key
	RSAPrivateKey, err := util.GetStoredRequestShareKey(u.vaultClient, details.ID, "rsa")
	if err != nil {
		log.Println("Error getting RSA private key")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_GET_RSA_PRIVATE_KEY))
	}

	var aesKey string
	if data.AESKeyEncrypted != nil {
		// Decrypt AES key
		AESKeyDecoded, err := base64.StdEncoding.DecodeString(*data.AESKeyEncrypted)
		if err != nil {
			log.Println("Error decoding AES key", err)
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_DECODE_AES_KEY))
		}
		aesKeyBytes, err := util.DecryptCipherTextRSA(AESKeyDecoded, RSAPrivateKey)
		if err != nil {
			log.Println("Error decrypting AES key", err)
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_DECRYPT_AES_KEY))
		}

		aesKey = base64.StdEncoding.EncodeToString(aesKeyBytes)

		err = util.StoreRequestShareKey(u.vaultClient, details.ID, aesKey, "aes")
		if err != nil {
			log.Println("Error storing AES key")
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_STORE_AES_KEY))
		}
	} else {
		aesKey, err = util.GetStoredRequestShareKey(u.vaultClient, details.ID, "aes")
		if err != nil {
			log.Println("Error getting AES key")
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_GET_AES_KEY))
		}
	}

	// Decrypt user profile
	decryptedProfile, err := u.decryptData(details.UserProfileJson, aesKey)
	if err != nil {
		log.Println("Error decrypting user profile")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_DECRYPT_USER_PROFILE))
	}

	var responseFiles []GetShareRequestFilesDomain
	for _, file := range details.Files {
		responseFiles = append(responseFiles, GetShareRequestFilesDomain{
			ShareRequestId: file.ShareRequestId,
			FileId:         file.FileId,
			FileName:       file.FileName,
			EncryptionType: file.EncryptionType,
		})
	}

	return &GetShareRequestDetailsByIdResponseDTO{
		ID:              details.ID,
		UserProfileJson: decryptedProfile,
		Files:           responseFiles,
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

	aes_key, err := util.GenerateAESKey()
	if err != nil {
		log.Println("Error generating AES key")
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_GENERATE_AES_KEY))
	}

	if data.Status == RequestStatus.Accepted {
		if err := u.createDetailShareRequest(shareRequest.RequestTo, shareRequest.ID, aes_key); err != nil {
			log.Println("Error creating detail share request")
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error: %d", e.ERROR_CREATE_DETAIL_SHARE_REQUEST))
		}
	}

	u.sendEmail(aes_key, shareRequest.RSAPublicKey, shareRequest.RequestBy)

	return &UpdateShareRequestResponseDTO{
		ID:     shareRequest.ID.String(),
		Status: shareRequest.Status,
	}, nil
}

func (u *useCase) createDetailShareRequest(userId, requestId uuid.UUID, aes_key string) error {
	tx := u.repository.BeginTransaction()
	defer tx.Rollback()

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
	errs = u.repository.UpdateShareRequestUserDataWithTransaction(tx, &ShareRequestModel{
		BaseModels: common.BaseModels{
			ID: requestId,
		},
		UserProfileJson: encryptedJsonProfile,
	})
	if errs != nil {
		tx.Rollback()
		log.Println("Error updating share request user data")
		return err
	}

	dataFile, err := u.fileUploadUseCase.GetAllFilesByUserId(userId.String())
	if err != nil {
		log.Println("Error getting all files by user id")
		return errors.New(err.Error())
	}

	var fileUploadModel []fileUpload.FileUploadModel
	var shareRequestFileModel []ShareRequestFileModel

	for _, file := range dataFile.Files {
		fileRes, err := u.fileUploadUseCase.DownloadFile(&fileUpload.FileDownloadRequestDTO{
			FileID: file.FileId,
			UserID: userId.String(),
		})

		if err != nil {
			log.Println("Error downloading file")
			return errors.New(err.Error())
		}

		encryptedFile, errs := util.EncryptPlainTextAESCBC(fileRes.FileBytes, aes_key)

		if errs != nil {
			log.Println("Error encrypting file")
			return errors.New(errs.Error())
		}

		// Save the encrypted file
		currentTime := time.Now().Format("2006010215")
		hash := sha256.New()
		hash.Write([]byte(file.FileName))
		hash.Write([]byte(time.Now().String()))
		fileName := file.FileName
		filePath := fmt.Sprintf("uploads/share/%s_%x_%s", currentTime, hash.Sum(nil), filepath.Ext(fileName))

		// save file to path
		if err := util.SaveFile(encryptedFile, filePath); err != nil {
			log.Println("Error saving file")
			return err
		}
		newFileId := uuid.New()
		fileUploadModel = append(fileUploadModel, fileUpload.FileUploadModel{
			BaseModels: common.BaseModels{
				ID: newFileId,
			},
			FilePath:       filePath,
			FileName:       fileName,
			KeyId:          requestId.String(),
			EncryptionType: "aes",
		})

		shareRequestFileModel = append(shareRequestFileModel, ShareRequestFileModel{
			ShareRequestId: requestId,
			FileId:         newFileId,
		})
	}

	if err := u.repository.BatchInsertFileUploadWithTransaction(tx, fileUploadModel); err != nil {
		log.Println("Error batch insert file upload", err)
		tx.Rollback()
		return err
	}

	if err := u.repository.BatchInsertShareRequestFileWithTransaction(tx, shareRequestFileModel); err != nil {
		log.Println("Error batch insert share request file")
		tx.Rollback()
		return err
	}
	tx.Commit()
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
		log.Println("Error decoding data", err)
		return "", fmt.Errorf("failed to decode data: %w", err)
	}

	decryptedData, err := util.DecryptCipherTextAESCBC(decodedData, key)
	if err != nil {
		log.Println("Error decrypting data", err)
		return "", fmt.Errorf("failed to decrypt data: %w", err)
	}

	return string(decryptedData), nil
}

func (u *useCase) sendEmail(aesKey string, rsaKey string, userId uuid.UUID) error {
	aesKeyDecoded, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		log.Println("Error decoding AES key", err)
	}

	// Get Email from user id
	profile, errs := u.profileUseCase.GetProfile(userId)
	if errs != nil {
		log.Println("Error getting profile", err)
		return errors.New(err.Error())
	}

	// Encrypt AES key using RSA public key
	encryptedAESKey, err := util.EncryptPlainTextRSA(aesKeyDecoded, rsaKey)
	if err != nil {
		log.Println("Error encrypting AES key", err)
		return err
	}

	encodedAESKey := base64.StdEncoding.EncodeToString(encryptedAESKey)

	err = email.SendEmail(*profile.Email, "Share Request", fmt.Sprintf("AES Key: %s", encodedAESKey))
	if err != nil {
		log.Println("Error sending email", err)
		return err
	}

	return nil
}
