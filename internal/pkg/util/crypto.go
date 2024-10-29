package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rc4"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
)

// GenerateAESKey generates a key for AES with a given bit size (16, 24, or 32 bytes for AES-128, AES-192, AES-256).
func GenerateAESKey() (string, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// GenerateRC4Key generates a key for RC4 (recommended key size is between 5 and 256 bytes).
func GenerateRC4Key() (string, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// GenerateDESKey generates an 8-byte key for DES encryption.
func GenerateDESKey() (string, error) {
	key := make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

type RSAKeyPair struct {
	PublicKey string
	PrivateKey string
}

// GenerateRSAKeyPair generates an RSA key pair with the given bit size.
func GenerateRSAKeyPair() (*RSAKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		PublicKey: base64.StdEncoding.EncodeToString(publicKeyBytes),
		PrivateKey: base64.StdEncoding.EncodeToString(privateKeyBytes),
	}, nil
}

// EncryptPlainTextRSA encrypts plaintext bytes using RSA-OAEP encryption
func EncryptPlainTextRSA(plainText []byte, publicKey string) ([]byte, error) {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Use SHA-256 as the hashing function for RSA-OAEP
	hash := sha256.New()
	encryptedData, err := rsa.EncryptOAEP(hash, rand.Reader, pubKey.(*rsa.PublicKey), plainText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	return encryptedData, nil
}

// DecryptCipherTextRSA decrypts ciphertext bytes using RSA-OAEP decryption
func DecryptCipherTextRSA(encryptedData []byte, privateKey string) ([]byte, error) {
	privKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Use SHA-256 as the hashing function for RSA-OAEP
	hash := sha256.New()
	decryptedData, err := rsa.DecryptOAEP(hash, rand.Reader, privKey, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decryptedData, nil
}

// encryptFile encrypts file bytes using AES encryption
func EncryptPlainTextAESGCM(plainText []byte, encryptionKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey); 
	
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
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
func DecryptCipherTextAESGCM(encryptedData []byte, encryptionKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey); 
	
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(key))

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

// EncryptPlainTextAESCBC encrypts plaintext bytes using AES in CBC mode
func EncryptPlainTextAESCBC(plainText []byte, encryptionKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Generate a random IV
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	// Pad the plaintext to be a multiple of the block size
	paddedPlainText := pad(plainText, block.BlockSize())

	// Encrypt the data
	ciphertext := make([]byte, len(paddedPlainText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedPlainText)

	// Prepend the IV to the ciphertext
	return append(iv, ciphertext...), nil
}

// DecryptCipherTextAESCBC decrypts ciphertext bytes using AES in CBC mode
func DecryptCipherTextAESCBC(encryptedData []byte, encryptionKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Ensure that the encrypted data is long enough
	if len(encryptedData) < block.BlockSize() {
		return nil, fmt.Errorf("encrypted data is too short")
	}

	// Extract the IV from the encrypted data
	iv := encryptedData[:block.BlockSize()]
	encryptedData = encryptedData[block.BlockSize():]

	// Decrypt the data
	decrypted := make([]byte, len(encryptedData))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, encryptedData)

	// Unpad the decrypted data
	return unpad(decrypted), nil
}

func EncryptPlainTextRC4(plainText []byte, encryptionKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey); 
	
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	cipherText := make([]byte, len(plainText))
	cipher.XORKeyStream(cipherText, plainText)

	return cipherText, nil
}

func DecryptCipherTextRC4(cipherText []byte, encryptionKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey); 
	
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	plainText := make([]byte, len(cipherText))
	cipher.XORKeyStream(plainText, cipherText)

	return plainText, nil
}

func EncryptPlainTextDES(plainText []byte, encryptionKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey); 
	
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create DES cipher: %w", err)
	}

	plainText = pad(plainText, block.BlockSize())
	encryptedData := make([]byte, des.BlockSize+len(plainText))
	iv := encryptedData[:des.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptedData[des.BlockSize:], plainText)

	return encryptedData, nil
}

func DecryptCipherTextDES(encryptedData []byte, encryptionKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey); 
	
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create DES cipher: %w", err)
	}

	if len(encryptedData) < des.BlockSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	iv := encryptedData[:des.BlockSize]
	encryptedData = encryptedData[des.BlockSize:]
	if len(encryptedData)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("encrypted data is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptedData, encryptedData)

	plainText := unpad(encryptedData)
	return plainText, nil
}

// Padding for block ciphers
func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// Unpadding for block ciphers
func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}