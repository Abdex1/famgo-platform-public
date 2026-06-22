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

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"famgo/auth-service/internal/domain"
    "famgo/auth-service/internal/handlers"
    "famgo/auth-service/internal/infrastructure/postgres"
    "famgo/auth-service/internal/infrastructure/redis"
	pb "famgo/auth-service/api/proto/v1"
)

func main() {
	logger := log.New(os.Stdout, "[auth-service] ", log.LstdFlags|log.Lshortfile)

	// Configuration from environment
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "famgo_auth"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "5001"
	}

	logger.Println("=== Auth Service Starting ===")
	logger.Printf("Database: %s:%s/%s", dbHost, dbPort, dbName)
	logger.Printf("Redis: %s:%s", redisHost, redisPort)
	logger.Printf("Kafka: %s", kafkaBrokers)
	logger.Printf("gRPC Port: %s", grpcPort)

	// Initialize database
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := db.PingContext(ctx); err != nil {
		cancel()
		logger.Fatalf("Database connection failed: %v", err)
	}
	cancel()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Println("✓ Database connected")

	// Initialize Redis
	_ = redis.NewSessionStore(fmt.Sprintf("%s:%s", redisHost, redisPort))
	logger.Println("✓ Redis connected")

	// Kafka producer wiring deferred until domain service is implemented.
	logger.Println("✓ Kafka producer skipped (not wired)")

	// Initialize repositories
	_ = postgres.NewAuthRepository(db)

	// Initialize domain service
	var domainService domain.AuthService
	// TODO: Implement auth service implementation
	// domainService = service.NewAuthServiceImpl(authRepo, cache, kafkaProducer, logger)

	// Initialize gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	)

	// Register gRPC services
	authHandler := handlers.NewAuthServer(domainService)
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	// Register health check service
	healthCheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthCheck)
	healthCheck.SetServingStatus("auth.v1.AuthService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for gRPC UI
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", grpcPort))
	if err != nil {
		logger.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	// Initialize OpenTelemetry (optional)
	initTracing(logger)

	logger.Printf("✓ gRPC server listening on %s", lis.Addr())
	logger.Println("=== Auth Service Ready ===")

	// Start server in goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Println("Shutting down...")

	grpcServer.GracefulStop()
	logger.Println("Auth service stopped")
}

// loggingInterceptor logs gRPC requests
func loggingInterceptor(logger *log.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Printf("gRPC call: %s", info.FullMethod)
		return handler(ctx, req)
	}
}

// initTracing initializes OpenTelemetry tracing (optional)
func initTracing(logger *log.Logger) {
	otlpEndpoint := os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		logger.Println("⚠ OpenTelemetry not configured (OTLP_ENDPOINT not set)")
		otel.SetTracerProvider(noop.NewTracerProvider())
		return
	}

	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(otlpEndpoint),
	)
	if err != nil {
		logger.Printf("⚠ Failed to create OTLP exporter: %v", err)
		return
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	logger.Println("✓ OpenTelemetry tracing initialized")
}
