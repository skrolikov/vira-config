package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl          string // для vira-id и других, кто использует DBUrl
	PostgresDSN    string // приоритетный параметр для Postgres (vira-dev)
	ViraIDEndpoint string // URL vira-id (vira-dev)
	KafkaAddr      string
	JwtSecret      string
	Port           string
	DevPort        string
	JwtTTL         time.Duration
	JwtRefreshTTL  time.Duration
	RedisAddr      string
	RedisDB        int
}

func Load() *Config {
	_ = godotenv.Load(".env") // загружаем .env, если есть

	cfg := &Config{}

	// Существующие обязательные
	cfg.DBUrl = mustEnv("DB_URL")
	cfg.JwtSecret = mustEnv("JWT_SECRET")

	// Новые: если POSTGRES_DSN не задан, падаем на DBUrl
	cfg.PostgresDSN = getEnv("POSTGRES_DSN", cfg.DBUrl)
	// Для vira-dev:
	cfg.ViraIDEndpoint = getEnv("VIRA_ID_ENDPOINT", "")
	cfg.DevPort = getEnv("DEV_PORT", "8080")

	// Опциональные
	cfg.Port = getEnv("PORT", "8080")
	cfg.KafkaAddr = getEnv("KAFKA_ADDR", "redpanda:9092")
	cfg.RedisAddr = getEnv("REDIS_ADDR", "redis:6379")
	cfg.RedisDB = mustAtoi(getEnv("REDIS_DB", "0"))

	ttlMinutes := mustAtoi(getEnv("JWT_TTL_MINUTES", "15"))
	refreshDays := mustAtoi(getEnv("JWT_REFRESH_DAYS", "7"))
	cfg.JwtTTL = time.Duration(ttlMinutes) * time.Minute
	cfg.JwtRefreshTTL = time.Duration(refreshDays) * 24 * time.Hour

	log.Printf("✅ Конфиг загружен: PORT=%s, Redis=%s, TTL=%v\n", cfg.Port, cfg.RedisAddr, cfg.JwtTTL)
	return cfg
}

// mustEnv — паника если переменная отсутствует
func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("❌ Обязательная переменная окружения не задана: %s", key)
	}
	return val
}

// getEnv — с дефолтом
func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// mustAtoi — безопасное преобразование строк в int
func mustAtoi(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("❌ Не удалось преобразовать число: %s", val)
	}
	return i
}
