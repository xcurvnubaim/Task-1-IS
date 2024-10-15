package test

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
)

// UploadFile Gin handler
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("uploadedFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving the file"})
		return
	}

	fileBytes, err := util.FileToBytes(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading the file"})
		return
	}

	encryptedBytes, err := util.EncryptPlainTextAESGCM(fileBytes, "example key 1234")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting the file"})
		return
	}

	savePath := "./uploads/" + file.Filename

	if err := os.Mkdir("./uploads", 0755); err != nil && !os.IsExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating the uploads directory", "message": err.Error()})
		return
	}

	err = util.BytesToFile(encryptedBytes, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving the file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File upload and encryption successful!"})
}

// Helper function to create a multipart form file
func createMultipartFormFile(t *testing.T, fieldName, fileName string, fileContent []byte) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatalf("failed to write file content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	return body, writer.FormDataContentType()
}

// Unit test for the UploadFile function
func TestUploadFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fileContent := []byte("This is a test file.")
	fieldName := "uploadedFile"
	fileName := "testfile.txt"

	body, contentType := createMultipartFormFile(t, fieldName, fileName, fileContent)

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", contentType)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.POST("/upload", UploadFile)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"message":"File upload and encryption successful!"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestDecryptFile(t *testing.T) {
	TestUploadFile(t)
	fileBytes, err := util.ReadBytes("./uploads/testfile.txt")
	if err != nil {
		t.Error(err)
	}

	decryptedBytes, err := util.DecryptCipherTextAESGCM(fileBytes, "example key 1234")
	if err != nil {
		t.Error(err)
	}

	expected := "This is a test file."
	log.Println(string(decryptedBytes))
	assert.Equal(t, expected, string(decryptedBytes))
}