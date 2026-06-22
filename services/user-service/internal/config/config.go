// services/user-service/internal/config/config.go
// Configuration Loading

package config

import (
	"os"
	"strconv"
	"time"
)

// Config represents application configuration
type Config struct {
	// Database
	DatabaseHost            string
	DatabasePort            int
	DatabaseUser            string
	DatabasePassword        string
	DatabaseName            string
	DatabaseSSLMode         string
	DatabaseMaxConnections  int32
	DatabaseMinConnections  int32
	DatabaseConnMaxLifetime time.Duration
	DatabaseConnMaxIdleTime time.Duration

	// Redis
	RedisURL  string
	RedisDB   int
	CacheTTL  int32

	// Server
	HTTPPort string
	GRPCPort string

	// Logging
	LogLevel string

	// Shutdown
	ShutdownTimeout time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		// Database
		DatabaseHost:            getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:            getEnvInt("DATABASE_PORT", 5432),
		DatabaseUser:            getEnv("DATABASE_USER", "user"),
		DatabasePassword:        getEnv("DATABASE_PASSWORD", "password"),
		DatabaseName:            getEnv("DATABASE_NAME", "user_db"),
		DatabaseSSLMode:         getEnv("DATABASE_SSL_MODE", "disable"),
		DatabaseMaxConnections:  int32(getEnvInt("DATABASE_MAX_CONNECTIONS", 25)),
		DatabaseMinConnections:  int32(getEnvInt("DATABASE_MIN_CONNECTIONS", 5)),
		DatabaseConnMaxLifetime: getEnvDuration("DATABASE_CONN_MAX_LIFETIME", 1*time.Hour),
		DatabaseConnMaxIdleTime: getEnvDuration("DATABASE_CONN_MAX_IDLE_TIME", 10*time.Minute),

		// Redis
		RedisURL: getEnv("REDIS_URL", "localhost:6379"),
		RedisDB:  getEnvInt("REDIS_DB", 0),
		CacheTTL: int32(getEnvInt("CACHE_TTL", 3600)),

		// Server
		HTTPPort: getEnv("HTTP_PORT", "5003"),
		GRPCPort: getEnv("GRPC_PORT", "5004"),

		// Logging
		LogLevel: getEnv("LOG_LEVEL", "info"),

		// Shutdown
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
	}
}

// Helper functions
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultVal
}
