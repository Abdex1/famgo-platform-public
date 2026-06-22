# FamGo Go Service Template
## Production-Grade Microservice Foundation for Go Services

---

## Overview

This is the **standard Go service template** for all FamGo platform microservices. Use this for:

- ✅ auth-service
- ✅ ride-service
- ✅ dispatch-service
- ✅ pooling-service
- ✅ gps-service
- ✅ payment-service
- ✅ wallet-service
- ✅ safety-service
- ✅ fraud-service
- ✅ pricing-service

**NOT for:** websocket-gateway, notification-service (use Node.js template), ML services (use Python template)

---

## Architecture

```
cmd/service/
  └── main.go              # Entry point, gRPC + REST server setup

internal/
  ├── domain/              # Business logic (service interfaces, entities)
  ├── infrastructure/      # External dependencies (DB, cache, Kafka)
  ├── handlers/            # API handlers (gRPC + REST)
  ├── repositories/        # Data access layer
  └── config/              # Configuration management

api/proto/
  └── service.proto        # gRPC service definitions

migrations/
  └── 001_init.sql        # Database migrations

tests/
  ├── unit/               # Unit tests
  └── integration/        # Integration tests
```

---

## Quick Start

### 1. Create Service from Template

```bash
# Copy template
cp -r services/_template-go services/my-service
cd services/my-service

# Copy environment
cp .env.example .env
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Development

```bash
make dev              # Start development server
make test             # Run tests
make build            # Build binary
```

### 4. Docker

```bash
make docker-build     # Build image
make docker-run       # Run container
```

---

## Configuration

All services use 12-factor app configuration via environment variables:

```bash
# Service
SERVICE_NAME=my-service
PORT=5001 (gRPC), 5002 (REST)
LOG_LEVEL=info

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=famgo

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=my-service

# Observability
JAEGER_ENABLED=false
JAEGER_ENDPOINT=http://localhost:14268/api/traces

# gRPC
GRPC_PORT=5001
REST_PORT=8080
```

---

## Service Structure

### Domain Layer (`internal/domain/`)

Business logic, independent of infrastructure:

```go
// domain/service.go
type ServiceInterface interface {
    GetEntity(ctx context.Context, id string) (*Entity, error)
    CreateEntity(ctx context.Context, entity *Entity) error
    UpdateEntity(ctx context.Context, entity *Entity) error
    DeleteEntity(ctx context.Context, id string) error
}

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}
```

### Infrastructure Layer (`internal/infrastructure/`)

Database, cache, Kafka, etc:

```go
// infrastructure/postgres/repository.go
type PostgresRepository struct {
    db *sql.DB
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Entity, error) {
    // Implementation
}

// infrastructure/redis/cache.go
type RedisCache struct {
    client *redis.Client
}

// infrastructure/kafka/consumer.go
type KafkaConsumer struct {
    consumer sarama.ConsumerGroup
}

func (c *KafkaConsumer) Consume(ctx context.Context, topics []string) error {
    // Implementation
}
```

### Handlers Layer (`internal/handlers/`)

API endpoints (gRPC + REST):

```go
// handlers/grpc.go - gRPC implementation
type ServiceServer struct {
    service domain.ServiceInterface
}

func (s *ServiceServer) GetEntity(ctx context.Context, req *pb.GetEntityRequest) (*pb.Entity, error) {
    entity, err := s.service.GetEntity(ctx, req.Id)
    // Convert to protobuf and return
}

// handlers/rest.go - REST adapter
func (h *RestHandler) GetEntity(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    entity, err := h.service.GetEntity(r.Context(), id)
    json.NewEncoder(w).Encode(entity)
}
```

### Repositories Layer (`internal/repositories/`)

Data access abstractions:

```go
// repositories/interface.go
type Repository interface {
    Create(ctx context.Context, entity *Entity) error
    GetByID(ctx context.Context, id string) (*Entity, error)
    Update(ctx context.Context, entity *Entity) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context) ([]*Entity, error)
}
```

---

## Protocol Buffers (gRPC)

Define service contracts in `api/proto/service.proto`:

```protobuf
syntax = "proto3";

package famgo.service;

service ServiceAPI {
    rpc GetEntity(GetEntityRequest) returns (Entity);
    rpc CreateEntity(CreateEntityRequest) returns (Entity);
    rpc UpdateEntity(UpdateEntityRequest) returns (Entity);
    rpc DeleteEntity(DeleteEntityRequest) returns (Empty);
}

message Entity {
    string id = 1;
    string name = 2;
    int64 created_at = 3;
}

message GetEntityRequest {
    string id = 1;
}
```

Generate code:

```bash
protoc --go_out=. --go-grpc_out=. api/proto/service.proto
```

---

## Database Migrations

Use SQL migrations in `migrations/`:

```sql
-- migrations/001_init.sql

CREATE TABLE entities (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_entities_name ON entities(name);
```

Run migrations:

```bash
# Using migrate tool
migrate -path migrations -database "postgres://..." up

# Or in code
migrate.Up()
```

---

## Kafka Event Publishing

Publish events from your service:

```go
// internal/infrastructure/kafka/producer.go

type KafkaProducer struct {
    producer sarama.AsyncProducer
}

func (p *KafkaProducer) PublishEvent(ctx context.Context, topic string, event interface{}) error {
    data, _ := json.Marshal(event)
    
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.StringEncoder(data),
        Headers: []sarama.RecordHeader{
            {
                Key:   []byte("trace_id"),
                Value: []byte(extractTraceID(ctx)),
            },
        },
    }
    
    p.producer.Input() <- msg
    return nil
}
```

---

## Kafka Event Consumption

Subscribe to events:

```go
// internal/infrastructure/kafka/consumer.go

type EventConsumer struct {
    consumer sarama.ConsumerGroup
    handlers map[string]EventHandler
}

func (c *EventConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for {
        select {
        case message := <-claim.Messages():
            topic := message.Topic
            if handler, ok := c.handlers[topic]; ok {
                handler.Handle(message.Value)
            }
            session.MarkMessage(message, "")
        case <-session.Context().Done():
            return nil
        }
    }
}

// Register handlers
consumer.handlers["ride.created"] = &RideCreatedHandler{}
consumer.handlers["payment.completed"] = &PaymentCompletedHandler{}
```

---

## Observability

### Logging

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

logger.InfoContext(ctx, "processing request",
    slog.String("trace_id", traceID),
    slog.String("user_id", userID),
    slog.Duration("latency", duration),
)
```

### Tracing

```go
import "go.opentelemetry.io/otel"

tracer := otel.Tracer("my-service")

ctx, span := tracer.Start(ctx, "ProcessRequest")
defer span.End()

span.SetAttribute("user.id", userID)
span.RecordError(err)
```

### Metrics

```go
import "go.opentelemetry.io/otel/metric"

meter := otel.Meter("my-service")

counter, _ := meter.Int64Counter("requests_total")
counter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("method", "GET"),
    attribute.String("status", "200"),
))
```

---

## Testing

### Unit Tests

```go
// internal/domain/service_test.go

func TestGetEntity(t *testing.T) {
    mockRepo := &MockRepository{}
    service := domain.NewService(mockRepo)
    
    entity, err := service.GetEntity(context.Background(), "123")
    
    assert.NoError(t, err)
    assert.Equal(t, "123", entity.ID)
}
```

### Integration Tests

```go
// tests/integration/service_test.go

func TestCreateEntityIntegration(t *testing.T) {
    db := setupTestDB()
    repo := repositories.NewPostgresRepository(db)
    service := domain.NewService(repo)
    
    entity := &domain.Entity{ID: "123", Name: "Test"}
    err := service.CreateEntity(context.Background(), entity)
    
    assert.NoError(t, err)
    
    retrieved, _ := service.GetEntity(context.Background(), "123")
    assert.Equal(t, "Test", retrieved.Name)
}
```

---

## Docker

Multi-stage build in `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/app /app
EXPOSE 5001 8080
CMD ["/app"]
```

Build and run:

```bash
docker build -t famgo/my-service:latest .
docker run -p 5001:5001 -p 8080:8080 famgo/my-service:latest
```

---

## Dependencies

Core Go modules for FamGo services:

```go
require (
    google.golang.org/grpc v1.59.0
    google.golang.org/protobuf v1.31.0
    github.com/lib/pq v1.10.9              // PostgreSQL
    github.com/redis/go-redis/v9 v9.3.0    // Redis
    github.com/Shopify/sarama v1.38.1      // Kafka
    go.opentelemetry.io/otel v1.20.0       // Tracing
    go.opentelemetry.io/otel/exporters/jaeger v1.20.0
    go.uber.org/zap v1.26.0                // Logging
)
```

---

## Service Boundaries (STRICT)

### DO implement in your service:
- Your service's core business logic
- Event publishing to Kafka
- Event subscriptions to Kafka
- Your service's database tables
- Your service's gRPC methods
- Your service's REST endpoints (via adapter)

### DON'T implement:
- Other services' logic (use gRPC calls)
- Other services' databases (use gRPC)
- Authentication (use auth-service)
- Payment logic (use payment-service)
- GPS tracking (use gps-service)

---

## Health Checks

Implement gRPC health checking:

```go
import "google.golang.org/grpc/health"
import healthpb "google.golang.org/grpc/health/grpc_health_v1"

checker := health.NewServer()
checker.SetServingStatus("famgo.service.ServiceAPI", healthpb.HealthCheckResponse_SERVING)

healthpb.RegisterHealthServer(grpcServer, checker)
```

REST health endpoint:

```go
router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "ok",
        "timestamp": time.Now(),
    })
})
```

---

## Makefile

```makefile
.PHONY: dev test build docker-build docker-run clean

dev:
	go run cmd/service/main.go

test:
	go test ./...

build:
	CGO_ENABLED=0 go build -o bin/app cmd/service/main.go

docker-build:
	docker build -t famgo/$(SERVICE_NAME):latest .

docker-run:
	docker run -p 5001:5001 -p 8080:8080 famgo/$(SERVICE_NAME):latest

clean:
	rm -rf bin/
```

---

## Deployment

Use Kubernetes manifests in `infra/kubernetes/`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-service
  template:
    metadata:
      labels:
        app: my-service
    spec:
      containers:
      - name: my-service
        image: famgo/my-service:latest
        ports:
        - containerPort: 5001  # gRPC
        - containerPort: 8080  # REST
        env:
        - name: SERVICE_NAME
          value: my-service
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db.host
        livenessProbe:
          grpc:
            port: 5001
        readinessProbe:
          grpc:
            port: 5001
```

---

## Events This Service Publishes

Define what Kafka topics your service publishes to. Example for ride-service:

```
PUBLISHES:
- ride.created
- ride.started
- ride.completed
- ride.cancelled

SUBSCRIBES:
- payment.completed
- driver.location.updated
```

---

## Next Steps

1. Copy this template
2. Rename files and structs to match your service
3. Define your domain entities
4. Implement your business logic
5. Define your gRPC service (api/proto/service.proto)
6. Implement your handlers
7. Add database migrations
8. Add tests
9. Deploy!

---

**Build enterprise-grade Go services. Follow this template exactly.**
