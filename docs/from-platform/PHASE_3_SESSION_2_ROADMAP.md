# PHASE 3 SESSION 2: AUTH SERVICE IMPLEMENTATION GUIDE

**Duration**: 2-3 hours  
**Blocking**: All other services  
**Complexity**: Medium  
**Pattern**: DDD + gRPC  

---

## 🎯 OBJECTIVES

By end of Session 2:
- ✅ Auth Service can register new users
- ✅ Auth Service can login users with JWT tokens
- ✅ Auth Service validates tokens for other services
- ✅ All code follows DDD patterns
- ✅ Service emits auth events to Kafka
- ✅ Integrated with Jaeger tracing
- ✅ Tests pass (unit + integration)

---

## 📋 IMPLEMENTATION CHECKLIST

### Phase 1: Setup & Dependencies (30 min)

- [ ] Create `services/auth-service/go.mod` with dependencies
- [ ] Create `services/auth-service/cmd/main.go` entry point
- [ ] Create `services/auth-service/internal/config/config.go`
- [ ] Verify builds successfully: `go build ./cmd`

### Phase 2: Domain Layer (45 min)

- [ ] Create `internal/domain/entities/user.go`
  - User struct with ID, Email, Phone, PasswordHash, Role, Status
  - UserRole constants (rider, driver, admin, ops, fraud-agent, support, super-admin)

- [ ] Create `internal/domain/valueobjects/jwt_claims.go`
  - Claims struct with UserID, Role, SessionID, DeviceID
  - Extends jwt.RegisteredClaims

- [ ] Create `internal/domain/services/jwt_service.go`
  - GenerateAccessToken(user *User) → string
  - GenerateRefreshToken(user *User) → string
  - ValidateToken(token string) → *Claims

- [ ] Create `internal/domain/services/password_service.go`
  - HashPassword(password string) → string
  - VerifyPassword(hash, password string) → bool

- [ ] Create `internal/domain/services/rbac_service.go`
  - CheckPermission(role string, resource string) → bool
  - GetDefaultPermissions(role string) → []Permission

### Phase 3: Infrastructure Layer (45 min)

- [ ] Create `internal/infrastructure/repositories/user_repository.go`
  - CreateUser(ctx, user) → (*User, error)
  - GetUserByEmail(ctx, email) → (*User, error)
  - GetUserByID(ctx, id) → (*User, error)
  - UpdateUser(ctx, user) → error
  - Uses pgx from shared/database

- [ ] Create `internal/infrastructure/redis/session_store.go`
  - StoreSession(ctx, sessionID, userID, ttl) → error
  - GetSession(ctx, sessionID) → (userID, error)
  - DeleteSession(ctx, sessionID) → error

- [ ] Create `internal/infrastructure/redis/otp_store.go`
  - StoreOTP(ctx, email, otp, ttl) → error
  - VerifyOTP(ctx, email, otp) → bool
  - DeleteOTP(ctx, email) → error

- [ ] Create `internal/infrastructure/events/auth_events.go`
  - PublishLoginSucceeded(ctx, userID) → error
  - PublishLoginFailed(ctx, email, reason) → error
  - PublishUserRegistered(ctx, userID) → error
  - Uses event-bus from shared

### Phase 4: Application Layer (45 min)

- [ ] Create `internal/application/usecases/register_usecase.go`
  - Execute(ctx, RegisterRequest) → (RegisterResponse, error)
  - Validates email/password
  - Calls PasswordService.HashPassword
  - Calls UserRepository.CreateUser
  - Publishes UserRegistered event
  - Returns JWT tokens

- [ ] Create `internal/application/usecases/login_usecase.go`
  - Execute(ctx, LoginRequest) → (LoginResponse, error)
  - Gets user from UserRepository
  - Calls PasswordService.VerifyPassword
  - Calls JWTService.GenerateAccessToken
  - Publishes LoginSucceeded event
  - Returns JWT tokens

- [ ] Create `internal/application/usecases/refresh_usecase.go`
  - Execute(ctx, RefreshRequest) → (TokenResponse, error)
  - Validates refresh token
  - Generates new access token
  - Returns new token

### Phase 5: Interface Layer (45 min)

- [ ] Create `proto/auth.proto` with service definitions:
  ```protobuf
  service AuthService {
      rpc Register(RegisterRequest) returns (AuthResponse);
      rpc Login(LoginRequest) returns (AuthResponse);
      rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
      rpc RefreshToken(RefreshTokenRequest) returns (AuthResponse);
  }
  ```

- [ ] Generate protobuf code: `protoc --go_out=...`

- [ ] Create `internal/interfaces/grpc/auth_handler.go`
  - Implements AuthServiceServer
  - Register(ctx, req) → RegisterResponse
  - Login(ctx, req) → AuthResponse
  - ValidateToken(ctx, req) → ValidateTokenResponse
  - RefreshToken(ctx, req) → AuthResponse
  - Calls use cases
  - Logs with correlation IDs

### Phase 6: Bootstrap & Main (30 min)

- [ ] Update `cmd/main.go`:
  - Load config from environment
  - Initialize logger (Zap)
  - Initialize Jaeger tracer
  - Connect to PostgreSQL
  - Connect to Redis
  - Create gRPC server
  - Register auth service
  - Setup auth middleware
  - Enable gRPC reflection
  - Graceful shutdown on SIGTERM

### Phase 7: Testing (30 min)

- [ ] Create `internal/repositories/user_repository_test.go`
  - Test CreateUser
  - Test GetUserByEmail
  - Test GetUserByID

- [ ] Create `internal/usecases/register_usecase_test.go`
  - Test successful registration
  - Test validation errors
  - Test duplicate user

- [ ] Create `internal/usecases/login_usecase_test.go`
  - Test successful login
  - Test wrong password
  - Test non-existent user

- [ ] Integration test with docker-compose
  - Start services
  - Register user
  - Login user
  - Validate token

### Phase 8: Docker & Deployment (15 min)

- [ ] Create `Dockerfile` (multi-stage build)
- [ ] Test build: `docker build -t auth-service .`
- [ ] Update docker-compose if needed
- [ ] Test with docker-compose: `docker-compose up`

---

## 📂 FINAL STRUCTURE

```
services/auth-service/
├── cmd/
│   └── main.go                              (Entry point)
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   │   └── user.go                      (User entity + roles)
│   │   ├── valueobjects/
│   │   │   └── jwt_claims.go                (JWT claims)
│   │   ├── services/
│   │   │   ├── jwt_service.go               (Token generation)
│   │   │   ├── password_service.go          (Bcrypt operations)
│   │   │   └── rbac_service.go              (Permission checks)
│   │   └── events/
│   │       └── auth_events.go               (Domain events)
│   ├── application/
│   │   ├── dto/
│   │   │   ├── register_dto.go
│   │   │   └── login_dto.go
│   │   └── usecases/
│   │       ├── register_usecase.go          (Business logic)
│   │       ├── login_usecase.go
│   │       └── refresh_usecase.go
│   ├── infrastructure/
│   │   ├── repositories/
│   │   │   └── user_repository.go           (Database layer)
│   │   ├── redis/
│   │   │   ├── session_store.go
│   │   │   └── otp_store.go
│   │   └── events/
│   │       └── auth_events_publisher.go     (Kafka publisher)
│   ├── interfaces/
│   │   └── grpc/
│   │       ├── auth_handler.go              (gRPC endpoints)
│   │       └── health.go                    (Health check)
│   ├── config/
│   │   └── config.go                        (Config loading)
│   └── bootstrap/
│       └── bootstrap.go                     (DI container)
├── proto/
│   └── auth.proto                           (gRPC definitions)
├── tests/
│   ├── integration_test.go
│   └── fixtures.go
├── go.mod
├── go.sum
├── Dockerfile
├── .env.example
└── README.md
```

---

## 🔑 KEY PATTERNS TO FOLLOW

### 1. Use Case Pattern
```go
type RegisterUseCase struct {
    userRepo UserRepository
    passwordSvc PasswordService
    eventBus EventBus
}

func (u *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
    // 1. Validate
    // 2. Hash password
    // 3. Create user
    // 4. Publish event
    // 5. Generate tokens
    // 6. Return response
}
```

### 2. gRPC Handler Pattern
```go
type AuthHandler struct {
    registerUseCase RegisterUseCase
    loginUseCase LoginUseCase
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
    // Use correlation ID from context
    correlationID := utilities.GetCorrelationID(ctx)
    
    // Call use case
    result, err := h.registerUseCase.Execute(ctx, convertToDTO(req))
    
    // Convert to proto and return
    return convertToProto(result), nil
}
```

### 3. Repository Pattern
```go
type UserRepository interface {
    CreateUser(ctx context.Context, user *User) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    UpdateUser(ctx context.Context, user *User) error
}

type PostgresUserRepository struct {
    pool *pgxpool.Pool
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
    // Use prepared statements
    // Handle errors gracefully
    // Return created user
}
```

### 4. Event Publishing Pattern
```go
func (u *RegisterUseCase) publishUserRegistered(ctx context.Context, user *User) error {
    event := envelope.NewEventEnvelope(
        governance.EventAuthUserRegistered,
        "v1",
        "auth",
        "auth-service",
        user.ID,
        UserRegisteredPayload{UserID: user.ID, Email: user.Email},
    )
    
    return u.eventBus.Publish(ctx, governance.EventAuthUserRegistered, event)
}
```

---

## 📚 DEPENDENCIES TO ADD

```go
// go.mod
require (
    github.com/golang-jwt/jwt/v5 v5.1.0
    golang.org/x/crypto v0.17.0
    google.golang.org/grpc v1.60.0
    google.golang.org/protobuf v1.31.0
    go.uber.org/zap v1.26.0
)
```

---

## 🧪 TESTING COMMANDS

```bash
# Run unit tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run integration tests (with docker-compose up)
go test -tags integration ./tests/...

# Build
go build ./cmd

# Run locally
go run ./cmd/main.go

# Build docker image
docker build -t famgo/auth-service:latest .

# Run in docker-compose
cd infra/docker && docker-compose up auth-service
```

---

## ✅ SUCCESS CRITERIA

- [x] All 10 checklist items completed
- [x] Service builds without errors
- [x] Service starts with `docker-compose up`
- [x] Can register new user (grpcurl test)
- [x] Can login user and get JWT
- [x] Can validate token (other services can call)
- [x] Events published to Kafka
- [x] Logs appear in Loki
- [x] Traces visible in Jaeger
- [x] Unit tests pass (>80% coverage)

---

## 🚀 START IMPLEMENTING

**Estimated time**: 2-3 hours  
**Start with**: Phase 1 (Setup & Dependencies)  
**Build incrementally**: Each phase builds on previous  
**Test frequently**: After each phase  

Ready to go? Start with `services/auth-service/go.mod` and `cmd/main.go`!

---

**Previous Session**: PHASE_3_SESSION_1_COMPLETION.md ✅  
**This Guide**: PHASE_3_SESSION_2_ROADMAP.md  
**Next Sessions**: GPS → Ride → Dispatch → Flutter Apps  
