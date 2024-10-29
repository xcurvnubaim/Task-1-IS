package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
	"github.com/xcurvnubaim/Task-1-IS/internal/configs"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/e"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

type IAuthUseCase interface {
	RegisterUser(*RegisterUserRequestDTO) e.ApiError
	LoginUser(*LoginUserRequestDTO) (*LoginUserResponseDTO, e.ApiError)
	GetMe(uuid.UUID) (*GetMeResponseDTO, e.ApiError)
	HashPassword(string) (string, error)
	VerifyPassword(string, string) bool
	GenerateToken(PayloadToken) (string, error)
	CreateKey(string) error
	GetAllUser() (*GetAllUsersResponseDTO, e.ApiError)
	GetUserByUsername(string) (*GetUser, e.ApiError)
}

type authUseCase struct {
	vaultClient *api.Client
	authRepository IAuthRepository
}

func NewAuthUseCase(vaultClient *api.Client,authRepository IAuthRepository) *authUseCase {
	return &authUseCase{
		vaultClient, 
		authRepository,
	}
}

func (uc *authUseCase) RegisterUser(data *RegisterUserRequestDTO) e.ApiError {
	userCheck, _ := uc.authRepository.GetUserByUsername(data.Username)

	if userCheck != nil {
		return e.NewApiError(400, "Username already registered")
	}

	hashedPassword, err := uc.HashPassword(data.Password)
	if err != nil {
		log.Println(err.Error())
		return e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_BCRYPT_HASH_FAILED))
	}

	user := &RegisterUserDomain{
		Id:       uuid.New(),
		Username:    data.Username,
		Password: hashedPassword,
	}

	if err := uc.CreateKey(user.Id.String()); err != nil {
		log.Println(err.Error())
		return e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_GENERATE_ENCRYPTION_KEY_FAILED))
	}

	if err := uc.authRepository.RegisterUser(user); err != nil {
		log.Println(err.Error())
		return e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}
	return nil
}

func (uc *authUseCase) LoginUser(data *LoginUserRequestDTO) (*LoginUserResponseDTO, e.ApiError) {
	user, err := uc.authRepository.GetUserByUsername(data.Username)
	if err != nil {
		return nil, e.NewApiError(400, "User not found")
	}

	if !uc.VerifyPassword(user.Password, data.Password) {
		return nil, e.NewApiError(400, "Password is incorrect")
	}

	payloadToken := PayloadToken{
		ID:   user.ID,
		Role: user.Role,
	}

	token, errToken := uc.GenerateToken(payloadToken)
	if errToken != nil {
		log.Println(errToken.Error())
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", e.ERROR_GENERATE_TOKEN_FAILED))
	}

	return &LoginUserResponseDTO{
		Username: user.Username,
		Roles: user.Role,
		Token: token,
	}, nil
}

func (uc *authUseCase) GetMe(userID uuid.UUID) (*GetMeResponseDTO, e.ApiError) {
	user, err := uc.authRepository.GetUserByID(userID)
	if err != nil {
		return &GetMeResponseDTO{}, e.NewApiError(404, "User not found")
	}
	return &GetMeResponseDTO{
		Username: user.Username,
		Roles: user.Role,
	}, nil
}

func (uc *authUseCase) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (uc *authUseCase) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (uc *authUseCase) GenerateToken(payloadToken PayloadToken) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = payloadToken.ID
	claims["role"] = payloadToken.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Menggunakan secret key dari variabel lingkungan
	secretKey := configs.Config.JWT_SECRET
	if secretKey == "" {
		// Handle kasus dimana secret key tidak ditemukan
		return "", errors.New("secret key for JWT is not set")
	}
	//create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (uc *authUseCase) CreateKey(userID string) error {
	aes_key, err := util.GenerateAESKey(); if err != nil {
		log.Println(err.Error())
		return errors.New("failed to generate encryption aes key")
	}

	rc4_key, err := util.GenerateRC4Key(); if err != nil {
		return errors.New("failed to generate encryption rc4 key")
	}

	des_key, err := util.GenerateDESKey(); if err != nil {
		return errors.New("failed to generate encryption des key")
	}

	if err := util.StoreUserKey(uc.vaultClient, userID, aes_key, rc4_key, des_key); err != nil {
		return err
	}

	return nil
}

func (uc *authUseCase) GetAllUser() (*GetAllUsersResponseDTO, e.ApiError) {
	users, err := uc.authRepository.GetAllUser()
	if err != nil {
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}

	var response []GetUser
	for _, user := range users {
		response = append(response, GetUser{
			ID:       user.ID.String(),
			Username: user.Username,
		})
	}

	return &GetAllUsersResponseDTO{
		Users: response,
	}, nil
}

func (uc *authUseCase) GetUserByUsername(username string) (*GetUser, e.ApiError) {
	user, err := uc.authRepository.GetUserByUsername(username)
	if err != nil {
		if err.Code() == e.ERROR_USER_NOT_FOUND {
			return nil, e.NewApiError(404, "User not found")
		}
		return nil, e.NewApiError(500, fmt.Sprintf("Internal Server Error (%d)", err.Code()))
	}
	return &GetUser{
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}