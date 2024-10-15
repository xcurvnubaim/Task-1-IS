package util

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/xcurvnubaim/Task-1-IS/internal/configs"
)

func InitVault() (*api.Client, error) {
	configVault := api.DefaultConfig()
	configVault.Address = configs.Config.VAULT_ADDR
	client, err := api.NewClient(configVault)
	if err != nil {
		return nil, err
	}
	client.SetToken(configs.Config.VAULT_TOKEN)
	return client, nil
}

func StoreUserKey(client *api.Client, userID string, aesKey string, rc4Key string, desKey string) error {
	// Define the Vault path for the user's key
	// Create the data to store in Vault
	data := map[string]interface{}{
		"aes_key": aesKey,
		"rc4_key": rc4Key,
		"des_key": desKey,
	}

	// Write the key to the Vault server
	// _, err := client.Logical().Write(path, data)
	_, err := client.KVv2("secret").Put(context.Background(), userID, data)
	if err != nil {
		return fmt.Errorf("error writing to Vault: %w", err)
	}

	return nil
}

func GetUserKey(client *api.Client, userID string, keyType string) (string, error) {
	secret, err := client.KVv2("secret").Get(context.Background(), userID)

	if err != nil {
		return "", err
	}

	switch keyType {
	case "aes":
		key, ok := secret.Data["aes_key"].(string)
		if !ok {
			return "", fmt.Errorf("key not found")
		}
		return key, nil
	case "rc4":
		key, ok := secret.Data["rc4_key"].(string)
		if !ok {
			return "", fmt.Errorf("key not found")
		}
		return key, nil
	case "des":
		key, ok := secret.Data["des_key"].(string)
		if !ok {
			return "", fmt.Errorf("key not found")
		}
		return key, nil
	default:
		return "", fmt.Errorf("key type not found")
	}
}