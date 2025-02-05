package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPass             string
	DBName             string
	ExternalAPIBaseURL string
}

func LoadConfig() *Config {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("[WARN] No .env file found: %v", err)
	}

	return &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPass:             getEnv("DB_PASS", ""),
		DBName:             getEnv("DB_NAME", "songsdb"),
		ExternalAPIBaseURL: getEnv("EXTERNAL_API_BASE_URL", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
