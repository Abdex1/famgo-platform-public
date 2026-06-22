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

	"famgo/driver-service/internal/config"
	"famgo/driver-service/internal/handler"
	"famgo/driver-service/internal/repository"
	"famgo/driver-service/internal/service"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
	"github.com/Abdex1/FamGo-platform/shared/pkg/observability"
)

// PATTERN 2: Service Bootstrap
// 11-step initialization sequence for driver service foundation

func main() {
	// Step 1: Load Configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Initialize Logger
	log := logger.New(cfg.LogLevel)
	log.Info("starting driver-service", map[string]interface{}{
		"version": "1.0.0",
		"env":     cfg.Environment,
	})

	// Step 3: Initialize Observability
	metrics, err := observability.InitMetrics("driver-service")
	if err != nil {
		log.Error("failed to init metrics", map[string]interface{}{"error": err})
		os.Exit(1)
	}
	_ = metrics

	tracer, err := observability.InitTracer("driver-service", cfg.OTelEndpoint)
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

	// Step 5: Run Migrations (WEEK 1 FOUNDATION ONLY)
	if err := runMigrations(db, log); err != nil {
		log.Error("migration failed", map[string]interface{}{"error": err})
		os.Exit(1)
	}

	// Step 6: Initialize Repositories
	driverRepo := repository.NewDriverRepository(db)

	// Step 7: Initialize Services (FOUNDATION ONLY)
	driverService := service.NewDriverService(driverRepo, log)

	// Step 8: Initialize HTTP Router (PATTERN 1)
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	// Register handlers
	handlers := handler.NewHandler(driverService, log)
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

// runMigrations executes database migrations (WEEK 1 FOUNDATION)
func runMigrations(db *sqlx.DB, log logger.Logger) error {
	// WEEK 1: Foundation tables only
	// Full tables (documents, vehicles, verification) in WEEK 3

	createDriversTable := `
	CREATE TABLE IF NOT EXISTS drivers (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		auth_id UUID NOT NULL UNIQUE,
		license_number VARCHAR(100) UNIQUE,
		license_expiry DATE,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		verification_status VARCHAR(50) DEFAULT 'pending',
		date_joined DATE DEFAULT CURRENT_DATE,
		rating DECIMAL(3,2) DEFAULT 0,
		total_rides INT DEFAULT 0,
		total_earnings DECIMAL(12,2) DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_drivers_auth_id ON drivers(auth_id);
	CREATE INDEX IF NOT EXISTS idx_drivers_status ON drivers(status);
	`

	// Driver state machine: pending -> approved -> active -> suspended
	createDriverStatesTable := `
	CREATE TABLE IF NOT EXISTS driver_states (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		driver_id UUID NOT NULL REFERENCES drivers(id),
		current_state VARCHAR(50) NOT NULL,
		previous_state VARCHAR(50),
		reason VARCHAR(500),
		transition_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_driver_states_driver_id ON driver_states(driver_id);
	CREATE INDEX IF NOT EXISTS idx_driver_states_current ON driver_states(current_state);
	`

	if _, err := db.Exec(createDriversTable); err != nil {
		return fmt.Errorf("failed to create drivers table: %w", err)
	}

	if _, err := db.Exec(createDriverStatesTable); err != nil {
		return fmt.Errorf("failed to create driver_states table: %w", err)
	}

	log.Info("migrations completed (WEEK 1 FOUNDATION)", nil)
	log.Info("NOTE: Full verification, documents, and location tables in WEEK 3", nil)

	return nil
}
