// services/gps-service/internal/config/config.go
// GPS Service configuration with geolocation settings and performance tuning

package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all GPS service configuration
type Config struct {
	// Service metadata
	ServiceName string
	Environment string
	LogLevel    string
	Version     string

	// Server
	GRPCPort string
	GRPCHost string

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
	KafkaBrokers     []string
	KafkaGroupID     string
	KafkaTopicPrefix string
	KafkaConsumerLag time.Duration

	// GPS Location Configuration
	LocationUpdateIntervalMs   int           // How often drivers send location updates (milliseconds)
	LocationHistoryRetention   time.Duration // How long to keep location history
	LocationAccuracyThreshold  float64       // Minimum accuracy in meters
	LocationCleanupInterval    time.Duration // Interval for cleaning stale locations

	// Geolocation Search
	SearchRadiusKm            float64       // Search radius for nearby drivers (km)
	SearchRadiusMax           float64       // Maximum search radius (km)
	SearchRadiusMin           float64       // Minimum search radius (km)
	NearbyDriversLimit        int           // Max drivers to return in nearby search
	NearbyDriversCacheTTL     time.Duration // Cache TTL for nearby drivers

	// Geohashing
	GeohashPrecision int // Geohash precision for indexing (1-12)

	// ETA Calculation
	ETACalculationMethod  string        // "haversine" or "routing"
	ETABaseSpeedKmPerHour float64       // Default speed for ETA calc (km/h)
	ETAUpdateInterval     time.Duration // How often to recalculate ETA

	// Driver Status
	DriverOnlineTimeout      time.Duration // Time before marking driver as offline
	DriverLocationExpiry     time.Duration // Location data expiration time
	DriverStatusUpdateTTL    time.Duration // TTL for driver status in Redis
	MaxConcurrentLocations   int           // Max concurrent location updates to process

	// Performance
	BatchLocationUpdateSize  int           // Batch size for bulk location updates
	LocationBatchWaitTime    time.Duration // Wait time before flushing location batch
	PrefixRefreshInterval    time.Duration // Interval to refresh location indices
	CacheUpdateBatchSize     int           // Batch size for cache updates

	// Observability
	JaegerURL          string
	JaegerSamplerType  string
	JaegerSamplerParam float64
	MetricsPort        int

	// Timeouts
	RequestTimeout  time.Duration
	ShutdownTimeout time.Duration

	// WebSocket (for real-time updates)
	WebSocketPort       string
	WebSocketMaxClients int
	WebSocketReadDeadline time.Duration
	WebSocketWriteDeadline time.Duration
	WebSocketPingInterval time.Duration

	// Auth
	JWTSecret string
	JWTIssuer string

	// Audit
	AuditEnabled bool
	AuditLogTTL  time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		// Service metadata
		ServiceName: getEnv("SERVICE_NAME", "gps-service"),
		Environment: getEnv("ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Version:     getEnv("SERVICE_VERSION", "1.0.0"),

		// Server
		GRPCPort: getEnv("GRPC_PORT", "5002"),
		GRPCHost: getEnv("GRPC_HOST", "0.0.0.0"),

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
		KafkaBrokers:     strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaGroupID:     getEnv("KAFKA_GROUP_ID", "gps-service"),
		KafkaTopicPrefix: getEnv("KAFKA_TOPIC_PREFIX", "famgo"),
		KafkaConsumerLag: getEnvDuration("KAFKA_CONSUMER_LAG", 30*time.Second),

		// GPS Location Configuration
		LocationUpdateIntervalMs:  getEnvInt("LOCATION_UPDATE_INTERVAL_MS", 5000),
		LocationHistoryRetention: getEnvDuration("LOCATION_HISTORY_RETENTION", 30*24*time.Hour),
		LocationAccuracyThreshold: getEnvFloat("LOCATION_ACCURACY_THRESHOLD", 50.0),
		LocationCleanupInterval:  getEnvDuration("LOCATION_CLEANUP_INTERVAL", 1*time.Hour),

		// Geolocation Search
		SearchRadiusKm:        getEnvFloat("SEARCH_RADIUS_KM", 5.0),
		SearchRadiusMax:       getEnvFloat("SEARCH_RADIUS_MAX_KM", 25.0),
		SearchRadiusMin:       getEnvFloat("SEARCH_RADIUS_MIN_KM", 0.5),
		NearbyDriversLimit:    getEnvInt("NEARBY_DRIVERS_LIMIT", 50),
		NearbyDriversCacheTTL: getEnvDuration("NEARBY_DRIVERS_CACHE_TTL", 30*time.Second),

		// Geohashing
		GeohashPrecision: getEnvInt("GEOHASH_PRECISION", 8),

		// ETA Calculation
		ETACalculationMethod:  getEnv("ETA_CALCULATION_METHOD", "haversine"),
		ETABaseSpeedKmPerHour: getEnvFloat("ETA_BASE_SPEED_KM_PER_HOUR", 40.0),
		ETAUpdateInterval:     getEnvDuration("ETA_UPDATE_INTERVAL", 30*time.Second),

		// Driver Status
		DriverOnlineTimeout:   getEnvDuration("DRIVER_ONLINE_TIMEOUT", 2*time.Minute),
		DriverLocationExpiry:  getEnvDuration("DRIVER_LOCATION_EXPIRY", 5*time.Minute),
		DriverStatusUpdateTTL: getEnvDuration("DRIVER_STATUS_UPDATE_TTL", 10*time.Minute),
		MaxConcurrentLocations: getEnvInt("MAX_CONCURRENT_LOCATIONS", 1000),

		// Performance
		BatchLocationUpdateSize: getEnvInt("BATCH_LOCATION_UPDATE_SIZE", 100),
		LocationBatchWaitTime:   getEnvDuration("LOCATION_BATCH_WAIT_TIME", 500*time.Millisecond),
		PrefixRefreshInterval:   getEnvDuration("PREFIX_REFRESH_INTERVAL", 5*time.Minute),
		CacheUpdateBatchSize:    getEnvInt("CACHE_UPDATE_BATCH_SIZE", 500),

		// Observability
		JaegerURL:          getEnv("JAEGER_URL", "http://localhost:14268/api/traces"),
		JaegerSamplerType:  getEnv("JAEGER_SAMPLER_TYPE", "probabilistic"),
		JaegerSamplerParam: getEnvFloat("JAEGER_SAMPLER_PARAM", 0.1),
		MetricsPort:        getEnvInt("METRICS_PORT", 9090),

		// Timeouts
		RequestTimeout:  getEnvDuration("REQUEST_TIMEOUT", 30*time.Second),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),

		// WebSocket
		WebSocketPort:         getEnv("WEBSOCKET_PORT", "5003"),
		WebSocketMaxClients:   getEnvInt("WEBSOCKET_MAX_CLIENTS", 10000),
		WebSocketReadDeadline: getEnvDuration("WEBSOCKET_READ_DEADLINE", 60*time.Second),
		WebSocketWriteDeadline: getEnvDuration("WEBSOCKET_WRITE_DEADLINE", 60*time.Second),
		WebSocketPingInterval: getEnvDuration("WEBSOCKET_PING_INTERVAL", 30*time.Second),

		// Auth
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-min-32-characters"),
		JWTIssuer: getEnv("JWT_ISSUER", "famgo-platform"),

		// Audit
		AuditEnabled: getEnvBool("AUDIT_ENABLED", true),
		AuditLogTTL:  getEnvDuration("AUDIT_LOG_TTL", 90*24*time.Hour),
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
