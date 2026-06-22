# 🚀 PHASE 2: CORE BACKEND COHERENCE - IMPLEMENTATION CHECKLIST

## Executive Summary
**Goal**: Close all backend gaps and create unified integration layer  
**Timeline**: 2 weeks (80 hours)  
**Outcome**: Full backend coherence with unified APIs  

---

## WEEK 1: Foundation & Gateway Layer (40 hours)

### Monday: Database Coherence Audit (8 hours)

#### Database Audit Checklist
```sql
-- All 40+ tables must have:
-- ✓ id UUID PRIMARY KEY
-- ✓ created_at TIMESTAMP NOT NULL DEFAULT NOW()
-- ✓ updated_at TIMESTAMP NOT NULL DEFAULT NOW()
-- ✓ deleted_at TIMESTAMP (soft delete support)
-- ✓ created_by UUID NOT NULL
-- ✓ updated_by UUID NOT NULL

-- Critical queries to run:
SELECT tablename FROM pg_tables WHERE schemaname = 'public';
-- EXPECTED: 40+ tables

SELECT constraint_name FROM information_schema.table_constraints 
WHERE table_name = 'rides' AND constraint_type = 'PRIMARY KEY';
-- EXPECTED: rides_pkey on (id)

SELECT indexname FROM pg_indexes WHERE tablename = 'ride_locations';
-- EXPECTED: GiST index on (location) for PostGIS
```

**Files to create**:
- `database/audit_schema.sql` - Audit table structure
- `database/coherence_check.sql` - Validation queries
- `database/migrations/006_audit_trail.sql` - Add audit columns where missing
- `database/migrations/007_add_soft_delete.sql` - Add deleted_at where missing

**Deliverable**: All 40+ tables standardized

#### Audit Trail Implementation
```sql
-- Universal audit function
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  NEW.updated_by = CURRENT_USER_ID(); -- From auth context
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to all tables
CREATE TRIGGER rides_updated_at BEFORE UPDATE ON rides
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  
CREATE TRIGGER payments_updated_at BEFORE UPDATE ON payments
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- ... repeat for all tables
```

---

### Tuesday-Wednesday: API Gateway Setup (16 hours)

#### Kong Configuration
```yaml
# kong/config.yaml - Central API Gateway

_format_version: "2.1"
_transform: true

services:
  # Auth Service
  - name: auth-service
    url: http://auth-service:5000
    plugins:
      - name: jwt
        config:
          key_claim_name: "iss"
          secret_is_base64: false
  
  # GPS Service
  - name: gps-service
    url: http://gps-service:5001
  
  # Ride Service
  - name: ride-service
    url: http://ride-service:5002
  
  # Dispatch Service
  - name: dispatch-service
    url: http://dispatch-service:5003
  
  # Payment Service
  - name: payment-service
    url: http://payment-service:5006
  
  # Wallet Service
  - name: wallet-service
    url: http://wallet-service:5007
  
  # Safety Service
  - name: safety-service
    url: http://safety-service:5008
  
  # Fraud Service
  - name: fraud-service
    url: http://fraud-service:5009

routes:
  # Auth Routes
  - name: auth-login
    service: auth-service
    paths:
      - /api/v1/auth/login
    methods: [POST]
    plugins:
      - name: rate-limiting
        config:
          minute: 5
          policy: local
  
  - name: auth-register
    service: auth-service
    paths:
      - /api/v1/auth/register
    methods: [POST]
    plugins:
      - name: rate-limiting
        config:
          minute: 3
  
  # Ride Routes
  - name: ride-create
    service: ride-service
    paths:
      - /api/v1/rides
    methods: [POST]
    plugins:
      - name: jwt
      - name: rate-limiting
        config:
          minute: 100
  
  - name: ride-get
    service: ride-service
    paths:
      - /api/v1/rides/(.*)
    methods: [GET]
    plugins:
      - name: jwt
  
  # ... (continue for all 36+ endpoints)

plugins:
  # Global plugins
  - name: cors
    config:
      origins:
        - "http://localhost:3000"
        - "http://localhost:8081"
        - "https://app.famgo.et"
      credentials: true
      methods: [GET, POST, PUT, DELETE, PATCH, OPTIONS]
      headers: [Content-Type, Authorization, X-Request-ID]
  
  - name: request-id
    config:
      header_name: X-Request-ID
      generator: uuid
  
  - name: request-size-limiting
    config:
      allowed_payload_size: 128
  
  - name: response-transformer
    config:
      add:
        headers:
          - X-API-Version:1.0
          - X-Service:FamGo
```

**Files to create**:
- `backend/api-gateway/kong/kong.yml`
- `backend/api-gateway/kong/Dockerfile`
- `backend/api-gateway/kong/docker-compose.yml`
- `backend/api-gateway/kong/kong-init.sh`

**Deliverable**: Kong running, all 36+ endpoints routed

---

### Thursday-Friday: Event Schema Registry (16 hours)

#### Kafka Event Schema Setup
```yaml
# backend/kafka/schemas/ride.v1.yaml
name: ride.v1
version: 1
type: object
required:
  - event_id
  - event_type
  - aggregate_id
  - timestamp
  - data

properties:
  event_id:
    type: string
    pattern: '^evt_[0-9a-f]{32}$'
    description: Unique event ID
  
  event_type:
    type: string
    enum:
      - ride.created.v1
      - ride.accepted.v1
      - ride.started.v1
      - ride.completed.v1
      - ride.cancelled.v1
    description: Event type
  
  aggregate_id:
    type: string
    pattern: '^ride_[0-9a-f-]{36}$'
    description: Ride ID
  
  correlation_id:
    type: string
    description: Correlation ID for tracing
  
  timestamp:
    type: string
    format: date-time
    description: RFC3339 UTC timestamp
  
  data:
    type: object
    properties:
      ride_id: { type: string }
      rider_id: { type: string }
      driver_id: { type: string }
      pickup_location: { $ref: '#/definitions/Location' }
      dropoff_location: { $ref: '#/definitions/Location' }
      estimated_fare: { type: number }
      currency: { type: string, enum: [ETB] }
      distance_km: { type: number }
      duration_minutes: { type: integer }

definitions:
  Location:
    type: object
    properties:
      latitude: { type: number }
      longitude: { type: number }
      address: { type: string }
      timezone: { type: string }
```

**Event Schema Files to Create**:
- `backend/kafka/schemas/auth.v1.yaml`
- `backend/kafka/schemas/ride.v1.yaml`
- `backend/kafka/schemas/payment.v1.yaml`
- `backend/kafka/schemas/dispatch.v1.yaml`
- `backend/kafka/schemas/wallet.v1.yaml`
- `backend/kafka/schemas/safety.v1.yaml`
- `backend/kafka/schemas/fraud.v1.yaml`
- `backend/kafka/schemas/gps.v1.yaml`

**Schema Registry Service**:
```go
// backend/services/schema-registry-service/main.go
package main

import (
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"log"
)

func main() {
	// Initialize Confluent Schema Registry
	sr, err := schemaregistry.NewClient(
		schemaregistry.NewConfig().
			SetBaseURL("http://schema-registry:8081"),
	)
	if err != nil {
		log.Fatal(err)
	}
	
	// Register all schemas
	schemas := []string{
		"auth.v1", "ride.v1", "payment.v1", 
		"dispatch.v1", "wallet.v1", "safety.v1",
		"fraud.v1", "gps.v1",
	}
	
	for _, schema := range schemas {
		// Load and register each schema
		registerSchema(sr, schema)
	}
	
	log.Println("All schemas registered successfully")
}
```

**Deliverable**: All 8 event types have versioned schemas

---

## WEEK 2: Integration Layer & Documentation (40 hours)

### Monday-Tuesday: Unified API Client (16 hours)

#### Shared Go Library
```go
// backend/shared/go/client/api_client.go
package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
	
	"github.com/cenkalti/backoff/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
	tracer     trace.Tracer
	logger     Logger
}

type APIResponse struct {
	Success bool                   `json:"success"`
	Data    interface{}            `json:"data,omitempty"`
	Error   *APIError              `json:"error,omitempty"`
	Meta    map[string]interface{} `json:"meta"`
}

type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Retry with exponential backoff
func (c *APIClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	_, span := c.tracer.Start(ctx, "APIClient.Do")
	defer span.End()
	
	var resp *http.Response
	var err error
	
	backoffPolicy := backoff.NewExponentialBackOff()
	backoffPolicy.MaxElapsedTime = 30 * time.Second
	
	operation := func() error {
		resp, err = c.httpClient.Do(req.WithContext(ctx))
		if err != nil {
			c.logger.Error("API call failed", "error", err)
			return err
		}
		
		if resp.StatusCode == http.StatusTooManyRequests {
			c.logger.Warn("Rate limited, will retry")
			return fmt.Errorf("rate limited")
		}
		
		if resp.StatusCode >= 500 {
			c.logger.Warn("Server error, will retry", "status", resp.StatusCode)
			return fmt.Errorf("server error %d", resp.StatusCode)
		}
		
		return nil
	}
	
	err = backoff.Retry(operation, backoffPolicy)
	return resp, err
}

// Get performs a GET request with standard error handling
func (c *APIClient) Get(ctx context.Context, endpoint string, result interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	// Add JWT token from context
	if token := ctx.Value("auth_token"); token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	
	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	
	if !apiResp.Success {
		return fmt.Errorf("API error: %s - %s", apiResp.Error.Code, apiResp.Error.Message)
	}
	
	// Marshal result into struct
	resultBytes, _ := json.Marshal(apiResp.Data)
	json.Unmarshal(resultBytes, result)
	
	return nil
}
```

**Files to create**:
- `backend/shared/go/client/api_client.go`
- `backend/shared/go/client/errors.go`
- `backend/shared/go/client/interceptors.go`
- `backend/shared/go/client/telemetry.go`
- `backend/shared/go/client/README.md`

**Deliverable**: Unified client library ready for all backend services

---

### Wednesday: REST API Wrapper (12 hours)

#### gRPC-to-REST Conversion
```go
// backend/services/api-wrapper/main.go
package main

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	// Create gRPC client for each service
	authConn, _ := grpc.Dial("auth-service:5000", grpc.WithInsecure())
	rideConn, _ := grpc.Dial("ride-service:5002", grpc.WithInsecure())
	paymentConn, _ := grpc.Dial("payment-service:5006", grpc.WithInsecure())
	
	// Create mux for REST handlers
	mux := runtime.NewServeMux()
	
	// Register gRPC handlers
	auth.RegisterAuthServiceHandler(context.Background(), mux, authConn)
	ride.RegisterRideServiceHandler(context.Background(), mux, rideConn)
	payment.RegisterPaymentServiceHandler(context.Background(), mux, paymentConn)
	
	// Start HTTP server on port 8080
	http.ListenAndServe(":8080", mux)
}
```

**OpenAPI Generation**:
```bash
# Generate OpenAPI docs from gRPC
protoc --openapiv2_out=. --openapiv2_opt=output_format=yaml *.proto

# Merge all OpenAPI files
npx openapi-merge \
  --input-dir backend/services/*/openapi.yaml \
  --output backend/shared/openapi-merged.yaml
```

**Files to create**:
- `backend/services/api-wrapper/main.go`
- `backend/services/api-wrapper/Dockerfile`
- `backend/shared/openapi/openapi-merged.yaml`
- `backend/shared/postman/FamGo-API.postman_collection.json`

**Deliverable**: REST API with OpenAPI docs + Postman collection

---

### Thursday-Friday: Documentation & Testing (12 hours)

#### API Documentation
```markdown
# FamGo API Documentation

## Authentication
All requests must include JWT token:
\`\`\`
Authorization: Bearer <jwt_token>
\`\`\`

## Standard Response Format
\`\`\`json
{
  "success": true,
  "data": {},
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "request_id": "req_123456",
    "version": "1.0"
  }
}
\`\`\`

## Error Codes
### Auth Errors
- INVALID_CREDENTIALS (401)
- TOKEN_EXPIRED (401)
- UNAUTHORIZED (403)

### Ride Errors
- RIDE_NOT_FOUND (404)
- INVALID_LOCATION (400)
- LOCATION_OUTSIDE_SERVICE_AREA (400)

### Payment Errors
- PAYMENT_FAILED (402)
- INSUFFICIENT_BALANCE (402)
- PROVIDER_ERROR (500)

## Endpoints

### Authentication
- POST /api/v1/auth/login
- POST /api/v1/auth/register
- POST /api/v1/auth/refresh
- POST /api/v1/auth/logout

### Rides
- POST /api/v1/rides (create ride)
- GET /api/v1/rides/:id (get ride details)
- PUT /api/v1/rides/:id/cancel (cancel ride)
- GET /api/v1/rides/:id/tracking (WebSocket, realtime)
- GET /api/v1/rides (list rides with pagination)

### Payments
- POST /api/v1/payments (initiate payment)
- GET /api/v1/payments/:id (get payment status)
- POST /api/v1/payments/:id/refund (refund payment)

### Wallet
- GET /api/v1/wallet/balance (get balance)
- POST /api/v1/wallet/deposit (add funds)
- POST /api/v1/wallet/withdraw (remove funds)
- GET /api/v1/wallet/transactions (transaction history)

### Driver (GPS)
- POST /api/v1/drivers/location (update location)
- GET /api/v1/drivers/nearby (find nearby drivers)
- POST /api/v1/drivers/status (online/offline)

### Dispatch
- POST /api/v1/dispatch/match (find drivers for ride)
- GET /api/v1/dispatch/:id (get dispatch details)

### Safety
- POST /api/v1/safety/sos (report emergency)
- GET /api/v1/safety/:id (get incident details)
- PUT /api/v1/safety/:id/escalate (escalate incident)

### Fraud
- POST /api/v1/fraud/check (check fraud score for ride)
- GET /api/v1/fraud/:id (get fraud check details)

```

**Files to create**:
- `backend/shared/docs/API_GUIDE.md`
- `backend/shared/docs/ERROR_CODES.md`
- `backend/shared/docs/AUTHENTICATION.md`
- `backend/shared/docs/WEBHOOKS.md`
- `backend/shared/docs/WEBSOCKET.md`

#### Contract Tests
```go
// backend/test/integration/contract_test.go
package integration

import (
	"testing"
	"github.com/pact-foundation/pact-go/v2/consumer"
)

func TestAuthServiceContract(t *testing.T) {
	mockProvider := consumer.NewMockProvider("AuthService", "ConsumerApp")
	defer mockProvider.Close()
	
	// Define expected interaction
	mockProvider.
		AddInteraction().
		GivenAState("user exists").
		UponReceiving("a request to login").
		WithRequest("POST", "/login", map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}).
		WillRespondWith(200, map[string]interface{}{
			"token": "jwt_token_here",
		})
	
	// Verify contract
	if err := mockProvider.Verify(t); err != nil {
		t.Fatal(err)
	}
}
```

**Deliverable**: Complete API documentation, Postman collection, contract tests

---

## INTEGRATION POINTS VALIDATION

### Checklist
```
WEEK 1 DELIVERABLES:
✅ Database coherence (all 40+ tables standardized)
✅ API Gateway (Kong) routing all 36+ endpoints
✅ Kafka event schema registry (8 event types)
✅ Error code standardization (across all services)
✅ Timestamp standardization (RFC3339 UTC)

WEEK 2 DELIVERABLES:
✅ Unified API client (Go library)
✅ REST API wrapper (gRPC-to-REST)
✅ OpenAPI documentation (100% coverage)
✅ Postman collection (all endpoints)
✅ Contract tests (inter-service)
✅ Integration tests (end-to-end)

READY FOR PHASE 3:
✅ Mobile apps can connect
✅ All APIs use same format
✅ All errors consistent
✅ All events validated
✅ All documentation complete
```

---

## Commands to Execute (Phase 2)

```bash
# Week 1: Database & Gateway
docker-compose -f infra/docker/docker-compose.yml up postgres redis kafka kong

# Run database migrations
make migrate-db

# Apply Kong configuration
kubectl apply -f backend/api-gateway/kong/kong.yml

# Week 2: Integration Layer
make build-api-client
make generate-openapi
make run-api-wrapper

# Testing
make test-contracts
make test-integration
make run-postman-collection
```

---

**Phase 2 Completion**: All backend services coherent and integrated  
**Estimated Hours**: 80 hours (2 weeks)  
**Next Phase**: Mobile apps and frontend dashboards
