package common

import "github.com/gin-gonic/gin"


type AuthHandler struct {
	authUseCase IUseCase
	app         *gin.Engine
}

func NewAuthHandler(app *gin.Engine, authUseCase IUseCase, prefixApi string) {
	authHandler := &AuthHandler{
		app:         app,
		authUseCase: authUseCase,
	}

	authHandler.Routes(prefixApi)
}

func (ah *AuthHandler) Routes(prefix string) {
	_ = ah.app.Group(prefix)
	{
		// authentication.POST("/register", ah.Register)
		// authentication.POST("/login", ah.Login)

		// authentication.Use(middleware.AuthenticateJWT())
		// {
		// 	authentication.GET("/me", ah.GetMe)
		// }
	}
}
