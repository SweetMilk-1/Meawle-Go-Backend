package config

import (
	"os"
	"strconv"
)

// Config представляет конфигурацию приложения
type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
	LogLevel  string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", ":8080"),
		DBPath:    getEnv("DB_PATH", "app.db"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv получает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt получает целочисленное значение переменной окружения или значение по умолчанию
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}