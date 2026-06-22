// services/dispatch-service/internal/config/config.go
// Dispatch Service configuration with matching algorithm parameters

package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all dispatch service configuration
type Config struct {
	// Service metadata
	ServiceName string
	Environment string
	LogLevel    string
	Version     string

	// Server
	GRPCPort string
	GRPCHost string
	HTTPPort string
	HTTPHost string

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
	RedisURL    string
	RedisDB     int
	RedisMaxAge time.Duration

	// Kafka
	KafkaEnabled     bool
	KafkaBrokers     []string
	KafkaGroupID     string
	KafkaTopicPrefix string
	KafkaConsumerLag time.Duration

	// Matching Algorithm Parameters
	ProximityWeight           float64       // Weight for distance (0.40 default)
	AcceptanceRateWeight      float64       // Weight for acceptance rate (0.30 default)
	RatingWeight              float64       // Weight for rating (0.20 default)
	AvailabilityWeight        float64       // Weight for online status (0.10 default)
	SearchRadiusKm            float64       // Initial search radius (5 km default)
	MaxSearchRadiusKm         float64       // Maximum search radius (25 km default)
	MinSearchRadiusKm         float64       // Minimum search radius (0.5 km default)
	NearbyDriversLimit        int           // Max drivers to consider (100 default)
	TopMatchesLimit           int           // Max matches to return (5 default)
	MinAcceptanceRatePercent  float64       // Minimum acceptance rate to consider (50%)
	MinRating                 float64       // Minimum rating to consider (3.5 stars)
	MatchExpirySeconds        int           // How long a match is valid (60s default)
	MatchRequestTTL           time.Duration // How long to store match requests
	MaxConcurrentMatches      int           // Max concurrent matching operations

	// Observability
	JaegerURL          string
	JaegerSamplerType  string
	JaegerSamplerParam float64
	MetricsPort        int

	// Timeouts
	RequestTimeout  time.Duration
	ShutdownTimeout time.Duration

	// Auth
	JWTSecret string
	JWTIssuer string

	// Audit
	AuditEnabled bool
	AuditLogTTL  time.Duration

	// External Services (gRPC endpoints)
	GPSServiceURL  string
	RideServiceURL string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		// Service metadata
		ServiceName: getEnv("SERVICE_NAME", "dispatch-service"),
		Environment: getEnv("ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Version:     getEnv("SERVICE_VERSION", "1.0.0"),

		// Server
		GRPCPort: getEnv("GRPC_PORT", "5005"),
		GRPCHost: getEnv("GRPC_HOST", "0.0.0.0"),
		HTTPPort: getEnv("HTTP_PORT", "8085"),
		HTTPHost: getEnv("HTTP_HOST", "0.0.0.0"),

		// Database
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

		// Redis
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisDB:     getEnvInt("REDIS_DB", 0),
		RedisMaxAge: getEnvDuration("REDIS_MAX_AGE", 24*time.Hour),

		// Kafka
		KafkaEnabled:     getEnvBool("KAFKA_ENABLED", false),
		KafkaBrokers:     strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaGroupID:     getEnv("KAFKA_GROUP_ID", "dispatch-service"),
		KafkaTopicPrefix: getEnv("KAFKA_TOPIC_PREFIX", "famgo"),
		KafkaConsumerLag: getEnvDuration("KAFKA_CONSUMER_LAG", 30*time.Second),

		// Matching Algorithm Parameters
		ProximityWeight:          getEnvFloat("PROXIMITY_WEIGHT", 0.40),
		AcceptanceRateWeight:     getEnvFloat("ACCEPTANCE_RATE_WEIGHT", 0.30),
		RatingWeight:             getEnvFloat("RATING_WEIGHT", 0.20),
		AvailabilityWeight:       getEnvFloat("AVAILABILITY_WEIGHT", 0.10),
		SearchRadiusKm:           getEnvFloat("SEARCH_RADIUS_KM", 5.0),
		MaxSearchRadiusKm:        getEnvFloat("MAX_SEARCH_RADIUS_KM", 25.0),
		MinSearchRadiusKm:        getEnvFloat("MIN_SEARCH_RADIUS_KM", 0.5),
		NearbyDriversLimit:       getEnvInt("NEARBY_DRIVERS_LIMIT", 100),
		TopMatchesLimit:          getEnvInt("TOP_MATCHES_LIMIT", 5),
		MinAcceptanceRatePercent: getEnvFloat("MIN_ACCEPTANCE_RATE_PERCENT", 50.0),
		MinRating:                getEnvFloat("MIN_RATING", 3.5),
		MatchExpirySeconds:       getEnvInt("MATCH_EXPIRY_SECONDS", 60),
		MatchRequestTTL:          getEnvDuration("MATCH_REQUEST_TTL", 5*time.Minute),
		MaxConcurrentMatches:     getEnvInt("MAX_CONCURRENT_MATCHES", 1000),

		// Observability
		JaegerURL:          getEnv("JAEGER_URL", "http://localhost:14268/api/traces"),
		JaegerSamplerType:  getEnv("JAEGER_SAMPLER_TYPE", "probabilistic"),
		JaegerSamplerParam: getEnvFloat("JAEGER_SAMPLER_PARAM", 0.1),
		MetricsPort:        getEnvInt("METRICS_PORT", 9092),

		// Timeouts
		RequestTimeout:  getEnvDuration("REQUEST_TIMEOUT", 30*time.Second),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),

		// Auth
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-min-32-characters"),
		JWTIssuer: getEnv("JWT_ISSUER", "famgo-platform"),

		// Audit
		AuditEnabled: getEnvBool("AUDIT_ENABLED", true),
		AuditLogTTL:  getEnvDuration("AUDIT_LOG_TTL", 90*24*time.Hour),

		// External Services
		GPSServiceURL:  getEnv("GPS_SERVICE_URL", "localhost:5002"),
		RideServiceURL: getEnv("RIDE_SERVICE_URL", "localhost:5004"),
	}
}

// Helper functions
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
