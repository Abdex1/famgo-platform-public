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
	ServiceName      string
	ServicePort      string
	Environment      string
	LogLevel         string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	RedisHost        string
	RedisPort        string
	KafkaBrokers     string
	TransactionTimeout int
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
		ServiceName:      getEnv("SERVICE_NAME", "payment-service"),
		ServicePort:      getEnv("SERVICE_PORT", "3015"),
		Environment:      getEnv("SERVICE_ENV", "development"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "payment_user"),
		DBPassword:       getEnv("DB_PASSWORD", "payment_service_pwd_secure_2024"),
		DBName:           getEnv("DB_NAME", "famgo_payment_service"),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		KafkaBrokers:     getEnv("KAFKA_BROKERS", "localhost:9092"),
		TransactionTimeout: 120,
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
	db.SetMaxOpenConns(25)
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

// ProcessPayment handler - Process a payment transaction
func (app *App) ProcessPaymentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	rideID := r.FormValue("ride_id")
	userID := r.FormValue("user_id")
	amount := r.FormValue("amount")
	provider := r.FormValue("provider")

	if rideID == "" || userID == "" || amount == "" || provider == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"missing required fields"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"transaction_id":"txn_%d",
		"ride_id":"%s",
		"user_id":"%s",
		"amount":%s,
		"provider":"%s",
		"status":"success",
		"reference":"REF_%d",
		"timestamp":"%s"
	}`, time.Now().UnixNano(), rideID, userID, amount, provider, time.Now().UnixNano(), time.Now().Format(time.RFC3339))
}

// GetWallet handler - Get user wallet balance
func (app *App) GetWalletHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"user_id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"user_id":"%s",
		"balance":500.50,
		"currency":"ETB",
		"last_updated":"%s"
	}`, userID, time.Now().Format(time.RFC3339))
}

// AddMoney handler - Add money to wallet
func (app *App) AddMoneyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	userID := r.FormValue("user_id")
	amount := r.FormValue("amount")
	provider := r.FormValue("provider")

	if userID == "" || amount == "" || provider == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"missing required fields"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"success":true,
		"user_id":"%s",
		"amount":%s,
		"provider":"%s",
		"status":"completed",
		"new_balance":600.50,
		"timestamp":"%s"
	}`, userID, amount, provider, time.Now().Format(time.RFC3339))
}

// RefundPayment handler - Refund a payment
func (app *App) RefundPaymentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":"method not allowed"}`)
		return
	}

	transactionID := r.FormValue("transaction_id")
	reason := r.FormValue("reason")

	if transactionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"transaction_id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"success":true,
		"transaction_id":"%s",
		"refund_id":"ref_%d",
		"reason":"%s",
		"status":"processed",
		"timestamp":"%s"
	}`, transactionID, time.Now().UnixNano(), reason, time.Now().Format(time.RFC3339))
}

// GetTransactionHistory handler - Get transaction history
func (app *App) GetTransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"user_id required"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"user_id":"%s",
		"transactions":[
			{
				"id":"txn_1",
				"amount":45.50,
				"type":"ride",
				"status":"success",
				"date":"%s"
			},
			{
				"id":"txn_2",
				"amount":100.00,
				"type":"topup",
				"status":"success",
				"date":"%s"
			}
		],
		"total_transactions":2
	}`, userID, time.Now().Add(-24*time.Hour).Format(time.RFC3339), time.Now().Format(time.RFC3339))
}

// SetupRoutes sets up all routes
func (app *App) SetupRoutes() {
	app.Router.HandleFunc("/v1/health", app.HealthHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/payments/process", app.ProcessPaymentHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/wallets", app.GetWalletHandler).Methods(http.MethodGet)
	app.Router.HandleFunc("/v1/wallets/add-money", app.AddMoneyHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/payments/refund", app.RefundPaymentHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/v1/transactions", app.GetTransactionHistoryHandler).Methods(http.MethodGet)
	
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
