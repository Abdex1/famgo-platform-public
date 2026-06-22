package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port          string
	Environment   string
	LogLevel      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	OTelEndpoint  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8082"),
		Environment:  getEnv("ENV", "development"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "famgo_driver"),
		OTelEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
	}
	if cfg.DBHost == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
