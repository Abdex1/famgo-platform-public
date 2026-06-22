# 📚 WEEK 1-4: COMPLETE FOUNDATION IMPLEMENTATION GUIDE

**Status:** Step 3 Complete - Now Beginning Week 1-4 Implementation  
**Phase:** Foundation & Core Services  
**Duration:** 4 weeks / 20 working days / ~160 hours  
**Team:** 3-4 engineers  

---

## WEEK 1: FOUNDATION PHASE (40 Hours)

### Days 1-2: Auth Service Deep Review & Plan (8 Hours)

#### Task 1: Examine Current Auth Service
**File:** `services/auth-service/`

**Status Check:**
- JWT service: ✅ Exists (jwt_service.go)
- Password service: ✅ Exists (password_service.go)
- RBAC service: ✅ Exists (rbac_service.go)
- Database layer: ✅ Exists (user_repository.go)
- gRPC: ✅ Defined (auth.proto)

**What's Missing:**
- Database migrations: ❌
- Input validation: ❌
- Comprehensive tests: ❌
- Observability: ❌
- Kubernetes manifests: ❌

#### Deliverables for Days 1-2:
1. **IMPLEMENTATION_PLAN.md** - Created in auth-service directory
2. **Architecture Document** - Describe current implementation
3. **Gap Analysis** - List missing components
4. **Backup** - Copy current version to backups/

#### Git Commit Template:
```bash
git add services/auth-service/IMPLEMENTATION_PLAN.md
git commit -m "docs: auth-service implementation plan and gap analysis

- Analyzed current JWT, RBAC, password services
- Identified gaps: migrations, validation, tests, observability
- Created implementation roadmap (4 priorities)
- Estimated total effort: 30-40 hours"
```

---

### Days 2-3: Database Migrations Implementation (8 Hours)

#### Task: Create Database Schema with Migrations

**Files to Create:**
1. `services/auth-service/db/migrations/000001_create_initial_schema.up.sql`
2. `services/auth-service/db/migrations/000001_create_initial_schema.down.sql`

#### Implementation Checklist:

**Create Tables:**
- [ ] users table (with role ENUM)
- [ ] sessions table (for session tracking)
- [ ] otp_codes table (for OTP MFA)
- [ ] roles table (role definitions)
- [ ] permissions table (permission definitions)
- [ ] role_permissions junction table (mapping)
- [ ] audit_logs table (audit trail)

**Add Indexes:**
- [ ] users(email, phone, role, status)
- [ ] sessions(user_id, expires_at, device_id)
- [ ] otp_codes(user_id, code, expires_at)
- [ ] audit_logs(user_id, created_at, resource_type)

**Add Constraints:**
- [ ] NOT NULL on required fields
- [ ] UNIQUE on email, phone
- [ ] CHECK constraints on role/status
- [ ] FOREIGN KEY constraints

#### Code to Create:

```sql
-- File: db/migrations/000001_create_initial_schema.up.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  phone VARCHAR(20) UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  role VARCHAR(50) NOT NULL DEFAULT 'passenger',
  status VARCHAR(50) NOT NULL DEFAULT 'active',
  phone_verified BOOLEAN DEFAULT false,
  email_verified BOOLEAN DEFAULT false,
  two_fa_enabled BOOLEAN DEFAULT false,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- [Add 6 more tables as shown in db/migrations above]

-- Add indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
-- [Add remaining indexes]
```

#### Git Commit:
```bash
git add services/auth-service/db/migrations/
git commit -m "feat: auth-service database schema migrations

- Create users, sessions, otp_codes tables
- Create roles and permissions tables
- Add audit_logs for compliance
- Include foreign keys and indexes
- Support soft-delete pattern"
```

---

### Days 3-4: Input Validation Layer (8 Hours)

#### Files to Create:
1. `services/auth-service/validation.go` - Request validation
2. Update `services/auth-service/auth_handler.go` - Use validation

#### Validation Rules to Implement:

```go
// SignupRequest
Email:     "required,email"
Password:  "required,min=8,max=128,strong"
FirstName: "required,min=2,max=100"
LastName:  "required,min=2,max=100"
Phone:     "required,e164"

// LoginRequest
Email:     "required,email"
Password:  "required,min=8"

// PasswordResetRequest
Email:          "required,email"
NewPassword:    "required,min=8,max=128,strong"
OTPCode:        "required,len=6,numeric"

// UpdateProfileRequest
FirstName: "required,min=2,max=100"
LastName:  "required,min=2,max=100"
Phone:     "required,e164"
```

#### Custom Validators to Add:
- Strong password (uppercase, numbers, special chars)
- Valid phone number (E.164 format)
- Valid email domain

#### Git Commit:
```bash
git add services/auth-service/validation.go
git add services/auth-service/auth_handler.go
git commit -m "feat: input validation for auth-service

- Add validator for all requests
- Implement custom validators (strong password, phone)
- Handle validation errors gracefully
- Return proper HTTP status codes"
```

---

### Days 4-5: Comprehensive Testing (12 Hours)

#### Files to Create:
1. `services/auth-service/auth_handler_test.go` - Integration tests
2. `services/auth-service/validation_test.go` - Validation tests
3. `services/auth-service/jwt_service_test.go` - JWT tests

#### Test Cases to Implement:

**Signup Tests (6 cases):**
- [ ] Successful signup with valid data
- [ ] Duplicate email rejection
- [ ] Weak password rejection
- [ ] Invalid email format rejection
- [ ] Missing required fields rejection
- [ ] Invalid phone format rejection

**Login Tests (4 cases):**
- [ ] Successful login
- [ ] Invalid credentials
- [ ] Account not found
- [ ] Account disabled

**JWT Tests (5 cases):**
- [ ] Token generation
- [ ] Token validation
- [ ] Token expiration
- [ ] Token refresh
- [ ] Invalid token rejection

**RBAC Tests (3 cases):**
- [ ] Admin can access admin endpoints
- [ ] Passenger cannot access driver endpoints
- [ ] Proper error on unauthorized access

**Total: 18+ test cases with 80%+ coverage**

#### Git Commit:
```bash
git add services/auth-service/*_test.go
git commit -m "test: comprehensive test suite for auth-service

- Add 18+ test cases for signup/login/JWT/RBAC
- Achieve 80%+ code coverage
- Test error cases and edge cases
- Include integration tests with mock database"
```

---

### Day 5: Observability Integration (4 Hours)

#### Files to Create/Update:
1. `services/auth-service/telemetry.go` - Tracing & metrics setup
2. Update `services/auth-service/main.go` - Initialize telemetry
3. Update `services/auth-service/auth_handler.go` - Add trace spans

#### Observability Components:

**Jaeger Tracing:**
```go
- InitTracer() function
- Trace spans for each handler
- Custom attributes (user_id, request_id, etc.)
```

**Prometheus Metrics:**
```
- HTTP request counter
- HTTP request duration
- Auth-specific metrics (signup, login, failed attempts)
- JWT validation metrics
```

**Structured Logging:**
```go
- Request logging
- Error logging
- Audit logging for important events
```

#### Git Commit:
```bash
git add services/auth-service/telemetry.go
git add services/auth-service/main.go
git commit -m "feat: opentelemetry observability integration

- Add Jaeger tracing initialization
- Implement Prometheus metrics export
- Add structured logging with context
- Enable distributed tracing across services"
```

---

### End of Week 1 Deliverables:

✅ **Auth Service Status:**
- Database migrations: 100% ✅
- Input validation: 100% ✅
- Tests: 80%+ coverage ✅
- Observability: 100% ✅
- Kubernetes manifests: NEXT (Week 2)

✅ **Code Quality:**
- All tests passing
- No linting errors
- 80%+ test coverage
- Documentation complete

✅ **Git History:**
- 5 clean commits
- Clear commit messages
- Proper attribution

---

## WEEK 2: KUBERNETES & CI/CD (40 Hours)

### Days 1-2: Kubernetes Manifests (12 Hours)

#### Files to Create:
1. `infra/kubernetes/base/auth-service.yaml` - Deployment
2. `infra/kubernetes/base/auth-service-service.yaml` - Service
3. `infra/kubernetes/base/auth-service-configmap.yaml` - ConfigMap
4. `infra/kubernetes/base/auth-service-secret.yaml` - Secrets
5. `infra/kubernetes/base/auth-service-hpa.yaml` - Autoscaling

#### Deployment Specification:

```yaml
# Requirements:
replicas: 3
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"

# Probes:
livenessProbe: /health/live (15s delay, 10s period)
readinessProbe: /health/ready (10s delay, 5s period)
startupProbe: /health (30s delay, 10s period)

# Autoscaling:
minReplicas: 3
maxReplicas: 10
targetCPU: 70%
targetMemory: 80%
```

#### Networking:
- Service type: ClusterIP
- Port: 80
- Target port: 8080
- Labels for monitoring

#### Security:
- RunAsNonRoot: true
- ReadOnlyRootFilesystem: true
- No privileged access

#### Git Commit:
```bash
git add infra/kubernetes/base/
git commit -m "infra: kubernetes manifests for auth-service

- Create deployment with 3 replicas
- Add health checks (liveness, readiness, startup)
- Configure HPA (3-10 replicas, 70% CPU target)
- Implement security constraints
- Setup resource requests/limits"
```

---

### Days 3-4: CI/CD Pipelines (16 Hours)

#### Files to Create:
1. `.github/workflows/build.yaml` - Build & test
2. `.github/workflows/docker.yaml` - Docker image build
3. `.github/workflows/deploy.yaml` - Kubernetes deployment
4. `.github/workflows/security.yaml` - Security scanning
5. `.github/workflows/lint.yaml` - Code quality

#### Build Workflow (.github/workflows/build.yaml):
```yaml
name: Build & Test
on: [push, pull_request]
jobs:
  - Checkout code
  - Setup Go 1.21
  - Run go mod download
  - Run go test ./... (coverage)
  - Run go build
  - Upload coverage to Codecov
```

#### Docker Workflow (.github/workflows/docker.yaml):
```yaml
name: Build & Push Docker Image
on: [push: main]
jobs:
  - Checkout
  - Setup Buildx
  - Login to registry
  - Build multi-stage image
  - Push to registry
  - Scan image for vulnerabilities
```

#### Deploy Workflow (.github/workflows/deploy.yaml):
```yaml
name: Deploy to Kubernetes
on: [push: main]
jobs:
  - Checkout
  - Setup kubeconfig
  - Apply manifests
  - Verify deployment
  - Wait for ready
  - Run smoke tests
```

#### Security Workflow (.github/workflows/security.yaml):
```yaml
name: Security Scan
on: [pull_request]
jobs:
  - SAST scanning
  - Dependency check
  - Container scanning
  - Secret scanning
```

#### Lint Workflow (.github/workflows/lint.yaml):
```yaml
name: Lint & Format Check
on: [pull_request]
jobs:
  - Go fmt check
  - Go vet
  - Golangci-lint
  - Markdown lint
```

#### Git Commit:
```bash
git add .github/workflows/
git commit -m "ci/cd: complete github actions workflows

- Add build & test workflow
- Add docker build & push workflow
- Add kubernetes deployment workflow
- Add security scanning
- Add linting workflow"
```

---

### Days 5: Testing & Verification (12 Hours)

#### Testing Stages:

**Unit Tests:**
- [ ] All handler functions tested
- [ ] Validation logic tested
- [ ] Error handling tested
- [ ] Coverage: 80%+

**Integration Tests:**
- [ ] Database connectivity
- [ ] JWT generation and validation
- [ ] RBAC enforcement
- [ ] Error scenarios

**E2E Tests:**
- [ ] Full signup flow
- [ ] Full login flow
- [ ] Token refresh flow
- [ ] Password reset flow

**Security Tests:**
- [ ] SQL injection protection
- [ ] Rate limiting
- [ ] CORS configuration
- [ ] TLS/SSL

#### Deployment Tests:
- [ ] Docker image builds
- [ ] Kubernetes manifests valid
- [ ] Services start healthy
- [ ] Health checks pass
- [ ] Metrics exported
- [ ] Traces collected

#### Git Commit:
```bash
git add tests/
git commit -m "test: end-to-end testing and validation

- Add 50+ integration tests
- Add 10+ deployment tests
- Verify all workflows
- Confirm observability working"
```

---

### End of Week 2 Deliverables:

✅ **Kubernetes Readiness:**
- Manifests created and validated ✅
- 3 replicas with autoscaling ✅
- Health checks configured ✅
- Security hardened ✅

✅ **CI/CD Pipeline:**
- 5 GitHub Actions workflows ✅
- Automated testing ✅
- Automated building ✅
- Automated deployment ✅
- Security scanning ✅

✅ **Deployment Verification:**
- All tests passing ✅
- Workflows green ✅
- Service healthy ✅
- Metrics flowing ✅

---

## WEEK 3-4: CORE SERVICES (80 Hours)

### Week 3: User Service + Ride Service (40 Hours)

#### Days 1-2: User Service Architecture (8 Hours)
- Create `services/user-service/` directory structure
- Implement user profile management
- Create user repository layer
- Add KYC verification workflow

#### Days 3-5: Ride Service Implementation (32 Hours)
- Create `services/ride-service/` directory structure
- Implement ride state machine
- Create ride lifecycle handlers
- Add event publishing
- Implement ride history
- Add rating/feedback system

#### Deliverables:
- [ ] User service: CRUD operations
- [ ] Ride service: Lifecycle management
- [ ] Database migrations: 20+ tables
- [ ] API contracts: gRPC + REST
- [ ] Tests: 80%+ coverage
- [ ] Kubernetes manifests: Both services

---

### Week 4: Dispatch Service + GPS Service (40 Hours)

#### Days 1-3: Dispatch Service Implementation (20 Hours)
- Implement matching algorithm
- Create driver ranking system
- Add ETA calculation
- Build supply-demand balancing

#### Days 4-5: GPS Service Real-time (20 Hours)
- Implement WebSocket server
- Redis GEO integration
- Real-time location tracking
- Passenger map updates

#### Deliverables:
- [ ] Dispatch service: Matching engine
- [ ] GPS service: Real-time tracking
- [ ] WebSocket gateway
- [ ] Redis integration
- [ ] Event streaming setup
- [ ] All tests and manifests

---

## END OF MONTH 1 CHECKPOINT

### Services Completed:
✅ Auth Service - 100% production-ready
✅ User Service - 80% complete
✅ Ride Service - 80% complete
✅ Dispatch Service - 80% complete
✅ GPS Service - 80% complete

### Infrastructure:
✅ Docker-compose running
✅ Kubernetes cluster ready
✅ CI/CD pipelines automated
✅ Monitoring active
✅ Logging aggregated

### Quality Metrics:
✅ 80%+ test coverage
✅ All security checks passing
✅ Zero critical vulnerabilities
✅ Performance meets targets
✅ Uptime: 99.9%+

---

## 📋 IMPLEMENTATION COMMANDS REFERENCE

### Setup Week 1:
```bash
# Create migrations
mkdir -p services/auth-service/db/migrations
# Create validation
touch services/auth-service/validation.go
# Create tests
touch services/auth-service/auth_handler_test.go
```

### Deploy Week 2:
```bash
# Create Kubernetes manifests
mkdir -p infra/kubernetes/base
# Create GitHub Actions
mkdir -p .github/workflows
```

### Implement Week 3-4:
```bash
# Create user service
cp -r services/_template services/user-service
# Create ride service
cp -r services/_template services/ride-service
```

---

## ⏱️ TIME BREAKDOWN

| Week | Phase | Hours | Team | Status |
|------|-------|-------|------|--------|
| 1 | Auth Service Foundation | 40 | 2 eng | Active |
| 2 | Kubernetes & CI/CD | 40 | 2-3 eng | Planning |
| 3 | User & Ride Services | 40 | 3 eng | Queued |
| 4 | Dispatch & GPS Services | 40 | 3 eng | Queued |

**Total:** 160 hours = 4 weeks (5 working days each)

---

**All planned, all documented, ready for execution!**

Next: Execute Week 1-4 following this guide.

