package fileUpload

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xcurvnubaim/Task-1-IS/internal/middleware"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/app"
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
			handler.GET("/", h.GetAllFilesByUserId)
			handler.POST("/upload", h.UploadFile)
			handler.GET("/download/:file_id", h.DownloadFile)
		}
	}
}

func (h *Handler) GetAllFilesByUserId(c *gin.Context) {
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
	_, errUuid := uuid.Parse(userIDStr)
	if errUuid != nil {
		c.JSON(400, app.NewErrorResponse("Invalid user ID", nil))
		return
	}

	// Save the file to the server
	res, err := h.useCase.GetAllFilesByUserId(userIDStr)
	if err != nil {
		c.JSON(500, app.NewErrorResponse("Internal Server Error", nil))
		return
	}

	c.JSON(200, app.NewSuccessResponse("Files retrieved successfully", res))
}



func (h *Handler) UploadFile(c *gin.Context) {
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

	var fileUpload FileUploadRequestDTO
	fileUpload.UserID = parsedID.String()

	err := c.ShouldBind(&fileUpload)

	if err != nil {
		errMsg := err.Error()
		c.JSON(400, app.NewErrorResponse("Validation Error", &errMsg))
		return
	}

	// Save the file to the server
	res, err := h.useCase.UploadFile(&fileUpload)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, app.NewSuccessResponse("File uploaded successfully", res))
}


func (h *Handler) DownloadFile(c *gin.Context) {
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

	var fileDownload FileDownloadRequestDTO
	
	fileDownload.UserID = parsedID.String()

	fileId := c.Param("file_id")

	fileDownload.FileID = fileId

	// Save the file to the server
	res, err := h.useCase.DownloadFile(&fileDownload)
	if err != nil {
		c.JSON(500, err)
		return
	}

	// c.JSON(200, app.NewSuccessResponse("File downloaded successfully", res))
	// Set the headers to indicate file download
    c.Header("Content-Disposition", "attachment; filename="+res.FileName)
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Length", string(len(res.FileBytes)))

	// Write the byte data to the response
    if _, err := c.Writer.Write(res.FileBytes); err != nil {
		log.Println(err.Error())
        c.JSON(500, res)
        return
    }

    // Optionally, end the response (not strictly necessary)
    c.Writer.Flush()
}