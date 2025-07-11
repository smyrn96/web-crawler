package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_DSN     string
	AUTH_TOKEN string
)

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	DB_DSN = os.Getenv("DB_DSN")
	AUTH_TOKEN = os.Getenv("AUTH_TOKEN")

	if DB_DSN == "" || AUTH_TOKEN == "" {
		log.Fatal("Missing required environment variables: DB_DSN or AUTH_TOKEN")
	}
}
