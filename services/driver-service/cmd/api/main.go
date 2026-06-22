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
	GPSUpdateInterval int
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
		ServiceName:       getEnv("SERVICE_NAME", "driver-service"),
		ServicePort:       getEnv("SERVICE_PORT", "3002"),
		Environment:       getEnv("SERVICE_ENV", "development"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "driver_user"),
		DBPassword:        getEnv("DB_PASSWORD", "driver_service_pwd_secure_2024"),
		DBName:            getEnv("DB_NAME", "famgo_driver_service"),
		RedisHost:         getEnv("REDIS_HOST", "localhost"),
		RedisPort:         getEnv("REDIS_PORT", "6379"),
		KafkaBrokers:      getEnv("KAFKA_BROKERS", "localhost:9092"),
		GPSUpdateInterval: 5,
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

// GetDriver handler - Get driver details
func (app *App) GetDriverHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	driverID := r.URL.Query().Get("id")
	if driverID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"driver id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"id":"%s",
		"name":"Driver Name",
		"phone":"+251999999999",
		"status":"online",
		"rating":4.8,
		"total_trips":250,
		"vehicle_type":"sedan",
		"license_number":"ET-XXXX-XXXX"
	}`, driverID)
}

// UpdateLocation handler - Update driver location (GPS)
func (app *App) UpdateLocationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	driverID := r.FormValue("driver_id")
	latitude := r.FormValue("latitude")
	longitude := r.FormValue("longitude")
	accuracy := r.FormValue("accuracy")

	if driverID == "" || latitude == "" || longitude == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"missing required fields"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"success":true,
		"driver_id":"%s",
		"latitude":%s,
		"longitude":%s,
		"accuracy":%s,
		"timestamp":"%s"
	}`, driverID, latitude, longitude, accuracy, time.Now().Format(time.RFC3339))
}

// AcceptRide handler - Driver accepts a ride
func (app *App) AcceptRideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	driverID := r.FormValue("driver_id")
	rideID := r.FormValue("ride_id")

	if driverID == "" || rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"missing required fields"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"success":true,
		"ride_id":"%s",
		"driver_id":"%s",
		"status":"accepted",
		"estimated_pickup_minutes":5,
		"timestamp":"%s"
	}`, rideID, driverID, time.Now().Format(time.RFC3339))
}

// GetMetrics handler - Get driver metrics/stats
func (app *App) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	driverID := r.URL.Query().Get("id")
	if driverID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"driver id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"driver_id":"%s",
		"total_trips":250,
		"today_trips":5,
		"rating":4.8,
		"acceptance_rate":92.5,
		"completion_rate":98.3,
		"earnings_today":245.50,
		"earnings_month":7450.00,
		"earnings_year":89400.00,
		"total_distance_km":15600.5,
		"timestamp":"%s"
	}`, driverID, time.Now().Format(time.RFC3339))
}

// GoOffline handler - Driver goes offline
func (app *App) GoOfflineHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	driverID := r.FormValue("driver_id")
	if driverID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"driver id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"success":true,
		"driver_id":"%s",
		"status":"offline",
		"timestamp":"%s"
	}`, driverID, time.Now().Format(time.RFC3339))
}

// SetupRoutes sets up all routes
func (app *App) SetupRoutes() {
	app.Router.HandleFunc("/v1/health", app.HealthHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/drivers", app.GetDriverHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/drivers/location", app.UpdateLocationHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/drivers/accept-ride", app.AcceptRideHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/drivers/metrics", app.GetMetricsHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/drivers/offline", app.GoOfflineHandler).Methods(http.MethodPost)
	
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
