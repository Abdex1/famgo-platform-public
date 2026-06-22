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
	DatabaseURL   string
	KafkaBroker   string
	Port          string
	Environment   string
}

// App holds application dependencies
type App struct {
	DB            *gorm.DB
	Config        Config
	Logger        *log.Logger
	KafkaWriter   *kafka.Writer
	KafkaReaders  map[string]*kafka.Reader
}

func main() {
	_ = godotenv.Load(".env")

	config := Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://famgo:[REDACTED]@localhost:5432/famgo"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		Port:        getEnv("PORT", "3001"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	logger := log.New(os.Stdout, "[USER-SERVICE] ", log.LstdFlags)

	// Connect to database
	db, err := connectDatabase(config.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Println("✓ Connected to PostgreSQL database")

	// Initialize Kafka producer
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(config.KafkaBroker),
		Topic:    "user.events", // Topic will be set dynamically
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()
	logger.Println("✓ Kafka producer initialized")

	// Initialize app
	app := &App{
		DB:          db,
		Config:      config,
		Logger:      logger,
		KafkaWriter: kafkaWriter,
		KafkaReaders: make(map[string]*kafka.Reader),
	}

	// Start event consumers
	go app.startEventConsumers()

	// Start server
	logger.Printf("Starting User Service on port %s...", config.Port)
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

// startEventConsumers starts consuming events from Kafka
func (app *App) startEventConsumers() {
	// Will consume auth events to sync user data
	app.Logger.Println("Event consumers started")
}

// StartServer starts the HTTP server
func (app *App) StartServer() error {
	app.Logger.Println("User Service is ready")
	app.Logger.Println("Endpoints:")
	app.Logger.Println("  GET    /v1/users/:id")
	app.Logger.Println("  PUT    /v1/users/:id")
	app.Logger.Println("  GET    /v1/users/:id/profile")
	app.Logger.Println("  PUT    /v1/users/:id/profile")
	app.Logger.Println("  POST   /v1/users/:id/preferences")
	app.Logger.Println("  GET    /v1/users/:id/history")

	// Placeholder: Keep server running
	select {}
}
