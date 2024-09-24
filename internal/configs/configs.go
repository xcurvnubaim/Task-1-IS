package configs

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

type ConfigEnv struct {
	DatabaseURL string
	DatabaseName string
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePassword string
	DatabaseDialect string

	ENV_MODE string
	APP_PORT string

	JWT_SECRET string
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
	Config.DatabaseName = os.Getenv("DB_NAME")
	Config.DatabaseHost = os.Getenv("DB_HOST")
	Config.DatabasePort = os.Getenv("DB_PORT")
	Config.DatabaseUser = os.Getenv("DB_USER")
	Config.DatabasePassword = os.Getenv("DB_PASSWORD")
	Config.DatabaseDialect = os.Getenv("DB_DIALECT")
	
	Config.ENV_MODE = os.Getenv("ENV_MODE")
	Config.APP_PORT = os.Getenv("APP_PORT")

	Config.JWT_SECRET = os.Getenv("JWT_SECRET")

	return nil
}
