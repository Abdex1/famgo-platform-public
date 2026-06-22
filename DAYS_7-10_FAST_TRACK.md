# 🚀 DAYS 7-10 FAST-TRACK EXECUTION PLAN

**Current Status:** 81% Complete (65 of 80 hours)  
**Remaining:** 15 hours (Days 7-10)  
**Target:** 100% Complete by Day 10  

---

## ⚡ RIDE SERVICE: FAST-TRACK BUILD (Days 7-9, 12 hours)

**Already Created:**
- ✅ Domain entities with state machine
- ✅ Domain errors

**Copy from User Service and Adapt (Use Template Pattern):**

```
services/ride-service/
├── internal/domain/
│   ├── entities.go ✅
│   ├── errors.go ✅
│   └── ride_service.go (2 hrs - copy user_service pattern)
├── internal/application/
│   ├── commands.go (2 hrs - 5 commands: Create, Assign, Start, Complete, Cancel)
│   ├── queries.go (1 hr - copy pattern)
│   └── interfaces.go (0.5 hr - copy pattern)
├── internal/infrastructure/
│   ├── postgres_repo.go (2 hrs - copy user repo pattern)
│   └── redis_cache.go (1 hr - copy cache pattern)
├── internal/transport/
│   └── http_handler.go (2 hrs - copy HTTP pattern)
├── internal/bootstrap/
│   └── container.go (1 hr - copy DI pattern)
├── internal/config/
│   └── config.go (copy from user-service)
├── cmd/
│   └── main.go (copy from user-service, port 5005)
├── db/migrations/
│   ├── 001_create_ride_schema.up.sql (1 hr)
│   └── 001_create_ride_schema.down.sql
├── deployments/
│   ├── deployment.yaml (copy + change port)
│   ├── service.yaml (copy + change port)
│   └── hpa.yaml (copy)
├── Dockerfile (copy + change port 5005)
└── tests/
    └── unit/ride_service_test.go (1 hr)
```

**Key Differences from User Service:**
1. Domain: State machine (Requested → Completed) instead of status
2. Commands: CreateRide, AssignDriver, StartRide, CompleteRide, CancelRide
3. Database: rides, ride_status_history tables
4. HTTP: POST /api/ride/request, GET /api/ride/{rideId}, etc.

**Estimated Time Breakdown:**
- Domain layer: 2.5 hrs
- Application layer: 3.5 hrs
- Infrastructure layer: 3 hrs
- Transport + Bootstrap: 3 hrs

---

## 🔗 DAYS 8-10: WIRING & PRODUCTION (8 hours)

### Event-Driven Workflows (3 hours)
```
ride.requested → dispatch-service → driver.assigned → ride-service → ride.assigned
                                                    → pricing-service → price.calculated
                                                    → wallet-service → payment.authorized
```

**Consumers to Implement:**
- Ride Service: Consumes driver.assigned, price.calculated, payment.authorized
- Dispatch Service: Consumes ride.requested
- Pricing Service: Consumes ride.requested
- Wallet Service: Consumes ride.requested

### gRPC Cross-Service Calls (2 hours)
```go
ride-service calls:
  - gps-service.GetNearbyDrivers(lat, lon, radius)
  - pricing-service.CalculateFare(distance, time)
  - user-service.GetUser(userID)
  - user-service.GetDriver(driverID)
```

### Saga Orchestration (2 hours)
```
CreateRideSaga:
  Step 1: Create ride (ride-service)
  Step 2: Calculate fare (pricing-service)
  Step 3: Search drivers (dispatch-service)
  Step 4: Pre-authorize payment (wallet-service)
  
  Compensation:
    - Cancel ride
    - Release authorization
```

### Production Verification (1 hour)
- Metrics: All services export /metrics
- Traces: All requests traced end-to-end
- Logs: Structured JSON logs aggregated
- Health: All health checks passing
- Security: All requirements met

---

## 📋 CHECKLIST: DAYS 7-10

### Day 7-8: Ride Domain + App (6 hours)
- [ ] Domain layer (entities, service, errors)
- [ ] Application layer (5 commands, 3 queries)
- [ ] Tests for domain logic

### Day 8-9: Ride Infra + Transport (6 hours)
- [ ] PostgreSQL repositories
- [ ] Redis caching
- [ ] HTTP handlers
- [ ] Bootstrap + Config
- [ ] Database migrations
- [ ] Dockerfile + K8s

### Day 9-10: Wiring + Production (8 hours)
- [ ] Event consumers (all services)
- [ ] gRPC client calls
- [ ] Saga orchestration
- [ ] Metrics verification
- [ ] Trace verification
- [ ] Log verification
- [ ] Health check verification
- [ ] Security verification

---

## 🎯 QUICK-START FILES

**Copy These Exactly (No Changes Needed):**
```
user-service/internal/config/ → ride-service/internal/config/
user-service/cmd/main.go → ride-service/cmd/main.go (change port to 5005)
user-service/internal/bootstrap/ → ride-service/internal/bootstrap/
user-service/Dockerfile → ride-service/Dockerfile (change port to 5005)
user-service/deployments/ → ride-service/deployments/ (change port to 5005)
```

**Adapt These (Copy + Change Domain):**
```
user-service/internal/application/ → ride-service/internal/application/
  (Copy command/query patterns, change to: CreateRide, AssignDriver, StartRide, CompleteRide, CancelRide)

user-service/internal/infrastructure/ → ride-service/internal/infrastructure/
  (Copy repo + cache patterns, change to: RideRepository, RideCache)

user-service/internal/transport/ → ride-service/internal/transport/
  (Copy HTTP handlers pattern, change to: /api/ride/request, /api/ride/{rideId}, etc.)
```

**Create New (State Machine Logic):**
```
ride-service/internal/domain/ride_service.go
  (State transition validation, lifecycle management)

ride-service/db/migrations/001_create_ride_schema.up.sql
  (Tables: rides, ride_status_history)
```

---

## ✨ SUCCESS CRITERIA: DAYS 7-10

**Day 9 EOD:**
- ✅ Ride Service 100% production-ready
- ✅ All tests passing
- ✅ All endpoints working
- ✅ Deployment ready

**Day 10 EOD:**
- ✅ Event workflows end-to-end
- ✅ gRPC communication verified
- ✅ Saga orchestration working
- ✅ All metrics exposed
- ✅ All traces propagated
- ✅ All logs aggregated
- ✅ All health checks passing
- ✅ Full production deployment ready

---

## 📊 FINAL PHASE STATUS

```
WEEKS 3-4: 10 DAYS
├─ Days 1-4: Audit (32 hrs) ✅ 100%
├─ Days 5-6: GPS (16 hrs) ✅ 100%
├─ Days 6-7: User (12 hrs) ✅ 100%
├─ Days 7-9: Ride (12 hrs) ⏳ 0% → 100%
└─ Days 8-10: Wiring + Prod (8 hrs) ⏳ 0% → 100%
```

**Target:** 100% complete by end of Day 10 = **80 of 80 hours**

---

**READY TO BUILD RIDE SERVICE (Days 7-9)**  
**THEN WIRE & HARDEN (Days 8-10)**  
**TARGET: 100% COMPLETE BY DAY 10**

