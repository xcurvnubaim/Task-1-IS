package test

import (
	"testing"

	"github.com/xcurvnubaim/Task-1-IS/internal/configs"
	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/util"
)

func init() {
	if err := configs.Setup("../.env"); err != nil {
		panic(err)
	}
}

func TestInitVault(t *testing.T) {
	client, err := util.InitVault()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if client == nil {
		t.Errorf("Client is nil")
	}
}

func TestStoreUserKey(t *testing.T) {
	client, err := util.InitVault()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if client == nil {
		t.Errorf("Client is nil")
	}

	aesKey, _ := util.GenerateAESKey()
	rc4Key, _ := util.GenerateRC4Key()
	desKey, _ := util.GenerateDESKey()

	err = util.StoreUserKey(client, "123", aesKey, rc4Key, desKey)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetUserKey(t *testing.T) {
	client, err := util.InitVault()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if client == nil {
		t.Errorf("Client is nil")
	}

	aesKey, _ := util.GenerateAESKey()
	rc4Key, _ := util.GenerateRC4Key()
	desKey, _ := util.GenerateDESKey()

	err = util.StoreUserKey(client, "123", aesKey, rc4Key, desKey)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("AES Key: %v", aesKey)
	key, err := util.GetUserKey(client, "123", "aes")
	t.Logf("AES Key: %v", key)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if key != aesKey {
		t.Errorf("Key is not equal")
	}
}
