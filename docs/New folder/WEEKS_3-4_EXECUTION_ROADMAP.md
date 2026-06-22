# 🗺️ WEEKS 3-4 EXECUTION ROADMAP: REPOSITORY-FIRST DEVELOPMENT

**Status:** Comprehensive execution plan  
**Timeline:** 10 working days (80 hours)  
**Mandate:** Complete existing services, not build new ones  
**Outcome:** Coherent, production-ready mobility platform

---

## 📊 EXECUTION OVERVIEW

### Phase 1: Repository Audit & Documentation (Days 1-4, 32 hours)

**Primary Objective:** Understand existing architecture before building

**Days 1-2: Contract & Package Governance (16 hours)**

```
Day 1 Morning (4 hours):
  ✅ Task 1.1: Contract Audit
     - Read shared/contracts/events/catalog/
     - Document all event types
     - Create EVENT_CATALOG.md
     - Document: producers, consumers, versions
     
  ✅ Task 1.2: Topic Audit
     - Read shared/contracts/events/topics/
     - List all Kafka topics
     - Document: retention, replication, policies
     - Create TOPIC_REGISTRY.md

Day 1 Afternoon (4 hours):
  ✅ Task 1.3: Envelope & Versioning Audit
     - Read shared/contracts/events/envelopes/
     - Read shared/contracts/events/versions/
     - Document: event structure
     - Document: versioning strategy
     - Document: schema migration
     - Create EVENT_STRUCTURE.md

Day 2 Morning (4 hours):
  ✅ Task 2.1: Package Audit
     - Read packages/kafka-sdk/
     - Read packages/event-bus/
     - Read packages/telemetry/
     - Read packages/redis-platform/
     - Verify: Usage patterns in services
     - Create PACKAGE_USAGE_GUIDE.md

Day 2 Afternoon (4 hours):
  ✅ Task 2.2: Verify No Duplication
     - Search: for "kafka.Dial" in services/
     - Search: for custom telemetry in services/
     - Search: for custom event-bus in services/
     - Verify: ALL use platform packages
     - Document: Any violations found
```

**Deliverables:**
- `EVENT_CATALOG.md` (list of all events)
- `TOPIC_REGISTRY.md` (list of all topics)
- `EVENT_STRUCTURE.md` (event format, versioning)
- `PACKAGE_USAGE_GUIDE.md` (how to use each SDK)
- `DUPLICATION_SCAN_REPORT.md` (violations found)

---

**Days 3-4: Service Audit & Reference Architecture (16 hours)**

```
Day 3 Morning (4 hours):
  ✅ Task 3.1: Auth Service Deep Audit
     - Read: services/auth-service/internal/domain/
     - Document: User entity, Role aggregate, Permissions
     - Verify: ZERO external dependencies
     - Document: Domain services and value objects
     
  ✅ Task 3.2: Auth Service - Application Layer
     - Read: services/auth-service/internal/application/
     - Document: Commands (Register, Login, Refresh)
     - Document: Queries (GetUser, VerifyToken)
     - Document: Command/Query handlers
     - Verify: Only uses domain and infrastructure interfaces

Day 3 Afternoon (4 hours):
  ✅ Task 3.3: Auth Service - Infrastructure Layer
     - Read: services/auth-service/internal/infrastructure/
     - Document: PostgreSQL repos
     - Document: Redis implementations
     - Document: JWT service
     - Document: OTP service
     - Verify: Uses platform abstractions

  ✅ Task 3.4: Auth Service - Transport Layer
     - Read: services/auth-service/internal/transport/
     - Document: HTTP handlers
     - Document: gRPC handlers
     - Document: WebSocket handlers
     - Verify: No business logic leaks

Day 4 Morning (4 hours):
  ✅ Task 3.5: Auth Service - Complete Assessment
     - Document: Complete architecture as REFERENCE_ARCHITECTURE.md
     - Document: Patterns to replicate in other services
     - Create: auth-service/ARCHITECTURE.md
     - Highlight: Clean separation of concerns

Day 4 Afternoon (4 hours):
  ✅ Task 3.6: Platform & Gateway Audit
     - Read: platform/event-bus/
     - Read: platform/saga/
     - Read: platform/feature-flags/
     - Read: services/api-gateway/
     - Document: Usage patterns
     - Create: PLATFORM_ABSTRACTIONS.md
```

**Deliverables:**
- `REFERENCE_ARCHITECTURE.md` (Auth service as reference)
- `ARCHITECTURE_PATTERNS.md` (patterns to replicate)
- `PLATFORM_ABSTRACTIONS.md` (platform layer overview)
- `API_GATEWAY_CONFIGURATION.md` (gateway setup)
- `REPOSITORY_AUDIT_COMPLETE.md` (summary)

---

### Phase 2: Service Completion (Days 5-9, 40 hours)

**Priority 1: GPS Service (Days 5-6, 16 hours)**

GPS service is foundation for all location-based features.

```
Day 5 Morning (4 hours):
  ✅ Task 5.1: GPS Domain Layer
     - Create: services/gps-service/internal/domain/
     - Define: DriverLocation entity
     - Define: Trip aggregate
     - Define: Geofence entity
     - Define: RoutePoint value object
     - Implement: LocationService (pure domain logic)
     - Implement: Methods for geofence detection
     - Reference: auth-service domain pattern
     - Verify: ZERO external dependencies

  ✅ Task 5.2: GPS Application Layer
     - Create: services/gps-service/internal/application/
     - Commands: UpdateDriverLocationCommand
     - Handlers: UpdateDriverLocationHandler
     - Queries: GetDriverLocationQuery
     - Reference: Auth service application pattern

Day 5 Afternoon (4 hours):
  ✅ Task 5.3: GPS Infrastructure Layer
     - Create: services/gps-service/internal/infrastructure/
     - PostgreSQL repos (using platform/database/)
     - Redis repos (using packages/redis-platform/)
     - Event publishing (using packages/event-bus/)
     - Reference: Auth service infrastructure pattern

  ✅ Task 5.4: GPS Transport Layer
     - Create: services/gps-service/internal/transport/
     - HTTP handlers
     - gRPC handlers (from api/proto/gps.proto)
     - WebSocket handlers
     - Health checks

Day 6 Morning (4 hours):
  ✅ Task 5.5: GPS Events & Tests
     - Events: MUST use shared/contracts/events/
     - DriverLocationUpdatedEvent
     - GeofenceEnteredEvent
     - GeofenceExitedEvent
     - Tests: Unit tests for LocationService
     - Tests: Integration tests for handlers
     - Coverage: >80%

Day 6 Afternoon (4 hours):
  ✅ Task 5.6: GPS Database & Deployment
     - Database migrations for gps-service
     - Dockerfile (multi-stage build)
     - Kubernetes manifests (Deployment, Service)
     - Health checks in manifests
     - Test local build: docker build services/gps-service
```

**Deliverables:**
- `services/gps-service/internal/domain/entities.go`
- `services/gps-service/internal/application/commands.go`
- `services/gps-service/internal/infrastructure/repos.go`
- `services/gps-service/internal/transport/handlers.go`
- `services/gps-service/db/migrations/001_create_gps_schema.sql`
- `services/gps-service/Dockerfile`
- `services/gps-service/deployments/kubernetes.yaml`
- `services/gps-service/tests/unit/location_service_test.go`
- `services/gps-service/README.md` with architecture

---

**Priority 2: User Service (Day 6-7, 12 hours)**

User service owns driver and passenger profiles.

```
Day 6 Afternoon / Day 7 Morning (4-6 hours):
  ✅ Task 6.1: User Domain Layer
     - Entities: User, DriverProfile, PassengerProfile, Device
     - Aggregates: User aggregate with profiles
     - Domain services: ProfileService
     - Value objects: PhoneNumber, Email, Rating
     - Reference: Auth service domain pattern

  ✅ Task 6.2: User Application Layer
     - Commands: CreateUserCommand, UpdateProfileCommand
     - Queries: GetUserQuery, GetProfileQuery
     - Handlers for all commands/queries
     - Reference: Auth service pattern

  ✅ Task 6.3: User Infrastructure & Transport
     - Repos: PostgreSQL for users/profiles
     - Redis for session/profile cache
     - HTTP & gRPC handlers
     - Health checks

Day 7 Afternoon (6 hours):
  ✅ Task 6.4: User Database & Tests
     - Schema migrations
     - Tests (>80% coverage)
     - Dockerfile
     - Kubernetes manifests
```

**Deliverables:**
- `services/user-service/` (complete service)

---

**Priority 3: Ride Service (Day 7-9, 12 hours)**

Ride service owns ride lifecycle and state management.

```
Day 7 Afternoon / Day 8 Morning (6-8 hours):
  ✅ Task 7.1: Ride Domain Layer
     - Aggregate: Ride (with state machine)
     - States: Requested, Searching, Assigned, DriverArriving, Started, Completed, Cancelled
     - Entities: Passenger, Pickup, Dropoff
     - Domain services: RideStateService (state transitions)
     - Events: RideRequested, RideAssigned, RideStarted, RideCompleted, RideCancelled
     - Value objects: Location, Money, Duration

  ✅ Task 7.2: Ride Application Layer
     - Commands: CreateRideCommand, AcceptRideCommand, StartRideCommand, CompleteRideCommand
     - Queries: GetRideQuery, GetRideHistoryQuery
     - Handlers with state machine transitions
     - Saga: RideCreationSaga (through platform/saga/)

  ✅ Task 7.3: Ride Infrastructure & Transport
     - PostgreSQL repos with state history
     - Event publishing
     - gRPC & HTTP handlers
     - Health checks

Day 8 Afternoon / Day 9 Morning (4-6 hours):
  ✅ Task 7.4: Ride Database & Tests
     - Schema migrations (with ride_status_history)
     - Tests (>80% coverage)
     - Dockerfile
     - Kubernetes manifests
```

**Deliverables:**
- `services/ride-service/` (complete service)

---

### Phase 3: Wiring Services Together (Days 8-9, 16 hours)

**Objective:** Services communicate through events and gRPC only

```
Day 8 Morning (4 hours):
  ✅ Task 8.1: Event-Driven Workflows
     - Workflow: Ride creation
       1. User calls ride-service/CreateRide
       2. ride-service publishes RideRequested event (to shared/contracts/events)
       3. dispatch-service subscribes to RideRequested
       4. dispatch-service publishes DriverAssigned event
       5. ride-service subscribes to DriverAssigned
       6. ride-service updates state
     
     - Verify: ALL events through shared/contracts/events
     - Verify: idempotency (same event = same result)
     - Verify: DLQ handling on failures

  ✅ Task 8.2: Cross-Service gRPC
     - GPS service provides: GetDriverLocationRPC
     - Ride service calls: gps-service.GetLocation(driverID)
     - Pricing service calls: pricing-service.CalculateFare(pickupLocation, dropoffLocation)
     - Verify: Service discovery working
     - Verify: Timeouts set
     - Verify: Circuit breakers active

Day 8 Afternoon (4 hours):
  ✅ Task 8.3: Saga Orchestration
     - Using: platform/saga/
     - Saga: RideCreationSaga
       Step 1: Create ride
       Step 2: Calculate price
       Step 3: Request driver
       Step 4: Update ride with driver
       Compensate: Cancel if any step fails
     
     - Test: Happy path
     - Test: Failure scenarios

  ✅ Task 8.4: Event Replay & Idempotency
     - Event: same event published twice
     - Verify: same result
     - Test: DLQ recovery
     - Test: Retry policies
```

---

### Phase 4: Production Readiness (Days 9-10, 24 hours)

```
Day 9 Morning (4 hours):
  ✅ Task 9.1: Metrics & Observability
     - Every service exposes Prometheus metrics:
       - request_count (counter)
       - request_duration_seconds (histogram)
       - request_errors_total (counter)
       - {service}_{entity}_created_total
     - Test: curl http://localhost:8080/metrics
     - Verify: Prometheus scraping

Day 9 Afternoon (4 hours):
  ✅ Task 9.2: Traces & Logs
     - OpenTelemetry initialization in every service
     - Trace propagation across services
     - Structured logs (JSON format)
     - Test: Jaeger UI shows cross-service traces
     - Test: Loki shows all logs

  ✅ Task 9.3: Health Checks
     - Every service: GET /health (liveness)
     - Every service: GET /ready (readiness)
     - Every service: GET /startup (startup probe)
     - Test: kubectl describe pod {service-pod}

Day 10 Morning (8 hours):
  ✅ Task 10.1: Deployment Validation
     - Verify: All services have Dockerfile
     - Verify: All services have Kubernetes manifests
     - Verify: All services have Helm charts
     - Test: kubectl apply -f services/*/deployments/
     - Test: All pods running
     - Test: Health checks passing

Day 10 Afternoon (8 hours):
  ✅ Task 10.2: Integration Testing
     - Test: Full ride workflow
       1. Create user (user-service)
       2. Request ride (ride-service)
       3. Get nearby drivers (gps-service)
       4. Calculate price (pricing-service)
       5. Assign driver (dispatch-service)
       6. Start trip (ride-service)
       7. Complete trip (ride-service)
     - Test: Payment flows through payment-service
     - Test: All events published correctly
     - Verify: No orphaned requests
```

---

## 📋 DELIVERABLES BY DAY

| Day | Phase | Hours | Deliverable | Status |
|-----|-------|-------|-------------|--------|
| 1-2 | Audit | 16 | EVENT_CATALOG, TOPIC_REGISTRY, PACKAGE_GUIDE | ⏳ |
| 3-4 | Audit | 16 | REFERENCE_ARCHITECTURE, PLATFORM_ABSTRACTIONS | ⏳ |
| 5-6 | GPS | 16 | GPS service complete, tests, deployment | ⏳ |
| 6-7 | User | 12 | User service complete, tests, deployment | ⏳ |
| 7-9 | Ride | 12 | Ride service complete, tests, deployment | ⏳ |
| 8-9 | Wiring | 16 | Event workflows, gRPC, sagas working | ⏳ |
| 9-10 | Production | 24 | Metrics, traces, logs, health, deployment | ⏳ |
| **Total** | **All** | **80** | **Production-ready platform** | **⏳** |

---

## 🎯 CRITICAL SUCCESS FACTORS

### Repository Integrity: Must Maintain 100%

✅ **DO:**
- Use shared/contracts/events for all events
- Use packages/kafka-sdk for all Kafka operations
- Use packages/event-bus for all publishing
- Use packages/telemetry for all observability
- Use platform/ abstractions for saga, feature flags, caching
- Use auth-service patterns in all new services
- Document all domain ownership

❌ **DO NOT:**
- Create parallel event contracts
- Implement custom Kafka clients
- Implement custom telemetry
- Implement custom event-bus
- Create service-local events
- Implement custom saga logic
- Violate domain boundaries
- Create duplicate implementations

### Service Architecture: All Identical

✅ Every service must have:
- internal/domain/ (pure business logic, zero external deps)
- internal/application/ (commands, queries, handlers)
- internal/infrastructure/ (repos, external clients)
- internal/transport/ (HTTP, gRPC, WebSocket handlers)
- db/migrations/ (schema)
- api/proto/ (gRPC contracts)
- tests/ (unit, integration, contract)
- Dockerfile
- kubernetes/
- health checks

### Production Requirements: 100% Met

✅ Every service must:
- Export Prometheus metrics
- Propagate traces (OpenTelemetry)
- Output structured logs (JSON)
- Have health checks (/health, /ready, /startup)
- Have JWT validation
- Have RBAC authorization
- Have input validation
- Have error handling
- Have retry policies
- Have timeout configuration
- Be deployable via kubectl

---

## 🚀 FINAL CHECKLIST: WEEKS 3-4 COMPLETE

### Repository Audit: 100%
- [x] All contracts documented
- [x] All packages audited
- [x] All platform abstractions understood
- [x] Auth service documented as reference
- [x] Service ownership clear
- [x] Database boundaries defined

### Service Completion: 100%
- [x] GPS service complete (production-ready)
- [x] User service complete (production-ready)
- [x] Ride service complete (production-ready)
- [x] Dispatch service wired
- [x] Pricing service wired
- [x] Payment service wired

### Architecture Alignment: 100%
- [x] All services use shared contracts
- [x] All services use platform abstractions
- [x] All services follow reference pattern
- [x] No duplicate implementations
- [x] No service boundary violations
- [x] Event-driven architecture working

### Production Readiness: 100%
- [x] All services observable (metrics, traces, logs)
- [x] All services secure (JWT, RBAC, audit)
- [x] All services reliable (retries, timeouts, circuit breakers)
- [x] All services deployable (docker, k8s, helm)
- [x] All services documented (README, API, architecture)

### Integration: 100%
- [x] Full ride workflow end-to-end
- [x] Event-driven communication working
- [x] gRPC cross-service communication working
- [x] Saga orchestration working
- [x] All tests passing
- [x] All services healthy

---

## 🏁 WEEKS 3-4 COMPLETE: NEXT PHASE

**After Weeks 3-4:**

- ✅ Repository fully audited
- ✅ Core services (GPS, User, Ride) complete
- ✅ Services properly wired
- ✅ Event-driven architecture working
- ✅ Production-ready infrastructure
- ✅ Observability fully integrated
- ✅ Security hardened
- ✅ Deployment automated

**Ready for:**
- Week 5: Dispatch, Pricing, Payment (20 hours each)
- Week 6-7: Remaining services (Wallet, Pooling, Fraud, Safety, etc.)
- Week 8-16: Advanced features, mobile apps, ML services

---

**WEEKS 3-4 EXECUTION ROADMAP ESTABLISHED** ✅

Comprehensive plan with specific tasks, deliverables, and success criteria.

Ready to implement repository-first development.

