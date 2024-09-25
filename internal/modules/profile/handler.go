package profile

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/middleware"
	"github.com/xcurvnubaim/Task-1-IS/internal/modules/common"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/app"
)

type ProfileHandler struct {
	profileUseCase IProfileUseCase
	app            *gin.Engine
}

func NewProfileHandler(app *gin.Engine, profileUseCase IProfileUseCase, prefixApi string) {
	authHandler := &ProfileHandler{
		app:            app,
		profileUseCase: profileUseCase,
	}

	authHandler.Routes(prefixApi)
}

func (h *ProfileHandler) Routes(prefix string) {
	profile := h.app.Group(prefix)
	{
		// authentication.POST("/register", ah.Register)
		// authentication.POST("/login", ah.Login)

		profile.Use(middleware.AuthenticateJWT())
		{
			profile.POST("/", h.CreateProfile)
			profile.GET("/", h.GetProfile)
		}
	}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	// Assuming userID is extracted from the context or request
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(400, app.NewErrorResponse("User ID not found in context", nil))
		return
	}

	// Perform a type assertion to ensure userID is a string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(400, app.NewErrorResponse("Invalid user ID type", nil))
		return
	}

	// Convert string to UUID
	parsedID, errUuid := uuid.Parse(userIDStr)
	if errUuid != nil {
		c.JSON(400, app.NewErrorResponse("Invalid user ID", nil))
		return
	}

	var profile CreateProfileRequestDTO
	profile.Id = parsedID
	profile.Fullname = c.PostForm("fullname")
	profile.ProfilePicturePath = fmt.Sprintf("%s/profile/%s.jpg", common.UploadFolder, parsedID.String())

	file, err := c.FormFile("profile_picture"); if err != nil {
		c.JSON(400, app.NewErrorResponse("Profile picture is required", nil))
		return
	}

	c.SaveUploadedFile(file, profile.ProfilePicturePath)

	res, errP := h.profileUseCase.CreateProfile(profile)
	if errP != nil {
		// Assuming err has a method Code() to retrieve HTTP status code
		var errMsg = errP.Error()
		c.JSON(errP.Code(), app.NewErrorResponse("Failed to create user profile", &errMsg))
		return
	}

	c.JSON(200, app.NewSuccessResponse("User profile created successfully", res))

}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Assuming userID is extracted from the context or request
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(400, app.NewErrorResponse("User ID not found in context", nil))
		return
	}

	// Perform a type assertion to ensure userID is a string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(400, app.NewErrorResponse("Invalid user ID type", nil))
		return
	}

	// Convert string to UUID
	parsedID, errUuid := uuid.Parse(userIDStr)
	if errUuid != nil {
		c.JSON(400, app.NewErrorResponse("Invalid user ID", nil))
		return
	}

	// Now call GetMe with the UUID
	user, errP := h.profileUseCase.GetProfile(parsedID)
	if errP != nil {
		// Assuming err has a method Code() to retrieve HTTP status code
		var errMsg = errP.Error()
		c.JSON(errP.Code(), app.NewErrorResponse("Failed to get user profile data", &errMsg))
		return
	}

	c.JSON(200, app.NewSuccessResponse("User profile data retrieved successfully", user))
}
