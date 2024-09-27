package test

import (
	"testing"

	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
)

func TestCrypto(t *testing.T) {
	var plainText = []byte("Hello World")

	encryptedText, err := util.EncryptPlainText(plainText)
	if err != nil {
		t.Error(err)
	}

	decryptedText, err := util.DecryptCipherText(encryptedText)
	if err != nil {
		t.Error(err)
	}

	if string(plainText) != string(decryptedText) {
		t.Error("Decrypted text is not equal to plain text")
	}
}
