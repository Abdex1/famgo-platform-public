# 📋 WEEKS 3-4: DELIVERY GOVERNANCE SPECIFICATION

**Status:** Repository-First Development Mandate  
**Phase:** Core Services Completion (Weeks 3-4, 80 hours)  
**Primary Objective:** Turn the repository into a complete, coherent, observable, secure, deployable mobility platform  
**Timeline:** 2 weeks / 10 working days  

---

## ⚠️ CRITICAL MANDATE

**DO NOT BUILD NEW SERVICES FIRST.**

The repository already contains the skeleton of an **enterprise mobility platform**:

- ✅ Platform layer (security, gateway, shared, packages)
- ✅ Core services (auth, dispatch, wallet, gps, ride, pricing, payment, etc.)
- ✅ Contract layer (events, schemas, protobufs)
- ✅ Infrastructure layer (docker, kubernetes, terraform)
- ✅ Observability stack (prometheus, grafana, loki, jaeger)

**Your objective is NOT to build more code.**

**Your objective is to COMPLETE, AUDIT, and WIRE existing services.**

**ANY implementation that introduces parallel systems is a regression and will be rejected.**

---

## 📊 REPOSITORY AUDIT: EXISTING ARCHITECTURE

### Layer 1: Shared Contracts (Foundation)

**Location:** `shared/contracts/`

**What Exists:**
```
shared/contracts/
├── events/
│   ├── catalog/              (Event types registry)
│   ├── common/               (Common event fields)
│   ├── dlq/                  (Dead-letter queue policies)
│   ├── driver/               (Driver events)
│   ├── envelopes/            (Event wrappers)
│   ├── idempotency/          (Idempotency keys)
│   ├── payment/              (Payment events)
│   ├── rating/               (Rating events)
│   ├── ride/                 (Ride events)
│   ├── schemas/              (Event schemas)
│   ├── topics/               (Kafka topics)
│   ├── trip/                 (Trip events)
│   ├── retry/                (Retry policies)
│   ├── versions/             (Event versioning)
│   ├── policies/             (Event policies)
├── grpc/                     (gRPC service contracts)
├── rest/                     (REST API contracts)
├── kafka/                    (Kafka topic contracts)
└── websocket/                (WebSocket contracts)
```

**Critical Rule:**
- ✅ All events must originate from `shared/contracts/events`
- ❌ NO service may define local event types
- ❌ NO parallel event contracts
- ❌ NO service-specific topics

---

### Layer 2: Packages (Reusable SDKs)

**Location:** `packages/`

**What Exists:**
```
packages/
├── api-client/              (REST client SDK)
├── auth-client/             (Auth service client)
├── config/                  (Configuration management)
├── event-bus/               (Event bus SDK)
├── feature-flags/           (Feature flag SDK)
├── geo-utils/               (Geospatial utilities)
├── grpc-clients/            (gRPC client generation)
├── i18n/                    (Internationalization)
├── kafka-sdk/               (Kafka wrapper)
├── maps-sdk/                (Maps/GPS utilities)
├── payment-sdk/             (Payment gateway SDK)
├── redis-platform/          (Redis wrapper)
├── telemetry/               (OpenTelemetry SDK)
├── types/                   (Shared types)
├── ui-kit/                  (UI components)
├── ui-theme/                (UI theme)
├── vault-sdk/               (Vault client)
└── websocket-sdk/           (WebSocket client)
```

**Critical Rule:**
- ✅ All services MUST use these SDKs
- ❌ NO service may import `kafka` directly
- ❌ NO service may create duplicate clients
- ❌ NO service may define custom telemetry

---

### Layer 3: Platform (Core Infrastructure)

**Location:** `platform/`

**What Exists:**
```
platform/
├── api-policies/            (API gateway policies)
├── cache/                   (Caching strategies)
├── cdc/                     (Change data capture)
├── config/                  (Configuration system)
├── database/                (Database abstractions)
├── domain/                  (Shared domain models)
├── event-bus/               (Event bus implementation)
├── event-governance/        (Event governance rules)
├── feature-flags/           (Feature flag system)
├── kafka/                   (Kafka platform)
├── orchestration/           (Service orchestration)
├── outbox/                  (Outbox pattern)
├── resilience/              (Circuit breakers, retries)
├── runtime/                 (Runtime configuration)
├── saga/                    (Saga orchestration)
├── security/                (Security platform)
├── service-mesh/            (Service mesh config)
```

**Critical Rule:**
- ✅ ALL services use platform implementations
- ❌ NO custom event bus
- ❌ NO custom telemetry stack
- ❌ NO custom feature flag system

---

### Layer 4: Security Layer

**Location:** `platform/security/` + `security/`

**What Exists:**
- JWT service
- OTP provider
- Vault integration
- RBAC policies
- Certificate management
- Audit logging

**Critical Rule:**
- ✅ Auth service is source of truth
- ❌ NO service may implement its own auth
- ❌ NO password storage outside auth-service

---

### Layer 5: Gateway Layer

**Location:** `gateway/` + `services/api-gateway/`

**What Exists:**
- Kong API gateway
- JWT proxy validation
- Rate limiting policies
- Request routing
- Response transformation

**Critical Rule:**
- ✅ All APIs must pass through gateway
- ❌ NO service exposed publicly directly
- ❌ NO duplicate rate limiting

---

### Layer 6: Core Services (Partially Implemented)

**Services Inventory:**

```
services/
├── auth-service/            ✅ MATURE (ref architecture)
├── user-service/            ⏳ STUB (needs completion)
├── ride-service/            ⏳ STUB (needs completion)
├── dispatch-service/        ⏳ STUB (needs completion)
├── gps-service/             ⏳ STUB (needs completion)
├── pricing-service/         ⏳ STUB (needs completion)
├── payment-service/         ⏳ STUB (needs completion)
├── wallet-service/          ⏳ STUB (needs completion)
├── pooling-service/         ⏳ STUB (needs completion)
├── fraud-service/           ⏳ STUB (needs completion)
├── safety-service/          ⏳ STUB (needs completion)
├── notification-service/    ⏳ STUB (needs completion)
├── analytics-service/       ⏳ STUB (needs completion)
├── subscription-service/    ⏳ STUB (needs completion)
├── driver-service/          ⏳ STUB (needs completion)
├── smart-pickup-service/    ⏳ STUB (needs completion)
├── voice-booking-service/   ⏳ STUB (needs completion)
├── api-gateway/             ⏳ STUB (needs completion)
└── websocket-gateway/       ⏳ STUB (needs completion)
```

---

### Layer 7: Infrastructure (Production-Ready)

**Location:** `infra/`

**What Exists:**
- Docker: Dockerfile templates, multi-stage builds
- Kubernetes: Manifests, helm charts
- Terraform: Cloud infrastructure
- Kong: API gateway
- PostgreSQL: Database
- Redis: Caching
- Redpanda: Kafka
- Grafana: Dashboards
- Tempo: Tracing
- Loki: Logs
- Prometheus: Metrics

---

## 🎯 WEEKS 3-4 EXECUTION PLAN

### PHASE 1: REPOSITORY AUDIT & DOCUMENTATION (Days 1-2, 16 hours)

**Day 1: Contract Governance Audit**

```bash
1. Audit shared/contracts/events/
   - Document all event types
   - Document all topics
   - Create EVENT_CATALOG.md
   
2. Audit shared/contracts/grpc/
   - Document all services
   - Document all methods
   - Create GRPC_CATALOG.md
   
3. Audit shared/contracts/rest/
   - Document all endpoints
   - Document all schemas
   - Create REST_CATALOG.md
   
4. Audit shared/protobufs/
   - Document all proto definitions
   - Document versioning strategy
```

**Deliverable:** `REPOSITORY_CONTRACT_CATALOG.md`

**Day 2: Service Ownership Audit**

```bash
1. Audit each service:
   - Current maturity (stub, partial, mature)
   - Current ownership
   - Current dependencies
   - Current database tables
   - Current events published
   
2. Create SERVICE_OWNERSHIP_MATRIX.md
   
3. Create DEPENDENCY_GRAPH.md
   
4. Create DATA_OWNERSHIP_MATRIX.md
```

**Deliverable:** `SERVICE_OWNERSHIP_MATRIX.md`

---

### PHASE 2: REFERENCE ARCHITECTURE AUDIT (Days 3-4, 16 hours)

**Day 3: Auth Service Deep Audit**

Auth service is the reference architecture. Complete audit:

```bash
1. Audit internal/domain/
   - entities
   - aggregates
   - value objects
   - domain services
   
2. Audit internal/application/
   - use cases
   - commands
   - queries
   
3. Audit internal/infrastructure/
   - postgres implementation
   - redis implementation
   - jwt service
   - otp service
   
4. Audit api/
   - grpc handlers
   - rest handlers
   - websocket handlers
   
5. Document as REFERENCE_ARCHITECTURE.md
```

**Day 4: Platform & Package Audit**

```bash
1. Audit platform/
   - event-bus implementation
   - saga orchestration
   - feature flags
   - caching strategy
   
2. Audit packages/
   - kafka-sdk usage patterns
   - telemetry patterns
   - redis-platform patterns
   
3. Create PLATFORM_USAGE_GUIDE.md
```

---

### PHASE 3: SERVICE COMPLETION (Days 5-7, 24 hours)

**Priority 1: GPS Service (8 hours)**

GPS service is foundation for all location-based features.

**Must Complete:**
- Domain layer (Location, Trip, Geofence aggregates)
- Application layer (tracking, location update, geofence detection)
- Infrastructure (Redis Geo, PostGIS, WebSocket)
- API (gRPC, REST, WebSocket)
- Contracts from `shared/contracts/events/trip/`

**Day 5 Morning:** GPS Service

**Priority 2: User Service (8 hours)**

**Must Complete:**
- Domain layer (User, Profile, Device aggregates)
- Application layer (registration, profile update, device management)
- Infrastructure (PostgreSQL, Redis)
- API (gRPC, REST)
- Contracts from `shared/contracts/events/driver/`

**Day 6 Morning:** User Service

**Priority 3: Ride Service (8 hours)**

**Must Complete:**
- Domain layer (Ride aggregate, State Machine)
- Application layer (create, update, cancel, complete)
- Infrastructure (PostgreSQL, Event Publishing)
- API (gRPC, REST)
- Contracts from `shared/contracts/events/ride/`

**Day 6 Afternoon / Day 7 Morning:** Ride Service

---

### PHASE 4: WIRING SERVICES TOGETHER (Days 8-9, 16 hours)

**Day 8: Event-Driven Architecture**

```bash
1. Wire event publishing through shared/contracts/events
   - ride-service publishes RideRequested
   - dispatch-service consumes RideRequested
   - dispatch-service publishes DriverAssigned
   - ride-service consumes DriverAssigned
   
2. Implement idempotency patterns
3. Implement DLQ handling
4. Implement retry policies
```

**Day 9: Cross-Service Communication**

```bash
1. gRPC service discovery
   - auth-service provides JWT validation
   - gps-service provides location lookups
   - pricing-service provides fare calculations
   
2. Saga orchestration
   - Ride creation saga
   - Payment saga
```

---

### PHASE 5: PRODUCTION READINESS (Day 10, 8 hours)

**Day 10: Complete Production Checklist**

```bash
1. Metrics:
   - Every service exposes Prometheus metrics
   - Request count, latency, errors
   - Business metrics (rides, earnings, etc.)
   
2. Traces:
   - OpenTelemetry integration
   - Cross-service trace propagation
   - Tempo collection
   
3. Logs:
   - Structured logging (JSON)
   - Log levels (DEBUG, INFO, WARN, ERROR)
   - Loki collection
   
4. Health Checks:
   - /health (liveness)
   - /ready (readiness)
   - /startup (startup probe)
   
5. Security:
   - JWT validation on all endpoints
   - RBAC enforcement
   - Audit logging
   
6. Deployment:
   - Kubernetes manifests
   - Helm charts
   - CI/CD pipelines
```

---

## 📋 MANDATORY SERVICE STRUCTURE

Every service MUST contain:

```
service-name/
├── cmd/                     # Entrypoints
│   └── main.go
├── internal/
│   ├── domain/              # Business logic (zero external deps)
│   │   ├── entities.go
│   │   ├── aggregates.go
│   │   ├── value_objects.go
│   │   └── services.go
│   ├── application/         # Use cases
│   │   ├── commands.go
│   │   ├── queries.go
│   │   └── handlers.go
│   ├── infrastructure/      # External integrations
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   ├── kafka.go
│   │   └── grpc_clients.go
│   └── transport/           # API handlers
│       ├── http.go
│       ├── grpc.go
│       └── websocket.go
├── api/                     # API contracts
│   ├── proto/               # gRPC definitions
│   └── openapi.yaml         # REST API spec
├── db/                      # Database
│   ├── migrations/
│   └── schema.sql
├── config/                  # Configuration
│   ├── .env.example
│   ├── .env.local
│   ├── .env.development
│   ├── .env.staging
│   └── .env.production
├── tests/                   # Tests
│   ├── unit/
│   ├── integration/
│   └── contract/
├── deployments/             # Kubernetes
│   ├── deployment.yaml
│   ├── service.yaml
│   └── hpa.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 🏗️ DOMAIN OWNERSHIP RULES

Every domain must have clear ownership. NO CROSS-SERVICE DATABASE WRITES.

### Auth Domain (auth-service owns)

**Tables:**
- users
- roles
- permissions
- user_roles
- sessions
- refresh_tokens
- devices
- otp_requests
- audit_logs

**Events:**
- user.created
- user.verified
- login.success
- login.failed
- otp.sent
- otp.verified
- logout

**No other service may:**
- Read/write users table
- Create JWT tokens
- Manage OTP

---

### GPS Domain (gps-service owns)

**Redis:**
- driver:location:{driver_id}
- trip:locations:{trip_id}
- geofence:{geofence_id}

**PostgreSQL (PostGIS):**
- trip_routes (historical)
- geofences
- heatmaps

**Events:**
- driver.location.updated
- trip.started
- trip.completed
- geofence.entered
- geofence.exited

**No other service may:**
- Store location directly
- Update PostGIS data
- Create trip routes

---

### Ride Domain (ride-service owns)

**Tables:**
- rides
- ride_status_history
- ride_passengers
- ride_events

**State Machine:**
```
Requested → Searching → Assigned → DriverArriving 
→ Started → Completed/Cancelled
```

**Events:**
- ride.requested
- ride.assigned
- ride.started
- ride.completed
- ride.cancelled
- ride.disputed

**No other service may:**
- Create rides
- Update ride status
- Store ride history

---

### Dispatch Domain (dispatch-service owns)

**Tables:**
- driver_assignments
- driver_candidates
- match_history

**Events:**
- driver.searched
- driver.offered
- driver.accepted
- driver.rejected
- driver.timeout

**No other service may:**
- Create assignments
- Update match logic
- Store match history

---

### Pricing Domain (pricing-service owns)

**Tables:**
- price_calculations
- surge_multipliers
- promotions
- subscription_pricing

**Events:**
- price.calculated
- surge.applied
- promotion.applied

**No other service may:**
- Calculate fares
- Apply surge
- Apply promotions

---

### Wallet Domain (wallet-service owns)

**Tables:**
- wallets
- ledger_entries (IMMUTABLE)
- derived_balances (READ-ONLY)

**Critical Rule:** NO direct balance updates. ONLY ledger entries.

**Events:**
- wallet.created
- credit.added
- debit.posted
- balance.updated

---

### Payment Domain (payment-service owns)

**Tables:**
- payment_intents
- payment_transactions
- refunds

**Events:**
- payment.initiated
- payment.succeeded
- payment.failed
- refund.requested
- refund.completed

---

## 📊 EVENT GOVERNANCE RULES

All events must:

1. **Originate from shared/contracts/events only**
2. **Include metadata:**
   ```go
   type Event struct {
       EventID       string    // Unique
       EventType     string    // From catalog
       EventVersion  int       // Versioned
       AggregateID   string    // Domain object ID
       Timestamp     time.Time
       CorrelationID string    // Request trace
       CausationID   string    // Event that triggered this
       Data          []byte    // Payload
   }
   ```

3. **Be idempotent** (same event = same result)
4. **Support replay** (events can be replayed)
5. **Be versioned** (schema evolution)

**Forbidden:**
- ❌ Services publishing directly to Kafka
- ❌ Event types defined in services
- ❌ Events without event_id
- ❌ Non-idempotent event handlers

---

## 🔒 SECURITY RULES

### Authentication

- ✅ All requests go through JWT validation
- ✅ JWT validated at gateway
- ✅ Services trust JWT from gateway
- ❌ Services don't re-validate JWT (trust boundary)

### Authorization

- ✅ RBAC enforced per endpoint
- ✅ Roles defined in auth-service
- ✅ Permissions checked at service level
- ❌ Services don't grant permissions

### Secrets

- ✅ All secrets in Vault
- ✅ Loaded via environment variables
- ✅ Never in code
- ❌ Never in logs
- ❌ Never in config files

---

## 📈 OBSERVABILITY RULES

Every service MUST expose:

### Metrics (Prometheus)

```go
- request_count (counter)
- request_duration_seconds (histogram)
- request_errors_total (counter)
- {service}_{entity}_created_total (counter)
- {service}_{entity}_processing_seconds (histogram)
```

### Traces (Jaeger/Tempo)

```go
- Every request starts trace
- Trace propagated to dependencies
- Traces include service name, operation, duration
- Traces stored in Tempo
```

### Logs (Loki)

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "service": "ride-service",
  "operation": "create_ride",
  "user_id": "user123",
  "trace_id": "abc123",
  "correlation_id": "req456",
  "message": "Ride created",
  "duration_ms": 125
}
```

### Health Checks

```bash
GET /health          # Liveness probe
GET /ready           # Readiness probe
GET /startup         # Startup probe
```

---

## ✅ PRODUCTION READINESS CHECKLIST

A service is NOT ready for production unless it has:

### Functional
- [ ] All domain entities implemented
- [ ] All use cases implemented
- [ ] All aggregates implemented
- [ ] All commands/queries implemented

### Operational
- [ ] Metrics exposed (Prometheus)
- [ ] Logs structured (JSON)
- [ ] Traces propagated (Jaeger)
- [ ] Health checks implemented
- [ ] Alerts configured (Prometheus)

### Security
- [ ] JWT validation
- [ ] RBAC enforcement
- [ ] Audit logging
- [ ] Secrets in Vault
- [ ] Input validation

### Reliability
- [ ] Retries configured
- [ ] Timeouts set
- [ ] Circuit breakers
- [ ] Idempotency keys
- [ ] DLQ handling

### Infrastructure
- [ ] Dockerfile
- [ ] Kubernetes manifests
- [ ] Helm charts
- [ ] Health checks
- [ ] Autoscaling config

### Documentation
- [ ] README with architecture
- [ ] API documentation
- [ ] Database schema
- [ ] Event catalog
- [ ] Runbook

---

## 🚀 WEEKS 3-4 MILESTONES

### Week 3 Milestones

**Day 1-2:** Repository Audit Complete
- [ ] Contract catalog documented
- [ ] Service ownership matrix documented
- [ ] Dependency graph documented

**Day 3-4:** Reference Architecture Documented
- [ ] Auth service documented as reference
- [ ] Platform abstractions documented
- [ ] Package usage patterns documented

**Day 5-7:** Services Completed
- [ ] GPS service complete (24 hours)
- [ ] User service complete (8 hours)
- [ ] Ride service complete (8 hours)

### Week 4 Milestones

**Day 8-9:** Services Wired Together
- [ ] Event-driven architecture working
- [ ] Cross-service communication working
- [ ] Saga orchestration working

**Day 10:** Production Ready
- [ ] All services have metrics
- [ ] All services have traces
- [ ] All services have logs
- [ ] All services have health checks
- [ ] All services deployable

---

## 📊 SUCCESS CRITERIA

### Repository Integrity: 100%

- [x] No parallel auth systems
- [x] No parallel event contracts
- [x] No parallel SDKs
- [x] No parallel infrastructure

### Service Completeness: >90%

- [x] GPS service complete
- [x] User service complete
- [x] Ride service complete
- [x] Dispatch service wired
- [x] Pricing service wired
- [x] Payment service wired

### Production Readiness: 100%

- [x] All services have metrics
- [x] All services have traces
- [x] All services have logs
- [x] All services have health checks
- [x] All services deployable
- [x] All services secure

### Architecture Alignment: 100%

- [x] All services use platform abstractions
- [x] All services use shared contracts
- [x] All services use shared packages
- [x] All services follow domain governance
- [x] No cross-service database writes

---

## 🎯 FINAL OUTCOME

At the end of Weeks 3-4:

**What You Have:**
- ✅ Complete repository audit
- ✅ Contract governance defined
- ✅ Service ownership clear
- ✅ GPS service complete
- ✅ User service complete
- ✅ Ride service complete
- ✅ Services wired together
- ✅ Event-driven architecture working
- ✅ Cross-service communication working
- ✅ All services observable
- ✅ All services secure
- ✅ All services deployable

**What You DON'T Have:**
- ❌ Duplicate implementations
- ❌ Parallel auth systems
- ❌ Parallel event contracts
- ❌ Service boundary violations
- ❌ Hidden dependencies

**Result:** A coherent, enterprise-grade mobility platform ready for integration testing and production deployment.

---

**WEEKS 3-4: REPOSITORY-FIRST DEVELOPMENT MANDATE ESTABLISHED** ✅

