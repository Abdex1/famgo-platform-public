package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"famgo/auth-service/internal/config"
	"famgo/auth-service/internal/handler"
	"famgo/auth-service/internal/repository"
	"famgo/auth-service/internal/service"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
	"github.com/Abdex1/FamGo-platform/shared/pkg/observability"
)

// PATTERN 2: Service Bootstrap
// 11-step initialization sequence following extracted pattern

func main() {
	// Step 1: Load Configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Initialize Logger
	log := logger.New(cfg.LogLevel)
	log.Info("starting auth-service", map[string]interface{}{
		"version": "1.0.0",
		"env":     cfg.Environment,
	})

	// Step 3: Initialize Observability (Prometheus + OpenTelemetry)
	_, err = observability.InitMetrics("auth-service")
	if err != nil {
		log.Error("failed to init metrics", map[string]interface{}{"error": err})
		os.Exit(1)
	}

	tracer, err := observability.InitTracer("auth-service", cfg.OTelEndpoint)
	if err != nil {
		log.Error("failed to init tracer", map[string]interface{}{"error": err})
		os.Exit(1)
	}
	defer tracer.Shutdown(context.Background())

	// Step 4: Initialize Database Connection
	dbConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := sqlx.Open("postgres", dbConnStr)
	if err != nil {
		log.Error("failed to open database", map[string]interface{}{"error": err})
		os.Exit(1)
	}
	defer db.Close()

	// Verify DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Error("failed to ping database", map[string]interface{}{"error": err})
		os.Exit(1)
	}
	log.Info("database connected", nil)

	// Step 5: Database Migrations (Run here)
	if err := runMigrations(db, log); err != nil {
		log.Error("migration failed", map[string]interface{}{"error": err})
		os.Exit(1)
	}

	// Step 6: Initialize Repositories
	userRepo := repository.NewUserRepository(db)

	// Step 7: Initialize Services
	authService := service.NewAuthService(cfg, userRepo, log)

	// Step 8: Initialize HTTP Router (Pattern 1: HTTP Handlers)
	router := chi.NewRouter()

	// Middleware stack (Pattern 1)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	// Register handlers
	handlers := handler.NewHandler(authService, log)
	handlers.RegisterRoutes(router)

	// Health checks (Step 9)
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not_ready"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})

	// Step 10: Start HTTP Server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown channel
	errChan := make(chan error, 1)
	go func() {
		log.Info("listening", map[string]interface{}{"port": cfg.Port})
		errChan <- server.ListenAndServe()
	}()

	// Step 11: Graceful Shutdown Handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		if err != http.ErrServerClosed {
			log.Error("server error", map[string]interface{}{"error": err})
		}
	case sig := <-sigChan:
		log.Info("shutdown signal received", map[string]interface{}{"signal": sig.String()})

		// Graceful shutdown with 30-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error("shutdown error", map[string]interface{}{"error": err})
		}

		db.Close()
		log.Info("shutdown complete", nil)
	}
}

// runMigrations executes database migrations
func runMigrations(db *sqlx.DB, log logger.Logger) error {
	// Migration: Create users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		phone VARCHAR(20),
		role VARCHAR(50) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		email_verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	// Migration: Create OTP verification table
	createOTPTable := `
	CREATE TABLE IF NOT EXISTS otp_verification (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) NOT NULL,
		otp VARCHAR(6) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		attempts INT DEFAULT 0,
		verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_otp_email ON otp_verification(email);
	`

	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	if _, err := db.Exec(createOTPTable); err != nil {
		return fmt.Errorf("failed to create otp_verification table: %w", err)
	}

	log.Info("migrations completed", nil)
	return nil
}
