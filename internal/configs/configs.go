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

	EMAIL_FROM string
	SMTP_HOST string
	SMTP_PORT string
	SMTP_USER string
	SMTP_PASSWORD string
}

var Config = &ConfigEnv{}

func Setup(pathEnv string) error {
	// Load .env file
	err := godotenv.Load(pathEnv)
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

	Config.EMAIL_FROM = os.Getenv("EMAIL_FROM")
	Config.SMTP_HOST = os.Getenv("SMTP_HOST")
	Config.SMTP_PORT = os.Getenv("SMTP_PORT")
	Config.SMTP_USER = os.Getenv("SMTP_USER")
	Config.SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	
	return nil
}