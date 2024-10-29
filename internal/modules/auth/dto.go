package auth

type (
	LoginUserRequestDTO struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	LoginUserResponseDTO struct {
		Username string `json:"username"`
		Roles    string `json:"roles"`
		Token    string `json:"token"`
	}

	RegisterUserRequestDTO struct {
		Username        string `json:"username" binding:"min=5,required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	RegisterUserResponseDTO struct {
		Username string `json:"username"`
	}

	GetMeResponseDTO struct {
		Username string `json:"username"`
		Roles    string `json:"roles"`
	}

	GetUser struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	GetAllUsersResponseDTO struct {
		Users []GetUser `json:"users"`
	}
)
