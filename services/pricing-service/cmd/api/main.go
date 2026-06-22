package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

// AppConfig holds application configuration
type AppConfig struct {
	ServiceName    string
	ServicePort    string
	Environment    string
	LogLevel       string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	RedisHost      string
	RedisPort      string
	KafkaBrokers   string
}

// App holds application state
type App struct {
	Config *AppConfig
	DB     *sql.DB
	Router *mux.Router
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *AppConfig {
	return &AppConfig{
		ServiceName:   getEnv("SERVICE_NAME", "pricing-service"),
		ServicePort:   getEnv("SERVICE_PORT", "3014"),
		Environment:   getEnv("SERVICE_ENV", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "pricing_user"),
		DBPassword:    getEnv("DB_PASSWORD", "pricing_service_pwd_secure_2024"),
		DBName:        getEnv("DB_NAME", "famgo_pricing_service"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		KafkaBrokers:  getEnv("KAFKA_BROKERS", "localhost:9092"),
	}
}

// getEnv gets an environment variable or returns a default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// InitDB initializes database connection
func InitDB(config *AppConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("✓ Connected to database: %s@%s:%s/%s", config.DBUser, config.DBHost, config.DBPort, config.DBName)
	return db, nil
}

// Health check handler
func (app *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.DB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `{"status":"unhealthy","error":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"status":"healthy",
		"service":"%s",
		"environment":"%s",
		"version":"1.0.0",
		"timestamp":"%s"
	}`, app.Config.ServiceName, app.Config.Environment, time.Now().Format(time.RFC3339))
}

// EstimatePrice handler
func (app *App) EstimatePriceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	var request struct {
		RideType        string  `json:"ride_type"`
		DistanceMeters  float64 `json:"distance_meters"`
		ActiveRides     int     `json:"active_rides"`
		AvailableDriver int     `json:"available_drivers"`
		IsPool          bool    `json:"is_pool"`
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"failed to parse request"}`)
		return
	}

	// Parse JSON from body
	r.Header.Set("Content-Type", "application/json")
	if err := r.ParseForm(); err == nil {
		// Pricing calculation logic
		baseFare := 2.0
		distanceFare := request.DistanceMeters * 0.0012
		surgeFare := 1.0
		
		if request.ActiveRides > 0 && request.AvailableDriver > 0 {
			surgeFare = 1.0 + (float64(request.ActiveRides) / float64(request.AvailableDriver) * 0.25)
		}

		if request.IsPool {
			surgeFare = surgeFare * 0.7 // 30% discount for pooled rides
		}

		totalFare := (baseFare + distanceFare) * surgeFare

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"base_fare":%.2f,
			"distance_fare":%.2f,
			"surge_multiplier":%.2f,
			"total_fare":%.2f,
			"is_surge":"%v"
		}`, baseFare, distanceFare, surgeFare, totalFare, surgeFare > 1.1)
	}
}

// SetupRoutes sets up all routes
func (app *App) SetupRoutes() {
	app.Router.HandleFunc("/v1/health", app.HealthHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/pricing/estimate", app.EstimatePriceHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/metrics", app.MetricsHandler).Methods(http.MethodGet)
	
	log.Println("✓ Routes configured")
}

// MetricsHandler returns service metrics
func (app *App) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"service":"%s",
		"uptime_seconds":0,
		"requests_total":0,
		"requests_success":0,
		"requests_error":0
	}`, app.Config.ServiceName)
}

// Start starts the HTTP server
func (app *App) Start() error {
	server := &http.Server{
		Addr:         ":" + app.Config.ServicePort,
		Handler:      app.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Starting %s on port %s (%s environment)\n", 
			app.Config.ServiceName, app.Config.ServicePort, app.Config.Environment)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("\n📨 Received signal: %v\n", sig)
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
		return err
	}

	if err := app.DB.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	log.Println("✓ Service stopped successfully")
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load config
	config := LoadConfig()
	log.Printf("📋 Loading configuration from environment (%s)\n", config.Environment)

	// Initialize database
	db, err := InitDB(config)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v\n", err)
	}
	defer db.Close()

	// Create application
	app := &App{
		Config: config,
		DB:     db,
		Router: mux.NewRouter(),
	}

	// Setup routes
	app.SetupRoutes()

	// Start server
	if err := app.Start(); err != nil {
		log.Fatalf("❌ Failed to start application: %v\n", err)
	}
}
