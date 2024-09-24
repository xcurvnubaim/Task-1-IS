package configs

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

type ConfigEnv struct {
	DatabaseURL string
	DatabaseName string
	ENV_MODE string
	APP_PORT string
}

var Config = &ConfigEnv{}

func Setup() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return err
	}

	Config.DatabaseURL = os.Getenv("DATABASE_URL")
	Config.DatabaseName = os.Getenv("DATABASE_NAME")
	Config.ENV_MODE = os.Getenv("ENV_MODE")
	Config.APP_PORT = os.Getenv("APP_PORT")

	return nil
}
