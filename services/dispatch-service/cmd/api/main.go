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
	MatchingTimeout   int
	SearchRadius      int
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
		ServiceName:     getEnv("SERVICE_NAME", "dispatch-service"),
		ServicePort:     getEnv("SERVICE_PORT", "3011"),
		Environment:     getEnv("SERVICE_ENV", "development"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "dispatch_user"),
		DBPassword:      getEnv("DB_PASSWORD", "dispatch_service_pwd_secure_2024"),
		DBName:          getEnv("DB_NAME", "famgo_dispatch_service"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		KafkaBrokers:    getEnv("KAFKA_BROKERS", "localhost:9092"),
		MatchingTimeout: 30,
		SearchRadius:    5000,
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

// MatchDrivers handler - Find matching drivers for a ride
func (app *App) MatchDriversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	rideID := r.FormValue("ride_id")
	pickupLat := r.FormValue("pickup_lat")
	pickupLng := r.FormValue("pickup_lng")
	rideType := r.FormValue("ride_type")

	if rideID == "" || pickupLat == "" || pickupLng == "" || rideType == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"missing required fields"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"ride_id":"%s",
		"matching_drivers":[
			{
				"driver_id":"driver_1",
				"distance_meters":450,
				"estimated_arrival_seconds":120,
				"rating":4.9,
				"total_trips":500
			},
			{
				"driver_id":"driver_2",
				"distance_meters":680,
				"estimated_arrival_seconds":180,
				"rating":4.7,
				"total_trips":350
			},
			{
				"driver_id":"driver_3",
				"distance_meters":920,
				"estimated_arrival_seconds":240,
				"rating":4.6,
				"total_trips":200
			}
		],
		"algorithm_version":"v1.0.0",
		"matching_timestamp":"%s"
	}`, rideID, time.Now().Format(time.RFC3339))
}

// AssignDriver handler - Assign a driver to a ride
func (app *App) AssignDriverHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	rideID := r.FormValue("ride_id")
	driverID := r.FormValue("driver_id")

	if rideID == "" || driverID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"missing required fields"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"success":true,
		"ride_id":"%s",
		"driver_id":"%s",
		"status":"assigned",
		"assignment_timestamp":"%s"
	}`, rideID, driverID, time.Now().Format(time.RFC3339))
}

// CancelDispatch handler - Cancel dispatch/matching
func (app *App) CancelDispatchHandler(w http.ResponseWriter, r *http.Request) {
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
		"status":"dispatch_cancelled",
		"reason":"%s",
		"timestamp":"%s"
	}`, rideID, reason, time.Now().Format(time.RFC3339))
}

// GetDispatchStatus handler - Get dispatch status
func (app *App) GetDispatchStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rideID := r.URL.Query().Get("ride_id")
	if rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"ride_id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"ride_id":"%s",
		"status":"driver_assigned",
		"assigned_driver":"driver_1",
		"matching_duration_seconds":15,
		"search_radius_meters":5000,
		"last_update":"%s"
	}`, rideID, time.Now().Format(time.RFC3339))
}

// GetMetrics handler - Get dispatch metrics
func (app *App) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"total_dispatches":10500,
		"successful_matches":10200,
		"failed_matches":300,
		"average_matching_time_seconds":18,
		"success_rate":97.14,
		"active_rides_being_matched":25,
		"timestamp":"%s"
	}`, time.Now().Format(time.RFC3339))
}

// SetupRoutes sets up all routes
func (app *App) SetupRoutes() {
	app.Router.HandleFunc("/v1/health", app.HealthHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/dispatch/match", app.MatchDriversHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/dispatch/assign", app.AssignDriverHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/dispatch/cancel", app.CancelDispatchHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/dispatch/status", app.GetDispatchStatusHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/dispatch/metrics", app.GetMetricsHandler).Methods(http.MethodGet)
	
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
