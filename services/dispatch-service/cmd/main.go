// services/dispatch-service/cmd/main.go

package main



import (

	"context"

	"fmt"

	"net"

	"net/http"

	"os"

	"os/signal"

	"syscall"

	"time"



	"github.com/gorilla/mux"

	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"

	"google.golang.org/grpc"



	dispatchv1 "github.com/Abdex1/FamGo-platform/services/dispatch-service/api/proto/v1"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/saga"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/usecases"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/config"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/ports"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/services"

	grpcHandler "github.com/Abdex1/FamGo-platform/services/dispatch-service/interfaces/grpc"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/infrastructure/clients"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/infrastructure/events"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/infrastructure/observability"

	redisinfra "github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/infrastructure/redis"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/infrastructure/repositories"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/infrastructure/security"

	rest "github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/interfaces/rest"

)



func main() {

	cfg := config.Load()



	logger, _ := zap.NewProduction()

	defer logger.Sync()



	rootCtx := context.Background()

	if _, err := observability.Bootstrap(rootCtx, cfg.ServiceName, cfg.Environment); err != nil {

		logger.Warn("telemetry bootstrap failed", zap.Error(err))

	}



	ctx, cancel := context.WithTimeout(rootCtx, 30*time.Second)

	defer cancel()



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



	dispatchRepo := repositories.NewDispatchRepository(pgPool)

	sessionRepo := repositories.NewMatchingSessionRepository(pgPool)

	matchRepo := repositories.NewMatchRepository(pgPool)



	gpsDiscovery, err := clients.NewGPSDriverDiscovery(cfg.GPSServiceURL)

	if err != nil {

		logger.Fatal("failed to initialize GPS discovery client", zap.Error(err))

	}

	defer gpsDiscovery.Close()



	matchingService := services.NewMatchingService(

		cfg.ProximityWeight,

		cfg.AcceptanceRateWeight,

		cfg.RatingWeight,

		cfg.AvailabilityWeight,

		cfg.MinAcceptanceRatePercent,

		cfg.MinRating,

		cfg.MaxSearchRadiusKm,

	)



	candidateEngine := services.NewDriverCandidateEngine(gpsDiscovery)

	matchingEngine := services.NewMatchingEngine(candidateEngine, matchingService, cfg.TopMatchesLimit)

	assignmentEngine := services.NewAssignmentEngine(time.Duration(cfg.MatchExpirySeconds)*time.Second, ports.NoOpPoolingHook{})

	timeoutService := services.NewTimeoutService(cfg.MatchRequestTTL)



	var eventPublisher ports.DispatchEventPublisher = events.NoOpPublisher{}

	var kafkaWriter *events.RawKafkaPublisher

	if cfg.KafkaEnabled {

		kafkaWriter = events.NewRawKafkaPublisher(cfg.KafkaBrokers)

		eventPublisher = events.NewKafkaPublisher(

			kafkaWriter.Publish,

			cfg.ServiceName,

			cfg.Environment,

		)

		logger.Info("Kafka publisher enabled", zap.Strings("brokers", cfg.KafkaBrokers))

	}



	dispatchUseCases := usecases.NewDispatchUseCases(

		dispatchRepo,

		sessionRepo,

		matchRepo,

		matchingEngine,

		assignmentEngine,

		timeoutService,

		eventPublisher,

		cfg.SearchRadiusKm,

		cfg.MaxSearchRadiusKm,

	)



	consumerCtx, consumerCancel := context.WithCancel(context.Background())

	defer consumerCancel()



	if cfg.KafkaEnabled {

		sagaHandler := saga.NewDispatchSagaHandler(dispatchUseCases)

		rideConsumer := events.NewRideCreatedConsumer(

			cfg.KafkaBrokers,

			cfg.KafkaGroupID,

			sagaHandler,

			logger,

		)

		rideConsumer.Start(consumerCtx)

		logger.Info("Kafka ride.created.v1 consumer started", zap.String("group_id", cfg.KafkaGroupID))

	}



	tokenValidator := security.NewTokenValidator(cfg.JWTSecret, cfg.JWTIssuer)

	var rateLimiter security.RateLimiter

	if redisClient, err := redisinfra.NewFromURL(cfg.RedisURL, cfg.RedisDB); err != nil {

		logger.Warn("redis unavailable; rate limiting disabled", zap.Error(err))

	} else {

		defer redisClient.Close()

		rateLimiter = security.NewRedisRateLimiter(redisClient.Raw())

		logger.Info("redis connected for rate limiting")

	}

	securityMiddleware := security.NewHTTPSecurityMiddleware(tokenValidator, rateLimiter)



	restHandler := rest.NewDispatchHandler(

		dispatchUseCases,

		candidateEngine,

		cfg.SearchRadiusKm,

		cfg.NearbyDriversLimit,

	)

	router := mux.NewRouter()

	restHandler.RegisterRoutes(router)



	httpServer := &http.Server{

		Addr:         fmt.Sprintf("%s:%s", cfg.HTTPHost, cfg.HTTPPort),

		Handler:      observability.HTTPMiddleware(cfg.ServiceName, securityMiddleware.Wrap(router)),

		ReadTimeout:  cfg.RequestTimeout,

		WriteTimeout: cfg.RequestTimeout,

		IdleTimeout:  60 * time.Second,

	}



	go func() {

		logger.Info("starting HTTP server", zap.String("port", cfg.HTTPPort))

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {

			logger.Fatal("HTTP server error", zap.Error(err))

		}

	}()



	grpcServer := grpc.NewServer(observability.GRPCServerOptions()...)

	dispatchHandler := grpcHandler.NewDispatchHandler(dispatchUseCases)

	dispatchv1.RegisterDispatchServiceServer(grpcServer, dispatchHandler)



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

	logger.Info("shutting down Dispatch service")



	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)

	defer shutdownCancel()



	consumerCancel()

	_ = httpServer.Shutdown(shutdownCtx)

	grpcServer.GracefulStop()

	pgPool.Close()

	if kafkaWriter != nil {

		_ = kafkaWriter.Close()

	}



	select {

	case <-shutdownCtx.Done():

		logger.Info("shutdown timeout exceeded")

	default:

		logger.Info("Dispatch service stopped")

	}

}

