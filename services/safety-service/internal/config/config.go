// services/safety-service/internal/config/config.go
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServiceName       string
	Environment       string
	LogLevel          string
	Version           string
	GRPCPort          string
	GRPCHost          string
	DatabaseHost      string
	DatabasePort      int
	DatabaseUser      string
	DatabasePassword  string
	DatabaseName      string
	DatabaseSSLMode   string
	DatabaseMaxConns  int32
	DatabaseMinConns  int32
	RedisURL          string
	KafkaBrokers      []string
	KafkaGroupID      string
	KafkaTopicPrefix  string
	SOSTimeoutSec     int
	EscalationDelaySec int
	JWTSecret         string
	JaegerURL         string
	ShutdownTimeout   time.Duration
}

func Load() *Config {
	return &Config{
		ServiceName:      getEnv("SERVICE_NAME", "safety-service"),
		Environment:      getEnv("ENV", "development"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		Version:          getEnv("SERVICE_VERSION", "1.0.0"),
		GRPCPort:         getEnv("GRPC_PORT", "5008"),
		GRPCHost:         getEnv("GRPC_HOST", "0.0.0.0"),
		DatabaseHost:     getEnv("DB_HOST", "localhost"),
		DatabasePort:     getEnvInt("DB_PORT", 5432),
		DatabaseUser:     getEnv("DB_USER", "app_user"),
		DatabasePassword: getEnv("DB_PASSWORD", "app_password"),
		DatabaseName:     getEnv("DB_NAME", "famgo_platform"),
		DatabaseSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		DatabaseMaxConns: int32(getEnvInt("DB_MAX_CONNECTIONS", 32)),
		DatabaseMinConns: int32(getEnvInt("DB_MIN_CONNECTIONS", 10)),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		KafkaBrokers:     strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaGroupID:     getEnv("KAFKA_GROUP_ID", "safety-service"),
		KafkaTopicPrefix: getEnv("KAFKA_TOPIC_PREFIX", "famgo"),
		SOSTimeoutSec:    getEnvInt("SOS_TIMEOUT_SEC", 300),
		EscalationDelaySec: getEnvInt("ESCALATION_DELAY_SEC", 30),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key-min-32-characters"),
		JaegerURL:        getEnv("JAEGER_URL", "http://localhost:14268/api/traces"),
		ShutdownTimeout:  getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		return fallback
	}
	return val
}
