package config

import (
	"log"

	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	APIServer string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBURL:     os.Getenv("DB_URL"),
		APIServer: os.Getenv("API_SERVER"),
	}
}
