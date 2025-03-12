package configs

import (
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

	BASE_URL string
}

var Config = &ConfigEnv{}

func Setup(pathEnv string) {
	// Load .env file
	// Try loading .env, but ignore errors if running in Docker
	_ = godotenv.Load(pathEnv)
	
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
	
	Config.BASE_URL = os.Getenv("BASE_URL")
}
