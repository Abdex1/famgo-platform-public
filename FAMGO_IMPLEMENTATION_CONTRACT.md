# 🔐 FAMGO IMPLEMENTATION CONTRACT

## CONTEXT: You Are NOT Designing Architecture

**This is NOT a green-field project.**  
**This is NOT a new development from scratch.**  
**This is IMPLEMENTATION COMPLETION of an existing enterprise platform.**

You have been given:
- ✅ Approved architecture (do NOT redesign)
- ✅ Existing service scaffolds (do NOT rebuild)
- ✅ Established patterns (do NOT introduce new ones)
- ✅ Real databases and infrastructure (do NOT mock)
- ✅ Event contracts (do NOT create local duplicates)
- ✅ Shared packages (do NOT reimplement)

---

## CRITICAL RULES - MEMORIZE IMMEDIATELY

### DO NOT

- ❌ Redesign architecture
- ❌ Change folder structure
- ❌ Introduce new patterns
- ❌ Generate placeholder code
- ❌ Generate TODOs
- ❌ Generate demo code
- ❌ Generate mock implementations
- ❌ Generate fake repositories
- ❌ Generate fake Kafka producers
- ❌ Generate fake Redis clients
- ❌ Create new patterns parallel to existing ones
- ❌ Leave empty methods
- ❌ Leave stubs or placeholder structs
- ❌ Leave TODO comments

### ALWAYS

- ✅ Finish existing implementation
- ✅ Use existing interfaces
- ✅ Use existing contracts
- ✅ Use existing events
- ✅ Use existing DTOs
- ✅ Use existing domain models
- ✅ Use existing observability systems
- ✅ Use existing security patterns
- ✅ Complete all business logic
- ✅ Complete validation
- ✅ Complete error handling
- ✅ Complete logging
- ✅ Complete tracing
- ✅ Complete metrics
- ✅ Complete tests (≥80% coverage)

---

## REPOSITORY CONSTRAINTS - EVERY SERVICE

Every service **must**:

- Own its own database (no cross-service writes)
- Communicate via APIs or events only
- Emit telemetry (traces, metrics, logs)
- Support OpenTelemetry
- Support Kafka (producer + consumer)
- Support Redis where required
- Support graceful shutdown
- Support health checks (/healthz)
- Support readiness checks (/readyz)
- Support Kubernetes deployment

---

## REQUIRED DELIVERABLES - WHEN IMPLEMENTING A FILE

### 1. Complete ALL Business Logic
- No partial implementations
- No deferred logic
- No "implement later" patterns
- All use cases fully realized

### 2. Complete Validation
- Request validation (input guards)
- Response validation (output guarantees)
- State validation (invariants)
- Constraint validation (business rules)

### 3. Complete Error Handling
- Typed errors (not generic)
- Error categorization (user, system, validation)
- Error recovery strategies
- Error logging with context
- No silent failures

### 4. Complete Logging
- Structured JSON logs
- Log levels used correctly (debug, info, warn, error)
- Context propagation (trace_id, user_id, request_id)
- Actionable messages (not debug noise)

### 5. Complete Tracing
- OpenTelemetry spans for every operation
- Span attributes (operation name, user, resource)
- Span events (state changes)
- Trace context propagated to downstream services

### 6. Complete Metrics
- Counter metrics (attempts, successes, failures)
- Histogram metrics (latency, duration)
- Gauge metrics (queue size, connection count)
- Exemplars linked to traces
- Cardinality explosion prevented

### 7. Complete Tests
- Unit tests (all public functions)
- Integration tests (with real dependencies)
- API tests (request/response validation)
- Error scenario tests (all error paths)
- Performance tests (SLA verification)
- ≥80% code coverage

---

## CODE QUALITY REQUIREMENTS

### Every Public Function

- [ ] Must be unit tested
- [ ] Must have clear contract (inputs/outputs)
- [ ] Must validate inputs
- [ ] Must handle errors
- [ ] Must emit logs
- [ ] Must emit traces
- [ ] Must emit metrics

### Every Handler

- [ ] Request validation (schema + business rules)
- [ ] Request binding (type safety)
- [ ] Response validation (contract verification)
- [ ] Response encoding (consistent format)
- [ ] Tracing (operation tracking)
- [ ] Metrics (operation counting)
- [ ] Structured logging (context enrichment)
- [ ] Error response (HTTP status + error code)

### Every Repository

- [ ] Transactions (ACID guarantees)
- [ ] Retries (transient failure handling)
- [ ] Context propagation (cancellation support)
- [ ] Connection pooling (resource management)
- [ ] Query optimization (indexes verified)
- [ ] Audit logging (access tracking)
- [ ] Error classification (constraint vs system)

### Every Kafka Consumer

- [ ] Idempotency (process-once guarantee)
- [ ] Retry handling (exponential backoff)
- [ ] DLQ support (dead letter queue)
- [ ] Offset management (progress tracking)
- [ ] Consumer group coordination
- [ ] Message ordering (per partition)
- [ ] Graceful shutdown (in-flight handling)

### Every Kafka Producer

- [ ] Schema validation (message format)
- [ ] Tracing propagation (trace context headers)
- [ ] Error handling (broker failures)
- [ ] Retry logic (transient failures)
- [ ] Compression (payload optimization)
- [ ] Partitioning (key selection)
- [ ] Monitoring (delivery metrics)

---

## COVERAGE TARGET: ≥80%

Coverage measured by:
- Line coverage (execution)
- Branch coverage (logic paths)
- Function coverage (entry points)

Do NOT stop until:
- All code paths tested
- All error cases covered
- All happy paths verified
- Integration verified
- Performance verified
- Security verified

---

## IMPLEMENTATION SEQUENCE

### Wave 1: Foundation Services (4-6 weeks)

1. **auth-service** (CURRENT)
   - JWT + refresh tokens
   - Email/phone verification
   - Password reset + MFA
   - Rate limiting + audit logs
   - OpenTelemetry integration

2. **dispatch-service**
   - Driver discovery + ranking
   - Assignment + reassignment
   - Matching sessions
   - Kafka event integration
   - Performance-optimized

### Wave 2: Core Services (4-6 weeks)

3. **ride-service**
4. **notification-service**
5. **geo-service**

### Wave 3: Financial (4-6 weeks)

6. **wallet-service**
7. **payment-service**
8. **subscription-service**

### Wave 4: Advanced (4-6 weeks)

9. **safety-service**
10. **pooling-service**
11. **analytics-service**

---

## HOW TO WORK WITH THIS CONTRACT

### For Each Service

**NEVER start with:** "Build this service from scratch"  
**ALWAYS start with:**

```markdown
## SERVICE: {name}

**Architecture:** Already exists (reference: internal/domain/)
**Interfaces:** Already defined (reference: internal/interfaces/)
**Events:** Already in shared/contracts/events/
**DTOs:** Already in shared/
**Observability:** Already in packages/telemetry/

### Current State Assessment

**What exists:** [List actual files/implementations]
**What's incomplete:** [List empty/stub files]
**What's missing:** [List required but absent]

### Implementation Tasks

**Task 1:** [Complete file X with responsibility Y]
**Task 2:** [Add missing functionality Z]
**Task 3:** [Wire A to B]

### Testing Strategy

**Unit tests:** [coverage target]
**Integration tests:** [scenario coverage]
**API tests:** [endpoint coverage]

### Acceptance Criteria

- [ ] All files have implementations (no empty stubs)
- [ ] All tests pass (including new tests)
- [ ] Coverage ≥80%
- [ ] Service compiles successfully
- [ ] All handlers have validation, logging, tracing, metrics
- [ ] All repositories have transactions and error handling
- [ ] All Kafka consumers have idempotency and DLQ
```

---

## EXECUTION PROTOCOL

### Step 1: Understand (NOT Design)
- Read existing implementation
- Understand patterns used
- Review interfaces defined
- Study domain models

### Step 2: Complete (NOT Create)
- Fill in stub implementations
- Connect existing pieces
- Add missing validators
- Implement missing handlers

### Step 3: Integrate (NOT Invent)
- Wire to existing packages
- Use existing events
- Emit to existing bus
- Store in existing database

### Step 4: Observe (NOT Mock)
- Add tracing to real operations
- Emit metrics for real actions
- Log real decisions
- Never use mock/fake producers

### Step 5: Test (NOT Verify Later)
- Write tests as you go
- Test error paths first
- Test integration paths
- Achieve 80%+ coverage
- Test doesn't halt at 80% - verify completeness

### Step 6: Deliver (NOT Demo)
- Service compiles (go build)
- Tests pass (go test ./...)
- Coverage sufficient (go test -cover)
- Deployment manifest works (kubectl apply -f)
- Health checks respond (/healthz + /readyz)

---

## SUCCESS CRITERIA

A service is **production-ready** when:

- ✅ 100% of code paths have tests or are verifiably safe
- ✅ 80%+ line coverage minimum
- ✅ All business logic is implemented (no stubs)
- ✅ All error cases are handled (no panics)
- ✅ All inputs are validated (request guards)
- ✅ All outputs are validated (response contracts)
- ✅ All operations are logged (structured, contextual)
- ✅ All operations are traced (OpenTelemetry)
- ✅ All operations are metered (Prometheus)
- ✅ All database operations are transactional
- ✅ All Kafka consumers are idempotent
- ✅ Service compiles cleanly (no warnings)
- ✅ Tests pass completely (100% pass rate)
- ✅ Integration verified (real dependencies)
- ✅ Performance verified (SLA targets met)
- ✅ Security verified (no vulnerabilities)
- ✅ Deployment verified (manifests functional)

---

## FAILURE MODES - WHAT STOPS YOU

You will be **STOPPED** if:

- ❌ Code doesn't compile (go build fails)
- ❌ Tests don't pass (go test fails)
- ❌ Coverage <80% (go test -cover)
- ❌ Empty files remain (no implementations)
- ❌ TODO comments exist (no deferred work)
- ❌ Mock implementations remain (must use real)
- ❌ Stubs exist (no partial implementations)
- ❌ Handlers lack validation (input guards missing)
- ❌ Operations lack tracing (no spans)
- ❌ Operations lack metrics (no counters)
- ❌ Operations lack logging (no context)
- ❌ Error handling incomplete (silent failures)
- ❌ Databases not transactional (race conditions)
- ❌ Kafka not idempotent (duplicate processing)
- ❌ Service won't gracefully shutdown
- ❌ Health checks not implemented
- ❌ Deployment manifest missing

---

## THIS CONTRACT IS BINDING

You accept this contract by proceeding with implementation.

By proceeding, you commit to:

1. Never designing architecture (only implementing approved design)
2. Never leaving stubs or TODOs (only complete implementations)
3. Never using fake/mock implementations (only real integrations)
4. Never reducing coverage below 80% (only increasing)
5. Never departing from established patterns (only following)
6. Delivering production-ready code (only code ready to deploy)

---

**FAMGO PLATFORM**  
**IMPLEMENTATION COMPLETION MODE: ACTIVE**  
**STANDARDS: PRODUCTION QUALITY ONLY**

---
