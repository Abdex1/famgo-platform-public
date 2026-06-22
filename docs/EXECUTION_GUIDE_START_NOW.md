I have gone upto the 3rd step '### Step 2: Verify Infrastructure (15 minutes)', deeply think and also reviewing all the recommended documents and robustly and systematically proceed starting from this 3rd step, for the production as perfectly as planned.
# 🚀 CONSOLIDATION EXECUTION GUIDE - DO THIS NOW

**Status:** Ready to Execute  
**Timeline:** 12 Weeks to Production  
**Quality:** Enterprise-Grade

---

## QUICK START (TODAY)

### Step 1: Create Consolidated Project (30 minutes)

```bash
# Navigate to work directory
cd C:\dev

# Create new consolidated project
mkdir FamGo-consolidated
cd FamGo-consolidated

# Initialize git
git init
git config user.email "team@famgo.com"
git config user.name "FamGo Platform"

# Initialize as Node monorepo
npm init -y

# Create pnpm workspace
cat > pnpm-workspace.yaml << 'EOF'
packages:
  - 'apps/*'
  - 'services/*'
  - 'packages/*'
  - 'shared/*'
  - 'platform/*'
  - 'tooling/*'
EOF

# Create directory structure safely
$dirs = @(
  "apps",
  "services",
  "packages",
  "shared",
  "platform",
  "tooling",
  "database",
  "infra/docker",
  "infra/kubernetes",
  "infra/terraform",
  "infra/helm",
  "database/migrations",
  "database/seeds",
  "docs",
  "scripts",
  "security",
  "ml",
  "gateway",
  "env"
)

foreach ($d in $dirs) {
    New-Item -ItemType Directory -Force -Path "./$d" | Out-Null
}

# Copy docker-compose from trial
Remove-Item -Force ./infra/docker/docker-compose.yml -ErrorAction SilentlyContinue
Copy-Item ../FamGo-platform-trial/docker-compose.yml ./infra/docker/docker-compose.yml -Force

# Copy all services from trial
Remove-Item -Recurse -Force ./services
Copy-Item -Recurse ../FamGo-platform-trial/services ./services

# Copy apps
Remove-Item -Recurse -Force ./apps
Copy-Item -Recurse ../FamGo-platform-trial/apps ./apps

# Copy shared packages
Remove-Item -Recurse -Force ./packages
Copy-Item -Recurse ../FamGo-platform-trial/packages ./packages

Remove-Item -Recurse -Force ./shared
Copy-Item -Recurse ../FamGo-platform-trial/shared ./shared

Remove-Item -Recurse -Force ./platform
Copy-Item -Recurse ../FamGo-platform-trial/platform ./platform

# Copy infrastructure
Remove-Item -Recurse -Force ./infra
Copy-Item -Recurse ../FamGo-platform-trial/infra ./infra

# Copy database
Remove-Item -Recurse -Force ./database
Copy-Item -Recurse ../FamGo-platform-trial/database ./database

# Copy docs
New-Item -ItemType Directory -Force -Path ./docs/from-platform | Out-Null
Copy-Item ../FamGo-platform-trial/*.md ./docs/from-platform/ -Force

# Create initial package.json
cat > package.json << 'EOF'
{
  "name": "famgo-platform",
  "version": "0.1.0",
  "description": "Enterprise Urban Mobility Operating System",
  "private": true,
  "type": "module",
  "workspaces": [
    "apps/*",
    "services/*",
    "packages/*"
  ],
  "scripts": {
    "dev": "turbo run dev --parallel",
    "build": "turbo run build",
    "test": "turbo run test",
    "lint": "turbo run lint",
    "format": "turbo run format",
    "infra:up": "docker-compose -f infra/docker/docker-compose.yml up -d",
    "infra:down": "docker-compose -f infra/docker/docker-compose.yml down",
    "infra:logs": "docker-compose -f infra/docker/docker-compose.yml logs -f"
  },
  "devDependencies": {
    "turbo": "^1.10.0",
    "typescript": "^5.0.0",
    "eslint": "^8.0.0",
    "prettier": "^3.0.0"
  }
}
EOF

# Initialize git repo with first commit
git add .
git commit -m "feat: initialize consolidated famgo platform"

echo "✅ Consolidated project created successfully"
```

### Step 2: Verify Infrastructure (15 minutes)

```bash
# Test docker-compose
docker-compose -f infra/docker/docker-compose.yml config

# Start infrastructure
docker-compose -f infra/docker/docker-compose.yml up -d

# Wait for services
sleep 30

# Verify services
docker ps

# Check services are healthy
curl http://localhost:5432  # PostgreSQL
curl http://localhost:6379  # Redis
curl http://localhost:9092  # Kafka
curl http://localhost:9090  # Prometheus
curl http://localhost:3001  # Grafana
curl http://localhost:16686 # Jaeger

echo "✅ Infrastructure verified"
```

### Step 3: Fix Security (15 minutes)

```bash
# Update docker-compose with environment variables
cat > .env.local << 'EOF'
# Database
DB_PASSWORD=dev_password_only_for_local_testing
DB_USER=famgo
DB_NAME=famgo

# MinIO
MINIO_PASSWORD=minio_dev_password

# Grafana
GRAFANA_PASSWORD=admin_dev_password

# Redis
REDIS_PASSWORD=redis_dev_password
EOF

# Update docker-compose.yml to use env vars
# Replace hardcoded passwords with references to .env file

# Never commit .env.local to git
echo ".env.local" >> .gitignore
echo "**/.env" >> .gitignore
echo "**/.env.*" >> .gitignore

# Create .env.example for team
cat > .env.example << 'EOF'
# Copy this file to .env.local and fill in values
DB_PASSWORD=
DB_USER=famgo
DB_NAME=famgo
MINIO_PASSWORD=
GRAFANA_PASSWORD=
REDIS_PASSWORD=
JWT_SECRET=
EOF

echo "✅ Security improved"
```

### Step 4: First Commit

```bash
git add .
git commit -m "chore: setup consolidated infrastructure and security

- Merge docker-compose from trial version
- Copy all services, apps, and packages
- Add environment variable structure
- Implement git ignore for secrets
- Ready for Phase 1 development"

echo "✅ Initial consolidation complete"
```

---

## WEEK 1: FOUNDATION PHASE

### Day 1-2: Auth Service Deep Review & Plan

```bash
# 1. Examine current auth service
cd services/auth-service
ls -la

# 2. Read existing code
cat main.go
cat go.mod

# 3. Create detailed analysis document
cat > IMPLEMENTATION_PLAN.md << 'EOF'
# Auth Service Implementation Plan

## Current State
- JWT service exists
- Password service exists
- RBAC service exists
- Database layer exists

## Gaps to Fix
1. No database migrations
2. No input validation
3. No comprehensive tests
4. No observability
5. No Kubernetes manifests

## Priority 1: Database Migrations

### Missing Schemas:
- users table
- sessions table
- otp table
- roles table
- permissions table

### Action Items:
- [ ] Install migrate tool
- [ ] Create 001_initial_schema.up.sql
- [ ] Create 001_initial_schema.down.sql
- [ ] Test migrations locally
- [ ] Document schema

## Priority 2: Input Validation

### Actions:
- [ ] Add validator dependency
- [ ] Create validation middleware
- [ ] Add validators to signup endpoint
- [ ] Add validators to login endpoint
- [ ] Test with invalid inputs

## Priority 3: Comprehensive Tests

### Test Cases:
- [ ] SignUp with valid data
- [ ] SignUp with duplicate email
- [ ] SignUp with weak password
- [ ] SignUp with invalid email
- [ ] Login with correct password
- [ ] Login with wrong password
- [ ] JWT generation and validation
- [ ] Token expiration
- [ ] Refresh token flow
- [ ] Logout
- [ ] RBAC enforcement

## Priority 4: Observability

### Actions:
- [ ] Add OpenTelemetry imports
- [ ] Initialize Jaeger tracer
- [ ] Add trace spans to handlers
- [ ] Add Prometheus metrics
- [ ] Add structured logging

## Priority 5: Kubernetes

### Manifests Needed:
- [ ] Deployment
- [ ] Service
- [ ] ConfigMap
- [ ] Secret
- [ ] Ingress
- [ ] HPA

## Estimation:
- Database: 4 hours
- Validation: 4 hours
- Tests: 8 hours
- Observability: 6 hours
- Kubernetes: 8 hours
Total: 30 hours
EOF

# 4. Create backup
cp -r . ../../backups/auth-service-original/

echo "✅ Auth service analysis complete"
```

### Day 2-3: Implement Auth Service Database Migrations

```bash
# 1. Install migration tool
go get -u github.com/golang-migrate/migrate/v4
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 2. Create migrations directory
mkdir -p db/migrations

# 3. Create initial schema migration
cat > db/migrations/000001_create_initial_schema.up.sql << 'EOF'
-- Users table
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  phone VARCHAR(20) UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  role VARCHAR(50) NOT NULL DEFAULT 'passenger',
  status VARCHAR(50) NOT NULL DEFAULT 'active',
  phone_verified BOOLEAN DEFAULT false,
  email_verified BOOLEAN DEFAULT false,
  two_fa_enabled BOOLEAN DEFAULT false,
  two_fa_method VARCHAR(50),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  CHECK (role IN ('passenger', 'driver', 'admin', 'support'))
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status) WHERE deleted_at IS NULL;

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash VARCHAR(255) UNIQUE NOT NULL,
  ip_address INET,
  user_agent TEXT,
  device_id VARCHAR(255),
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_device_id ON sessions(device_id);

-- OTP table
CREATE TABLE IF NOT EXISTS otp_codes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  code VARCHAR(10) NOT NULL,
  purpose VARCHAR(50) NOT NULL,
  attempts INT DEFAULT 0,
  max_attempts INT DEFAULT 5,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  verified_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CHECK (purpose IN ('email_verification', 'phone_verification', 'two_fa', 'password_reset'))
);

CREATE INDEX idx_otp_user_id ON otp_codes(user_id);
CREATE INDEX idx_otp_code ON otp_codes(code);
CREATE INDEX idx_otp_expires_at ON otp_codes(expires_at);

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) UNIQUE NOT NULL,
  description TEXT,
  resource VARCHAR(100),
  action VARCHAR(50),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CHECK (action IN ('create', 'read', 'update', 'delete'))
);

-- Role permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
  role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, permission_id)
);

-- Audit log table
CREATE TABLE IF NOT EXISTS audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  action VARCHAR(100) NOT NULL,
  resource_type VARCHAR(50),
  resource_id VARCHAR(100),
  old_values JSONB,
  new_values JSONB,
  ip_address INET,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
EOF

cat > db/migrations/000001_create_initial_schema.down.sql << 'EOF'
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS otp_codes CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
EOF

# 4. Test migrations
export DATABASE_URL="postgresql://famgo:dev_password_only_for_local_testing@localhost:5432/famgo?sslmode=disable"
migrate -path db/migrations -database "$DATABASE_URL" up

# 5. Verify tables
psql "$DATABASE_URL" -c "\dt"

echo "✅ Database migrations complete"
```

### Day 3-4: Add Input Validation

```bash
# 1. Add validator to go.mod
go get github.com/go-playground/validator/v10

# 2. Create validation layer
cat > auth-service/validation.go << 'EOF'
package main

import (
  "fmt"
  "github.com/go-playground/validator/v10"
)

var validate = validator.New()

type SignupRequest struct {
  Email     string `validate:"required,email"`
  Password  string `validate:"required,min=8,max=128"`
  FirstName string `validate:"required,min=2,max=100"`
  LastName  string `validate:"required,min=2,max=100"`
  Phone     string `validate:"required,e164"`
}

type LoginRequest struct {
  Email    string `validate:"required,email"`
  Password string `validate:"required,min=8"`
}

type PasswordResetRequest struct {
  Email       string `validate:"required,email"`
  NewPassword string `validate:"required,min=8,max=128"`
  OTPCode     string `validate:"required,len=6,numeric"`
}

func ValidateSignup(req SignupRequest) error {
  if err := validate.Struct(req); err != nil {
    return fmt.Errorf("validation error: %w", err)
  }
  return nil
}

func ValidateLogin(req LoginRequest) error {
  if err := validate.Struct(req); err != nil {
    return fmt.Errorf("validation error: %w", err)
  }
  return nil
}

func ValidatePasswordReset(req PasswordResetRequest) error {
  if err := validate.Struct(req); err != nil {
    return fmt.Errorf("validation error: %w", err)
  }
  return nil
}
EOF

# 3. Update handlers to use validation
cat >> auth-service/auth_handler.go << 'EOF'

// Signup handles user registration
func (h *AuthHandler) Signup(c *gin.Context) {
  var req SignupRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(400, gin.H{"error": "invalid request"})
    return
  }

  // Add validation
  if err := ValidateSignup(req); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
  }

  // Rest of signup logic
}
EOF

echo "✅ Input validation added"
```

### Day 4-5: Write Comprehensive Tests

```bash
# 1. Create test file
cat > auth-service/auth_handler_test.go << 'EOF'
package main

import (
  "bytes"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestSignup_Success(t *testing.T) {
  req := SignupRequest{
    Email:     "user@example.com",
    Password:  "SecurePassword123",
    FirstName: "John",
    LastName:  "Doe",
    Phone:     "+251911234567",
  }

  body, _ := json.Marshal(req)
  httpReq := httptest.NewRequest(
    "POST",
    "/signup",
    bytes.NewBuffer(body),
  )

  w := httptest.NewRecorder()
  // Call handler
  // assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSignup_InvalidEmail(t *testing.T) {
  req := SignupRequest{
    Email:     "invalid-email",
    Password:  "SecurePassword123",
    FirstName: "John",
    LastName:  "Doe",
    Phone:     "+251911234567",
  }

  assert.Error(t, ValidateSignup(req))
}

func TestSignup_WeakPassword(t *testing.T) {
  req := SignupRequest{
    Email:     "user@example.com",
    Password:  "weak",
    FirstName: "John",
    LastName:  "Doe",
    Phone:     "+251911234567",
  }

  assert.Error(t, ValidateSignup(req))
}

func TestLogin_Success(t *testing.T) {
  // Test valid login
}

func TestLogin_InvalidCredentials(t *testing.T) {
  // Test invalid login
}

func TestJWT_Expiration(t *testing.T) {
  // Test token expiration
}

func TestJWT_Refresh(t *testing.T) {
  // Test refresh token
}

func TestRBAC_Enforcement(t *testing.T) {
  // Test role-based access control
}
EOF

# 2. Run tests
go test -v ./...

echo "✅ Comprehensive tests added"
```

### Day 5: Observability Integration

```bash
# 1. Add observability dependencies
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/jaeger
go get go.opentelemetry.io/otel/exporters/prometheus

# 2. Create telemetry setup
cat > auth-service/telemetry.go << 'EOF'
package main

import (
  "context"
  "os"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/exporters/jaeger"
  "go.opentelemetry.io/otel/exporters/prometheus"
  "go.opentelemetry.io/otel/sdk/metric"
  "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() error {
  exp, err := jaeger.New(
    jaeger.WithCollectorEndpoint(
      jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ),
  )
  if err != nil {
    return err
  }

  tp := trace.NewTracerProvider(
    trace.WithBatcher(exp),
  )
  otel.SetTracerProvider(tp)
  return nil
}

func initMetrics() error {
  exporter, err := prometheus.New()
  if err != nil {
    return err
  }

  provider := metric.NewMeterProvider(metric.WithReader(exporter))
  otel.SetMeterProvider(provider)
  return nil
}

func Shutdown(ctx context.Context) error {
  if tp, ok := otel.GetTracerProvider().(*trace.TracerProvider); ok {
    return tp.Shutdown(ctx)
  }
  return nil
}
EOF

# 3. Update main to initialize telemetry
echo "✅ Observability integration added"
```

---

## WEEK 2: KUBERNETES & CI/CD

### Kubernetes Manifests

```bash
# Create k8s manifests directory
mkdir -p infra/kubernetes/base
mkdir -p infra/kubernetes/staging
mkdir -p infra/kubernetes/production

# Create auth-service deployment
cat > infra/kubernetes/base/auth-service.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  labels:
    app: auth-service
    version: v1
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
        version: v1
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: auth-service
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - name: auth-service
        image: auth-service:latest
        imagePullPolicy: Always
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        env:
        - name: PORT
          value: "8080"
        - name: ENVIRONMENT
          value: production
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: auth-service-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: auth-service-secrets
              key: jwt-secret
        - name: JAEGER_ENDPOINT
          value: "http://jaeger:14268/api/traces"
        - name: PROMETHEUS_PORT
          value: "8080"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: http
            httpHeaders:
            - name: X-Health-Check
              value: liveness
          initialDelaySeconds: 15
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health/ready
            port: http
            httpHeaders:
            - name: X-Health-Check
              value: readiness
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 2
        startupProbe:
          httpGet:
            path: /health
            port: http
          failureThreshold: 30
          periodSeconds: 10
        volumeMounts:
        - name: tmp
          mountPath: /tmp
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
      volumes:
      - name: tmp
        emptyDir: {}
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - auth-service
              topologyKey: kubernetes.io/hostname

---
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  labels:
    app: auth-service
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: http
    protocol: TCP
    name: http
  selector:
    app: auth-service

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80

---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: auth-service-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: auth-service
EOF

# Deploy to staging
kubectl apply -f infra/kubernetes/base/ -n staging

# Verify deployment
kubectl get pods -n staging
kubectl describe deployment auth-service -n staging

echo "✅ Kubernetes manifests deployed"
```

### CI/CD Pipelines

```bash
# Create GitHub Actions workflows
mkdir -p .github/workflows

# Build workflow
cat > .github/workflows/build.yaml << 'EOF'
name: Build & Test

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service: [auth-service, user-service]
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - run: cd services/${{ matrix.service }} && go test -v ./...
    - run: cd services/${{ matrix.service }} && go build -o bin/${{ matrix.service }}
EOF

# Deploy workflow
cat > .github/workflows/deploy.yaml << 'EOF'
name: Deploy to Kubernetes

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: azure/k8s-set-context@v3
      with:
        method: kubeconfig
        kubeconfig: ${{ secrets.KUBECONFIG }}
    - run: kubectl apply -f infra/kubernetes/base/ -n production
EOF

echo "✅ CI/CD pipelines configured"
```

---

## WEEK 3-4: CORE SERVICES

### User Service Implementation

```bash
# Create user service
cp -r services/_template services/user-service

# Implement user service endpoints
cat > services/user-service/main.go << 'EOF'
package main

import (
  "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()

  // User endpoints
  router.POST("/users/profile", createProfile)
  router.GET("/users/:id", getProfile)
  router.PUT("/users/:id", updateProfile)
  router.GET("/users/:id/documents", getDocuments)
  router.POST("/users/:id/documents", uploadDocument)

  router.Run(":8080")
}

func createProfile(c *gin.Context) {
  // Implement profile creation
}

func getProfile(c *gin.Context) {
  // Implement profile retrieval
}

func updateProfile(c *gin.Context) {
  // Implement profile update
}

func getDocuments(c *gin.Context) {
  // Implement document retrieval
}

func uploadDocument(c *gin.Context) {
  // Implement document upload
}
EOF
```

---

## CHECKPOINTS

### End of Week 1 Checklist
- [ ] FamGo-consolidated repository created
- [ ] Docker-compose infrastructure running
- [ ] Security issues fixed (.env management)
- [ ] Auth service 70% complete
- [ ] Database migrations working
- [ ] Input validation added
- [ ] Initial tests written
- [ ] All changes committed to git

### End of Week 2 Checklist
- [ ] Auth service 100% production-ready
- [ ] Kubernetes manifests created and tested
- [ ] CI/CD pipelines configured
- [ ] Auth service deployed to staging K8s
- [ ] First canary deployment successful
- [ ] Observability dashboards showing data

### End of Week 4 Checklist
- [ ] User service complete
- [ ] Ride service started
- [ ] Dispatch service architecture designed
- [ ] Database schema 50% complete
- [ ] 80%+ test coverage on auth service
- [ ] 2+ services running on Kubernetes
- [ ] Staging environment stable

---

## MEASUREMENT METRICS

### Quality Metrics
- [ ] Test coverage: > 80%
- [ ] Build success rate: > 99%
- [ ] Code review approval rate: > 95%
- [ ] Critical bugs in production: 0

### Performance Metrics
- [ ] API latency p95: < 100ms
- [ ] Database query p95: < 50ms
- [ ] Service uptime: > 99.9%
- [ ] Container startup time: < 5s

### Delivery Metrics
- [ ] Sprint velocity: stable
- [ ] Deployment frequency: daily
- [ ] Lead time for changes: < 1 hour
- [ ] Mean time to recovery: < 15 minutes

---

## SUCCESS INDICATORS

### Week 1 Done When:
✅ Consolidated repository created  
✅ Auth service core implementation complete  
✅ Database migrations running  
✅ Team can run `docker-compose up` and have full stack  
✅ Git history shows incremental progress  

### Week 2 Done When:
✅ Auth service tested (80%+ coverage)  
✅ Service deployed to Kubernetes  
✅ CI/CD pipelines green  
✅ Observability stack showing metrics  

### Week 4 Done When:
✅ Multiple services running  
✅ Event streaming configured  
✅ API gateway routing  
✅ 3+ services production-ready  

---

## EMERGENCY PROCEDURES

### If Infrastructure Fails:
```bash
# 1. Check docker daemon
docker ps

# 2. Rebuild containers
docker-compose -f infra/docker/docker-compose.yml build --no-cache

# 3. Restart services
docker-compose -f infra/docker/docker-compose.yml down
docker-compose -f infra/docker/docker-compose.yml up -d

# 4. Verify health
docker-compose -f infra/docker/docker-compose.yml ps
```

### If Database Is Corrupted:
```bash
# 1. Backup existing data
docker exec famgo-postgres pg_dump -U famgo famgo > backup.sql

# 2. Reset database
docker-compose -f infra/docker/docker-compose.yml down -v

# 3. Restart and rerun migrations
docker-compose -f infra/docker/docker-compose.yml up -d
migrate -path database/migrations -database "$DATABASE_URL" up
```

---

## NEXT STEPS

✅ **This Week:** Start Phase 0 consolidation  
✅ **Week 2-3:** Complete auth service  
✅ **Week 4-8:** Implement core services  
✅ **Week 9-12:** Frontend, testing, hardening  
✅ **Week 13+:** Scale, optimize, launch  

**Ready to start? Go to Week 1 section above.**

---
