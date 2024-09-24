package auth

type (
	LoginUserRequestDTO struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	LoginUserResponseDTO struct {
		Email string `json:"email"`
		Roles string `json:"roles"`
		Token string `json:"token"`
	}

	RegisterUserRequestDTO struct {
		Email           string `json:"email" binding:"email,required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	RegisterUserResponseDTO struct {
		Email string `json:"email"`
	}

	GetMeResponseDTO struct {
		Email string `json:"email"`
		Roles string `json:"roles"`
	}
)
