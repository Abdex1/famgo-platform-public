// services/ride-service/cmd/main.go
// Ride Service Entry Point

package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/bootstrap"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/config"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	// Initialize database
	db, err := initializeDatabase(cfg, logger)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Redis
	redisClient := initializeRedis(cfg, logger)
	defer redisClient.Close()

	// Create application container
	container := bootstrap.NewAppContainer(db, redisClient, logger)
	defer container.Cleanup()

	// Setup HTTP routes
	router := setupRoutes(container)

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", zap.Error(err))
		}
	}()

	logger.Info("starting ride service",
		zap.String("version", "1.0.0"),
		zap.Int("port", cfg.Port),
		zap.String("environment", cfg.Environment))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server error", zap.Error(err))
	}
}

func initializeDatabase(cfg *config.Config, logger *zap.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Info("database connected",
		zap.String("host", cfg.Database.Host),
		zap.String("name", cfg.Database.Name))

	return db, nil
}

func initializeRedis(cfg *config.Config, logger *zap.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatal("redis connection failed", zap.Error(err))
	}

	logger.Info("redis connected",
		zap.String("host", cfg.Redis.Host),
		zap.Int("port", cfg.Redis.Port))

	return client
}

func setupRoutes(container *bootstrap.AppContainer) *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))

	// Health checks
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"alive"}`))
	})

	router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})

	// API Routes
	router.Route("/rides", func(r chi.Router) {
		// Create ride
		r.Post("/", container.HTTPServer.CreateRide)

		// Get specific ride
		r.Get("/{rideID}", container.HTTPServer.GetRide)

		// Ride state transitions
		r.Post("/{rideID}/assign", container.HTTPServer.AssignDriver)
		r.Post("/{rideID}/start", container.HTTPServer.StartRide)
		r.Post("/{rideID}/complete", container.HTTPServer.CompleteRide)
		r.Post("/{rideID}/cancel", container.HTTPServer.CancelRide)
	})

	// Passenger routes
	router.Route("/passengers", func(r chi.Router) {
		r.Get("/{passengerID}/rides", container.HTTPServer.GetPassengerRides)
	})

	// Driver routes
	router.Route("/drivers", func(r chi.Router) {
		r.Get("/{driverID}/rides", container.HTTPServer.GetDriverRides)
	})

	return router
}
