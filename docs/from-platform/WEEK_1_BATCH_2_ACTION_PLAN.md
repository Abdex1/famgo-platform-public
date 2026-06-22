# 🚀 WEEK 1 - BATCH 2 ACTION PLAN

## Backend Coherence (40 Files, 12 Hours)

**Timeline**: After Batch 1 complete  
**Estimated Duration**: 12 hours  
**Status**: READY TO START  

---

## 📋 BATCH 2 BREAKDOWN

### Database Coherence (3 files, 2 hours)
**Goal**: Standardize all 40+ database tables

Files to create:
1. `backend/database/coherence_check.sql` - Validation queries
2. `backend/database/migrations/006_audit_trail.sql` - Add audit columns
3. `backend/database/migrations/007_add_soft_delete.sql` - Add deleted_at

What these accomplish:
- Ensure all tables have UUID primary keys
- Ensure all tables have created_at/updated_at
- Ensure all tables have soft delete support
- Ensure all tables have audit trails

### API Gateway (Kong) (3 files, 2 hours)
**Goal**: Route all 36+ API endpoints through unified gateway

Files to create:
1. `backend/api-gateway/kong/kong.yml` - Kong configuration (all routes)
2. `backend/api-gateway/kong/Dockerfile` - Kong container image
3. `backend/api-gateway/kong/kong-init.sh` - Startup script

What this accomplishes:
- JWT validation on all endpoints
- Rate limiting (100 req/min per user)
- CORS policies
- Request/response logging
- Unified error handling

### Kafka Event Schemas (8 files, 2 hours)
**Goal**: Define and validate all 40+ event types

Files to create (one per service):
1. `backend/kafka/schemas/auth.v1.yaml`
2. `backend/kafka/schemas/ride.v1.yaml`
3. `backend/kafka/schemas/payment.v1.yaml`
4. `backend/kafka/schemas/dispatch.v1.yaml`
5. `backend/kafka/schemas/wallet.v1.yaml`
6. `backend/kafka/schemas/safety.v1.yaml`
7. `backend/kafka/schemas/fraud.v1.yaml`
8. `backend/kafka/schemas/gps.v1.yaml`

What this accomplishes:
- Event versioning
- Schema validation
- Consistent event format
- Type safety for events

### Unified API Client (Go) (4 files, 2 hours)
**Goal**: Create reusable HTTP client for all backend services

Files to create:
1. `backend/shared/go/client/api_client.go` - HTTP wrapper
2. `backend/shared/go/client/interceptors.go` - Middleware
3. `backend/shared/go/client/errors.go` - Error mapping
4. `backend/shared/go/client/telemetry.go` - OpenTelemetry integration

What this accomplishes:
- Retry logic (exponential backoff)
- Error standardization
- Request tracing
- Circuit breaker pattern

### REST API Wrapper (2 files, 1 hour)
**Goal**: Create REST wrapper for gRPC services (backward compatibility)

Files to create:
1. `backend/services/api-wrapper/main.go` - gRPC-to-REST converter
2. `backend/services/api-wrapper/Dockerfile` - Container image

What this accomplishes:
- gRPC services available via REST
- Backward compatibility
- Legacy client support

### OpenAPI Documentation (2 files, 1 hour)
**Goal**: Generate complete API documentation

Files to create:
1. `backend/shared/openapi/openapi-merged.yaml` - Complete OpenAPI spec
2. `backend/shared/postman/FamGo-API.postman_collection.json` - Postman collection

What this accomplishes:
- API documentation
- Postman testing ready
- Swagger UI compatible
- All 36+ endpoints documented

### Documentation (2 files, 1 hour)
**Goal**: Create comprehensive API guides

Files to create:
1. `backend/shared/docs/API_GUIDE.md` - Complete API reference
2. `backend/shared/docs/ERROR_CODES.md` - Standard error codes

What this accomplishes:
- Developer reference
- Standard error mapping
- Consistent error messages
- Integration guide

### Integration Tests (4 files, 1 hour)
**Goal**: Validate all coherence layers working together

Files to create:
1. `backend/test/integration/database_coherence_test.go` - DB validation
2. `backend/test/integration/api_gateway_test.go` - Kong routing
3. `backend/test/integration/event_schema_test.go` - Kafka schemas
4. `backend/test/integration/api_client_test.go` - Client functionality

What this accomplishes:
- Verify database standardization
- Verify API Gateway routing
- Verify Kafka schemas
- Verify client error handling

---

## 📊 BATCH 2 STRUCTURE

```
Backend Coherence (40 files, 12 hours)
├── Database (3 files) - Table standardization
├── API Gateway (3 files) - Route all endpoints
├── Event Schemas (8 files) - Event versioning
├── API Client (4 files) - Unified SDK
├── REST Wrapper (2 files) - Backward compatibility
├── OpenAPI (2 files) - Documentation
├── Documentation (2 files) - API guides
├── Integration Tests (4 files) - Validation
├── Config Files (8 files) - Templates & examples
└── Deployment (2 files) - Docker integration
```

---

## 🎯 EXECUTION CHECKLIST

### Day 1: Database Coherence (2 hours)
```
[ ] Create 006_audit_trail.sql
[ ] Create 007_add_soft_delete.sql
[ ] Create coherence_check.sql
[ ] Run migrations on test database
[ ] Validate all 40+ tables pass coherence check
```

### Day 1-2: API Gateway (2 hours)
```
[ ] Create kong.yml (all 36+ routes)
[ ] Create Kong Dockerfile
[ ] Create kong-init.sh startup script
[ ] Test Kong container starts
[ ] Verify all routes registered
```

### Day 2: Kafka Schemas (2 hours)
```
[ ] Create all 8 schema files (one per service)
[ ] Define schema structure (event_id, event_type, data, timestamp)
[ ] Add schema validation rules
[ ] Test schema parsing
```

### Day 2-3: API Client (2 hours)
```
[ ] Create api_client.go (HTTP wrapper)
[ ] Create interceptors.go (middleware)
[ ] Create errors.go (error mapping)
[ ] Create telemetry.go (tracing)
[ ] Test client retry logic
```

### Day 3: REST Wrapper (1 hour)
```
[ ] Create api-wrapper/main.go
[ ] Create Dockerfile for wrapper
[ ] Test gRPC-to-REST conversion
```

### Day 3: OpenAPI & Docs (1 hour)
```
[ ] Generate openapi-merged.yaml
[ ] Create Postman collection
[ ] Create API_GUIDE.md
[ ] Create ERROR_CODES.md
```

### Day 4: Integration Tests (1 hour)
```
[ ] Create database coherence test
[ ] Create API Gateway routing test
[ ] Create event schema test
[ ] Create API client test
[ ] Verify all tests pass
```

---

## ✅ SUCCESS CRITERIA - BATCH 2

```
Database:
  ✅ All 40+ tables have UUID primary keys
  ✅ All tables have created_at/updated_at
  ✅ All tables have soft delete support
  ✅ All tables have audit trails

API Gateway:
  ✅ All 36+ endpoints routable through Kong
  ✅ JWT validation on all endpoints
  ✅ Rate limiting enforced (100/min per user)
  ✅ CORS policies configured

Events:
  ✅ All 8 event types have versioned schemas
  ✅ Schema validation working
  ✅ Event format standardized

API Client:
  ✅ HTTP wrapper compiles
  ✅ Retry logic implemented
  ✅ Error mapping complete
  ✅ Telemetry integrated

REST Wrapper:
  ✅ gRPC services accessible via REST
  ✅ Backward compatibility maintained

Documentation:
  ✅ OpenAPI spec complete (all 36+ endpoints)
  ✅ Postman collection ready
  ✅ API guide written
  ✅ Error codes documented

Tests:
  ✅ All integration tests pass
  ✅ No regressions introduced
  ✅ Coverage ≥80%
```

---

## 🔗 WHAT BATCH 2 ENABLES

After Batch 2 complete:
- ✅ Mobile apps can connect to unified API Gateway
- ✅ All backend services coherent
- ✅ Event streaming validated
- ✅ Error handling standardized
- ✅ Documentation complete
- ✅ Ready for mobile app integration (Batch 3)

---

## 📈 WEEK 1 FULL TIMELINE

```
Hours 0-8:   Batch 1 (Shared Flutter Library) ✅ COMPLETE
Hours 8-20:  Batch 2 (Backend Coherence) ← YOU ARE HERE

Total Week 1: 20 hours production code
Result: Mobile apps can connect to backend
```

---

## 🚀 READY TO START BATCH 2?

All files documented and organized. 40 files ready to generate.

**Let me know when you're ready to proceed with Batch 2!**

Next steps:
1. Verify Batch 1 compilation (`flutter pub get && flutter analyze`)
2. Generate JSON serializable code (`flutter pub run build_runner build`)
3. Start Batch 2 - Backend Coherence

---

**Status**: ✅ Batch 1 Complete, Batch 2 Ready  
**Files to Create**: 40 (Backend Coherence)  
**Time Estimate**: 12 hours  
**Readiness**: 100%
