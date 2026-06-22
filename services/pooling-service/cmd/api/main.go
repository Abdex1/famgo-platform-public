package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Abdex1/FamGo-platform/services/pooling-service/internal/domain/services"
	"github.com/Abdex1/FamGo-platform/services/pooling-service/internal/infrastructure/postgres"
	"github.com/Abdex1/FamGo-platform/services/pooling-service/internal/interfaces/rest"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "famgo_user")
	dbPassword := getEnv("DB_PASSWORD", "secure_password")
	dbName := getEnv("DB_NAME", "famgo_platform")
	servicePort := getEnv("SERVICE_PORT", "3013")

	// Database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL")

	// Initialize repositories
	poolRepo := postgres.NewPoolRepository(db)

	// Initialize pooling engine
	poolingEngine := services.NewPoolingEngine()

	// Initialize handlers
	poolingHandler := rest.NewPoolingHandler(poolRepo, poolingEngine)

	// Setup router
	router := mux.NewRouter()
	poolingHandler.RegisterRoutes(router)

	// Start server
	server := &http.Server{
		Addr:    ":" + servicePort,
		Handler: router,
	}

	log.Printf("🚀 Pooling Service starting on port %s\n", servicePort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
