package config

import (
	"fmt"
	"os"
)

// Config holds all service configuration
type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Auth
	JWTSecret       string
	JWTExpiry       int // minutes
	RefreshExpiry   int // days
	OTPExpiry       int // minutes

	// Email (Brevo)
	BrevoAPIKey string
	BrevoSender string

	// Logging
	LogLevel string

	// Observability
	OTelEndpoint string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "famgo"),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "famgo"),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		JWTExpiry:       15, // 15 minutes
		RefreshExpiry:   7,  // 7 days
		OTPExpiry:       10, // 10 minutes
		BrevoAPIKey:     getEnv("BREVO_API_KEY", ""),
		BrevoSender:     getEnv("BREVO_SENDER", "noreply@famgo.local"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		OTelEndpoint:    getEnv("OTEL_ENDPOINT", "localhost:4317"),
	}

	// Validate required fields
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}
	if cfg.BrevoAPIKey == "" {
		return nil, fmt.Errorf("BREVO_API_KEY environment variable is required")
	}

	return cfg, nil
}

// getEnv retrieves environment variable or default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
