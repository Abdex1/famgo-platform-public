# 📦 PACKAGES USAGE GUIDE: SDK Inventory & Patterns

**Status:** Comprehensive audit of `packages/`  
**Created:** Day 1 Afternoon of Weeks 3-4 Audit Phase  
**Repository:** github.com/Abdex1/FamGo-platform  
**Purpose:** How to use each SDK (NEVER implement duplicates)

---

## 🎯 CRITICAL RULE

✅ **MANDATORY:** All services MUST use these packages  
❌ **FORBIDDEN:** Services must NEVER implement custom versions

---

## 📦 PACKAGES QUICK REFERENCE

| Package | Purpose | NEVER DO | DO THIS |
|---------|---------|----------|---------|
| kafka-sdk | Kafka wrapper | Raw kafka-go imports | `kafkasdk.NewClient()` |
| event-bus | Event publishing | Direct Kafka publishing | `eventbus.PublishIdempotent()` |
| telemetry | Metrics/traces/logs | Custom prometheus/logging | `telemetry.NewMetrics()` |
| redis-platform | Redis wrapper | Raw redis-go imports | `redisplatform.NewClient()` |
| auth-client | Auth service calls | HTTP calls to auth-service | `authclient.ValidateToken()` |
| grpc-clients | Service-to-service gRPC | Manual gRPC dialup | `grpcclient.NewGPSClient()` |
| payment-sdk | Payment processing | Direct Stripe/PayPal API | `paymentsdk.CreatePaymentIntent()` |
| vault-sdk | Secrets management | Env vars / config files | `vaultsdk.GetSecret()` |
| websocket-sdk | WebSocket connections | Raw gorilla/websocket | `websocketsdk.NewConn()` |

---

## 1️⃣ KAFKA-SDK: Kafka Wrapper

**Location:** `packages/kafka-sdk/`  
**When to Use:** Any Kafka producer/consumer operations

### ❌ WRONG - DO NOT USE:
```go
import "github.com/segmentio/kafka-go"
reader := kafka.NewReader(kafka.ReaderConfig{...})
writer := kafka.NewWriter(kafka.WriterConfig{...})
```

### ✅ CORRECT - USE THIS:
```go
import "github.com/Abdex1/FamGo-platform/packages/kafka-sdk"

client := kafkasdk.NewClient(&kafkasdk.Config{
    Brokers:  []string{"kafka:9092"},
    GroupID:  "ride-service-group",
    LogLevel: "info",
})

// Produce
err := client.Produce(ctx, &kafkasdk.Message{
    Topic: "ride.lifecycle",
    Key:   rideID,
    Value: rideEvent,
})

// Consume
for msg := range client.Consume(ctx, []string{"ride.lifecycle"}) {
    processMessage(msg)
    client.Commit(ctx, msg)
}

defer client.Close()
```

---

## 2️⃣ EVENT-BUS: Event Publishing

**Location:** `packages/event-bus/`  
**When to Use:** Publishing events to other services

### ❌ WRONG - DO NOT USE:
```go
// Direct Kafka
kafka.Produce("ride.lifecycle", event)

// Custom event-bus
type EventBus struct { ... }
eventBus.Publish(event)
```

### ✅ CORRECT - USE THIS:
```go
import (
    "github.com/Abdex1/FamGo-platform/packages/event-bus"
    "github.com/Abdex1/FamGo-platform/shared/contracts/events"
)

bus := eventbus.NewEventBus(&eventbus.Config{
    KafkaClient: kafkaClient,
    Telemetry:   telemetrySDK,
})

// Publish with idempotency
err := bus.PublishIdempotent(ctx, events.RideRequestedEvent{
    EventID:       uuid.New().String(),
    EventType:     events.EventTypeRideRequested,
    EventVersion:  events.VersionV1,
    AggregateID:   rideID,
    Timestamp:     time.Now(),
    CorrelationID: correlationID,
    Data: map[string]interface{}{
        "rider_id":   riderID,
        "pickup":     pickupLoc,
        "dropoff":    dropoffLoc,
    },
})

// Subscribe
handler := func(ctx context.Context, evt *events.Event) error {
    // Process event
    return nil
}
bus.Subscribe(ctx, events.EventTypeRideAssigned, handler)
```

---

## 3️⃣ TELEMETRY: Metrics, Traces, Logs

**Location:** `packages/telemetry/`  
**When to Use:** Every service for observability

### ❌ WRONG - DO NOT USE:
```go
// Direct Prometheus
import "github.com/prometheus/client_golang/prometheus"
counter := prometheus.NewCounter(...)

// Raw logging
log.Println("something happened")

// Raw OpenTelemetry
import "go.opentelemetry.io/otel"
tracer := otel.Tracer("my-service")
```

### ✅ CORRECT - USE THIS:
```go
import "github.com/Abdex1/FamGo-platform/packages/telemetry"

// Initialize
metrics := telemetry.NewMetrics("ride-service")
logger := telemetry.NewLogger("ride-service")
tracer := telemetry.NewTracer("ride-service")

// Record metrics
defer func(start time.Time) {
    duration := time.Since(start)
    metrics.RecordLatency("CreateRide", duration)
    if err != nil {
        metrics.RecordError("CreateRide", err)
    }
}(time.Now())

// Structured logging (JSON output)
logger.Info("creating ride", map[string]interface{}{
    "rider_id": riderID,
    "trace_id": traceID,
})

// Tracing
ctx, span := tracer.Start(ctx, "CreateRide")
defer span.End()

// Propagate trace to other services
headers := tracer.InjectHeaders(ctx)
```

**Metrics Exported at /metrics:**
```
ride_service_requests_total{method="CreateRide"} 1234
ride_service_request_duration_seconds{method="CreateRide"} 0.125
ride_service_errors_total{method="CreateRide"} 5
ride_service_rides_created_total 1000
```

---

## 4️⃣ REDIS-PLATFORM: Redis Wrapper

**Location:** `packages/redis-platform/`  
**When to Use:** Caching, real-time data, sessions

### ❌ WRONG - DO NOT USE:
```go
import "github.com/go-redis/redis/v8"
client := redis.NewClient(&redis.Options{...})
val, _ := client.Get(ctx, "some_random_key").Result()
```

### ✅ CORRECT - USE THIS:
```go
import "github.com/Abdex1/FamGo-platform/packages/redis-platform"

redis := redisplatform.NewClient(&redisplatform.Config{
    Host: "redis",
    Port: 6379,
})

// Key naming: {service}:{entity}:{id}
err := redis.SetEX(ctx, "gps:location:driver_123", locationData, 5*time.Minute)
val, err := redis.Get(ctx, "gps:location:driver_123")

// Delete
err := redis.Delete(ctx, "gps:location:driver_123")

// Cache with function
result, err := redis.CacheFunc(ctx, "ride:estimate:123", 10*time.Minute, func() (interface{}, error) {
    return calculateFare(ctx)
})

// Increment counter
err := redis.Increment(ctx, "ride:counter:2024-01", 1)

// Geo operations (for location)
err := redis.GeoAdd(ctx, "drivers:locations", lat, lon, "driver_123")
nearby, err := redis.GeoRadius(ctx, "drivers:locations", lat, lon, 5000) // 5km
```

---

## 5️⃣ AUTH-CLIENT: Auth Service Client

**Location:** `packages/auth-client/`  
**When to Use:** JWT validation, user lookups, RBAC checks

### ❌ WRONG - DO NOT USE:
```go
// HTTP call without abstraction
import "net/http"
resp, _ := http.Post("http://auth-service/verify", ...)

// Custom JWT validation
import "github.com/dgrijalva/jwt-go"
// Custom code to validate JWT
```

### ✅ CORRECT - USE THIS:
```go
import "github.com/Abdex1/FamGo-platform/packages/auth-client"

authClient := authclient.NewClient(&authclient.Config{
    AuthServiceURL: "http://auth-service",
    CacheTTL:       5 * time.Minute,
})

// Validate JWT token (with caching)
claims, err := authClient.ValidateToken(ctx, token)
userID := claims.UserID

// Check permissions
permitted, _ := authClient.HasPermission(ctx, userID, "ride:create")

// Get user
user, _ := authClient.GetUser(ctx, userID)

// Get roles
roles, _ := authClient.GetUserRoles(ctx, userID)
```

---

## 6️⃣ GRPC-CLIENTS: Service-to-Service gRPC

**Location:** `packages/grpc-clients/`  
**When to Use:** Calling other services (GPS, Pricing, Dispatch, etc.)

### ❌ WRONG - DO NOT USE:
```go
import "google.golang.org/grpc"
conn, _ := grpc.Dial("gps-service:50051")
client := gps.NewGPSServiceClient(conn)

// Or HTTP without abstraction
http.Get("http://gps-service/api/locations")
```

### ✅ CORRECT - USE THIS:
```go
import "github.com/Abdex1/FamGo-platform/packages/grpc-clients"

gpsClient := grpcclient.NewGPSClient(&grpcclient.Config{
    ServiceURL: "gps-service:50051",
    Timeout:    5 * time.Second,
})

pricingClient := grpcclient.NewPricingClient(&grpcclient.Config{
    ServiceURL: "pricing-service:50051",
    Timeout:    3 * time.Second,
})

// Call services
location, _ := gpsClient.GetDriverLocation(ctx, &GetLocationRequest{
    DriverID: driverID,
})

fare, _ := pricingClient.CalculateFare(ctx, &FareRequest{
    Pickup:  pickupLoc,
    Dropoff: dropoffLoc,
})
```

---

## 7️⃣ PAYMENT-SDK: Payment Gateway

**Location:** `packages/payment-sdk/`  
**When to Use:** Processing payments, refunds

### ❌ WRONG - DO NOT USE:
```go
// Direct Stripe
import "github.com/stripe/stripe-go"
// Custom Stripe integration

// Direct PayPal
// Custom PayPal integration
```

### ✅ CORRECT - USE THIS:
```go
import "github.com/Abdex1/FamGo-platform/packages/payment-sdk"

paymentService := paymentsdk.NewPaymentService(&paymentsdk.Config{
    Provider: "stripe",
    APIKey:   os.Getenv("PAYMENT_API_KEY"),
})

// Create payment intent
intent, _ := paymentService.CreatePaymentIntent(ctx, &paymentsdk.PaymentRequest{
    Amount:      fare,
    Currency:    "USD",
    CustomerID:  userID,
    Description: "Ride fare for " + rideID,
})

// Process payment
result, _ := paymentService.ProcessPayment(ctx, &paymentsdk.ProcessRequest{
    IntentID:        intent.ID,
    PaymentMethodID: pmID,
})

// Refund
refund, _ := paymentService.RefundPayment(ctx, result.TransactionID, refundAmount)
```

---

## 8️⃣ VAULT-SDK: Secrets Management

**Location:** `packages/vault-sdk/`  
**When to Use:** Getting secrets (API keys, passwords)

### ❌ WRONG - DO NOT USE:
```go
// Environment variables
dbPassword := os.Getenv("DB_PASSWORD")

// Config files
// password: "secret123" in config.yaml

// Direct Vault
import vault "github.com/hashicorp/vault/api"
```

### ✅ CORRECT - USE THIS:
```go
import "github.com/Abdex1/FamGo-platform/packages/vault-sdk"

vaultClient := vaultsdk.NewVaultClient(&vaultsdk.Config{
    VaultAddr: os.Getenv("VAULT_ADDR"),
    RoleName:  "ride-service",
})

// Get database credentials
dbCreds, _ := vaultClient.GetDatabaseCredentials(ctx)
dbURL := fmt.Sprintf("postgres://%s:%s@db:5432/famgo", 
    dbCreds.Username, 
    dbCreds.Password)

// Get API keys
stripeKey, _ := vaultClient.GetSecret(ctx, "stripe/api-key")

// Rotate secrets
_ = vaultClient.RotateSecret(ctx, "db/password")
```

---

## 9️⃣ WEBSOCKET-SDK: WebSocket Client

**Location:** `packages/websocket-sdk/`  
**When to Use:** Real-time connections (GPS tracking, notifications)

### ❌ WRONG - DO NOT USE:
```go
import "github.com/gorilla/websocket"
dialer := websocket.Dialer{}
conn, _ := dialer.Dial("ws://gateway:8080/ws", nil)
```

### ✅ CORRECT - USE THIS:
```go
import "github.com/Abdex1/FamGo-platform/packages/websocket-sdk"

wsClient := websocketsdk.NewClient(&websocketsdk.Config{
    URL:     "ws://gateway:8080/ws",
    Timeout: 10 * time.Second,
})

// Connect
conn, _ := wsClient.Connect(ctx)

// Send message
_ = conn.SendJSON(ctx, map[string]interface{}{
    "type": "location_update",
    "data": locationData,
})

// Receive message
msg, _ := conn.ReceiveJSON(ctx)

// Close
conn.Close()
```

---

## 🎯 USAGE PATTERNS BY SERVICE LAYER

### Domain Layer (internal/domain/)
❌ NO package imports  
❌ Pure logic only

### Application Layer (internal/application/)
✅ `event-bus` - For publishing events
✅ `telemetry` - For recording operation metrics

### Infrastructure Layer (internal/infrastructure/)
✅ `kafka-sdk` - For low-level Kafka operations
✅ `redis-platform` - For caching/storage
✅ `auth-client` - For user lookups
✅ `grpc-clients` - For service calls
✅ `payment-sdk` - For payment operations
✅ `vault-sdk` - For secrets retrieval

### Transport Layer (internal/transport/)
✅ `telemetry` - For recording HTTP/gRPC metrics
✅ `websocket-sdk` - For WebSocket connections

---

## ✅ VERIFICATION CHECKLIST

**For Each Service:**

- [ ] NO raw kafka-go imports (use kafka-sdk)
- [ ] NO raw redis imports (use redis-platform)
- [ ] NO raw prometheus imports (use telemetry)
- [ ] NO direct HTTP calls to other services (use grpc-clients)
- [ ] NO custom JWT validation (use auth-client)
- [ ] NO secrets in environment variables (use vault-sdk)
- [ ] NO direct payment gateway API (use payment-sdk)
- [ ] NO raw WebSocket imports (use websocket-sdk)
- [ ] ALL events from shared/contracts/events (use event-bus)

---

## 🚀 NEXT STEPS

1. Days 2-4: Continue audit of remaining packages
2. Days 5-9: Use templates when building services
3. Days 9-10: Verify all services use correct packages

---

**PACKAGE_USAGE_GUIDE AUDIT COMPLETE** ✅

Repository: github.com/Abdex1/FamGo-platform  
All 9 core packages documented with usage patterns.  
Ready for service implementation.

