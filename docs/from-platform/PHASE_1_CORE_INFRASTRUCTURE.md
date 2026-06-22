# PHASE 1 - CORE INFRASTRUCTURE SETUP
## Enterprise Foundation: Database, Auth, Gateway, Event Bus

**Objective**: Establish the foundational infrastructure layer supporting all 18+ microservices.

**Timeline**: 2-3 weeks
**Deliverables**: 5 core systems + integration tests

---

## Phase 1 Breakdown

### 1.1 Database Layer Setup ✓ Planned

**Objective**: PostgreSQL + PostGIS initialization

**Tasks**:
1. Create PostgreSQL container with PostGIS extension
2. Create initial schema migration (001_initial_schema.sql)
3. Create databases:
   - `famgo` (transactional)
   - `famgo_analytics` (analytics)
4. Set up pgvector extension for embeddings
5. Run migrations

**Files to Create**:
```
database/migrations/
  └── 001_initial_schema.sql    # Core tables: users, drivers, rides, bookings
      ├── users
      ├── drivers
      ├── vehicles
      ├── ride_requests
      ├── ride_sessions
      ├── bookings
      ├── feedback
      ├── wallet_transactions
      └── audit_logs
```

**Commands**:
```bash
# Start PostgreSQL
docker-compose -f infra/docker/docker-compose.yml up postgres

# Run migrations
pnpm run db:migrate

# Verify
psql -U famgo -d famgo -c "\dt"
```

---

### 1.2 Authentication Service ✓ Planned

**Objective**: Rewrite FastAPI auth layer in Go with enterprise security

**Source**: `C:\dev\FamGo\backend\app\routes\auth.py` + `backend/app/utils/jwt_handler.py`

**New Structure**:
```
services/auth-service/
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   │   ├── user.go         # User entity
│   │   │   ├── session.go      # Session entity
│   │   │   ├── device.go       # Device fingerprinting
│   │   │   └── refresh_token.go
│   │   ├── services/
│   │   │   ├── jwt_service.go
│   │   │   ├── otp_service.go  # One-time password
│   │   │   ├── password_service.go
│   │   │   └── rbac_service.go # Role-based access control
│   │   └── events/
│   │       └── auth_events.go  # Events emitted
│   ├── application/
│   │   ├── commands/
│   │   │   ├── register.go
│   │   │   ├── login.go
│   │   │   └── refresh_token.go
│   │   └── queries/
│   │       └── get_user.go
│   ├── infrastructure/
│   │   ├── postgres/
│   │   │   └── user_repository.go
│   │   ├── redis/
│   │   │   └── session_store.go
│   │   ├── security/
│   │   │   ├── jwt_manager.go
│   │   │   ├── token_hasher.go
│   │   │   └── device_fingerprint.go
│   │   └── vault/
│   │       └── client.go       # Secrets from HashiCorp Vault
│   └── interfaces/
│       ├── rest/
│       │   ├── handlers/
│       │   │   └── auth_handler.go
│       │   └── routes/
│       │       └── routes.go
│       └── grpc/
│           └── auth_service.pb.go
├── cmd/
│   └── api/
│       └── main.go
├── migrations/
│   └── 001_initial_schema.sql
└── go.mod
```

**Key Features to Implement**:
- JWT token generation & validation
- OTP (for payment security)
- Device fingerprinting (fraud prevention)
- Session management
- RBAC (rider, driver, admin, operator, support)
- Password hashing (bcrypt)
- Token refresh
- Audit logging

**Endpoints**:
```
POST   /v1/auth/register        # Register user
POST   /v1/auth/login           # Login
POST   /v1/auth/refresh         # Refresh token
POST   /v1/auth/logout          # Logout
POST   /v1/auth/otp/request     # Request OTP
POST   /v1/auth/otp/verify      # Verify OTP
GET    /v1/auth/me              # Get current user
PUT    /v1/auth/profile         # Update profile
POST   /v1/auth/change-password # Change password
```

**Tests Required**:
- Unit: JWT generation, validation, expiry
- Integration: PostgreSQL + Redis interaction
- E2E: Full auth flow

---

### 1.3 API Gateway (Kong) ✓ Planned

**Objective**: Centralized API routing, rate limiting, auth

**Configuration**:
```yaml
# infra/docker/kong/kong.yml
services:
  - name: auth-service
    url: http://auth-service:3000
    routes:
      - /v1/auth
    plugins:
      - rate-limiting (100 req/min)
      - jwt-auth
      - request-id
  
  - name: user-service
    url: http://user-service:3000
    routes:
      - /v1/users
    plugins:
      - rate-limiting (200 req/min)
      - jwt-auth
      - rbac-enforcement
  
  - name: ride-service
    url: http://ride-service:3000
    routes:
      - /v1/rides
    plugins:
      - rate-limiting (500 req/min)
      - jwt-auth
```

**Features**:
- Request routing to microservices
- Rate limiting per endpoint
- JWT validation
- Request/Response transformation
- CORS handling
- API versioning

---

### 1.4 Event Bus (Kafka) Setup ✓ Planned

**Objective**: Event streaming infrastructure

**Topics to Create**:
```bash
# Authentication events
kafka-topics --create --topic auth.user.registered
kafka-topics --create --topic auth.user.logged_in
kafka-topics --create --topic auth.token.refreshed

# Ride events
kafka-topics --create --topic ride.created
kafka-topics --create --topic ride.matched
kafka-topics --create --topic ride.started
kafka-topics --create --topic ride.completed
kafka-topics --create --topic ride.cancelled

# Driver events
kafka-topics --create --topic driver.registered
kafka-topics --create --topic driver.online
kafka-topics --create --topic driver.offline
kafka-topics --create --topic driver.location.updated

# Payment events
kafka-topics --create --topic payment.requested
kafka-topics --create --topic payment.completed
kafka-topics --create --topic payment.failed

# Notification events
kafka-topics --create --topic notification.send.sms
kafka-topics --create --topic notification.send.push
```

**Event Governance**:
```
shared/contracts/events/
├── schemas/
│   ├── auth_events.go
│   ├── ride_events.go
│   ├── driver_events.go
│   └── payment_events.go
├── governance/
│   ├── naming_conventions.md
│   ├── retention_policies.md
│   └── schema_versioning.md
└── catalog.json              # Event catalog
```

**Event Example**:
```go
// shared/contracts/events/auth_events.go
type UserRegisteredEvent struct {
    EventID       string    `json:"event_id"`
    EventType     string    `json:"event_type"`      // "auth.user.registered"
    AggregateID   string    `json:"aggregate_id"`    // user_id
    Timestamp     time.Time `json:"timestamp"`
    Version       int       `json:"version"`
    
    UserID        string    `json:"user_id"`
    Email         string    `json:"email"`
    Role          string    `json:"role"`            // "rider", "driver", "admin"
    RegisteredAt  time.Time `json:"registered_at"`
}
```

---

### 1.5 Redis Cache & GEO Setup ✓ Planned

**Objective**: In-memory cache + geospatial queries

**Configuration**:
```
# For sessions
SETEX session:${sessionID} 3600 ${userData}

# For geospatial driver indexing
GEOADD drivers:geo 13.361389 38.115556 "driver:123"  # Palermo
GEOADD drivers:geo 15.087269 37.502669 "driver:456"  # Catania
GEORADIUSBYMEMBER drivers:geo driver:123 50 km

# For real-time presence
SET driver:123:online "true" EX 300
```

**Data Structures**:
- Sessions (hash + TTL)
- GEO index of drivers
- Rate limit counters
- OTP codes (TTL)
- Cached user profiles

---

## Phase 1 Execution Checklist

### Week 1: Database & Auth Service
- [ ] PostgreSQL container + PostGIS setup
- [ ] Initial schema migration
- [ ] Auth service skeleton Go code
- [ ] JWT implementation
- [ ] PostgreSQL repository layer
- [ ] Unit tests

### Week 2: Gateway & Event Bus
- [ ] Kong gateway configuration
- [ ] Service registration
- [ ] Rate limiting setup
- [ ] Kafka topics creation
- [ ] Event schemas defined
- [ ] Integration tests

### Week 3: Testing & Integration
- [ ] E2E auth flow tests
- [ ] Gateway routing verification
- [ ] Kafka producer/consumer tests
- [ ] Documentation
- [ ] Production readiness review

---

## Key Metrics (Success Criteria)

✓ **Database**:
- PostgreSQL container healthy
- All migrations run successfully
- PostGIS queries working

✓ **Auth Service**:
- JWT tokens generated and validated
- 95% code coverage in tests
- Login/Register endpoints responding

✓ **Gateway**:
- All services registered
- Rate limits enforced
- <50ms latency overhead

✓ **Event Bus**:
- Kafka topics created
- Producer/consumer working
- Event contracts defined

✓ **Redis**:
- Cache hits >80%
- GEO queries <5ms
- Session management working

---

## Next Phase Preparation

Once Phase 1 complete:
- Services can use auth via Kong gateway
- Services can emit events to Kafka
- Cache layer available for all services
- Database schema ready for new services

**Phase 2 Focus**: Core domain services (Ride, Driver, User)

---

## Files to Create/Modify

**New Files**:
```
database/migrations/001_initial_schema.sql
services/auth-service/go.mod
services/auth-service/cmd/api/main.go
services/auth-service/internal/domain/entities/user.go
services/auth-service/internal/domain/services/jwt_service.go
services/auth-service/internal/infrastructure/postgres/user_repository.go
services/auth-service/internal/infrastructure/redis/session_store.go
services/auth-service/internal/infrastructure/security/jwt_manager.go
services/auth-service/internal/interfaces/rest/handlers/auth_handler.go
services/auth-service/internal/interfaces/rest/routes/routes.go
shared/contracts/events/auth_events.go
shared/contracts/events/ride_events.go
shared/contracts/grpc/auth_service.proto
infra/docker/docker-compose.yml
```

**Already Created** (Phase 0):
```
.github/workflows/ci.yml
package.json
tsconfig.json
turbo.json
MIGRATION_MAPPING.md
```

---

## Command Reference

```bash
# Start infrastructure
docker-compose -f infra/docker/docker-compose.yml up -d

# View logs
docker-compose -f infra/docker/docker-compose.yml logs -f postgres

# Create auth service
mkdir -p services/auth-service/{cmd/api,internal/{domain,application,infrastructure,interfaces}}

# Initialize Go module
cd services/auth-service && go mod init github.com/FamGo/platform/auth-service

# Run tests
go test ./...

# Build
go build -o bin/auth-service cmd/api/main.go

# Run auth service
./bin/auth-service

# Verify Kong
curl http://localhost:8001/services/

# Kafka topics
kafka-topics --create --topic auth.user.registered --bootstrap-server localhost:9092
```

---

**Phase 1 Owner**: Backend Lead
**Phase 1 Status**: Ready to execute
**Estimated Start**: After Phase 0 ✓
