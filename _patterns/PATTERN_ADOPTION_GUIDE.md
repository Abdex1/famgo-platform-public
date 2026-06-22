# 📚 PATTERN ADOPTION GUIDE
## How to Use Extracted Patterns in FamGo Services

**Status:** Week 0 Phase 2 - COMPLETE  
**Patterns:** 8 extracted, documented, ready to use  
**Location:** `_patterns/`  
**Authority:** Must be used by all services (enforced by governance)

---

## PATTERN LIBRARY CONTENTS

```
_patterns/
├── 01-http-patterns/
│   ├── README.md (HTTP routing, middleware, handlers)
│   └── examples/ (example code snippets)
│
├── 02-service-bootstrap/
│   ├── README.md (service initialization, shutdown)
│   └── examples/
│
├── 03-kafka-patterns/
│   ├── README.md (producer, consumer, events)
│   └── examples/
│
├── 04-state-machines/
│   ├── README.md (driver states, trip states)
│   └── examples/
│
├── 05-data-access/
│   ├── README.md (database pooling, repositories)
│   └── examples/
│
├── 06-payment-gateway/
│   ├── README.md (gateway abstraction, providers)
│   └── examples/
│
├── 07-testing/
│   ├── README.md (mocks, table-driven tests, integration tests)
│   └── examples/
│
├── 08-observability/
│   ├── README.md (metrics, logging, tracing)
│   └── examples/
│
└── PATTERN_ADOPTION_GUIDE.md (this file)
```

---

## HOW TO USE THESE PATTERNS

### For Service Developers

**Step 1: Read Pattern Documentation**
- Open the relevant pattern in `_patterns/{number}-{name}/README.md`
- Understand the pattern and why it's used
- Review code examples

**Step 2: Apply Pattern to Your Service**
- Copy the pattern code to your service
- Adapt names to your domain (e.g., DriverService instead of generic Service)
- Preserve the structure and behavior

**Step 3: Integrate with FamGo Standards**
- Pattern + shared/ libraries = complete solution
- Follow ADOPTION_RULES.md and ARCHITECTURE_GUARDRAILS.md
- Use patterns within your service's architecture

**Step 4: Test**
- Follow testing pattern (Pattern 7)
- Write mocks for dependencies
- Use table-driven tests for coverage

---

## PATTERN MAPPING TO SERVICES

### Pattern 1: HTTP Handler Patterns
**Used by:**
- All 19 services that expose HTTP endpoints
- Includes: auth, user, driver, ride, dispatch, etc.

**When:**
- Implementing any HTTP endpoint
- Adding middleware
- Error handling

### Pattern 2: Service Bootstrap
**Used by:**
- All 19 services
- In: cmd/main.go

**When:**
- Service startup
- Database/Redis/Kafka initialization
- Graceful shutdown

### Pattern 3: Kafka Patterns
**Used by:**
- Services publishing events: trip, driver, payment, ride, etc.
- Services consuming events: dispatch, payment, notification, etc.

**When:**
- Publishing domain events
- Consuming Kafka events
- Building event-driven features

### Pattern 4: State Machines
**Used by:**
- Driver service (driver states)
- Trip service (trip states)
- Payment service (payment states)

**When:**
- Managing entity lifecycle
- Validating state transitions
- Implementing domain state logic

### Pattern 5: Data Access
**Used by:**
- All services with database access
- In: internal/repository/

**When:**
- Querying database
- Creating/updating entities
- Implementing repository pattern

### Pattern 6: Payment Gateway
**Used by:**
- Payment service

**When:**
- Processing payments
- Handling webhooks
- Supporting multiple providers

### Pattern 7: Testing
**Used by:**
- All 19 services
- In: internal/*_test.go files

**When:**
- Writing unit tests
- Writing integration tests
- Achieving 80%+ coverage

### Pattern 8: Observability
**Used by:**
- All 19 services

**When:**
- Exporting metrics
- Structured logging
- Distributed tracing

---

## CHECKLIST: IMPLEMENTING A SERVICE

### Pre-Implementation
- [ ] Read ADOPTION_RULES.md (governance rules)
- [ ] Read ARCHITECTURE_GUARDRAILS.md (boundaries)
- [ ] Write comparison document (MODULE_COMPARISON_TEMPLATE.md)
- [ ] Get governance approval

### During Implementation
- [ ] Use Pattern 1 for HTTP endpoints
- [ ] Use Pattern 2 for service bootstrap (cmd/main.go)
- [ ] Use Pattern 3 for Kafka events (if applicable)
- [ ] Use Pattern 4 for state machines (if applicable)
- [ ] Use Pattern 5 for database access
- [ ] Use Pattern 6 for payments (payment service only)
- [ ] Use Pattern 7 for tests (80%+ coverage)
- [ ] Use Pattern 8 for observability

### Post-Implementation
- [ ] All patterns applied correctly
- [ ] Integration with shared/ libraries
- [ ] Tests passing (80%+ coverage)
- [ ] Documentation complete
- [ ] Production readiness checklist (PRODUCTION_ACCEPTANCE_CHECKLIST.md)
- [ ] Governance approval

---

## ANTI-PATTERNS: WHAT NOT TO DO

### ❌ DO NOT
- Invent your own patterns (use extracted patterns)
- Skip using patterns for "consistency"
- Copy Uber code without extracting patterns
- Bypass shared/ libraries for custom implementation
- Skip testing requirements
- Ignore observability needs

### ✅ DO
- Use extracted patterns consistently
- Adapt patterns to your domain
- Integrate with shared/ libraries
- Follow testing requirements
- Implement observability from start

---

## PATTERN EVOLUTION

### When Pattern Becomes Obsolete
1. Document why (code comment)
2. Propose alternative
3. Governance board approval
4. Update pattern documentation
5. Communicate to all services

### When New Pattern Needed
1. Identify gap
2. Extract from proven implementation
3. Document thoroughly
4. Add to pattern library
5. Communicate availability

---

## SUPPORT

### Questions About Pattern?
- Refer to pattern README
- Check examples in examples/ directory
- Review existing service implementation
- Ask tech lead

### Pattern Not Working?
1. Verify you followed the pattern exactly
2. Check integration with shared/ libraries
3. Review test examples
4. Ask for guidance

---

## FINAL CHECKLIST: WEEK 0 PHASE 2 COMPLETE

```
✅ Pattern 1: HTTP Handlers - Extracted, documented, ready
✅ Pattern 2: Service Bootstrap - Extracted, documented, ready
✅ Pattern 3: Kafka - Extracted, documented, ready
✅ Pattern 4: State Machines - Extracted, documented, ready
✅ Pattern 5: Data Access - Extracted, documented, ready
✅ Pattern 6: Payment Gateway - Extracted, documented, ready
✅ Pattern 7: Testing - Extracted, documented, ready
✅ Pattern 8: Observability - Extracted, documented, ready

✅ Pattern Library Created (_patterns/)
✅ Examples Prepared
✅ Adoption Guide Complete
✅ Team Ready for Week 1
```

---

**Status:** Week 0 Phase 2 ✅ COMPLETE  
**All 8 Patterns:** EXTRACTED, DOCUMENTED, READY FOR USE  
**Next:** Week 1 Implementation (Auth, User, Driver services)

---
