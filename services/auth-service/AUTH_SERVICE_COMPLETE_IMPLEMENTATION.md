# 🔐 AUTH-SERVICE: COMPLETE IMPLEMENTATION SPECIFICATION

**Status:** Production Implementation Phase  
**Timeline:** Days 1-5 (40 hours)  
**Target:** 100% Production Ready  
**Coverage Target:** ≥80% with all tests passing  

---

## EXECUTIVE SUMMARY

Auth-service is the **foundation service** for FamGo platform. It must be **100% complete** before proceeding to Wave 2 services.

**Current Assessment:**
- ✅ Foundation code exists (main.go, validation.go, telemetry.go)
- ✅ Dependencies configured (go.mod complete)
- ⏳ 60% infrastructure skeleton in place
- ❌ Internal implementations mostly empty (need completion)
- ❌ Tests incomplete or missing
- ❌ Some handlers not implemented

**This Document:** Complete task breakdown to production readiness.

---

## PART 1: ARCHITECTURE UNDERSTANDING

### Service Responsibility
Auth-service owns:
- User registration (email + phone)
- User login (JWT auth)
- Token management (access + refresh)
- Email verification
- Phone verification
- Password reset with OTP
- MFA support
- Session management
- Rate limiting
- Audit logging

### Service Boundaries
Auth-service **DOES NOT**:
- Manage user profiles (user-service owns)
- Process rides (ride-service owns)
- Handle payments (payment-service owns)
- Store location (gps-service owns)

### Event-Driven Integration
Auth-service **publishes**:
- `user.registered.v1`
- `user.email_verified.v1`
- `user.phone_verified.v1`
- `user.password_reset.v1`
- `user.mfa_enabled.v1`

Auth-service **consumes**:
- (None - pure producer for now)

---

## PART 2: CURRENT STATE ANALYSIS

### Files That Exist ✅

```
services/auth-service/
├── cmd/main.go                      ✅ Main bootstrap (mostly complete)
├── validation.go                    ✅ Validation layer (complete, using go-playground/validator)
├── telemetry.go                     ⏳ Telemetry (framework exists, needs integration)
├── go.mod                           ✅ Dependencies configured
├── internal/
│   ├── config/                      ❌ Empty (need implementation)
│   ├── domain/
│   │   ├── auth.go                  ⏳ Domain logic (partial)
│   │   ├── entities/                ❌ Empty (need User, Token entities)
│   │   ├── events/                  ❌ Empty (need event definitions)
│   │   ├── services/                ❌ Empty (need domain services)
│   │   └── valueobjects/            ❌ Empty (need JWT, Password VOs)
│   ├── application/                 ❌ Empty (need use cases)
│   ├── handler/                     ❌ Empty (need HTTP handlers)
│   ├── infrastructure/              ❌ Empty (need external integrations)
│   ├── repository/
│   │   └── repository.go            ⏳ Repository interface (partial)
│   ├── bootstrap/                   ❌ Empty (need DI setup)
│   └── service/                     ❌ Empty (need application services)
└── Dockerfile                       ✅ Container definition
```

### Critical Missing Implementations ❌

| Component | Status | Priority | Effort |
|-----------|--------|----------|--------|
| Config loader | Empty | CRITICAL | 2h |
| User entity | Missing | CRITICAL | 2h |
| User repository | Stub | CRITICAL | 4h |
| JWT service | Partial | CRITICAL | 3h |
| Password service | Missing | CRITICAL | 2h |
| RBAC service | Partial | CRITICAL | 2h |
| OTP service | Missing | CRITICAL | 3h |
| HTTP handlers | Empty | CRITICAL | 6h |
| Middleware (validation, auth, logging) | Missing | HIGH | 4h |
| Tests | Minimal | HIGH | 8h |
| Database migrations | Partial | MEDIUM | 2h |
| Event publishing | Missing | MEDIUM | 3h |

---

## PART 3: COMPLETE IMPLEMENTATION TASKS

### TASK 1: Configuration Layer (2 hours)

**File:** `services/auth-service/internal/config/config.go`

**Responsibility:** Load and validate environment configuration

**Implementation Requirements:**

```go
type Config struct {
    // Server
    Port        string
    Environment string
    
    // Database
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBMaxConn  int
    DBIdleConn int
    
    // JWT
    JWTSecret           string
    JWTExpirationHours  int
    RefreshExpiration   int
    
    // OTP
    OTPLength           int
    OTPExpirationMins   int
    OTPMaxAttempts      int
    
    // Observability
    OTelEndpoint string
    LogLevel     string
    
    // External Services
    SMSProvider      string
    SMSAPIKey        string
    EmailProvider    string
    EmailAPIKey      string
    
    // Security
    RateLimitPerMin   int
    PasswordMinLength int
    SessionTimeout    int
}

// Load() error - loads from environment with validation
// Validate() error - validates all required fields
```

**Tests Required:**
- Load valid config
- Load with missing required fields (error)
- Load with invalid values (error)
- Validate all constraints

---

### TASK 2: Domain Entities (2 hours)

**File:** `services/auth-service/internal/domain/entities/user.go`

**Responsibility:** Define User domain entity with business logic

```go
type User struct {
    ID            uuid.UUID
    Email         string
    PasswordHash  string
    Phone         string
    FirstName     string
    LastName      string
    Role          Role          // PASSENGER, DRIVER, ADMIN
    Status        UserStatus    // PENDING, ACTIVE, SUSPENDED
    EmailVerified bool
    PhoneVerified bool
    MFAEnabled    bool
    MFAMethod     string        // SMS, EMAIL, AUTHENTICATOR
    LastLogin     *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// Methods
func (u *User) VerifyPassword(plaintext string) bool
func (u *User) IsActive() bool
func (u *User) CanLogin() error
func (u *User) VerifyEmail()
func (u *User) VerifyPhone()
func (u *User) EnableMFA(method string)
func (u *User) UpdateLastLogin()
```

**Tests Required:**
- User creation
- Password verification
- Status checks
- MFA operations
- Validation logic

---

### TASK 3: Value Objects (2 hours)

**File:** `services/auth-service/internal/domain/valueobjects/jwt.go`

**Responsibility:** JWT token generation and validation

```go
type TokenPair struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int
}

type Claims struct {
    UserID string
    Email  string
    Role   string
    // Standard JWT claims
    jwt.RegisteredClaims
}

type JWTService struct {
    secret     string
    expiration time.Duration
}

func (js *JWTService) GenerateAccessToken(userID, email, role string) (string, error)
func (js *JWTService) GenerateRefreshToken(userID string) (string, error)
func (js *JWTService) VerifyToken(token string) (*Claims, error)
func (js *JWTService) RefreshAccessToken(refreshToken string) (string, error)
func (js *JWTService) RevokeToken(token string) error
```

**Tests Required:**
- Token generation
- Token verification
- Token expiration
- Token refresh
- Token revocation
- Invalid token handling

---

### TASK 4: User Repository (4 hours)

**File:** `services/auth-service/internal/repository/user_repository.go`

**Responsibility:** Database persistence for User entity

```go
type UserRepository interface {
    // Create new user (with email uniqueness constraint)
    Create(ctx context.Context, user *User) error
    
    // Find operations
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByPhone(ctx context.Context, phone string) (*User, error)
    
    // Update operations (use transactions)
    Update(ctx context.Context, user *User) error
    UpdateStatus(ctx context.Context, id uuid.UUID, status UserStatus) error
    UpdateLastLogin(ctx context.Context, id uuid.UUID) error
    UpdateMFASettings(ctx context.Context, id uuid.UUID, enabled bool, method string) error
    
    // Verification operations
    VerifyEmail(ctx context.Context, userID uuid.UUID) error
    VerifyPhone(ctx context.Context, userID uuid.UUID) error
    
    // Query operations
    ListByRole(ctx context.Context, role Role, limit, offset int) ([]*User, error)
    ListActive(ctx context.Context, limit, offset int) ([]*User, error)
    
    // Cleanup
    Delete(ctx context.Context, id uuid.UUID) error // Soft delete
}

// Implementation: PostgresUserRepository
type PostgresUserRepository struct {
    db *sqlx.DB
}

// All methods with:
// - Transaction wrapping where needed
// - Proper error handling
// - Context cancellation support
// - Query optimization
// - Audit logging
```

**Critical Implementation Details:**

```go
// Example: Create with transaction
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    // Insert user
    query := `INSERT INTO users (...) VALUES (...) RETURNING id`
    err = tx.QueryRowContext(ctx, query, ...).Scan(...)
    if err != nil {
        if isUniqueConstraintError(err) {
            return ErrEmailAlreadyExists
        }
        return err
    }
    
    // Log audit event
    err = r.logAudit(ctx, tx, user.ID, "USER_CREATED", nil)
    if err != nil {
        return err
    }
    
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }
    
    return nil
}
```

**Tests Required:**
- Create user (success + duplicate email error)
- Find operations (success + not found)
- Update operations
- Verification operations
- Transaction rollback on error
- Concurrent access handling

---

### TASK 5: Domain Services (2 hours)

**File:** `services/auth-service/internal/domain/services/auth_service.go`

**Responsibility:** Core authentication business logic

```go
type AuthDomainService struct {
    userRepo UserRepository
    jwtSvc   JWTService
    passwordSvc PasswordService
    otpSvc   OTPService
}

// Registration
func (s *AuthDomainService) RegisterUser(ctx context.Context, email, password, phone, firstName, lastName string) (*User, error) {
    // Validate input
    // Check email doesn't exist
    // Hash password
    // Create user entity
    // Persist to repository
    // Publish user.registered event
    // Return user
}

// Login
func (s *AuthDomainService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
    // Find user by email
    // Verify password
    // Check user is active
    // Generate tokens
    // Update last login
    // Return token pair
}

// Token refresh
func (s *AuthDomainService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
    // Verify refresh token
    // Extract user ID
    // Generate new access token
    // Return new token
}

// Password reset
func (s *AuthDomainService) RequestPasswordReset(ctx context.Context, email string) error {
    // Find user
    // Generate OTP
    // Send OTP via email/SMS
    // Store OTP with expiration
    // Return
}

// Complete password reset
func (s *AuthDomainService) CompletePasswordReset(ctx context.Context, email, otp, newPassword string) error {
    // Verify OTP
    // Validate new password
    // Update user password
    // Invalidate OTP
    // Publish event
}

// Email verification
func (s *AuthDomainService) RequestEmailVerification(ctx context.Context, userID uuid.UUID) error
func (s *AuthDomainService) VerifyEmail(ctx context.Context, userID uuid.UUID, otp string) error

// MFA
func (s *AuthDomainService) EnableMFA(ctx context.Context, userID uuid.UUID, method string) error
func (s *AuthDomainService) VerifyMFA(ctx context.Context, userID uuid.UUID, code string) error
```

**Tests Required:**
- Registration flow (success + validation errors)
- Login flow (success + invalid credentials)
- Token refresh
- Password reset flow
- Email verification
- MFA operations

---

### TASK 6: HTTP Handlers (6 hours)

**File:** `services/auth-service/internal/handler/auth_handler.go`

**Responsibility:** HTTP request/response handling

```go
type AuthHandler struct {
    authSvc AuthService
    log     Logger
    metrics AuthMetrics
}

// POST /auth/signup
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "signup")
    defer span.End()
    
    // Parse request
    var req SignupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.errorResponse(w, http.StatusBadRequest, "invalid_request", err.Error())
        return
    }
    
    // Validate request
    if err := ValidateSignup(req); err != nil {
        span.AddEvent("validation_failed")
        h.log.Warn("signup validation failed", map[string]interface{}{"error": err})
        h.errorResponse(w, http.StatusBadRequest, "validation_error", err.Error())
        return
    }
    
    // Call service
    user, err := h.authSvc.Register(ctx, req.Email, req.Password, req.Phone, req.FirstName, req.LastName)
    if err != nil {
        span.AddEvent("registration_failed")
        h.log.Error("signup failed", map[string]interface{}{"error": err})
        
        if errors.Is(err, ErrEmailAlreadyExists) {
            h.errorResponse(w, http.StatusConflict, "email_exists", "Email already registered")
            return
        }
        h.errorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to register user")
        return
    }
    
    // Record metrics
    h.metrics.SignupSuccess.Add(ctx, 1)
    h.log.Info("user registered", map[string]interface{}{"user_id": user.ID})
    
    // Return response
    h.successResponse(w, http.StatusCreated, map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
    })
}

// POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) { ... }

// POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) { ... }

// POST /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) { ... }

// POST /auth/password-reset/request
func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) { ... }

// POST /auth/password-reset/verify
func (h *AuthHandler) VerifyPasswordReset(w http.ResponseWriter, r *http.Request) { ... }

// POST /auth/email/verify/request
func (h *AuthHandler) RequestEmailVerification(w http.ResponseWriter, r *http.Request) { ... }

// POST /auth/email/verify
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) { ... }

// POST /auth/mfa/enable
func (h *AuthHandler) EnableMFA(w http.ResponseWriter, r *http.Request) { ... }

// POST /auth/mfa/verify
func (h *AuthHandler) VerifyMFA(w http.ResponseWriter, r *http.Request) { ... }
```

**Every Handler Must Have:**
- ✅ Request validation (using ValidateXxx functions)
- ✅ Request binding (json.Decoder or similar)
- ✅ Response validation (verify response structure)
- ✅ Error handling (with proper HTTP status codes)
- ✅ Logging (structured, contextual)
- ✅ Tracing (OpenTelemetry spans)
- ✅ Metrics (counters, histograms)

**Tests Required:**
- Successful flow for each handler
- Validation error cases
- Business logic error cases (email exists, invalid credentials)
- Missing required fields
- Invalid data types
- Response structure validation

---

### TASK 7: Middleware (4 hours)

**Files:** 
- `services/auth-service/internal/handler/middleware/auth_middleware.go`
- `services/auth-service/internal/handler/middleware/rate_limit.go`
- `services/auth-service/internal/handler/middleware/logging.go`

**Authentication Middleware:**
```go
func AuthMiddleware(jwtSvc JWTService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from header
            token := extractToken(r)
            if token == "" {
                http.Error(w, "missing token", http.StatusUnauthorized)
                return
            }
            
            // Verify token
            claims, err := jwtSvc.VerifyToken(token)
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }
            
            // Add to context
            ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

**Rate Limiting Middleware:**
```go
func RateLimitMiddleware(limiter RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            clientIP := getClientIP(r)
            if !limiter.Allow(clientIP) {
                http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

**Tests Required:**
- Middleware chaining
- Token validation
- Rate limiting enforcement
- Logging output

---

### TASK 8: Comprehensive Tests (8 hours)

**Test Coverage Requirements:**

```
handlers/        - 80%+ coverage (all happy paths + error cases)
repository/      - 90%+ coverage (all database operations)
domain/services/ - 85%+ coverage (all business logic)
middleware/      - 80%+ coverage (all middleware functions)
integration/     - 70%+ coverage (end-to-end flows)
```

**Test Files to Create:**

1. `services/auth-service/internal/handler/auth_handler_test.go` (800+ lines)
   - Signup endpoint tests (success, validation errors, duplicate email)
   - Login endpoint tests (success, invalid credentials)
   - Token refresh tests
   - Password reset tests
   - Email verification tests
   - MFA tests

2. `services/auth-service/internal/repository/user_repository_test.go` (600+ lines)
   - Create user tests
   - Find operations (by ID, email, phone)
   - Update operations
   - Verification operations
   - Concurrent access tests
   - Transaction rollback tests

3. `services/auth-service/internal/domain/services/auth_service_test.go` (500+ lines)
   - Registration flow
   - Login flow
   - Token operations
   - Password reset flow
   - Email/phone verification
   - MFA operations

4. `services/auth-service/tests/integration_test.go` (400+ lines)
   - Full signup → login → verify email flow
   - Password reset flow
   - Token refresh flow
   - Rate limiting enforcement
   - Concurrent user registration

**Test Infrastructure:**
```go
// TestDB setup for integration tests
func setupTestDB(t *testing.T) *sqlx.DB {
    // Connect to test PostgreSQL
    // Run migrations
    // Return connection
}

// Mock implementations for external services
type MockSMSProvider struct { ... }
type MockEmailProvider struct { ... }
```

---

### TASK 9: Database Migrations (2 hours)

**File:** `services/auth-service/internal/migrations/migrations.go`

**Schema Requirements:**

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'passenger',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_method VARCHAR(50),
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP -- Soft delete
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status) WHERE deleted_at IS NULL;

-- OTP verification
CREATE TABLE otp_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    otp_code VARCHAR(6) NOT NULL,
    otp_type VARCHAR(50), -- EMAIL, PHONE, PASSWORD_RESET
    expires_at TIMESTAMP NOT NULL,
    attempts INT DEFAULT 0,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_otp_user_id ON otp_verifications(user_id);
CREATE INDEX idx_otp_expires ON otp_verifications(expires_at);

-- Session tracking
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    refresh_token_hash VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

-- Audit log
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID,
    actor_type VARCHAR(50),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    changes JSONB,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_actor ON audit_logs(actor_id);
CREATE INDEX idx_audit_action ON audit_logs(action);
CREATE INDEX idx_audit_created ON audit_logs(created_at);
```

---

### TASK 10: Event Publishing (3 hours)

**File:** `services/auth-service/internal/events/event_publisher.go`

**Events to Publish:**

```go
// From shared/contracts/events/auth_events.go
type UserRegistered struct {
    UserID    uuid.UUID
    Email     string
    Phone     string
    FirstName string
    LastName  string
    Role      string
    Timestamp time.Time
}

type UserEmailVerified struct {
    UserID    uuid.UUID
    Email     string
    Timestamp time.Time
}

type UserPhoneVerified struct {
    UserID    uuid.UUID
    Phone     string
    Timestamp time.Time
}

type PasswordResetCompleted struct {
    UserID    uuid.UUID
    Email     string
    Timestamp time.Time
}

type MFAEnabled struct {
    UserID    uuid.UUID
    Method    string
    Timestamp time.Time
}

// Publisher implementation
type EventPublisher struct {
    kafkaProducer KafkaProducer
    tracer        Tracer
    logger        Logger
}

func (p *EventPublisher) PublishUserRegistered(ctx context.Context, event UserRegistered) error {
    // Validate event
    // Serialize to JSON
    // Set trace context headers
    // Publish to Kafka topic "user.registered.v1"
    // Log event
    // Emit metric
}
```

---

## PART 4: TESTING STRATEGY

### Coverage Targets

| Component | Target | Priority |
|-----------|--------|----------|
| Handlers | 80%+ | CRITICAL |
| Repository | 90%+ | CRITICAL |
| Domain Services | 85%+ | CRITICAL |
| Middleware | 80%+ | HIGH |
| Integration | 70%+ | HIGH |

### Test Execution

```bash
# Run all tests with coverage
go test -v -cover -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out

# Run integration tests
go test -v -tags=integration ./tests/...

# Benchmark tests
go test -bench=. -benchmem ./...
```

### CI/CD Validation

```yaml
# All tests must pass
# Coverage must be ≥80%
# No SQL errors
# No race conditions (go test -race)
# All handlers return proper HTTP status codes
# All errors logged with trace context
```

---

## PART 5: DEPLOYMENT CHECKLIST

### Pre-Deployment

- [ ] All tests passing (go test ./...)
- [ ] Coverage ≥80% (go test -cover)
- [ ] No race conditions (go test -race)
- [ ] Code compiles cleanly (go build ./cmd)
- [ ] No lint warnings (golangci-lint)
- [ ] Database migrations run successfully
- [ ] Kubernetes manifests valid (kubeval)
- [ ] Health checks respond (/healthz + /readyz)
- [ ] Environment variables documented
- [ ] Secrets configured in deployment

### Deployment Manifest

```yaml
# services/auth-service/k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: famgo/auth-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: auth-config
              key: db_host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: auth-secrets
              key: db_password
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

---

## PART 6: SUCCESS CRITERIA

### Code Quality
- ✅ All files implement their responsibility (no stubs)
- ✅ No TODO comments
- ✅ No empty methods
- ✅ All error cases handled
- ✅ All inputs validated
- ✅ All outputs verified

### Testing
- ✅ Coverage ≥80%
- ✅ All tests passing
- ✅ No race conditions
- ✅ Integration tests verify full flows
- ✅ Load tests verify performance

### Observability
- ✅ All operations traced (OpenTelemetry)
- ✅ All operations metered (Prometheus)
- ✅ All operations logged (structured JSON)
- ✅ Logs include trace context
- ✅ Health checks functional

### Deployment
- ✅ Dockerfile builds successfully
- ✅ Container runs without errors
- ✅ Kubernetes manifests valid
- ✅ Database migrations apply cleanly
- ✅ Service starts with graceful shutdown support

### Production Readiness
- ✅ Handles concurrent users
- ✅ Rate limiting enforced
- ✅ Audit logs captured
- ✅ JWT tokens validated
- ✅ Password hashed securely
- ✅ OTP flow secure
- ✅ Session management working
- ✅ Error responses appropriate (no sensitive data leaking)

---

## PART 7: IMPLEMENTATION ORDER

### Day 1-2: Foundation (8 hours)
1. Config loader
2. Database migrations
3. Domain entities
4. Value objects (JWT, Password)

### Day 3: Repository (6 hours)
1. User repository interface
2. Postgres implementation
3. Query optimization
4. Transaction handling

### Day 4: Handlers (8 hours)
1. Auth handlers
2. Middleware
3. Route registration
4. Error responses

### Day 5: Tests & Deployment (8 hours)
1. Unit tests
2. Integration tests
3. Kubernetes manifests
4. Documentation

---

## REFERENCE: EXISTING CODE TO PRESERVE

**DO NOT MODIFY:**
- cmd/main.go (bootstrap exists)
- validation.go (already complete)
- telemetry.go (framework exists)
- go.mod (dependencies configured)

**COMPLETE/EXTEND:**
- internal/config (currently empty)
- internal/domain (add missing pieces)
- internal/repository (complete implementations)
- internal/handler (add all handlers)

**CREATE NEW:**
- internal/middleware
- internal/events
- internal/bootstrap
- tests/

---

**END OF SPECIFICATION**

This document is your complete blueprint for implementing auth-service to production quality. Follow it exactly. No deviations.

**Next Action:** Proceed to file-by-file implementation starting with Task 1 (Config Layer).
