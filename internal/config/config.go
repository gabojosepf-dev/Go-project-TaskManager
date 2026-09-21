package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	ServerPort  string
	FrontendURL string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env vars")
	}

	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "tasksdb"),
		JWTSecret:   getEnv("JWT_SECRET", "12345678901234567890123456789012"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
