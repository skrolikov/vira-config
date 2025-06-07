package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl         string
	JwtSecret     string
	Port          string
	JwtTTL        time.Duration
	JwtRefreshTTL time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		DBUrl:         os.Getenv("DB_URL"),
		JwtSecret:     os.Getenv("JWT_SECRET"),
		Port:          port,
		JwtTTL:        15 * time.Minute,
		JwtRefreshTTL: 7 * 24 * time.Hour,
	}

	log.Println("Конфиг загружен.")
	return cfg
}
