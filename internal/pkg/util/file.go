package util

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveFile(fileBytes []byte, path string) error {
	// Create the destination file
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Write the file
	if _, err = dst.Write(fileBytes); err != nil {
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
