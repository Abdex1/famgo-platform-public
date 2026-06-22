package bootstrap

import (
	"database/sql"
	"time"

	"github.com/Abdex1/FamGo-platform/packages/event-bus"
	"github.com/Abdex1/FamGo-platform/packages/redis-platform"
	"github.com/Abdex1/FamGo-platform/packages/telemetry"
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/application"
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/domain"
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/infrastructure"
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/transport"
)

// Container holds all service dependencies
type Container struct {
	// Repositories
	LocationRepo   domain.LocationRepository
	TripRepo       domain.TripRepository
	GeofenceRepo   domain.GeofenceRepository

	// Domain Services
	LocationService *domain.LocationService

	// Application Handlers
	UpdateLocationHandler   *application.UpdateDriverLocationHandler
	GetLocationHandler      *application.GetDriverLocationHandler
	GetNearbyDriversHandler *application.GetNearbyDriversHandler

	// Transport
	HTTPHandler *transport.HTTPHandler

	// Infrastructure
	DB    *sql.DB
	Redis redis_platform.RedisClient

	// Observability
	Metrics telemetry.Metrics
	Logger  telemetry.Logger
}

// NewContainer creates and initializes all dependencies
func NewContainer(db *sql.DB, redis redis_platform.RedisClient) *Container {
	// Initialize observability
	metrics := telemetry.NewMetrics("gps-service")
	logger := telemetry.NewLogger("gps-service")

	// Initialize repositories
	locationRepo := infrastructure.NewPostgresLocationRepository(db)
	tripRepo := infrastructure.NewPostgresTripRepository(db)
	geofenceRepo := infrastructure.NewPostgresGeofenceRepository(db)

	// Initialize domain service
	locationService := domain.NewLocationService()

	// Initialize caches
	locationCache := infrastructure.NewRedisLocationCache(redis, 5*time.Minute)
	tripCache := infrastructure.NewRedisTripCache(redis, 5*time.Minute)
	geofenceCache := infrastructure.NewRedisGeofenceCache(redis, 1*time.Hour)
	driverCache := infrastructure.NewRedisDriverCache(redis)

	// Initialize event bus (from packages)
	eventBus := event_bus.NewEventBus(&event_bus.Config{
		KafkaClient: nil, // Would be initialized from config
	})

	// Initialize application handlers
	updateLocationHandler := application.NewUpdateDriverLocationHandler(
		locationRepo,
		geofenceRepo,
		eventBus,
		locationService,
		metrics,
		logger,
	)

	getLocationHandler := application.NewGetDriverLocationHandler(
		locationRepo,
		metrics,
		logger,
	)

	getNearbyDriversHandler := application.NewGetNearbyDriversHandler(
		locationRepo,
		locationService,
		metrics,
		logger,
	)

	// Initialize HTTP handler
	httpHandler := transport.NewHTTPHandler(
		updateLocationHandler,
		getLocationHandler,
		getNearbyDriversHandler,
		metrics,
		logger,
	)

	// Return container
	return &Container{
		LocationRepo:   locationRepo,
		TripRepo:       tripRepo,
		GeofenceRepo:   geofenceRepo,
		LocationService: locationService,
		UpdateLocationHandler:   updateLocationHandler,
		GetLocationHandler:      getLocationHandler,
		GetNearbyDriversHandler: getNearbyDriversHandler,
		HTTPHandler:            httpHandler,
		DB:                     db,
		Redis:                  redis,
		Metrics:                metrics,
		Logger:                 logger,
	}
}
