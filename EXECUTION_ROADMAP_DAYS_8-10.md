# 📋 EXECUTION ROADMAP: Days 8-10 (Wiring + Production Readiness)

**Phase:** Final Service Integration + Production Readiness  
**Duration:** 40 hours (Days 8-10)  
**Standard:** WEEKS_3-4_EXECUTION_ROADMAP.md  
**Status:** READY TO BEGIN

---

## PHASE 3: DAYS 8-9 WIRING (16 HOURS)

**Objective:** Connect all services via event-driven + gRPC patterns with resilience

### Day 8 Morning (4 hours): Event-Driven Workflows

**Task 1: Ride Service Event Subscriptions** (1 hour)
- [ ] Subscribe to `DriverAssigned` events from dispatch-service
- [ ] Subscribe to `PaymentProcessed` events from payment-service
- [ ] Implement event handlers for state transitions
- [ ] Test event deserialization from shared/contracts/events

**Task 2: GPS Service Event Publishing** (1 hour)
- [ ] Publish `DriverLocationUpdated` events (packages/event-bus)
- [ ] Publish `DriverOnline` / `DriverOffline` events
- [ ] Implement idempotency (deduplication via event-id)
- [ ] Test event serialization

**Task 3: User Service Event Publishing** (1 hour)
- [ ] Publish `UserRegistered` events
- [ ] Publish `UserProfileUpdated` events
- [ ] Wire EventPublisher to all handlers
- [ ] Test event publishing flow

**Task 4: Event Flow Testing** (1 hour)
- [ ] Create integration test: RideRequested → dispatch subscription
- [ ] Create integration test: DriverAssigned → ride subscription
- [ ] Create integration test: RideCompleted → payment subscription
- [ ] Verify message ordering and delivery

**Deliverables:**
- [ ] All services publishing events
- [ ] All services subscribing to required events
- [ ] Event flow tests passing (3+ workflows)
- [ ] Document: `EVENT_WORKFLOW_VERIFICATION.md`

---

### Day 8 Afternoon (4 hours): gRPC Cross-Service Communication

**Task 5: Create gRPC Client Stubs** (1.5 hours)
- [ ] Generate proto stubs for: Pricing, GPS, Dispatch services
- [ ] Create gRPC clients in packages/grpc-clients
- [ ] Implement proper error handling and timeouts
- [ ] Add circuit breaker patterns

**Task 6: Implement gRPC Calls from Ride Service** (1.5 hours)
- [ ] CreateRideCommand calls: `pricing.CalculateFare()`
- [ ] CreateRideCommand calls: `dispatch.RequestDriver()`
- [ ] AssignDriverCommand calls: `gps.SubscribeLocation()`
- [ ] All calls with 5-30s timeouts + retry logic

**Task 7: Implement gRPC Calls from Dispatch Service** (1 hour)
- [ ] `gps.GetNearbyDrivers()` for driver search
- [ ] `ride.GetRideDetails()` for ride info
- [ ] Implement caching layer for frequently called endpoints

**Deliverables:**
- [ ] All gRPC clients generated and working
- [ ] All cross-service calls implemented
- [ ] Integration tests for gRPC flows (3+ scenarios)
- [ ] Document: `GRPC_INTEGRATION_VERIFIED.md`

---

### Day 9 Morning (4 hours): Saga Orchestration + Resilience

**Task 8: Implement RideCreationSaga** (2 hours)
- [ ] Define saga steps: CreateRide → RequestDriver → CalculateFare → AssignDriver
- [ ] Implement compensation: CancelRide, ReleaseDriver, RefundPayment
- [ ] Use platform/saga for state management
- [ ] Implement timeout handling (30s total, 5s per step)

**Task 9: Circuit Breaker + Resilience** (2 hours)
- [ ] Add circuit breakers to: pricing, gps, dispatch, payment
- [ ] Configure: 50% failure threshold, 30s timeout, 60s open state
- [ ] Implement fallback for GPS location (default coordinates)
- [ ] Implement fallback for pricing (default fare calculation)
- [ ] Test circuit breaker state transitions

**Deliverables:**
- [ ] Saga implementation complete and tested
- [ ] Circuit breakers on all external calls
- [ ] Fallback strategies documented
- [ ] Document: `SAGA_AND_RESILIENCE_VERIFICATION.md`

---

### Day 9 Afternoon (4 hours): Service Discovery + Deployment Prep

**Task 10: Service Discovery Setup** (2 hours)
- [ ] Configure Kubernetes service discovery (DNS)
- [ ] Update service endpoints in bootstrap containers
- [ ] Test service-to-service connectivity
- [ ] Document: Load balancing strategy, service targets

**Task 11: Deployment Preparation** (2 hours)
- [ ] Verify all services in Kubernetes manifests
- [ ] Update image tags to DHI versions
- [ ] Set resource limits (CPU: 200m→1000m, Memory: 256Mi→512Mi)
- [ ] Configure HPA scaling policies
- [ ] Test rolling updates with 0 downtime

**Deliverables:**
- [ ] Service discovery working end-to-end
- [ ] Kubernetes manifests validated
- [ ] Rolling update tested
- [ ] Document: `DEPLOYMENT_READINESS.md`

---

## PHASE 4: DAYS 9-10 PRODUCTION READINESS (24 HOURS)

**Objective:** Full observability, security, and integration testing

### Day 9 Afternoon → Day 10 Morning (8 hours): Observability

**Task 12: Prometheus Metrics** (2.5 hours)
- [ ] Export metrics on all HTTP endpoints:
  - `request_count` by method/path
  - `request_duration_seconds` histogram
  - `request_errors_total` by error type
- [ ] Domain metrics: rides created, completed, cancelled
- [ ] Infrastructure metrics: DB pool, Redis connections
- [ ] Integration with Prometheus (scrape config)
- [ ] Test metrics in Prometheus UI

**Task 13: Jaeger Trace Propagation** (2.5 hours)
- [ ] Add OpenTelemetry instrumentation to all services
- [ ] Propagate trace context via gRPC + events
- [ ] Sample traces at 10% rate (tunable)
- [ ] Test end-to-end trace: CreateRide → Dispatch → Driver
- [ ] Verify traces in Jaeger UI

**Task 14: Loki Structured Logging** (2 hours)
- [ ] Convert all logs to structured JSON
- [ ] Include: timestamp, level, service, trace_id, operation, duration
- [ ] Log all domain events (Create, Assign, Start, Complete, Cancel)
- [ ] Log all external calls (gRPC, events)
- [ ] Verify logs aggregating in Loki

**Task 15: Grafana Dashboards** (1 hour)
- [ ] Create dashboard: Request latency (p50, p95, p99)
- [ ] Create dashboard: Error rates and types
- [ ] Create dashboard: Ride volume and completion rate
- [ ] Create dashboard: Service health (up/down)
- [ ] Create dashboard: Resource usage (CPU, memory, disk)

**Deliverables:**
- [ ] All metrics exporting to Prometheus
- [ ] All traces propagating to Jaeger
- [ ] All logs structured in Loki
- [ ] 5 Grafana dashboards created and verified
- [ ] Document: `OBSERVABILITY_COMPLETE.md`

---

### Day 10 Morning → Afternoon (8 hours): Security + Auth

**Task 16: JWT Validation** (2 hours)
- [ ] Implement auth middleware on all HTTP handlers
- [ ] Implement auth interceptor on all gRPC services
- [ ] Extract user_id and roles from JWT
- [ ] Reject invalid/expired tokens
- [ ] Test: valid token, expired token, invalid signature

**Task 17: RBAC Authorization** (2 hours)
- [ ] Define roles: PASSENGER, DRIVER, ADMIN
- [ ] Define permissions per endpoint
- [ ] Implement role checks on all handlers
- [ ] Test: authorized user, unauthorized user, missing role
- [ ] Audit log all RBAC decisions

**Task 18: Input Validation** (2 hours)
- [ ] Validate all request parameters (type, range, format)
- [ ] Validate geographic coordinates (valid lat/lon ranges)
- [ ] Validate fare amounts (non-negative, reasonable ranges)
- [ ] Reject oversized requests (payload > 1MB)
- [ ] Test: valid input, invalid input, edge cases

**Task 19: Audit Logging** (2 hours)
- [ ] Log all sensitive operations: Create ride, Assign driver, Complete ride
- [ ] Log all RBAC decisions (allow/deny)
- [ ] Log all authentication events (login, token refresh)
- [ ] Include: user, operation, result, timestamp, IP
- [ ] Store audit logs in persistent storage

**Deliverables:**
- [ ] JWT validation working on all endpoints
- [ ] RBAC policies enforced
- [ ] Input validation complete
- [ ] Audit logging active
- [ ] Document: `SECURITY_HARDENING_COMPLETE.md`

---

### Day 10 Afternoon (8 hours): Integration Testing + Finalization

**Task 20: Full Workflow Integration Tests** (3 hours)
- [ ] Test 1: CreateRide → GetRide → Assign → Start → Complete (happy path)
- [ ] Test 2: CreateRide → Assign → Cancel (cancellation flow)
- [ ] Test 3: Driver location updates via GPS service (real-time tracking)
- [ ] Test 4: Multi-ride concurrent execution (stress test)
- [ ] Verify all events published and consumed correctly
- [ ] Verify all metrics recorded

**Task 21: Failure Scenario Testing** (2 hours)
- [ ] Test: Pricing service unavailable (circuit breaker + fallback)
- [ ] Test: GPS service timeout (retry + circuit break)
- [ ] Test: Event message loss (DLQ recovery)
- [ ] Test: Database connection pool exhausted (graceful degradation)
- [ ] Test: Saga step failure (compensation triggered)

**Task 22: Performance Validation** (2 hours)
- [ ] Measure: CreateRide latency (target <100ms p95)
- [ ] Measure: GetRide latency (target <50ms p95)
- [ ] Measure: Concurrent ride throughput (target 1000+ rides/sec with 3 replicas)
- [ ] Measure: Memory usage per pod (target <512MB base)
- [ ] Load test with 10,000 concurrent users

**Task 23: Documentation + Runbooks** (1 hour)
- [ ] Create: Deployment runbook
- [ ] Create: Troubleshooting guide
- [ ] Create: Incident response procedures
- [ ] Create: Rollback procedures
- [ ] Create: Scaling guidelines

**Deliverables:**
- [ ] All integration tests passing (4+ workflows)
- [ ] All failure scenarios handled gracefully
- [ ] Performance targets met
- [ ] Production-ready documentation
- [ ] Document: `INTEGRATION_TESTING_COMPLETE.md`

---

## DELIVERABLES CHECKLIST (DAYS 8-10)

### Day 8
- [ ] Event-driven workflows working end-to-end
- [ ] gRPC cross-service communication functional
- [ ] Document: `EVENT_WORKFLOW_VERIFICATION.md`
- [ ] Document: `GRPC_INTEGRATION_VERIFIED.md`

### Day 9
- [ ] Saga orchestration implemented and tested
- [ ] Resilience patterns (circuit breaker, retry) working
- [ ] Service discovery and deployment preparation complete
- [ ] Full observability stack (Prometheus, Jaeger, Loki, Grafana)
- [ ] Document: `SAGA_AND_RESILIENCE_VERIFICATION.md`
- [ ] Document: `DEPLOYMENT_READINESS.md`
- [ ] Document: `OBSERVABILITY_COMPLETE.md`

### Day 10
- [ ] Security hardening complete (JWT, RBAC, audit)
- [ ] All integration tests passing
- [ ] Performance validated and documented
- [ ] Production-ready documentation
- [ ] Document: `SECURITY_HARDENING_COMPLETE.md`
- [ ] Document: `INTEGRATION_TESTING_COMPLETE.md`

---

## FINAL COMPLIANCE STATUS (Target)

**After Days 8-10:**

| Rule | Target | Status |
|------|--------|--------|
| 1: Events from shared/contracts | ✅ 100% | ON TRACK |
| 2: SDKs from packages | ✅ 100% | ON TRACK |
| 3: Platform abstractions | ✅ 100% | ON TRACK (saga, resilience) |
| 4: Reference architecture | ✅ 100% | ON TRACK |
| 5: No cross-service DB | ✅ 100% | ON TRACK |
| **OVERALL** | **✅ 100%** | **PRODUCTION READY** |

---

## SUCCESS CRITERIA

**Services Deployed:**
- ✅ Ride Service (100% complete)
- ✅ GPS Service (fully compliant)
- ✅ User Service (fully compliant)
- ✅ Dispatch Service (integrated)
- ✅ Pricing Service (integrated)
- ✅ Payment Service (integrated)

**Reliability:**
- ✅ All services accessible via gRPC + HTTP
- ✅ Event-driven workflows executing end-to-end
- ✅ Circuit breakers protecting against cascading failures
- ✅ Zero data loss on service failure (idempotency)

**Observability:**
- ✅ 100% of requests traced
- ✅ All metrics exported to Prometheus
- ✅ All logs aggregated in Loki
- ✅ All traces viewable in Jaeger

**Security:**
- ✅ All requests authenticated (JWT)
- ✅ All requests authorized (RBAC)
- ✅ All sensitive operations audited
- ✅ Input validation on all endpoints

**Performance:**
- ✅ CreateRide: <100ms p95
- ✅ GetRide: <50ms p95
- ✅ Throughput: 1000+ rides/sec
- ✅ Memory: <512MB per pod

---

## EXECUTION START

**Ready to begin Days 8-10?** 

**Prerequisites Met:**
- ✅ Audit Phase complete
- ✅ All critical violations fixed
- ✅ All services have compliant domain layer
- ✅ Ride Service 97% complete
- ✅ Transport layers all implemented

**Next Command:** Execute Day 8 Morning tasks (Event-Driven Workflows)

---

