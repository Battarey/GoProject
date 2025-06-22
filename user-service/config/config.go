package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	JWTSecret string
	SMTPHost  string
	SMTPPort  string
	SMTPUser  string
	SMTPPass  string
	FromEmail string
}

func LoadConfig() *Config {
	// Загружаем .env файл, если он есть
	_ = godotenv.Load()

	cfg := &Config{
		DBUrl:     os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		SMTPHost:  os.Getenv("SMTP_HOST"),
		SMTPPort:  os.Getenv("SMTP_PORT"),
		SMTPUser:  os.Getenv("SMTP_USER"),
		SMTPPass:  os.Getenv("SMTP_PASS"),
		FromEmail: os.Getenv("FROM_EMAIL"),
	}

	if cfg.DBUrl == "" || cfg.JWTSecret == "" {
		log.Fatal("DB_URL и JWT_SECRET должны быть заданы в .env или переменных окружения")
	}

	return cfg
}
