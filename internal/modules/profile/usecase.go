package profile

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/e"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
)

type IProfileUseCase interface {
	CreateProfile(CreateProfileRequestDTO) (*CreateProfileResponseDTO, e.ApiError)
	GetProfile(id uuid.UUID) (*GetProfileResponseDTO, e.ApiError)
	encryptFile(string, string, []byte) error
}

type profileUseCase struct {
	vaultClient *api.Client
	repository  IProfileRepository
}

func NewProfileUseCase(vaultClient *api.Client, repository IProfileRepository) *profileUseCase {
	return &profileUseCase{
		vaultClient: vaultClient,
		repository:  repository,
	}
}

func (uc *profileUseCase) CreateProfile(data CreateProfileRequestDTO) (*CreateProfileResponseDTO, e.ApiError) {
	aesKey, err := util.GetUserKey(uc.vaultClient, data.Id.String(), "aes")
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_GET_USER_KEY_FAILED))
	}

	var encryptedData CreateProfileRequestDTO

	// Check if Fullname is not empty before encrypting
	if data.Fullname != "" {
		if encryptedData.Fullname, err = uc.encryptData(data.Fullname, aesKey); err != nil {
			log.Println(err.Error())
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_ENCRYPT_DATA_FAILED))
		}
	}

	// Check if Email is not empty before encrypting
	if data.Email != "" {
		if encryptedData.Email, err = uc.encryptData(data.Email, aesKey); err != nil {
			log.Println(err.Error())
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_ENCRYPT_DATA_FAILED))
		}
	}

	// Check if Phone is not empty before encrypting
	if data.Phone != "" {
		if encryptedData.Phone, err = uc.encryptData(data.Phone, aesKey); err != nil {
			log.Println(err.Error())
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_ENCRYPT_DATA_FAILED))
		}
	}

	// Check if Address is not empty before encrypting
	if data.Address != "" {
		if encryptedData.Address, err = uc.encryptData(data.Address, aesKey); err != nil {
			log.Println(err.Error())
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_ENCRYPT_DATA_FAILED))
		}
	}

	// Handle optional field Nik
	if data.Nik != "" {
		if encryptedData.Nik, err = uc.encryptData(data.Nik, aesKey); err != nil {
			log.Println(err.Error())
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_ENCRYPT_DATA_FAILED))
		}
	}


	profile := NewProfile(
		data.Id,
		&encryptedData.Fullname,
		&encryptedData.Email,
		&encryptedData.Phone,
		&encryptedData.Address,
		&encryptedData.Nik,
		&data.ProfilePicturePath,
	)

	isExist, errs := uc.repository.GetProfileById(profile.UserId)
	if errs != nil {
		log.Println(errs.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", errs.Code()))
	}
	fmt.Println(isExist)
	if isExist.Email != "" {
		return nil, e.NewApiError(400, "Profile already exist")
	}

	// Encrypt the uploaded file
	if err := uc.encryptFile(data.Id.String(), data.ProfilePicturePath, data.ProfilePictureByte); err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_ENCRYPT_FILE_FAILED))
	}

	if err := uc.repository.CreateProfile(profile); err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}

	return &CreateProfileResponseDTO{
		Fullname:       &data.Fullname,
		Email:          &data.Email,
		Phone:          &data.Phone,
		Address:        &data.Address,
		Nik:            &data.Nik,
		ProfilePicture: &profile.ProfilePicture,
	}, nil
}

func (uc *profileUseCase) GetProfile(id uuid.UUID) (*GetProfileResponseDTO, e.ApiError) {
	aesKey, err := util.GetUserKey(uc.vaultClient, id.String(), "aes")
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_GET_USER_KEY_FAILED))
	}

	profile, errs := uc.repository.GetProfileById(id)
	if errs != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", errs.Code()))
	}

	var decryptedData GetProfileResponseDTO

	// Check if Email is not empty before decrypting
	if profile.Email != "" {
		if decryptedData.Email, err = uc.decryptData(profile.Email, aesKey); err != nil {
			log.Println(err.Error(), "email")
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_DECRYPT_DATA_FAILED))
		}
	}

	// Check if Fullname is not empty before decrypting
	if profile.FullName != "" {
		if decryptedData.Fullname, err = uc.decryptData(profile.FullName, aesKey); err != nil {
			log.Println(err.Error())
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_DECRYPT_DATA_FAILED))
		}
	}

	// Check if Phone is not empty before decrypting
	if profile.Phone != "" {
		if decryptedData.Phone, err = uc.decryptData(profile.Phone, aesKey); err != nil {
			log.Println(err.Error())
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_DECRYPT_DATA_FAILED))
		}
	}

	// Check if Address is not empty before decrypting
	if profile.Address != "" {
		if decryptedData.Address, err = uc.decryptData(profile.Address, aesKey); err != nil {
			log.Println(err.Error())
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_DECRYPT_DATA_FAILED))
		}
	}

	// Check if Nik is not empty before decrypting
	if profile.Nik != "" {
		if decryptedData.Nik, err = uc.decryptData(profile.Nik, aesKey); err != nil {
			log.Println(err.Error())
			return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_DECRYPT_DATA_FAILED))
		}
	}


	return &GetProfileResponseDTO{
		Email:          decryptedData.Email,
		Roles:          &profile.Roles,
		Fullname:       decryptedData.Fullname,
		Phone:          decryptedData.Phone,
		Address:        decryptedData.Address,
		Nik:            decryptedData.Nik,
		ProfilePicture: &profile.ProfilePicture,
	}, nil
}

// encryptFile encrypts the uploaded file using AES encryption
func (uc *profileUseCase) encryptFile(id string, filePath string, fileBytes []byte) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	aesKey, err := util.GetUserKey(uc.vaultClient, id, "aes")
	if err != nil {
		return fmt.Errorf("failed to get user key: %w", err)
	}

	cyperText, err := util.EncryptPlainTextAESGCM(fileBytes, aesKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %w", err)
	}

	// Write the encrypted data back to the same file (or a new file)
	if err := os.WriteFile(filePath+".enc", cyperText, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	return nil
}

func (uc *profileUseCase) encryptData(data string, key string) (string, error) {
	encryptedData, err := util.EncryptPlainTextAESCBC([]byte(data), key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func (uc *profileUseCase) decryptData(data string, key string) (*string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	decryptedData, err := util.DecryptCipherTextAESCBC(decodedData, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	decryptedDataString := string(decryptedData) 
	return &decryptedDataString, nil
}