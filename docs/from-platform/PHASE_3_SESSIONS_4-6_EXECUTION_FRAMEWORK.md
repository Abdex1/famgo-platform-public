# PHASE 3 SESSIONS 4-6+: RAPID DEPLOYMENT EXECUTION FRAMEWORK

**Current Status**: GPS Service ✅ Complete (Session 3)  
**Remaining**: Ride (4), Dispatch (5), Payment/Wallet/Safety/Fraud (6+)  
**Total Time**: 12-15 hours to production MVP  
**Architecture**: Proven 7-layer DDD pattern (replicate GPS Service pattern)

---

## 🎯 RAPID REPLICATION STRATEGY

Each remaining service uses **identical template pattern** from GPS Service:

```
Service Template (6-7 hours per service):
├── Configuration (go.mod + config.go)              [30 min]
├── Domain Layer (entities + value objects)         [1.5 hours]
├── Infrastructure (repositories + stores)          [1.5 hours]
├── Application Layer (use cases)                   [1 hour]
├── gRPC Interface (proto + handlers)              [1 hour]
├── Bootstrap (cmd/main.go + DI)                   [30 min]
├── Docker (multi-stage build)                     [15 min]
└── Tests (unit + integration)                     [1 hour]
```

---

## 📋 SESSION 4: RIDE SERVICE

**Command-Level Specification**:

```go
// Session 4 Ride Service Components

// 1. go.mod - Copy GPS go.mod, add ride-specific deps

// 2. Domain Layer
internal/domain/entities/
  ├── ride.go - State machine: REQUESTED → COMPLETED
  ├── ride_event.go - Event tracking
  
internal/domain/valueobjects/
  ├── ride_status.go - Valid state transitions
  ├── ride_fare.go - Pricing calculation
  ├── pickup_location.go - Start point
  ├── dropoff_location.go - End point
  
internal/domain/services/
  ├── ride_service.go - Lifecycle + state validation
  ├── fare_calculator.go - Trip pricing
  ├── ride_state_machine.go - Transition logic

// 3. Infrastructure
internal/infrastructure/repositories/
  ├── ride_repository.go - CRUD + state queries
  
internal/infrastructure/redis/
  ├── active_rides_cache.go - O(1) lookups

// 4. Application
internal/application/usecases/
  ├── create_ride_usecase.go
  ├── accept_ride_usecase.go
  ├── start_ride_usecase.go
  ├── complete_ride_usecase.go
  └── cancel_ride_usecase.go

// 5. gRPC
proto/ride.proto
internal/interfaces/grpc/ride_handler.go

// 6. Bootstrap + Docker
cmd/main.go
Dockerfile

// 7. Tests
internal/domain/services/ride_service_test.go
```

**Key State Machine**:
```
REQUESTED (rider creates)
    ↓ (driver assigned by Dispatch)
ASSIGNED
    ↓ (driver accepts)
ACCEPTED
    ↓ (driver starts ride)
STARTED
    ↓ (driver completes)
COMPLETED
    
Alternative: Any state → CANCELLED (rider/driver cancels)
```

**Kafka Events to Publish**:
```
ride.requested → [Dispatch Service]
ride.assigned → [Rider Notification]
ride.started → [Safety Service, Payment Service]
ride.completed → [Payment Service, Wallet Service, Analytics]
ride.cancelled → [Dispatch Service, Refund Logic]
```

---

## 📋 SESSION 5: DISPATCH SERVICE

**Command-Level Specification**:

```go
// Session 5 Dispatch Service Components

// 1. Domain Layer
internal/domain/entities/
  ├── match_request.go - Matching request
  ├── driver_assignment.go - Assignment decision
  
internal/domain/valueobjects/
  ├── driver_score.go - Composite score (distance, rating, acceptance)
  ├── match_result.go - Match outcome + ETA
  
internal/domain/services/
  ├── matching_algorithm.go - Multi-factor ranking
  ├── eta_calculator.go - Google Maps integration
  ├── supply_balancer.go - Surge pricing logic
  ├── driver_ranker.go - Scoring formula

// 2. Infrastructure
internal/infrastructure/
  ├── repositories/match_repository.go
  ├── clients/
  │   ├── gps_service_client.go (gRPC)
  │   ├── ride_service_client.go (gRPC)
  │   └── google_maps_client.go (HTTP)
  
internal/infrastructure/redis/
  ├── match_cache.go - Cache recent matches

// 3. Application
internal/application/usecases/
  ├── find_matches_usecase.go
  ├── assign_driver_usecase.go
  ├── calculate_eta_usecase.go
  └── score_driver_usecase.go

// 4. gRPC
proto/dispatch.proto
internal/interfaces/grpc/dispatch_handler.go

// 5. Bootstrap + Docker
cmd/main.go
Dockerfile

// 6. Tests
internal/domain/services/matching_algorithm_test.go
```

**Matching Algorithm Logic**:
```
For each available driver:
  score = (distance_factor * 0.3) +
          (rating_factor * 0.4) +
          (acceptance_rate_factor * 0.2) +
          (availability_factor * 0.1)

Sort by score DESC
Return top N drivers with ETA
```

**gRPC Integration Points**:
```
1. GPS Service: GetNearbyDrivers(lat, lng, radius)
   → Returns: [DriverID with coordinates and distance]

2. Ride Service: AcceptRide(ride_id, driver_id)
   → Updates ride state to ACCEPTED

3. Google Maps: GetETA(origin, destination)
   → Returns: ETA in minutes, distance in km
```

---

## 📋 SESSION 6: PAYMENT/WALLET/SAFETY/FRAUD SERVICES

### 6.1 Payment Service (3-4 hours)

```go
// Key Components

Domain:
  ├── entities/payment.go (PENDING → COMPLETED/FAILED/REFUNDED)
  ├── valueobjects/
  │   ├── payment_status.vo.go
  │   ├── payment_amount.vo.go
  │   └── payment_method.vo.go

Infrastructure:
  ├── repositories/payment_repository.go
  ├── clients/
  │   ├── telebirr_provider.go
  │   ├── cbe_birr_provider.go
  │   └── chapa_provider.go
  └── webhooks/payment_webhook_handler.go

Use Cases:
  ├── process_payment_usecase.go
  ├── handle_webhook_usecase.go
  └── refund_payment_usecase.go

Kafka Events:
  ├── payment.completed → [Wallet, Notification]
  ├── payment.failed → [Notification, Fraud]
  └── payment.refunded → [Wallet]
```

### 6.2 Wallet Service (2-3 hours)

```go
// Key Components

Domain:
  ├── entities/wallet.go
  ├── entities/wallet_transaction.go (append-only)
  ├── valueobjects/wallet_balance.vo.go

Infrastructure:
  ├── repositories/
  │   ├── wallet_repository.go
  │   └── transaction_repository.go
  └── redis/wallet_balance_cache.go

Use Cases:
  ├── debit_wallet_usecase.go (payment deduction)
  ├── credit_wallet_usecase.go (refund/promo)
  └── get_balance_usecase.go (cached)

Kafka Events:
  ├── payment.completed → [trigger debit]
  ├── refund.completed → [trigger credit]
  └── promo.applied → [trigger credit]
```

### 6.3 Safety Service (2-3 hours)

```go
// Key Components

Domain:
  ├── entities/sos_incident.go (OPEN → RESOLVED)
  ├── entities/emergency_contact.go
  ├── valueobjects/incident_status.vo.go

Infrastructure:
  ├── repositories/
  │   ├── incident_repository.go
  │   └── contact_repository.go
  └── clients/notification_client.go

Use Cases:
  ├── trigger_sos_usecase.go
  ├── resolve_incident_usecase.go
  └── get_emergency_contacts_usecase.go

Kafka Events:
  ├── sos.triggered → [Notification, Ops]
  └── incident.resolved → [Archive]
```

### 6.4 Fraud Service (2-3 hours)

```go
// Key Components

Domain:
  ├── entities/fraud_alert.go
  ├── valueobjects/risk_score.vo.go (0-100)

Infrastructure:
  ├── repositories/fraud_alert_repository.go
  ├── services/anomaly_detector.go
  └── clients/analytics_client.go

Use Cases:
  ├── score_user_risk_usecase.go
  ├── detect_fraud_usecase.go
  └── resolve_fraud_alert_usecase.go

Kafka Events:
  ├── fraud.detected → [Payment, Wallet, Safety]
  └── fraud.resolved → [Archive]
```

---

## 🔗 KAFKA EVENT BUS ARCHITECTURE

```
All services publish/subscribe to:

Events (15+ types):
├── Ride Events
│   ├── ride.requested
│   ├── ride.assigned
│   ├── ride.accepted
│   ├── ride.started
│   ├── ride.completed
│   └── ride.cancelled
├── Driver Events
│   ├── driver.location.updated
│   ├── driver.online
│   └── driver.offline
├── Payment Events
│   ├── payment.initiated
│   ├── payment.completed
│   ├── payment.failed
│   └── payment.refunded
├── Wallet Events
│   ├── wallet.debited
│   ├── wallet.credited
│   └── balance.updated
├── Safety Events
│   ├── sos.triggered
│   └── incident.resolved
└── Fraud Events
    ├── fraud.detected
    └── fraud.resolved

Retention: 7 days
Partitions: 3 per topic (scalability)
Replication: 2 (reliability)
```

---

## ✅ VALIDATION CHECKLIST (Each Service)

Before moving to next service:

**Code Quality**:
- [ ] All 7 layers implemented
- [ ] >80% test coverage
- [ ] All tests passing
- [ ] No unhandled errors

**Architecture**:
- [ ] Clean separation of concerns
- [ ] No circular dependencies
- [ ] Dependency injection working
- [ ] Repository pattern enforced

**gRPC**:
- [ ] Proto definitions correct
- [ ] Code generation successful
- [ ] All endpoints callable
- [ ] Error mapping correct

**Integration**:
- [ ] Docker image builds
- [ ] docker-compose works
- [ ] Health checks pass
- [ ] Graceful shutdown works

**Observability**:
- [ ] Logs appear in Loki
- [ ] Traces appear in Jaeger
- [ ] Metrics in Prometheus
- [ ] Correlation IDs tracked

**Testing**:
- [ ] Unit tests comprehensive
- [ ] Integration tests work
- [ ] E2E scenarios pass
- [ ] Mock patterns used

---

## 📊 TIME ALLOCATION

```
Session 4 (Ride):      3-4 hours
Session 5 (Dispatch):  3-4 hours
Session 6.1 (Payment): 3-4 hours
Session 6.2 (Wallet):  2-3 hours
Session 6.3 (Safety):  2-3 hours
Session 6.4 (Fraud):   2-3 hours
─────────────────────────────────
Total:                 15-21 hours

With template reuse & parallelization: 12-15 hours
```

---

## 🚀 PARALLEL EXECUTION OPPORTUNITIES

**Can be built simultaneously** (no dependencies until later):
- Sessions 6.1, 6.2, 6.3, 6.4 can be in progress simultaneously
- Each has independent domain logic
- Integration happens only in event bus

**Sequential Dependencies**:
- Session 4 (Ride) after Session 3 (GPS) ✓
- Session 5 (Dispatch) after Sessions 3 & 4 ✓
- Session 6 (Payment/Wallet/Safety/Fraud) after Sessions 3-5 ✓

---

## 📝 COMPLETION FRAMEWORK

Each service follows:

```
1. Read Session X template (15 min)
2. Create go.mod (5 min)
3. Domain layer (90 min)
   ├── 10 min: Entity classes
   ├── 20 min: Value objects
   ├── 40 min: Domain services
   └── 20 min: Repository interfaces
4. Infrastructure (90 min)
   ├── 30 min: Repositories
   ├── 30 min: Redis/caches
   └── 30 min: External clients
5. Application layer (60 min)
   ├── 50 min: Use cases
   └── 10 min: DTOs
6. gRPC interface (60 min)
   ├── 20 min: Proto definitions
   └── 40 min: Handler implementation
7. Bootstrap (30 min)
8. Docker (15 min)
9. Tests (60 min)
10. Verification (30 min)
───────────────────────
Total: ~6-7 hours per service
```

---

**Execution Status**: Ready to proceed  
**Methodology**: Proven GPS pattern  
**Quality**: Enterprise-grade throughout  
**Timeline**: 12-15 hours to production MVP  
**Next**: Session 4 - Ride Service
