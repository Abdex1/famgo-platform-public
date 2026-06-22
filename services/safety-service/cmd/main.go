// services/safety-service/cmd/main.go
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/safety-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/safety-service/internal/config"
	grpcHandler "github.com/Abdex1/FamGo-platform/services/safety-service/interfaces/grpc"
	"github.com/Abdex1/FamGo-platform/services/safety-service/internal/infrastructure/repositories"
	"github.com/Abdex1/FamGo-platform/services/safety-service/proto/safety"
)

func main() {
	cfg := config.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pgConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:[REDACTED]@%s:%d/%s?sslmode=%s",
		cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost,
		cfg.DatabasePort, cfg.DatabaseName, cfg.DatabaseSSLMode,
	))
	if err != nil {
		logger.Fatal("invalid postgres config", zap.Error(err))
	}

	pgConfig.MaxConns = cfg.DatabaseMaxConns
	pgConfig.MinConns = cfg.DatabaseMinConns

	pgPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		logger.Fatal("failed to create postgres pool", zap.Error(err))
	}
	defer pgPool.Close()

	if err := pgPool.Ping(ctx); err != nil {
		logger.Fatal("failed to ping postgres", zap.Error(err))
	}
	logger.Info("connected to PostgreSQL")

	sosRepo := repositories.NewSOSRepository(pgPool)
	safetyUseCases := usecases.NewSafetyUseCases(sosRepo)

	grpcServer := grpc.NewServer()
	safetyHandler := grpcHandler.NewSafetyHandler(safetyUseCases)
	safety.RegisterSafetyServiceServer(grpcServer, safetyHandler)

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("shutting down Safety service")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	grpcServer.GracefulStop()
	pgPool.Close()

	select {
	case <-shutdownCtx.Done():
		logger.Info("shutdown timeout exceeded")
	default:
		logger.Info("Safety service stopped")
	}
}
