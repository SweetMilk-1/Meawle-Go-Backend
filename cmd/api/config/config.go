package config

import (
	"log"
	"os"

	"meawle/internal/config"
)

// LoadConfig загружает конфигурацию приложения
func LoadConfig() *config.Config {
	cfg := config.Load()
	return cfg
}

// SetupLogger настраивает логгер приложения
func SetupLogger() *log.Logger {
	return log.New(os.Stdout, "[API] ", log.LstdFlags|log.Lshortfile)
}