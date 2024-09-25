package profile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/e"
)

type IProfileUseCase interface {
	CreateProfile(CreateProfileRequestDTO) (*CreateProfileResponseDTO, e.ApiError)
	GetProfile(id uuid.UUID) (*GetProfileResponseDTO, e.ApiError)
	encryptFile(filePath string) error
}

type profileUseCase struct {
	repository IProfileRepository
}

func NewProfileUseCase(repository IProfileRepository) *profileUseCase {
	return &profileUseCase{repository}
}

func (uc *profileUseCase) CreateProfile(data CreateProfileRequestDTO) (*CreateProfileResponseDTO, e.ApiError) {
	profile := &CreateProfileDomain{
		Id:             data.Id,
		Fullname:       data.Fullname,
		ProfilePicture: data.ProfilePicturePath,
	}

	isExist, err := uc.repository.GetProfileById(profile.Id)
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}
	fmt.Println(isExist)
	if isExist.Email != "" {
		return nil, e.NewApiError(400, "Profile already exist")
	}

	if err := uc.repository.CreateProfile(profile); err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}

	// Encrypt the uploaded file
	if err := uc.encryptFile(data.ProfilePicturePath); err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_ENCRYPT_FILE_FAILED))
	}

	return &CreateProfileResponseDTO{
		Fullname:       profile.Fullname,
		ProfilePicture: profile.ProfilePicture,
	}, nil
}

func (uc *profileUseCase) GetProfile(id uuid.UUID) (*GetProfileResponseDTO, e.ApiError) {
	profile, err := uc.repository.GetProfileById(id)
	if err != nil {
		log.Println(err.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}

	return &GetProfileResponseDTO{
		Email:          profile.Email,
		Roles:          profile.Roles,
		Fullname:       profile.FullName,
		ProfilePicture: profile.ProfilePicture,
	}, nil
}

// encryptFile encrypts the uploaded file using AES encryption
func (uc *profileUseCase) encryptFile(filePath string) error {
	// Read the original file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	var encryptionKey = "example key 1234"
	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a GCM (Galois/Counter Mode) instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a nonce for encryption
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	encryptedData := gcm.Seal(nonce, nonce, data, nil)

	// Write the encrypted data back to the same file (or a new file)
	if err := os.WriteFile(filePath + ".enc", encryptedData, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	return nil
}
