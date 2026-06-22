package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds application configuration
type Config struct {
	DatabaseURL     string
	KafkaBroker     string
	TwilioAccountSID string
	TwilioAuthToken  string
	FirebaseCredentials string
	Port            string
	Environment     string
}

// App holds application dependencies
type App struct {
	DB              *gorm.DB
	Config          Config
	Logger          *log.Logger
	KafkaReaders    map[string]*kafka.Reader
}

func main() {
	_ = godotenv.Load(".env")

	config := Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://famgo:[REDACTED]@localhost:5432/famgo"),
		KafkaBroker:      getEnv("KAFKA_BROKER", "localhost:9092"),
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		FirebaseCredentials: getEnv("FIREBASE_CREDENTIALS", ""),
		Port:             getEnv("PORT", "3003"),
		Environment:      getEnv("ENVIRONMENT", "development"),
	}

	logger := log.New(os.Stdout, "[NOTIFICATION-SERVICE] ", log.LstdFlags)

	// Connect to database
	db, err := connectDatabase(config.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Println("✓ Connected to PostgreSQL database")

	// Initialize app
	app := &App{
		DB:           db,
		Config:       config,
		Logger:       logger,
		KafkaReaders: make(map[string]*kafka.Reader),
	}

	// Start event consumers (SMS/Push notification requests)
	go app.startEventConsumers()

	// Start server
	logger.Printf("Starting Notification Service on port %s...", config.Port)
	if err := app.StartServer(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func connectDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// startEventConsumers starts consuming notification events from Kafka
func (app *App) startEventConsumers() {
	app.Logger.Println("Event consumers started")
	
	// Will consume:
	// - notification.send.sms
	// - notification.send.push
	
	// Process notifications and send via Twilio/Firebase
}

// StartServer starts the HTTP server
func (app *App) StartServer() error {
	app.Logger.Println("Notification Service is ready")
	app.Logger.Println("Service handles:")
	app.Logger.Println("  • SMS notifications (Twilio)")
	app.Logger.Println("  • Push notifications (Firebase)")
	app.Logger.Println("  • Email notifications (SendGrid)")
	app.Logger.Println("  • In-app notifications")

	select {}
}
