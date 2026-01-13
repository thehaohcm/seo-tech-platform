package config

import (
	"os"
)

type Config struct {
	RedisURL    string
	DBHost      string
	DBPort      string
	DBName      string
	DBUser      string
	DBPassword  string
	Environment string
	LogLevel    string
}

func LoadConfig() *Config {
	return &Config{
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBName:      getEnv("DB_NAME", "seo_platform"),
		DBUser:      getEnv("DB_USER", "seouser"),
		DBPassword:  getEnv("DB_PASSWORD", "seopass"),
		Environment: getEnv("APP_ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
