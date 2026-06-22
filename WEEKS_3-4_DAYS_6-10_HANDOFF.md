# 📋 WEEKS 3-4: DAYS 6-10 EXECUTION HANDOFF

**Status:** Days 1-6 Complete (60% of phase, 48 of 80 hours)  
**Next:** Days 6-10 Execution (40% remaining, 32 hours)  
**Repository:** github.com/Abdex1/FamGo-platform  

---

## ✅ HANDOFF: WHAT'S READY

### Completed (Days 1-6)
- ✅ **Audit Phase:** Complete (0 violations)
- ✅ **GPS Service:** 100% production-ready
- ✅ **Reference Patterns:** Established
- ✅ **Documentation:** Comprehensive
- ✅ **Repository Governance:** Defined

### For Days 6-10 (Copy GPS Pattern)
- Use `/services/gps-service/` as exact template
- All 4 layers already implemented
- All patterns verified and working
- All tests passing (>80% coverage)

---

## 🎯 DAYS 6-10 EXECUTION PLAN

### Days 6-7: USER SERVICE (12 hours)

**Follow GPS Pattern Exactly:**

1. **Copy Directory Structure**
   ```bash
   cp -r services/gps-service services/user-service
   ```

2. **Update Domain Layer** (`internal/domain/`)
   - Replace entities: User, DriverProfile, PassengerProfile (from USER_SERVICE_TEMPLATES.md)
   - Replace LocationService with UserService (business logic for user management)
   - Keep repository interfaces pattern

3. **Update Application Layer** (`internal/application/`)
   - Replace commands: CreateUserCommand, UpdateProfileCommand
   - Replace queries: GetUserQuery, GetProfileQuery
   - Follow exact handler pattern from GPS

4. **Update Infrastructure Layer** (`internal/infrastructure/`)
   - Replace repositories: UserRepository (users table)
   - Replace caches: UserCache (user profile caching)
   - Keep event publishing through platform/event-bus

5. **Update Transport Layer** (`internal/transport/`)
   - Replace endpoints: POST /api/user/register, GET /api/user/profile
   - Keep health checks pattern (live, ready, startup)
   - Keep metrics pattern

6. **Database & Deployment**
   - Create migrations (user tables, indexes)
   - Create unit tests (>80% coverage)
   - Create Dockerfile (copy GPS pattern)
   - Create Kubernetes manifests (copy GPS pattern)

**Database Tables:**
```sql
CREATE TABLE users (...)
CREATE TABLE driver_profiles (...)
CREATE TABLE passenger_profiles (...)
```

**Events:** Use from `shared/contracts/events/driver/` and `shared/contracts/events/rider/`

---

### Days 7-9: RIDE SERVICE (12 hours)

**Follow GPS Pattern + Add State Machine:**

1. **Copy GPS Directory** → `services/ride-service`

2. **Domain Layer: State Machine**
   - Ride aggregate with state machine
   - States: Requested → Searching → Assigned → DriverArriving → Started → Completed/Cancelled
   - Implement state transitions with validation

   ```go
   type RideStatus string
   const (
       RideStatusRequested RideStatus = "REQUESTED"
       RideStatusSearching RideStatus = "SEARCHING"
       RideStatusAssigned RideStatus = "ASSIGNED"
       RideStatusDriverArriving RideStatus = "DRIVER_ARRIVING"
       RideStatusStarted RideStatus = "STARTED"
       RideStatusCompleted RideStatus = "COMPLETED"
       RideStatusCancelled RideStatus = "CANCELLED"
   )

   // Allowed transitions
   func (r *Ride) CanTransitionTo(newStatus RideStatus) bool {
       // Validate state transitions
   }
   ```

3. **Application Layer: Ride Lifecycle**
   - CreateRideCommand → publishes RideRequested
   - StartRideCommand → publishes RideStarted
   - CompleteRideCommand → publishes RideCompleted
   - CancelRideCommand → publishes RideCancelled

4. **Infrastructure Layer: Ride Persistence**
   - RideRepository: Create, Get, Update, GetByDriver, GetByPassenger
   - RideCache: Cache active rides
   - Event publishing through platform/event-bus

5. **Transport Layer: Ride Endpoints**
   - POST /api/ride/request
   - GET /api/ride/{ride_id}
   - POST /api/ride/{ride_id}/start
   - POST /api/ride/{ride_id}/complete
   - POST /api/ride/{ride_id}/cancel

6. **Events:**
   - ride.requested
   - ride.assigned (consumed from dispatch)
   - ride.started
   - ride.completed
   - ride.cancelled

**Database Tables:**
```sql
CREATE TABLE rides (...)
CREATE TABLE ride_status_history (...)
CREATE TABLE ride_passengers (...)
```

---

### Days 8-9: WIRING SERVICES (16 hours)

**Event-Driven Workflows:**

1. **Ride Request → Dispatch Flow**
   ```
   User POST /api/ride/request
   → ride-service publishes ride.requested
   → dispatch-service consumes, searches drivers
   → dispatch-service publishes driver.assigned
   → ride-service consumes, updates state
   → User WebSocket receives assignment
   ```

2. **Location Updates → Ride Updates**
   ```
   GPS service updates location (gps-service publishes driver.location.updated)
   → ride-service consumes
   → ride-service calculates ETA
   → ride-service publishes ride.eta.updated
   → User receives update via WebSocket
   ```

3. **gRPC Cross-Service Calls**
   ```go
   ride-service calls:
   - gps-service.GetNearbyDrivers()
   - pricing-service.CalculateFare()
   - user-service.GetUser()
   
   gps-service calls:
   - user-service.GetDriver()
   
   dispatch-service calls:
   - gps-service.GetNearbyDrivers()
   - rating-service.GetDriverRating()
   ```

4. **Saga Orchestration: CreateRideSaga**
   ```
   Step 1: Create ride (ride-service)
   Step 2: Calculate fare (pricing-service)
   Step 3: Search drivers (dispatch-service)
   Step 4: Pre-authorize payment (wallet-service)
   
   Compensate on failure:
   - Cancel ride
   - Release authorization
   ```

5. **Event Idempotency Testing**
   - Publish same event twice
   - Verify same result (handlers are idempotent)
   - Test DLQ handling

---

### Days 9-10: PRODUCTION READINESS (16 hours)

**Complete Observability & Security:**

1. **Metrics (Prometheus)**
   - Every service exports /metrics
   - Record request count, latency, errors
   - Record business metrics (rides created, fares calculated)
   - Verify metrics flow to Prometheus

2. **Traces (Jaeger/Tempo)**
   - Every request starts trace
   - Trace propagated across services
   - Verify cross-service traces in Jaeger UI

3. **Logs (Loki)**
   - All logs structured JSON
   - Log levels: DEBUG, INFO, WARN, ERROR
   - Logs flow to Loki
   - Searchable by trace_id, correlation_id

4. **Health Checks**
   - Every service: /health, /ready, /startup
   - Kubernetes probes configured
   - Database connectivity checks
   - Redis connectivity checks

5. **Security Hardening**
   - JWT validation on all endpoints
   - RBAC checks
   - Audit logging for all mutations
   - Secrets in Vault

6. **Reliability**
   - Retries configured (exponential backoff)
   - Timeouts on all external calls
   - Circuit breakers on services
   - Idempotency guaranteed

7. **Deployment Verification**
   - All services build without errors
   - All services deploy to Kubernetes
   - All health checks passing
   - Integration test suite passing
   - Load test baseline captured

---

## 📋 QUICK CHECKLIST: Days 6-10

### Day 6-7: User Service
- [ ] Copy GPS directory → user-service
- [ ] Update domain layer (User, DriverProfile, PassengerProfile)
- [ ] Update application layer (Create, Update, Get)
- [ ] Update infrastructure (repos, caches)
- [ ] Update transport (HTTP endpoints)
- [ ] Database migrations
- [ ] Tests >80% coverage
- [ ] Dockerfile + K8s manifests

### Day 7-9: Ride Service
- [ ] Copy GPS directory → ride-service
- [ ] Implement state machine (Requested→...→Completed)
- [ ] Application layer with lifecycle
- [ ] Infrastructure (repos, caches)
- [ ] Transport (endpoints)
- [ ] Database migrations
- [ ] Tests >80% coverage
- [ ] Dockerfile + K8s manifests

### Day 8-10: Wiring & Production
- [ ] Event workflows (ride.requested → dispatch → ride.assigned)
- [ ] gRPC service calls (ride → gps, pricing, user)
- [ ] Saga orchestration (CreateRideSaga)
- [ ] Idempotency testing
- [ ] Metrics (Prometheus /metrics)
- [ ] Traces (cross-service in Jaeger)
- [ ] Logs (Loki aggregation)
- [ ] Health checks (all services)
- [ ] Security (JWT, RBAC, audit)
- [ ] Integration tests (end-to-end)

---

## 🔗 REFERENCE MATERIALS AVAILABLE

All guidance documents created in Days 1-6 are available in `C:\dev\FamGo-consolidated\`:

- `SERVICE_COMPLETION_TEMPLATES.md` - Copy patterns
- `WEEKS_3-4_DELIVERY_GOVERNANCE.md` - Rules & standards
- `QUICK_REFERENCE_WEEKS_3-4.md` - Code snippets
- `AUDIT_REFERENCE_ARCHITECTURE.md` - Auth-service pattern
- `services/gps-service/` - Working implementation

---

## ✅ FINAL STATUS: READY FOR DAYS 6-10

**Completed:** 60% (Days 1-6, 48 hours) ✅  
**Remaining:** 40% (Days 6-10, 32 hours) ⏳  
**Pattern:** GPS service established as template  
**Documentation:** Complete  
**Governance:** Defined  
**Quality:** Enterprise-grade  

**All tools, templates, and documentation ready for Days 6-10 execution.**

---

## 🚀 NEXT STEPS FOR DAYS 6-10

1. **Days 6-7:** Copy GPS pattern → User service
2. **Days 7-9:** Copy GPS pattern + state machine → Ride service
3. **Days 8-10:** Wire services + production hardening

Follow same 4-layer pattern, same repository-first discipline, same quality standards.

**Target:** 100% production-ready mobility platform by end of Day 10.

