// services/ride-service/internal/bootstrap/bootstrap.go
// Dependency Injection and Application Bootstrap

package bootstrap

import (
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/application"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/infrastructure"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/transport"
)

// AppContainer holds all application dependencies
type AppContainer struct {
	// Infrastructure
	DB     *sql.DB
	Redis  *redis.Client
	Logger *zap.Logger

	// Repositories
	RideRepo                    domain.RideRepository
	RideStatusHistoryRepo       domain.RideStatusHistoryRepository
	RideCache                   domain.RideCache

	// Domain Services
	RideService *domain.RideService

	// Event Publishing (COMPLIANT: using packages/event-bus)
	EventPublisher *application.EventPublisher

	// Command Handlers
	CreateRideHandler     *application.CreateRideHandler
	AssignDriverHandler   *application.AssignDriverHandler
	StartRideHandler      *application.StartRideHandler
	CompleteRideHandler   *application.CompleteRideHandler
	CancelRideHandler     *application.CancelRideHandler

	// Query Handlers
	GetRideHandler              *application.GetRideHandler
	GetPassengerRidesHandler    *application.GetPassengerRidesHandler
	GetDriverRidesHandler       *application.GetDriverRidesHandler
	GetActiveRidesHandler       *application.GetActiveRidesHandler

	// Transport
	HTTPServer *transport.HTTPServer
}

// NewAppContainer creates and initializes all application dependencies
func NewAppContainer(db *sql.DB, redis *redis.Client, logger *zap.Logger) *AppContainer {
	container := &AppContainer{
		DB:     db,
		Redis:  redis,
		Logger: logger,
	}

	// Initialize infrastructure layer
	container.initializeRepositories()

	// Initialize event publishing (COMPLIANT: uses packages/event-bus)
	container.initializeEventPublishing()

	// Initialize domain services
	container.initializeDomainServices()

	// Initialize application layer (handlers)
	container.initializeCommandHandlers()
	container.initializeQueryHandlers()

	// Initialize transport layer
	container.initializeTransport()

	return container
}

// initializeRepositories sets up all repository implementations
func (c *AppContainer) initializeRepositories() {
	// PostgreSQL repositories
	c.RideRepo = infrastructure.NewPostgresRideRepository(c.DB)
	c.RideStatusHistoryRepo = infrastructure.NewPostgresRideStatusHistoryRepository(c.DB)

	// Redis cache
	c.RideCache = infrastructure.NewRedisRideCache(c.Redis, 3600) // 1 hour default TTL
}

// initializeDomainServices sets up domain business logic
func (c *AppContainer) initializeDomainServices() {
	c.RideService = domain.NewRideService(c.Logger)
}

// initializeEventPublishing sets up event publishing
func (c *AppContainer) initializeEventPublishing() {
	// COMPLIANT: EventBus would be injected from packages/event-bus
	// For now, this is a placeholder - actual event-bus setup in bootstrap main.go
	c.EventPublisher = application.NewEventPublisher(
		nil, // EventBus from packages/event-bus (injected at runtime)
		c.Logger,
	)
}

// initializeCommandHandlers sets up command handlers
func (c *AppContainer) initializeCommandHandlers() {
	c.CreateRideHandler = application.NewCreateRideHandler(
		c.RideRepo,
		c.RideCache,
		c.RideService,
		c.EventPublisher,
		c.Logger,
	)

	c.AssignDriverHandler = application.NewAssignDriverHandler(
		c.RideRepo,
		c.RideCache,
		c.RideService,
		c.EventPublisher,
		c.Logger,
	)

	c.StartRideHandler = application.NewStartRideHandler(
		c.RideRepo,
		c.RideCache,
		c.RideService,
		c.EventPublisher,
		c.Logger,
	)

	c.CompleteRideHandler = application.NewCompleteRideHandler(
		c.RideRepo,
		c.RideCache,
		c.RideService,
		c.EventPublisher,
		c.Logger,
	)

	c.CancelRideHandler = application.NewCancelRideHandler(
		c.RideRepo,
		c.RideCache,
		c.RideService,
		c.EventPublisher,
		c.Logger,
	)
}

// initializeQueryHandlers sets up query handlers
func (c *AppContainer) initializeQueryHandlers() {
	c.GetRideHandler = application.NewGetRideHandler(
		c.RideRepo,
		c.RideCache,
		c.Logger,
	)

	c.GetPassengerRidesHandler = application.NewGetPassengerRidesHandler(
		c.RideRepo,
		c.RideCache,
		c.Logger,
	)

	c.GetDriverRidesHandler = application.NewGetDriverRidesHandler(
		c.RideRepo,
		c.RideCache,
		c.Logger,
	)

	c.GetActiveRidesHandler = application.NewGetActiveRidesHandler(
		c.RideRepo,
		c.RideCache,
		c.Logger,
	)
}

// initializeTransport sets up HTTP and gRPC servers
func (c *AppContainer) initializeTransport() {
	c.HTTPServer = transport.NewHTTPServer(
		c.CreateRideHandler,
		c.AssignDriverHandler,
		c.StartRideHandler,
		c.CompleteRideHandler,
		c.CancelRideHandler,
		c.GetRideHandler,
		c.GetPassengerRidesHandler,
		c.GetDriverRidesHandler,
		c.Logger,
	)
}

// Cleanup closes all resources
func (c *AppContainer) Cleanup() error {
	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			c.Logger.Error("failed to close database", zap.Error(err))
			return err
		}
	}

	if c.Redis != nil {
		if err := c.Redis.Close(); err != nil {
			c.Logger.Error("failed to close redis", zap.Error(err))
			return err
		}
	}

	c.Logger.Sync()
	return nil
}

// WithGracefulShutdown provides graceful shutdown with timeout
func (c *AppContainer) WithGracefulShutdown(timeout time.Duration) func() error {
	return func() error {
		done := make(chan error, 1)
		go func() {
			done <- c.Cleanup()
		}()

		select {
		case err := <-done:
			return err
		case <-time.After(timeout):
			c.Logger.Warn("graceful shutdown timeout exceeded")
			return nil
		}
	}
}
