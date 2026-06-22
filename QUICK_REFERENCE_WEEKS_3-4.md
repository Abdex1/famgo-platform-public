# 🎯 QUICK REFERENCE: WEEKS 3-4 DEVELOPMENT GUIDE

**Purpose:** Quick lookup for developers during implementation  
**Scope:** Essential rules, patterns, and commands  
**Format:** Copy-paste ready code and quick reference

---

## ⚡ THE FIVE RULES

### Rule 1: Events from Shared Contracts ONLY
```go
// ❌ WRONG - Service-local event
type RideCreatedEvent struct { ... }
kafka.Publish("ride.created", event)

// ✅ CORRECT - Use shared contracts
import "github.com/famgo/shared/contracts/events"
eventBus.Publish(ctx, events.RideCreatedEvent{
    EventID: uuid.New().String(),
    EventType: events.EventTypeRideCreated,
    Version: events.VersionV1,
    AggregateID: rideID,
    Timestamp: time.Now(),
    Data: rideData,
})
```

### Rule 2: SDKs from Packages ONLY
```go
// ❌ WRONG - Raw Kafka
import "github.com/segmentio/kafka-go"
r := kafka.NewReader(...)

// ✅ CORRECT - Use SDK
import "github.com/famgo/packages/kafka-sdk"
client := kafkasdk.NewClient(config)
```

### Rule 3: Platform Abstractions Required
```go
// ❌ WRONG - Custom event-bus
type EventBus struct { ... }

// ✅ CORRECT - Use platform
import "github.com/famgo/platform/event-bus"
eventBus := event-bus.New()
```

### Rule 4: Reference Architecture Pattern
```
Every service:
internal/
  ├── domain/          (Pure logic, zero external deps)
  ├── application/     (Commands, queries, handlers)
  ├── infrastructure/  (Repos, external clients)
  └── transport/       (HTTP, gRPC, WebSocket)
```

### Rule 5: No Cross-Service Database Writes
```go
// ❌ WRONG - ride-service writes to wallet
db.Exec("UPDATE wallets SET balance = balance - ?", amount)

// ✅ CORRECT - wallet-service owns wallet table
// ride-service calls:
paymentService.DeductFare(ctx, userID, amount)
// Inside paymentService: writes to its own tables
```

---

## 🏗️ SERVICE STRUCTURE TEMPLATE

### Directory Layout
```bash
services/gps-service/
├── cmd/
│   └── main.go                # Entrypoint
├── internal/
│   ├── domain/
│   │   ├── entities.go        # Aggregates, entities
│   │   ├── value_objects.go   # Value objects
│   │   └── services.go        # Domain services (pure logic)
│   ├── application/
│   │   ├── commands.go        # Commands
│   │   ├── queries.go         # Queries
│   │   ├── handlers.go        # Handlers
│   │   └── interfaces.go      # Repo interfaces
│   ├── infrastructure/
│   │   ├── postgres.go        # Postgres repos
│   │   ├── redis.go           # Redis repos
│   │   └── external.go        # External service clients
│   └── transport/
│       ├── http.go            # HTTP handlers
│       ├── grpc.go            # gRPC handlers
│       ├── websocket.go       # WebSocket handlers
│       └── health.go          # Health checks
├── api/
│   ├── proto/
│   │   └── gps.proto          # gRPC contract
│   └── openapi.yaml           # REST contract
├── db/
│   ├── migrations/
│   │   ├── 001_create_schema.up.sql
│   │   └── 001_create_schema.down.sql
│   └── schema.sql
├── config/
│   ├── .env.example
│   ├── .env.local
│   ├── .env.development
│   ├── .env.staging
│   └── .env.production
├── tests/
│   ├── unit/
│   │   └── *_test.go
│   ├── integration/
│   │   └── *_test.go
│   └── contract/
│       └── *_contract_test.go
├── deployments/
│   ├── Dockerfile
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── hpa.yaml
│   └── helm/
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## 🔌 DEPENDENCY INJECTION PATTERN

```go
// services/gps-service/internal/bootstrap/wire.go

package bootstrap

import (
    "database/sql"
    "github.com/famgo/packages/event-bus"
    "github.com/famgo/packages/telemetry"
    "github.com/famgo/packages/redis-platform"
)

// Wire up all dependencies
func NewGPSServer(
    db *sql.DB,
    redis redis-platform.RedisClient,
    eventBus event-bus.EventBus,
    metrics telemetry.Metrics,
) *transport.GPSServiceServer {
    
    // Infrastructure layer
    locationRepo := infrastructure.NewPostgresLocationRepo(db)
    
    // Domain layer
    locationService := &domain.LocationService{}
    
    // Application layer
    updateHandler := &application.UpdateDriverLocationHandler{
        locationRepo:    locationRepo,
        eventBus:        eventBus,
        locationService: locationService,
    }
    
    // Transport layer
    return &transport.GPSServiceServer{
        updateLocationHandler: updateHandler,
        metrics:               metrics,
    }
}
```

---

## 📝 COMMAND HANDLER PATTERN

```go
// services/gps-service/internal/application/commands.go

package application

type UpdateLocationCommand struct {
    DriverID  string  `validate:"required,uuid"`
    Latitude  float64 `validate:"required,min=-90,max=90"`
    Longitude float64 `validate:"required,min=-180,max=180"`
}

type UpdateLocationHandler struct {
    locationRepo LocationRepository
    eventBus     event-bus.EventBus
    service      *domain.LocationService
}

func (h *UpdateLocationHandler) Handle(
    ctx context.Context,
    cmd UpdateLocationCommand,
) error {
    // 1. Validate
    if err := validator.Validate(cmd); err != nil {
        return err
    }
    
    // 2. Get current state
    oldLoc, _ := h.locationRepo.GetDriverLocation(ctx, cmd.DriverID)
    
    // 3. Apply domain logic
    newLoc := &domain.DriverLocation{
        ID:        uuid.New().String(),
        DriverID:  cmd.DriverID,
        Latitude:  cmd.Latitude,
        Longitude: cmd.Longitude,
        UpdatedAt: time.Now(),
    }
    
    // 4. Persist
    if err := h.locationRepo.UpdateDriverLocation(ctx, newLoc); err != nil {
        return err
    }
    
    // 5. Publish event (THROUGH SHARED CONTRACTS)
    return h.eventBus.PublishIdempotent(ctx, events.DriverLocationUpdatedEvent{
        EventID:    uuid.New().String(),
        EventType:  events.EventTypeDriverLocationUpdated,
        Version:    events.VersionV1,
        AggregateID: cmd.DriverID,
        Timestamp:  time.Now(),
        Data: map[string]interface{}{
            "old_location": oldLoc,
            "new_location": newLoc,
        },
    })
}
```

---

## 📊 HEALTH CHECK PATTERN

```go
// services/gps-service/internal/transport/health.go

package transport

import (
    "context"
    "net/http"
    "time"
    "database/sql"
    "github.com/famgo/packages/redis-platform"
)

type HealthHandler struct {
    db    *sql.DB
    redis redis-platform.RedisClient
}

// GET /health - Liveness probe
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"alive"}`))
}

// GET /ready - Readiness probe
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
    w.Write([]byte(`{"status":"ready"}`))
}

// GET /startup - Startup probe
func (h *HealthHandler) Startup(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"started"}`))
}
```

---

## 📈 METRICS PATTERN

```go
// In any handler:
import "github.com/famgo/packages/telemetry"

type GPSServiceServer struct {
    metrics telemetry.Metrics
}

func (s *GPSServiceServer) UpdateLocation(ctx context.Context, req *Request) (*Response, error) {
    // Record request
    s.metrics.RecordRequest("UpdateLocation")
    
    // Do work
    if err := doWork(); err != nil {
        // Record error
        s.metrics.RecordError("UpdateLocation", err)
        return nil, err
    }
    
    // Record success
    s.metrics.RecordSuccess("UpdateLocation")
    return response, nil
}
```

---

## 🔍 EVENT PUBLISHING PATTERN

```go
// CORRECT: Use event-bus from packages
import (
    "github.com/famgo/packages/event-bus"
    "github.com/famgo/shared/contracts/events"
)

// In your handler:
err := h.eventBus.PublishIdempotent(ctx, events.RideCreatedEvent{
    EventID:    uuid.New().String(),
    EventType:  events.EventTypeRideCreated,
    Version:    events.VersionV1,
    AggregateID: rideID,
    Timestamp:  time.Now(),
    Data: map[string]interface{}{
        "pickup_location": pickupLoc,
        "dropoff_location": dropoffLoc,
    },
})
```

---

## 🧪 TEST PATTERN

```go
// tests/unit/location_service_test.go

package unit

import (
    "testing"
    "github.com/famgo/gps-service/internal/domain"
)

func TestLocationService_IsWithinGeofence(t *testing.T) {
    service := &domain.LocationService{}
    
    tests := []struct {
        name     string
        location domain.DriverLocation
        geofence domain.Geofence
        expected bool
    }{
        {
            name: "driver inside geofence",
            location: domain.DriverLocation{Latitude: 37.7749, Longitude: -122.4194},
            geofence: domain.Geofence{Latitude: 37.7749, Longitude: -122.4194, Radius: 100},
            expected: true,
        },
        {
            name: "driver outside geofence",
            location: domain.DriverLocation{Latitude: 37.8, Longitude: -122.4},
            geofence: domain.Geofence{Latitude: 37.7749, Longitude: -122.4194, Radius: 100},
            expected: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := service.IsWithinGeofence(tt.location, tt.geofence)
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

---

## 🚀 COMMON COMMANDS

```bash
# Audit repository
find shared/contracts/ -type f -name "*.go" | xargs grep "type.*Event"

# Verify no duplicate SDKs
find services/ -type f -name "*.go" | xargs grep "kafka.Dial"
find services/ -type f -name "*.go" | xargs grep "NewKafkaReader"

# Run tests for service
cd services/gps-service && go test -v -cover ./...

# Build Docker image
docker build -f services/gps-service/Dockerfile -t gps-service:latest services/gps-service

# Test Kubernetes manifest
kubectl apply -f services/gps-service/deployments/ --dry-run=client

# Deploy service
kubectl apply -f services/gps-service/deployments/

# Check health
kubectl get pods -l app=gps-service
kubectl describe pod <pod-name>
kubectl logs <pod-name>

# Port forward for testing
kubectl port-forward svc/gps-service 8080:80
curl http://localhost:8080/health

# View metrics
curl http://localhost:8080/metrics

# Check traces
# Visit Jaeger UI at localhost:16686
```

---

## 📋 CHECKLIST: Before Committing

- [ ] All events use shared/contracts/events
- [ ] All SDKs from packages/ (not raw libraries)
- [ ] All platform abstractions used
- [ ] Tests passing (>80% coverage)
- [ ] Code follows auth-service pattern
- [ ] Health checks implemented
- [ ] Metrics recorded
- [ ] Traces propagated
- [ ] Logs structured (JSON)
- [ ] No cross-service database writes
- [ ] No service boundary violations
- [ ] Dockerfile builds
- [ ] Kubernetes manifests apply
- [ ] README complete

---

## 🔗 QUICK LINKS

**Architecture Reference:**
```
services/auth-service/         ← Follow this pattern
```

**Event Registry:**
```
shared/contracts/events/catalog/events.go
```

**Package Registry:**
```
packages/kafka-sdk/
packages/event-bus/
packages/telemetry/
packages/redis-platform/
```

**Platform Registry:**
```
platform/event-bus/
platform/saga/
platform/feature-flags/
```

---

## ⚠️ COMMON MISTAKES TO AVOID

❌ **MISTAKE 1:** Publishing events directly to Kafka
✅ **FIX:** Use eventBus.Publish() from packages/event-bus

❌ **MISTAKE 2:** Creating service-local events
✅ **FIX:** Use events from shared/contracts/events

❌ **MISTAKE 3:** Writing to other service's database
✅ **FIX:** Call service via gRPC or event

❌ **MISTAKE 4:** Not following auth-service pattern
✅ **FIX:** Copy structure exactly (domain, application, infrastructure, transport)

❌ **MISTAKE 5:** No metrics, logs, or traces
✅ **FIX:** Add observability FIRST (not last)

---

**QUICK REFERENCE COMPLETE** ✅

Print this and use during implementation.

