// services/gps-service/cmd/main.go
// GPS Service main bootstrap

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/shared/database"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/config"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/services"
	grpcHandler "github.com/Abdex1/FamGo-platform/services/gps-service/interfaces/grpc"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/infrastructure/redis"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/infrastructure/repositories"
	"github.com/Abdex1/FamGo-platform/services/gps-service/proto/gps"
)

func main() {
	// Load config
	cfg := config.Load()

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize PostgreSQL
	pgConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
		cfg.DatabaseSSLMode,
	))
	if err != nil {
		logger.Fatal("invalid postgres config", zap.Error(err))
	}

	pgConfig.MaxConns = cfg.DatabaseMaxConnections
	pgConfig.MinConns = cfg.DatabaseMinConnections
	pgConfig.MaxConnLifetime = cfg.DatabaseConnMaxLifetime
	pgConfig.MaxConnIdleTime = cfg.DatabaseConnMaxIdleTime

	pgPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		logger.Fatal("failed to create postgres pool", zap.Error(err))
	}
	defer pgPool.Close()

	if err := pgPool.Ping(ctx); err != nil {
		logger.Fatal("failed to ping postgres", zap.Error(err))
	}
	logger.Info("connected to PostgreSQL")

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
		DB:   cfg.RedisDB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("failed to ping redis", zap.Error(err))
	}
	logger.Info("connected to Redis")

	// Initialize repositories and stores
	locationRepo := repositories.NewDriverLocationRepository(pgPool)
	geoIndexStore := redis.NewGeoIndexStore(redisClient, "gps", cfg.NearbyDriversCacheTTL)
	trackingStore := redis.NewDriverTrackingStore(redisClient, "gps", cfg.DriverStatusUpdateTTL)

	// Initialize domain services
	locationService := services.NewLocationService()
	redisGeoService := services.NewRedisGeoService()

	// Initialize use cases
	locationUseCases := usecases.NewLocationUseCases(
		locationRepo,
		geoIndexStore,
		trackingStore,
		locationService,
		redisGeoService,
	)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	gpsHandler := grpcHandler.NewGPSHandler(locationUseCases)
	gps.RegisterGPSServiceServer(grpcServer, gpsHandler)

	// Start gRPC server
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.GRPCHost, cfg.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		logger.Info("starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("shutting down GPS service")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	grpcServer.GracefulStop()
	pgPool.Close()
	redisClient.Close()

	select {
	case <-shutdownCtx.Done():
		logger.Info("shutdown timeout exceeded")
	default:
		logger.Info("GPS service stopped")
	}
}
