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

	"famgo/user-service/internal/config"
	"famgo/user-service/internal/handler"
	"famgo/user-service/internal/repository"
	"famgo/user-service/internal/service"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
	"github.com/Abdex1/FamGo-platform/shared/pkg/observability"
)

// PATTERN 2: Service Bootstrap
// 11-step initialization sequence

func main() {
	// Step 1: Load Configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Initialize Logger
	log := logger.New(cfg.LogLevel)
	log.Info("starting user-service", map[string]interface{}{
		"version": "1.0.0",
		"env":     cfg.Environment,
	})

	// Step 3: Initialize Observability
	metrics, err := observability.InitMetrics("user-service")
	if err != nil {
		log.Error("failed to init metrics", map[string]interface{}{"error": err})
		os.Exit(1)
	}
	_ = metrics

	tracer, err := observability.InitTracer("user-service", cfg.OTelEndpoint)
	if err != nil {
		log.Error("failed to init tracer", map[string]interface{}{"error": err})
		os.Exit(1)
	}
	defer tracer.Shutdown(context.Background())

	// Step 4: Initialize Database
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

	// Step 5: Run Migrations
	if err := runMigrations(db, log); err != nil {
		log.Error("migration failed", map[string]interface{}{"error": err})
		os.Exit(1)
	}

	// Step 6: Initialize Repositories
	userRepo := repository.NewUserRepository(db)

	// Step 7: Initialize Services
	userService := service.NewUserService(userRepo, log)

	// Step 8: Initialize HTTP Router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	// Register handlers
	handlers := handler.NewHandler(userService, log)
	handlers.RegisterRoutes(router)

	// Health checks
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

	errChan := make(chan error, 1)
	go func() {
		log.Info("listening", map[string]interface{}{"port": cfg.Port})
		errChan <- server.ListenAndServe()
	}()

	// Step 11: Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		if err != http.ErrServerClosed {
			log.Error("server error", map[string]interface{}{"error": err})
		}
	case sig := <-sigChan:
		log.Info("shutdown signal received", map[string]interface{}{"signal": sig.String()})

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
	createProfilesTable := `
	CREATE TABLE IF NOT EXISTS user_profiles (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		auth_id UUID NOT NULL UNIQUE,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		profile_picture_url VARCHAR(500),
		email_verified BOOLEAN DEFAULT FALSE,
		phone_verified BOOLEAN DEFAULT FALSE,
		rating DECIMAL(3,2) DEFAULT 0,
		total_rides INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_profiles_auth_id ON user_profiles(auth_id);
	`

	createPreferencesTable := `
	CREATE TABLE IF NOT EXISTS user_preferences (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES user_profiles(id),
		notification_email BOOLEAN DEFAULT TRUE,
		notification_sms BOOLEAN DEFAULT TRUE,
		language VARCHAR(10) DEFAULT 'en',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id)
	);
	`

	createAddressesTable := `
	CREATE TABLE IF NOT EXISTS user_addresses (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES user_profiles(id),
		type VARCHAR(50),
		address_line_1 VARCHAR(255),
		city VARCHAR(100),
		lat DECIMAL(10,8),
		lng DECIMAL(11,8),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON user_addresses(user_id);
	`

	if _, err := db.Exec(createProfilesTable); err != nil {
		return fmt.Errorf("failed to create user_profiles table: %w", err)
	}

	if _, err := db.Exec(createPreferencesTable); err != nil {
		return fmt.Errorf("failed to create user_preferences table: %w", err)
	}

	if _, err := db.Exec(createAddressesTable); err != nil {
		return fmt.Errorf("failed to create user_addresses table: %w", err)
	}

	log.Info("migrations completed", nil)
	return nil
}
