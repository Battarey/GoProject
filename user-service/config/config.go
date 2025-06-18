package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	JWTSecret string
}

func LoadConfig() *Config {
	// Загружаем .env файл, если он есть
	_ = godotenv.Load()

	cfg := &Config{
		DBUrl:     os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if cfg.DBUrl == "" || cfg.JWTSecret == "" {
		log.Fatal("DB_URL и JWT_SECRET должны быть заданы в .env или переменных окружения")
	}

	return cfg
}
