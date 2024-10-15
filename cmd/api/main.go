package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xcurvnubaim/Task-1-IS/internal/configs"
	"github.com/xcurvnubaim/Task-1-IS/internal/database"
	"github.com/xcurvnubaim/Task-1-IS/internal/middleware"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/auth"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/fileUpload"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/profile"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
)

func main() {
	// Setup configuration
	if err := configs.Setup(".env"); err != nil {
		panic(err)
	}

	// Setup for production
	if configs.Config.ENV_MODE == "production" {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("Production mode")
	}

	// Start the server
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	// Setup Database
	db, err  := database.New()
	if err != nil {
		panic(err)
	}

	vaultClient, err := util.InitVault()
	if err != nil {
		panic(err)
	}
	
	var authRepository auth.IAuthRepository = auth.NewAuthRepository(db)
	var authService auth.IAuthUseCase = auth.NewAuthUseCase(vaultClient, authRepository)
	auth.NewAuthHandler(r, authService, "/api/v1/auth")

	var profileRepository profile.IProfileRepository = profile.NewProfileRepository(db)
	var profileService profile.IProfileUseCase = profile.NewProfileUseCase(vaultClient, profileRepository)
	profile.NewProfileHandler(r, profileService, "/api/v1/profile")

	var FileUploadRepository fileUpload.IRepository = fileUpload.NewRepository(db)
	var FileUploadService fileUpload.IUseCase = fileUpload.NewuseCase(vaultClient, FileUploadRepository)
	fileUpload.NewHandler(r, FileUploadService, "/api/v1/file")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(":" + configs.Config.APP_PORT); err != nil {
		panic(err)
	}
}
