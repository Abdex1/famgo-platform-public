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
	ServiceName       string
	ServicePort       string
	Environment       string
	LogLevel          string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	RedisHost         string
	RedisPort         string
	KafkaBrokers      string
	RideTimeoutMinutes int
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
		ServiceName:       getEnv("SERVICE_NAME", "ride-service"),
		ServicePort:       getEnv("SERVICE_PORT", "3010"),
		Environment:       getEnv("SERVICE_ENV", "development"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "ride_user"),
		DBPassword:        getEnv("DB_PASSWORD", "ride_service_pwd_secure_2024"),
		DBName:            getEnv("DB_NAME", "famgo_ride_service"),
		RedisHost:         getEnv("REDIS_HOST", "localhost"),
		RedisPort:         getEnv("REDIS_PORT", "6379"),
		KafkaBrokers:      getEnv("KAFKA_BROKERS", "localhost:9092"),
		RideTimeoutMinutes: 30,
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

// CreateRide handler - Create a new ride
func (app *App) CreateRideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	userID := r.FormValue("user_id")
	pickupLat := r.FormValue("pickup_lat")
	pickupLng := r.FormValue("pickup_lng")
	dropoffLat := r.FormValue("dropoff_lat")
	dropoffLng := r.FormValue("dropoff_lng")
	rideType := r.FormValue("ride_type")

	if userID == "" || pickupLat == "" || pickupLng == "" || dropoffLat == "" || dropoffLng == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"missing required fields"}`)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{
		"ride_id":"ride_%d",
		"user_id":"%s",
		"status":"searching",
		"pickup_location":{"lat":%s,"lng":%s},
		"dropoff_location":{"lat":%s,"lng":%s},
		"ride_type":"%s",
		"estimated_fare":45.50,
		"created_at":"%s"
	}`, time.Now().UnixNano(), userID, pickupLat, pickupLng, dropoffLat, dropoffLng, rideType, time.Now().Format(time.RFC3339))
}

// GetRide handler - Get ride details
func (app *App) GetRideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rideID := r.URL.Query().Get("id")
	if rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"ride id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"ride_id":"%s",
		"user_id":"user_123",
		"driver_id":"driver_456",
		"status":"in_progress",
		"pickup_location":{"lat":9.0320,"lng":38.7469},
		"dropoff_location":{"lat":9.0265,"lng":38.7400},
		"current_location":{"lat":9.0300,"lng":38.7450},
		"ride_type":"economy",
		"fare":45.50,
		"created_at":"%s"
	}`, rideID, time.Now().Add(-30*time.Minute).Format(time.RFC3339))
}

// CancelRide handler - Cancel a ride
func (app *App) CancelRideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	rideID := r.FormValue("ride_id")
	reason := r.FormValue("reason")

	if rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"ride_id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"success":true,
		"ride_id":"%s",
		"status":"cancelled",
		"reason":"%s",
		"timestamp":"%s"
	}`, rideID, reason, time.Now().Format(time.RFC3339))
}

// CompleteRide handler - Complete a ride
func (app *App) CompleteRideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	rideID := r.FormValue("ride_id")
	if rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"ride_id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"success":true,
		"ride_id":"%s",
		"status":"completed",
		"distance_km":12.5,
		"final_fare":45.50,
		"duration_minutes":18,
		"completed_at":"%s"
	}`, rideID, time.Now().Format(time.RFC3339))
}

// RateRide handler - Rate a completed ride
func (app *App) RateRideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	rideID := r.FormValue("ride_id")
	rating := r.FormValue("rating")
	comment := r.FormValue("comment")

	if rideID == "" || rating == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"ride_id and rating required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"success":true,
		"ride_id":"%s",
		"rating":%s,
		"comment":"%s",
		"timestamp":"%s"
	}`, rideID, rating, comment, time.Now().Format(time.RFC3339))
}

// SetupRoutes sets up all routes
func (app *App) SetupRoutes() {
	app.Router.HandleFunc("/v1/health", app.HealthHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/rides", app.CreateRideHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/rides", app.GetRideHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/rides/cancel", app.CancelRideHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/rides/complete", app.CompleteRideHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/rides/rate", app.RateRideHandler).Methods(http.MethodPost)
	
	log.Println("✓ Routes configured")
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
