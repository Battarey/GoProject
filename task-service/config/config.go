package config

import (
	"os"
)

type Config struct {
	DBUrl     string
	JWTSecret string
	Port      string
}

func LoadConfig() *Config {
	return &Config{
		DBUrl:     getEnv("DB_URL", "host=db user=user password=password dbname=tasks_db port=5432 sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "supersecretkey"),
		Port:      getEnv("TASK_SERVICE_PORT", "50052"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
