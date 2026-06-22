# 🔐 FAMGO COPILOT IMPLEMENTATION CONTRACT

## CONTEXT: You Are NOT Designing Architecture

**This is NOT a green-field project.**  
**This is NOT new development from scratch.**  
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

### ALWAYS

- ✅ Finish existing implementation
- ✅ Use existing interfaces
- ✅ Use existing contracts
- ✅ Complete all business logic
- ✅ Complete validation
- ✅ Complete error handling
- ✅ Complete logging, tracing, metrics
- ✅ Complete tests (≥80% coverage)

---

## REPOSITORY CONSTRAINTS

Every service **must**:

- Own its own database
- Communicate via APIs or events only
- Emit telemetry (OpenTelemetry traces + Prometheus metrics)
- Support graceful shutdown
- Support health checks (/healthz) and readiness checks (/readyz)
- Support Kubernetes deployment

---

## CODE QUALITY REQUIREMENTS

### Every Public Function
- [ ] Unit tested
- [ ] Input validated
- [ ] Error handled
- [ ] Logged + traced + metered

### Every Handler
- [ ] Request validation
- [ ] Response validation
- [ ] Proper HTTP status codes
- [ ] Tracing + metrics + logging

### Every Repository
- [ ] Transactions (ACID)
- [ ] Retries for transients
- [ ] Context propagation
- [ ] Query optimization

---

## SUCCESS CRITERIA

A service is **production-ready** when:

- ✅ 80%+ code coverage
- ✅ All tests passing
- ✅ No stubs, TODOs, or empty methods
- ✅ All business logic complete
- ✅ All error cases handled
- ✅ All inputs validated
- ✅ All outputs validated
- ✅ All operations logged, traced, metered
- ✅ Compiles cleanly
- ✅ Deploys successfully

---

**FAMGO PLATFORM - IMPLEMENTATION COMPLETION MODE ACTIVE**  
**STANDARDS: PRODUCTION QUALITY ONLY**
