package test

import (
	"testing"

	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
)

func TestCryptoAES(t *testing.T) {
	var plainText = []byte("Hello World")

	key, _ := util.GenerateAESKey() 

	encryptedText, err := util.EncryptPlainTextAESGCM(plainText, string(key))
	if err != nil {
		t.Error(err)
	}

	decryptedText, err := util.DecryptCipherTextAESGCM(encryptedText,  string(key))
	if err != nil {
		t.Error(err)
	}
	if string(plainText) != string(decryptedText) {
		t.Error("Decrypted text is not equal to plain text")
	}
}

func TestCryptoRC4(t *testing.T) {
	var plainText = []byte("Hello World")

	encryptedText, err := util.EncryptPlainTextRC4(plainText, "example key 1234")
	if err != nil {
		t.Error(err)
	}

	decryptedText, err := util.DecryptCipherTextRC4(encryptedText, "example key 1234")
	if err != nil {
		t.Error(err)
	}
	if string(plainText) != string(decryptedText) {
		t.Error("Decrypted text is not equal to plain text")
	}
}

func TestCryptoDES(t *testing.T) {
	var plainText = []byte("Hello World")

	encryptedText, err := util.EncryptPlainTextDES(plainText, "12345678")
	if err != nil {
		t.Error(err)
	}

	decryptedText, err := util.DecryptCipherTextDES(encryptedText, "12345678")
	if err != nil {
		t.Error(err)
	}
	if string(plainText) != string(decryptedText) {
		t.Error("Decrypted text is not equal to plain text")
	}
}