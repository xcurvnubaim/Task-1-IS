package shareRequest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/middleware"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/app"
	CustomValidator "github.com/xcurvnubaim/Task-1-IS/internal/pkg/validator"
)

type Handler struct {
	useCase IUseCase
	app     *gin.Engine
}

func NewHandler(app *gin.Engine, useCase IUseCase, prefixApi string) {
	handler := &Handler{
		app:     app,
		useCase: useCase,
	}

	handler.Routes(prefixApi)
}

func (h *Handler) Routes(prefix string) {
	handler := h.app.Group(prefix)
	{
		// handler.POST("/", h.UploadFile)

		handler.Use(middleware.AuthenticateJWT())
		{
			handler.POST("/", h.CreateShareRequest)
			handler.PUT("/", h.UpdateShareRequest)
			handler.GET("/by-me", h.GetAllShareRequestByMe)
			handler.GET("/to-me", h.GetAllShareRequestToMe)
			handler.GET("/:id", h.GetShareRequestByID)
		}
	}
}

func (h *Handler) CreateShareRequest(c *gin.Context) {
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
	uid, errUuid := uuid.Parse(userIDStr)
	if errUuid != nil {
		c.JSON(400, app.NewErrorResponse("Invalid user ID", nil))
		return
	}

	var data CreateShareRequestDTO
	data.UserId = uid
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, app.NewErrorResponse("Invalid request body", nil))
		return
	}

	res, errM := h.useCase.CreateShareRequest(&data)
	if errM != nil {
		c.JSON(500, app.NewErrorResponse(errM.Error(), nil))
		return
	}

	c.JSON(200, app.NewSuccessResponse("Create Share request successfully", res))
}

func (h *Handler) UpdateShareRequest(c *gin.Context) {
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
	uid, errUuid := uuid.Parse(userIDStr)
	if errUuid != nil {
		c.JSON(400, app.NewErrorResponse("Invalid user ID", nil))
		return
	}

	var data UpdateShareRequestDTO
	data.UserId = uid

	if err := c.ShouldBindJSON(&data); err != nil {
		var errMessages = CustomValidator.FormatValidationErrors(err)
		c.JSON(400, app.NewErrorResponse("Invalid request body", &errMessages))
		return
	}

	res, errM := h.useCase.UpdateShareRequestStatus(&data)
	if errM != nil {
		var errMsg = errM.Error()
		c.JSON(500, app.NewErrorResponse("Internal Server Error", &errMsg))
		return
	}

	c.JSON(200, app.NewSuccessResponse("Share request updated successfully", res))
}

func (h *Handler) GetAllShareRequestByMe(c *gin.Context) {
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
	id, errUuid := uuid.Parse(userIDStr)
	if errUuid != nil {
		c.JSON(400, app.NewErrorResponse("Invalid user ID", nil))
		return
	}

	var data GetAllShareRequestDTO
	data.UserId = id

	// Save the file to the server
	res, err := h.useCase.GetAllShareRequestByUser(&data)
	if err != nil {
		c.JSON(500, app.NewErrorResponse("Internal Server Error", nil))
		return
	}

	c.JSON(200, app.NewSuccessResponse("Files retrieved successfully", res))
}

func (h *Handler) GetAllShareRequestToMe(c *gin.Context) {
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
	id, errUuid := uuid.Parse(userIDStr)
	if errUuid != nil {
		c.JSON(400, app.NewErrorResponse("Invalid user ID", nil))
		return
	}

	var data GetAllShareRequestDTO
	data.UserId = id

	// Save the file to the server
	res, err := h.useCase.GetAllShareRequestToUser(&data)
	if err != nil {
		c.JSON(500, app.NewErrorResponse("Internal Server Error", nil))
		return
	}

	c.JSON(200, app.NewSuccessResponse("Share request retrieved successfully", res))
}

func (h *Handler) GetShareRequestByID(c *gin.Context) {
	RequestId := c.Param("id")
	EncryptedAesKey := c.Query("aes_key")
	// Convert string to UUID
	id, errUuid := uuid.Parse(RequestId)
	if errUuid != nil {
		c.JSON(400, app.NewErrorResponse("Invalid request ID", nil))
		return
	}
	data := GetShareRequestDetailsByIdRequestDTO{
		Id:              id.String(),
		AESKeyEncrypted: &EncryptedAesKey,
	}
	res, err := h.useCase.GetShareRequestDetailsById(&data)
	if err != nil {
		c.JSON(500, app.NewErrorResponse("Internal Server Error", nil))
		return
	}
	c.JSON(200, app.NewSuccessResponse("Share request retrieved successfully", res))
}
