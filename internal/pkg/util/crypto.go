package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// encryptFile encrypts file bytes using AES encryption
func EncryptPlainText(plainText []byte) ([]byte, error) {
	// Read the original file

	var encryptionKey = "example key 1234"
	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a GCM (Galois/Counter Mode) instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a nonce for encryption
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	encryptedData := gcm.Seal(nonce, nonce, plainText, nil)

	return encryptedData, nil
}

// decryptFile decrypts file bytes using AES encryption
func DecryptCipherText(encryptedData []byte) ([]byte, error) {
	var encryptionKey = "example key 1234"
	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(encryptionKey))

	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a GCM (Galois/Counter Mode) instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Get the nonce size
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("encrypted data is too short")
	}

	// Extract the nonce from the encrypted data
	nonce, encryptedData := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt the data
	decryptedData, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decryptedData, nil

}
