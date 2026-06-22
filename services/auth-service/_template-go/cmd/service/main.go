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

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// Configuration from environment
	serviceName := getEnv("SERVICE_NAME", "template-service")
	port := getEnv("PORT", "5001")
	logLevel := getEnv("LOG_LEVEL", "info")

	// Initialize logger
	logger := NewLogger(serviceName, logLevel)
	logger.Infof("Starting %s on port %s", serviceName, port)

	// Setup gRPC server with interceptors
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server with middleware
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandler(func(p interface{}) (err error) {
					logger.Errorf("Panic recovered: %v", p)
					return status.Errorf(codes.Internal, "Internal server error")
				}),
			),
		),
		grpc.MaxConcurrentStreams(1000),
	}

	grpcServer := grpc.NewServer(opts...)

	// Register services here
	// pb.RegisterServiceServer(grpcServer, NewServiceHandler())

	// Health check
	registerHealthCheck(grpcServer)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down gracefully...")
		grpcServer.GracefulStop()
	}()

	// Start server
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Logger wrapper
type Logger struct {
	name string
}

func NewLogger(name, level string) *Logger {
	return &Logger{name: name}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] [%s] "+format, append([]interface{}{l.name}, args...)...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] [%s] "+format, append([]interface{}{l.name}, args...)...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf("[FATAL] [%s] "+format, append([]interface{}{l.name}, args...)...)
}

func registerHealthCheck(server *grpc.Server) {
	// Health check implementation
	// This would implement grpc.health.v1.Health service
}
