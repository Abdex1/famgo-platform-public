# 🎯 SERVICE COMPLETION TEMPLATES (Weeks 3-4)

**Status:** Reference templates for completing services  
**Focus:** Domain-driven architecture with platform abstractions  
**Use Case:** Complete GPS, User, and Ride services

---

## 🏗️ SERVICE TEMPLATE: DOMAIN LAYER

Every service domain layer MUST follow this structure.

### Example: GPS Service Domain

**File:** `services/gps-service/internal/domain/entities.go`

```go
package domain

import (
    "time"
    "github.com/google/uuid"
    "github.com/famgo/shared/contracts/events/common"
)

// DriverLocation represents current driver location
type DriverLocation struct {
    ID        string      // UUID
    DriverID  string      // Foreign key to driver
    Latitude  float64
    Longitude float64
    Accuracy  float32
    UpdatedAt time.Time
}

// Trip represents a trip being tracked
type Trip struct {
    ID        string
    RideID    string      // Foreign key to ride
    StartedAt time.Time
    Location  DriverLocation
    Route     []RoutePoint
    Status    TripStatus
}

type TripStatus string

const (
    TripStatusActive    TripStatus = "ACTIVE"
    TripStatusCompleted TripStatus = "COMPLETED"
    TripStatusCancelled TripStatus = "CANCELLED"
)

// Geofence represents a geographic boundary
type Geofence struct {
    ID        string
    Name      string
    Latitude  float64
    Longitude float64
    Radius    float32      // meters
    CreatedAt time.Time
}

// RoutePoint represents a point in the trip route
type RoutePoint struct {
    Latitude  float64
    Longitude float64
    Timestamp time.Time
}

// DomainService: LocationService
type LocationService struct {
    // NO external dependencies here - pure logic
}

func (s *LocationService) IsWithinGeofence(
    location DriverLocation,
    geofence Geofence,
) bool {
    // Pure domain logic - no I/O
    return distance(location, geofence) <= geofence.Radius
}

func (s *LocationService) CalculateDeviation(
    expected, actual DriverLocation,
) float32 {
    // Pure domain logic
    return distance(expected, actual)
}
```

---

## 🔄 APPLICATION LAYER TEMPLATE

**File:** `services/gps-service/internal/application/commands.go`

```go
package application

import (
    "context"
    "time"
    "github.com/google/uuid"
    "github.com/famgo/gps-service/internal/domain"
    "github.com/famgo/shared/contracts/events"
)

// Command: Update driver location
type UpdateDriverLocationCommand struct {
    DriverID  string
    Latitude  float64
    Longitude float64
    Accuracy  float32
}

// Handler: Update driver location
type UpdateDriverLocationHandler struct {
    locationRepo    LocationRepository        // Interface - depends on abstraction
    geofenceRepo    GeofenceRepository        // Interface - depends on abstraction
    eventBus        events.EventBus           // From packages/event-bus
    locationService *domain.LocationService
}

func (h *UpdateDriverLocationHandler) Handle(
    ctx context.Context,
    cmd UpdateDriverLocationCommand,
) error {
    // 1. Get driver location (from cache via Redis abstraction)
    oldLocation, _ := h.locationRepo.GetDriverLocation(ctx, cmd.DriverID)
    
    // 2. Create new location entity
    newLocation := &domain.DriverLocation{
        ID:        uuid.New().String(),
        DriverID:  cmd.DriverID,
        Latitude:  cmd.Latitude,
        Longitude: cmd.Longitude,
        Accuracy:  cmd.Accuracy,
        UpdatedAt: time.Now(),
    }
    
    // 3. Apply domain logic
    geofences, _ := h.geofenceRepo.GetAllGeofences(ctx)
    for _, geofence := range geofences {
        if h.locationService.IsWithinGeofence(*newLocation, geofence) {
            // Domain event - NO immediate side effects
            geofenceEntered := events.GeofenceEnteredEvent{
                DriverID:   cmd.DriverID,
                GeofenceID: geofence.ID,
            }
            
            // 4. Publish event through platform event-bus
            err := h.eventBus.Publish(ctx, geofenceEntered)
            if err != nil {
                return err
            }
        }
    }
    
    // 5. Persist location
    err := h.locationRepo.UpdateDriverLocation(ctx, newLocation)
    if err != nil {
        return err
    }
    
    // 6. Publish domain event through shared/contracts/events
    locationUpdated := events.DriverLocationUpdatedEvent{
        EventID:    uuid.New().String(),
        EventType:  events.EventTypeDriverLocationUpdated,
        Version:    events.VersionV1,
        AggregateID: cmd.DriverID,
        Timestamp:  time.Now(),
        Data: LocationUpdatedData{
            OldLatitude:  oldLocation.Latitude,
            OldLongitude: oldLocation.Longitude,
            NewLatitude:  newLocation.Latitude,
            NewLongitude: newLocation.Longitude,
        },
    }
    
    // Idempotent publishing through platform event-bus
    return h.eventBus.PublishIdempotent(ctx, locationUpdated)
}
```

**File:** `services/gps-service/internal/application/queries.go`

```go
package application

import (
    "context"
    "github.com/famgo/gps-service/internal/domain"
)

type GetDriverLocationQuery struct {
    DriverID string
}

type GetDriverLocationHandler struct {
    locationRepo LocationRepository  // Redis via packages/redis-platform
}

func (h *GetDriverLocationHandler) Handle(
    ctx context.Context,
    q GetDriverLocationQuery,
) (*domain.DriverLocation, error) {
    // Query: NO side effects, read-only
    return h.locationRepo.GetDriverLocation(ctx, q.DriverID)
}
```

---

## 🔌 INFRASTRUCTURE LAYER TEMPLATE

**File:** `services/gps-service/internal/infrastructure/postgres_repo.go`

```go
package infrastructure

import (
    "context"
    "database/sql"
    "github.com/famgo/gps-service/internal/domain"
    "github.com/famgo/gps-service/internal/application"
    "github.com/famgo/packages/database"  // Platform abstraction
)

// Implements LocationRepository interface
type PostgresLocationRepository struct {
    db *sql.DB
}

func (r *PostgresLocationRepository) GetDriverLocation(
    ctx context.Context,
    driverID string,
) (*domain.DriverLocation, error) {
    // Use platform database abstractions
    row := r.db.QueryRowContext(ctx,
        `SELECT id, driver_id, latitude, longitude, accuracy, updated_at
         FROM driver_locations WHERE driver_id = $1`,
        driverID)
    
    loc := &domain.DriverLocation{}
    err := row.Scan(
        &loc.ID,
        &loc.DriverID,
        &loc.Latitude,
        &loc.Longitude,
        &loc.Accuracy,
        &loc.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    
    return loc, nil
}

func (r *PostgresLocationRepository) UpdateDriverLocation(
    ctx context.Context,
    location *domain.DriverLocation,
) error {
    // Transactional update
    _, err := r.db.ExecContext(ctx,
        `INSERT INTO driver_locations (id, driver_id, latitude, longitude, accuracy, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6)
         ON CONFLICT (driver_id) DO UPDATE SET
            latitude = EXCLUDED.latitude,
            longitude = EXCLUDED.longitude,
            accuracy = EXCLUDED.accuracy,
            updated_at = EXCLUDED.updated_at`,
        location.ID,
        location.DriverID,
        location.Latitude,
        location.Longitude,
        location.Accuracy,
        location.UpdatedAt,
    )
    return err
}
```

**File:** `services/gps-service/internal/infrastructure/redis_repo.go`

```go
package infrastructure

import (
    "context"
    "encoding/json"
    "github.com/famgo/gps-service/internal/domain"
    "github.com/famgo/packages/redis-platform"  // Shared Redis platform
)

type RedisLocationRepository struct {
    redis redis-platform.RedisClient
}

func (r *RedisLocationRepository) GetDriverLocation(
    ctx context.Context,
    driverID string,
) (*domain.DriverLocation, error) {
    // Use platform Redis abstraction - not raw redis client
    val, err := r.redis.Get(ctx, "driver:location:"+driverID)
    if err != nil {
        return nil, err
    }
    
    var loc domain.DriverLocation
    err = json.Unmarshal(val, &loc)
    return &loc, err
}

func (r *RedisLocationRepository) CacheLocation(
    ctx context.Context,
    location *domain.DriverLocation,
    ttl time.Duration,
) error {
    // Use platform Redis TTL strategy
    data, _ := json.Marshal(location)
    return r.redis.SetEX(ctx, "driver:location:"+location.DriverID, data, ttl)
}
```

---

## 📡 TRANSPORT LAYER TEMPLATE

**File:** `services/gps-service/internal/transport/grpc_handler.go`

```go
package transport

import (
    "context"
    "github.com/famgo/api/proto/gps"
    "github.com/famgo/gps-service/internal/application"
    "github.com/famgo/packages/telemetry"  // Shared telemetry
)

type GPSServiceServer struct {
    gps.UnimplementedGPSServiceServer
    
    updateLocationHandler *application.UpdateDriverLocationHandler
    metrics               telemetry.Metrics
}

func (s *GPSServiceServer) UpdateDriverLocation(
    ctx context.Context,
    req *gps.UpdateLocationRequest,
) (*gps.LocationResponse, error) {
    // Extract from context (JWT already validated at gateway)
    driverID := extractDriverID(ctx)
    
    // Record metrics
    s.metrics.RecordRequest("UpdateDriverLocation")
    
    // Apply command
    cmd := application.UpdateDriverLocationCommand{
        DriverID:  driverID,
        Latitude:  req.Latitude,
        Longitude: req.Longitude,
        Accuracy:  req.Accuracy,
    }
    
    err := s.updateLocationHandler.Handle(ctx, cmd)
    if err != nil {
        s.metrics.RecordError("UpdateDriverLocation", err)
        return nil, err
    }
    
    return &gps.LocationResponse{
        Success: true,
        Message: "Location updated",
    }, nil
}
```

**File:** `services/gps-service/internal/transport/http_handler.go`

```go
package transport

import (
    "net/http"
    "encoding/json"
    "github.com/famgo/gps-service/internal/application"
    "github.com/famgo/packages/telemetry"  // Shared telemetry
)

type HTTPHandler struct {
    updateLocationHandler *application.UpdateDriverLocationHandler
    metrics              telemetry.Metrics
}

func (h *HTTPHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
    // JWT already validated by gateway
    driverID := extractDriverID(r.Context())
    
    var req struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
        Accuracy  float32 `json:"accuracy"`
    }
    
    json.NewDecoder(r.Body).Decode(&req)
    
    cmd := application.UpdateDriverLocationCommand{
        DriverID:  driverID,
        Latitude:  req.Latitude,
        Longitude: req.Longitude,
        Accuracy:  req.Accuracy,
    }
    
    err := h.updateLocationHandler.Handle(r.Context(), cmd)
    if err != nil {
        h.metrics.RecordError("UpdateLocation", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
```

---

## 📊 INFRASTRUCTURE: REPOSITORIES (Interface Definition)

**File:** `services/gps-service/internal/application/interfaces.go`

```go
package application

import (
    "context"
    "github.com/famgo/gps-service/internal/domain"
)

// Repository interfaces - what the application depends on
// Implementation details are in infrastructure layer

type LocationRepository interface {
    GetDriverLocation(ctx context.Context, driverID string) (*domain.DriverLocation, error)
    UpdateDriverLocation(ctx context.Context, location *domain.DriverLocation) error
    ListActiveLocations(ctx context.Context) ([]domain.DriverLocation, error)
}

type GeofenceRepository interface {
    GetGeofence(ctx context.Context, geofenceID string) (*domain.Geofence, error)
    GetAllGeofences(ctx context.Context) ([]domain.Geofence, error)
    CreateGeofence(ctx context.Context, geofence *domain.Geofence) error
}

type TripRepository interface {
    GetTrip(ctx context.Context, tripID string) (*domain.Trip, error)
    UpdateTrip(ctx context.Context, trip *domain.Trip) error
    AddRoutePoint(ctx context.Context, tripID string, point domain.RoutePoint) error
}
```

---

## 📦 EVENT CONTRACTS (Using shared/contracts/events)

**File:** Service MUST use `shared/contracts/events` - NO local events

```go
// ✅ CORRECT: Using shared contracts
import "github.com/famgo/shared/contracts/events"

type LocationUpdatedEvent struct {
    EventID       string
    EventType     string  // = events.EventTypeDriverLocationUpdated
    Version       int     // = events.VersionV1
    AggregateID   string
    Timestamp     time.Time
    CorrelationID string
    CausationID   string
    Data          LocationUpdatedData
}

// ❌ WRONG: Service-local event
type LocationUpdatedEvent struct {
    DriverID string
    Latitude float64
    // ...
}
```

---

## 🚀 BOOTSTRAP (Dependency Injection)

**File:** `services/gps-service/internal/bootstrap/wire.go`

```go
package bootstrap

import (
    "database/sql"
    "github.com/famgo/gps-service/internal/application"
    "github.com/famgo/gps-service/internal/infrastructure"
    "github.com/famgo/gps-service/internal/transport"
    "github.com/famgo/packages/event-bus"
    "github.com/famgo/packages/telemetry"
    "github.com/famgo/packages/redis-platform"
)

func NewUpdateLocationHandler(
    db *sql.DB,
    redis redis-platform.RedisClient,
    eventBus event-bus.EventBus,
) *application.UpdateDriverLocationHandler {
    // Infrastructure implementations
    locationRepo := infrastructure.NewPostgresLocationRepository(db)
    geofenceRepo := infrastructure.NewPostgresGeofenceRepository(db)
    
    // Domain service
    locationService := &domain.LocationService{}
    
    // Application handler
    return &application.UpdateDriverLocationHandler{
        locationRepo:    locationRepo,
        geofenceRepo:    geofenceRepo,
        eventBus:        eventBus,
        locationService: locationService,
    }
}

func NewGRPCServer(
    db *sql.DB,
    redis redis-platform.RedisClient,
    eventBus event-bus.EventBus,
    metrics telemetry.Metrics,
) *transport.GPSServiceServer {
    handler := NewUpdateLocationHandler(db, redis, eventBus)
    
    return &transport.GPSServiceServer{
        updateLocationHandler: handler,
        metrics:               metrics,
    }
}
```

---

## 🏥 HEALTH CHECKS

Every service MUST have health checks.

**File:** `services/gps-service/internal/transport/health.go`

```go
package transport

import (
    "context"
    "net/http"
    "database/sql"
    "github.com/famgo/packages/redis-platform"
    "github.com/famgo/packages/telemetry"
)

type HealthHandler struct {
    db      *sql.DB
    redis   redis-platform.RedisClient
    metrics telemetry.Metrics
}

// GET /health - Liveness probe (is service alive?)
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "alive"}`))
}

// GET /ready - Readiness probe (can service handle traffic?)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()
    
    // Check database
    if err := h.db.PingContext(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    
    // Check Redis
    if err := h.redis.Ping(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "ready"}`))
}

// GET /startup - Startup probe (did initialization complete?)
func (h *HealthHandler) Startup(w http.ResponseWriter, r *http.Request) {
    // Check all resources initialized
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "started"}`))
}
```

---

## 📈 OBSERVABILITY

Every service MUST record metrics, logs, and traces.

**File:** `services/gps-service/cmd/main.go`

```go
package main

import (
    "github.com/famgo/packages/telemetry"
    "github.com/famgo/gps-service/internal/bootstrap"
)

func main() {
    // Initialize observability FIRST
    metrics := telemetry.NewMetrics("gps-service")
    logs := telemetry.NewLogger("gps-service")
    tracer := telemetry.NewTracer("gps-service")
    
    defer tracer.Shutdown()
    
    logs.Info("Starting GPS service")
    
    // Initialize infrastructure
    db := setupDatabase()
    redis := setupRedis()
    eventBus := setupEventBus()
    
    // Create handlers with observability
    grpcServer := bootstrap.NewGRPCServer(db, redis, eventBus, metrics)
    
    // Start server with health checks
    healthHandler := &transport.HealthHandler{db: db, redis: redis, metrics: metrics}
    
    metrics.RecordServiceStart("gps-service")
    
    // Server...
}
```

---

## ✅ VALIDATION

Every service MUST validate inputs.

**File:** `services/gps-service/internal/application/validation.go`

```go
package application

import "github.com/go-playground/validator/v10"

func ValidateUpdateLocationCommand(cmd UpdateDriverLocationCommand) error {
    return validator.New().Struct(cmd)
}

// In handler:
func (h *Handler) Handle(ctx context.Context, cmd UpdateDriverLocationCommand) error {
    // Validate input FIRST
    if err := ValidateUpdateLocationCommand(cmd); err != nil {
        return err
    }
    
    // Then execute
    // ...
}
```

---

## 🧪 TESTING

Every service MUST have comprehensive tests.

**File:** `services/gps-service/tests/unit/location_service_test.go`

```go
package unit

import (
    "testing"
    "github.com/famgo/gps-service/internal/domain"
)

func TestLocationService_IsWithinGeofence(t *testing.T) {
    service := &domain.LocationService{}
    
    location := domain.DriverLocation{
        Latitude:  37.7749,
        Longitude: -122.4194,
    }
    
    geofence := domain.Geofence{
        Latitude:  37.7749,
        Longitude: -122.4194,
        Radius:    100, // 100 meters
    }
    
    result := service.IsWithinGeofence(location, geofence)
    
    if !result {
        t.Errorf("Expected true, got %v", result)
    }
}
```

---

## 🎯 SERVICE COMPLETION CHECKLIST

For each service (GPS, User, Ride):

- [ ] Domain layer complete (entities, aggregates, services)
- [ ] Application layer complete (commands, queries, handlers)
- [ ] Infrastructure layer complete (repos, external clients)
- [ ] Transport layer complete (HTTP, gRPC, WebSocket)
- [ ] Database schema and migrations
- [ ] API contracts (proto, OpenAPI)
- [ ] Event publishing (using shared/contracts/events)
- [ ] Health checks (live, ready, startup)
- [ ] Metrics recording (Prometheus)
- [ ] Trace propagation (Jaeger)
- [ ] Structured logging (JSON)
- [ ] Input validation
- [ ] Unit tests (>80% coverage)
- [ ] Integration tests
- [ ] Dockerfile
- [ ] Kubernetes manifests
- [ ] Helm charts
- [ ] README with architecture
- [ ] API documentation

---

**SERVICE COMPLETION TEMPLATES ESTABLISHED** ✅

All services must follow this structure.
All services must use platform abstractions.
All services must use shared contracts.

