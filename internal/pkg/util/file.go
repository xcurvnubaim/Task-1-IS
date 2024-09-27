package util

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveFile(file *multipart.FileHeader, path string) error {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(path)
	if err != nil {
		return err
	}

	// Copy the file
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil

}
func FileToBytes(file *multipart.FileHeader) ([]byte, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Read the file
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func BytesToFile(fileBytes []byte, path string) error {
	// Create the destination file
	dst, err := os.Create(path)
	if err != nil {
		return err
	}

	// Write the file
	if _, err = dst.Write(fileBytes); err != nil {
		return err
	}

	return nil
}

func ReadBytes(path string) ([]byte, error) {
	// Open the file
	src, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer src.Close()

	// Read the file

	fileBytes, err := io.ReadAll(src)

	if err != nil {
		return nil, err
	}

	return fileBytes, nil

}
