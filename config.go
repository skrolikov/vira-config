package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит все параметры конфигурации приложения
type Config struct {
	// Database
	DBUrl             string        `json:"db_url" env:"DB_URL"`
	DevPostgresDSN    string        `json:"dev_postgres_dsn" env:"DEV_POSTGRES_DSN"`
	WishPostgresDSN   string        `json:"wish_postgres_dsn" env:"WISH_POSTGRES_DSN"`
	DBMaxOpenConns    int           `json:"db_max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns    int           `json:"db_max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	DBConnMaxLifetime time.Duration `json:"db_conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME"`
	DBConnMaxIdleTime time.Duration `json:"db_conn_max_idle_time" env:"DB_CONN_MAX_IDLE_TIME"`

	// Server
	Port            string        `json:"port" env:"PORT"`
	DevPort         string        `json:"dev_port" env:"DEV_PORT"`
	WishPort        string        `json:"wish_port" env:"WISH_PORT"`
	ReadTimeout     time.Duration `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `json:"write_timeout" env:"WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `json:"idle_timeout" env:"IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`

	// JWT
	JwtSecret     string        `json:"-" env:"JWT_SECRET"` // Скрываем в логах
	JwtTTL        time.Duration `json:"jwt_ttl" env:"JWT_TTL"`
	JwtRefreshTTL time.Duration `json:"jwt_refresh_ttl" env:"JWT_REFRESH_TTL"`
	JwtIssuer     string        `json:"jwt_issuer" env:"JWT_ISSUER"`

	// Redis
	RedisAddr     string `json:"redis_addr" env:"REDIS_ADDR"`
	RedisDB       int    `json:"redis_db" env:"REDIS_DB"`
	RedisPassword string `json:"-" env:"REDIS_PASSWORD"` // Скрываем в логах
	RedisPoolSize int    `json:"redis_pool_size" env:"REDIS_POOL_SIZE"`

	// Kafka
	KafkaAddr          string `json:"kafka_addr" env:"KAFKA_ADDR"`
	KafkaConsumerGroup string `json:"kafka_consumer_group" env:"KAFKA_CONSUMER_GROUP"`

	// External services
	ViraIDEndpoint string `json:"vira_id_endpoint" env:"VIRA_ID_ENDPOINT"`

	// Feature flags
	EnableDebug   bool `json:"enable_debug" env:"ENABLE_DEBUG"`
	EnableSwagger bool `json:"enable_swagger" env:"ENABLE_SWAGGER"`

	// Logging
	LogLevel  string `json:"log_level" env:"LOG_LEVEL"`
	LogFormat string `json:"log_format" env:"LOG_FORMAT"`
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	// Загружаем .env файл если существует (для разработки)
	_ = godotenv.Load(".env", ".env.local")

	cfg := &Config{
		// Database defaults
		DBMaxOpenConns:    10,
		DBMaxIdleConns:    5,
		DBConnMaxLifetime: 30 * time.Minute,
		DBConnMaxIdleTime: 5 * time.Minute,

		// Server defaults
		Port:            "8080",
		WishPort:        "8082",
		DevPort:         "8083",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     30 * time.Second,
		ShutdownTimeout: 5 * time.Second,

		// JWT defaults
		JwtTTL:        15 * time.Minute,
		JwtRefreshTTL: 7 * 24 * time.Hour,
		JwtIssuer:     "vira-api",

		// Redis defaults
		RedisAddr:     "redis:6379",
		RedisDB:       0,
		RedisPoolSize: 10,

		// Kafka defaults
		KafkaAddr:          "redpanda:9092",
		KafkaConsumerGroup: "vira-api-group",

		// Logging defaults
		LogLevel:  "info",
		LogFormat: "json",
	}

	// Загружаем обязательные переменные
	cfg.DBUrl = mustGetEnv("DB_URL")
	cfg.JwtSecret = mustGetEnv("JWT_SECRET")

	// Загружаем опциональные переменные
	loadOptionalEnvs(cfg)

	// Валидация конфигурации
	validateConfig(cfg)

	logConfig(cfg)
	return cfg
}

// loadOptionalEnvs загружает необязательные переменные окружения
func loadOptionalEnvs(cfg *Config) {
	// Database
	cfg.DevPostgresDSN = getEnv("DEV_POSTGRES_DSN", cfg.DBUrl)
	cfg.WishPostgresDSN = getEnv("WISH_POSTGRES_DSN", cfg.DBUrl)
	cfg.DBMaxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", cfg.DBMaxOpenConns)
	cfg.DBMaxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", cfg.DBMaxIdleConns)
	cfg.DBConnMaxLifetime = getEnvAsDuration("DB_CONN_MAX_LIFETIME", cfg.DBConnMaxLifetime)
	cfg.DBConnMaxIdleTime = getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", cfg.DBConnMaxIdleTime)

	// Server
	cfg.Port = getEnv("PORT", cfg.Port)
	cfg.DevPort = getEnv("DEV_PORT", cfg.DevPort)
	cfg.WishPort = getEnv("WISH_PORT", cfg.WishPort)
	cfg.ReadTimeout = getEnvAsDuration("READ_TIMEOUT", cfg.ReadTimeout)
	cfg.WriteTimeout = getEnvAsDuration("WRITE_TIMEOUT", cfg.WriteTimeout)
	cfg.IdleTimeout = getEnvAsDuration("IDLE_TIMEOUT", cfg.IdleTimeout)
	cfg.ShutdownTimeout = getEnvAsDuration("SHUTDOWN_TIMEOUT", cfg.ShutdownTimeout)

	// JWT
	cfg.JwtTTL = getEnvAsDuration("JWT_TTL", cfg.JwtTTL)
	cfg.JwtRefreshTTL = getEnvAsDuration("JWT_REFRESH_TTL", cfg.JwtRefreshTTL)
	cfg.JwtIssuer = getEnv("JWT_ISSUER", cfg.JwtIssuer)

	// Redis
	cfg.RedisAddr = getEnv("REDIS_ADDR", cfg.RedisAddr)
	cfg.RedisDB = getEnvAsInt("REDIS_DB", cfg.RedisDB)
	cfg.RedisPassword = getEnv("REDIS_PASSWORD", "")
	cfg.RedisPoolSize = getEnvAsInt("REDIS_POOL_SIZE", cfg.RedisPoolSize)

	// Kafka
	cfg.KafkaAddr = getEnv("KAFKA_ADDR", cfg.KafkaAddr)
	cfg.KafkaConsumerGroup = getEnv("KAFKA_CONSUMER_GROUP", cfg.KafkaConsumerGroup)

	// External services
	cfg.ViraIDEndpoint = getEnv("VIRA_ID_ENDPOINT", "")

	// Feature flags
	cfg.EnableDebug = getEnvAsBool("ENABLE_DEBUG", false)
	cfg.EnableSwagger = getEnvAsBool("ENABLE_SWAGGER", false)

	// Logging
	cfg.LogLevel = getEnv("LOG_LEVEL", cfg.LogLevel)
	cfg.LogFormat = getEnv("LOG_FORMAT", cfg.LogFormat)
}

// validateConfig проверяет валидность конфигурации
func validateConfig(cfg *Config) {
	if cfg.JwtTTL >= cfg.JwtRefreshTTL {
		log.Fatal("❌ JWT refresh TTL должен быть больше чем JWT TTL")
	}

	if cfg.DBMaxIdleConns > cfg.DBMaxOpenConns {
		log.Fatal("❌ DB_MAX_IDLE_CONNS не может быть больше DB_MAX_OPEN_CONNS")
	}
}

// logConfig логирует загруженную конфигурацию (без секретов)
func logConfig(cfg *Config) {
	log.Println("✅ Конфигурация загружена:")
	log.Printf("• Порт сервера: %s (dev: %s, wish: %s)", cfg.Port, cfg.DevPort, cfg.WishPort)
	log.Printf("• Таймауты: read=%v, write=%v, idle=%v", cfg.ReadTimeout, cfg.WriteTimeout, cfg.IdleTimeout)
	log.Printf("• JWT: TTL=%v, RefreshTTL=%v, Issuer=%s", cfg.JwtTTL, cfg.JwtRefreshTTL, cfg.JwtIssuer)
	log.Printf("• Redis: %s (DB: %d, Pool: %d)", cfg.RedisAddr, cfg.RedisDB, cfg.RedisPoolSize)
	log.Printf("• Kafka: %s (Group: %s)", cfg.KafkaAddr, cfg.KafkaConsumerGroup)
	log.Printf("• DB: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v",
		cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
}

// mustGetEnv возвращает значение переменной окружения или паникует если не задана
func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("❌ Обязательная переменная окружения не задана: %s", key)
	}
	return val
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// getEnvAsInt возвращает целочисленное значение переменной окружения
func getEnvAsInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			log.Printf("⚠️ Неверное значение %s=%s, используем значение по умолчанию %d", key, val, fallback)
			return fallback
		}
		return i
	}
	return fallback
}

// getEnvAsDuration возвращает значение переменной окружения как time.Duration
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		d, err := time.ParseDuration(val)
		if err != nil {
			log.Printf("⚠️ Неверное значение %s=%s, используем значение по умолчанию %v", key, val, fallback)
			return fallback
		}
		return d
	}
	return fallback
}

// getEnvAsBool возвращает булево значение переменной окружения
func getEnvAsBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		switch strings.ToLower(val) {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		default:
			log.Printf("⚠️ Неверное значение %s=%s, используем значение по умолчанию %v", key, val, fallback)
		}
	}
	return fallback
}
