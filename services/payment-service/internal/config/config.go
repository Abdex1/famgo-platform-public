// services/payment-service/internal/config/config.go
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServiceName string
	Environment string
	LogLevel    string
	Version     string
	GRPCPort    string
	GRPCHost    string

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

	RedisURL    string
	RedisDB     int
	KafkaBrokers     []string
	KafkaGroupID     string
	KafkaTopicPrefix string

	// Payment Providers
	TelebirrAPIKey       string
	TelebirrAPISecret    string
	TelebirrBaseURL      string
	CBEBirrAPIKey        string
	CBEBirrAPISecret     string
	CBEBirrBaseURL       string
	ChapaAPIKey          string
	ChapaAPISecret       string
	ChapaBaseURL         string
	WebhookSecret        string
	WebhookURL           string

	// Payment Configuration
	MinPaymentAmount     float64
	MaxPaymentAmount     float64
	PaymentTimeoutSec    int
	RefundTimeoutSec     int
	MaxRetries           int
	RetryDelayMs         int

	JWTSecret    string
	JWTIssuer    string
	JaegerURL    string
	MetricsPort  int
	RequestTimeout  time.Duration
	ShutdownTimeout time.Duration
	AuditEnabled bool
	AuditLogTTL  time.Duration
}

func Load() *Config {
	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "payment-service"),
		Environment: getEnv("ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Version:     getEnv("SERVICE_VERSION", "1.0.0"),
		GRPCPort:    getEnv("GRPC_PORT", "5006"),
		GRPCHost:    getEnv("GRPC_HOST", "0.0.0.0"),

		DatabaseHost:            getEnv("DB_HOST", "localhost"),
		DatabasePort:            getEnvInt("DB_PORT", 5432),
		DatabaseUser:            getEnv("DB_USER", "app_user"),
		DatabasePassword:        getEnv("DB_PASSWORD", "app_password"),
		DatabaseName:            getEnv("DB_NAME", "famgo_platform"),
		DatabaseSSLMode:         getEnv("DB_SSL_MODE", "disable"),
		DatabaseMaxConnections:  int32(getEnvInt("DB_MAX_CONNECTIONS", 32)),
		DatabaseMinConnections:  int32(getEnvInt("DB_MIN_CONNECTIONS", 10)),
		DatabaseConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 15*time.Minute),
		DatabaseConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),

		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisDB:     getEnvInt("REDIS_DB", 0),
		KafkaBrokers:     strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaGroupID:     getEnv("KAFKA_GROUP_ID", "payment-service"),
		KafkaTopicPrefix: getEnv("KAFKA_TOPIC_PREFIX", "famgo"),

		TelebirrAPIKey:   getEnv("TELEBIRR_API_KEY", ""),
		TelebirrAPISecret: getEnv("TELEBIRR_API_SECRET", ""),
		TelebirrBaseURL:  getEnv("TELEBIRR_BASE_URL", "https://api.telebirr.com"),
		CBEBirrAPIKey:    getEnv("CBE_BIRR_API_KEY", ""),
		CBEBirrAPISecret: getEnv("CBE_BIRR_API_SECRET", ""),
		CBEBirrBaseURL:   getEnv("CBE_BIRR_BASE_URL", "https://api.cbebirr.com"),
		ChapaAPIKey:      getEnv("CHAPA_API_KEY", ""),
		ChapaAPISecret:   getEnv("CHAPA_API_SECRET", ""),
		ChapaBaseURL:     getEnv("CHAPA_BASE_URL", "https://api.chapa.co"),
		WebhookSecret:    getEnv("WEBHOOK_SECRET", ""),
		WebhookURL:       getEnv("WEBHOOK_URL", "https://api.famgo.et/webhooks/payment"),

		MinPaymentAmount:     getEnvFloat("MIN_PAYMENT_AMOUNT", 10.0),
		MaxPaymentAmount:     getEnvFloat("MAX_PAYMENT_AMOUNT", 100000.0),
		PaymentTimeoutSec:    getEnvInt("PAYMENT_TIMEOUT_SEC", 30),
		RefundTimeoutSec:     getEnvInt("REFUND_TIMEOUT_SEC", 60),
		MaxRetries:           getEnvInt("MAX_RETRIES", 3),
		RetryDelayMs:         getEnvInt("RETRY_DELAY_MS", 1000),

		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-min-32-characters"),
		JWTIssuer:    getEnv("JWT_ISSUER", "famgo-platform"),
		JaegerURL:    getEnv("JAEGER_URL", "http://localhost:14268/api/traces"),
		MetricsPort:  getEnvInt("METRICS_PORT", 9093),
		RequestTimeout:  getEnvDuration("REQUEST_TIMEOUT", 30*time.Second),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		AuditEnabled: getEnvBool("AUDIT_ENABLED", true),
		AuditLogTTL:  getEnvDuration("AUDIT_LOG_TTL", 90*24*time.Hour),
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

func getEnvBool(key string, fallback bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func getEnvFloat(key string, fallback float64) float64 {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseFloat(valStr, 64)
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
