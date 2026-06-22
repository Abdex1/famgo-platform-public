# WEEK 2 IMPLEMENTATION ROADMAP — REALIGNMENT EXECUTION

**Status:** Week 2 Ready to Execute  
**Date:** 2025-01-15  
**Duration:** 7 days  
**Objectives:** Create correct templates + infrastructure for 100% specification compliance

---

## DAILY BREAKDOWN

### DAY 1-2: Go Service Template

**Deliverable:** Production-grade Go microservice template for 10 core services

**Location:** `services/_template-go/`

**Structure:**
```
_template-go/
├── cmd/service/main.go
├── internal/domain/service.go
├── internal/infrastructure/postgres/repository.go
├── internal/infrastructure/redis/cache.go
├── internal/infrastructure/kafka/consumer.go
├── internal/handlers/grpc.go
├── internal/handlers/rest.go
├── api/proto/service.proto
├── migrations/001_init.sql
├── tests/unit/service_test.go
├── Dockerfile
├── go.mod
├── Makefile
└── README.md
```

**Files Created Today:**
- ✅ cmd/service/main.go (gRPC server setup)
- ✅ internal/domain/service.go (business logic interface)
- ✅ Makefile (14 Go commands)
- ✅ go.mod (all dependencies)
- ✅ README.md (comprehensive guide)

**Testing:**
```bash
cd services/_template-go
go mod download
go run cmd/service/main.go
# Should start gRPC server on :5001
```

**Status:** ✅ 70% Complete (need infrastructure files)

---

### DAY 2-3: Python FastAPI Template

**Deliverable:** Production-grade FastAPI service template for 5 ML services

**Location:** `services/_template-python/`

**Structure:**
```
_template-python/
├── app/main.py
├── app/models/
├── app/services/
├── app/routes/
├── ml/
│   ├── pipelines/
│   ├── models/
│   └── preprocessing.py
├── migrations/
├── tests/unit/
├── tests/integration/
├── Dockerfile
├── requirements.txt
├── Makefile
└── README.md
```

**Files Created Today:**
- ✅ app/main.py (FastAPI server + Kafka consumer)
- ✅ requirements.txt (all Python dependencies)

**Features:**
- FastAPI with async/await
- PostgreSQL connection pooling
- Redis caching
- Kafka producer + consumer
- Health checks + readiness checks
- OpenTelemetry tracing
- Background tasks

**Testing:**
```bash
cd services/_template-python
pip install -r requirements.txt
python -m app.main
# Should start FastAPI server on :8080
```

**Status:** ✅ 70% Complete

---

### DAY 3-4: Kong API Gateway Setup

**Deliverable:** Kong configuration replacing NestJS as API gateway

**Location:** `infra/kong/`

**Structure:**
```
infra/kong/
├── kong.yml (service + route definitions)
├── docker-compose.yml
├── plugins/
│   ├── jwt.yml
│   ├── rate-limit.yml
│   └── cors.yml
└── README.md
```

**Kong Configuration (kong.yml):**
```yaml
_format_version: "3.0"
_transform: true

services:
  - name: auth-service
    protocol: grpc
    host: auth-service
    port: 5001
    
  - name: ride-service
    protocol: grpc
    host: ride-service
    port: 5001
    
  - name: gps-service
    protocol: grpc
    host: gps-service
    port: 5001
    
  # ... all 18 services

routes:
  - name: auth-login
    service: auth-service
    protocols: [http, https]
    methods: [POST]
    paths: [/auth/login]
    
  - name: ride-create
    service: ride-service
    protocols: [http, https]
    methods: [POST]
    paths: [/rides]

plugins:
  - name: jwt
    service: auth-service
    config:
      key_claim_name: sub
      secret_is_base64: false
      
  - name: rate-limiting
    config:
      minute: 1000
      policy: local
```

**Docker Compose:**
```yaml
version: '3'

services:
  kong-database:
    image: postgres:15
    environment:
      POSTGRES_DB: kong
      POSTGRES_USER: kong
      POSTGRES_PASSWORD: kong
    ports:
      - "5432:5432"
    
  kong:
    image: kong:3.4
    environment:
      KONG_DATABASE: postgres
      KONG_PG_HOST: kong-database
      KONG_PG_USER: kong
      KONG_PG_PASSWORD: kong
      KONG_PROXY_ACCESS_LOG: /dev/stdout
      KONG_ADMIN_ACCESS_LOG: /dev/stdout
    ports:
      - "8000:8000"   # Proxy
      - "8443:8443"   # Proxy SSL
      - "8001:8001"   # Admin
    depends_on:
      - kong-database
    
  kong-admin-ui:
    image: pgbi/kong-dashboard:latest
    ports:
      - "8002:8080"
```

**Start Kong:**
```bash
cd infra/kong
docker-compose up

# Admin UI: http://localhost:8002
# Proxy: http://localhost:8000
```

**Status:** ⏳ Not Started

---

### DAY 5: Kafka Topics Creation

**Deliverable:** All 15 Kafka topics with proper configuration

**Location:** `infra/kafka/`

**Topics (Per Specification):**
```
ride.created
ride.matching.started
ride.driver.assigned
ride.started
ride.completed
ride.cancelled

driver.location.updated

pool.created
pool.updated

pricing.calculated

payment.completed
payment.failed

wallet.transaction.created

safety.sos.triggered

fraud.detected

notification.send
```

**Create Topics Script (infra/kafka/setup-topics.sh):**
```bash
#!/bin/bash

KAFKA_BROKERS="localhost:9092"

# Function to create topic
create_topic() {
    local topic=$1
    echo "Creating topic: $topic"
    kafka-topics.sh --create \
        --bootstrap-server $KAFKA_BROKERS \
        --topic $topic \
        --partitions 3 \
        --replication-factor 1 \
        --config retention.ms=604800000 \
        --if-not-exists
}

# Create all topics
create_topic "ride.created"
create_topic "ride.matching.started"
create_topic "ride.driver.assigned"
create_topic "ride.started"
create_topic "ride.completed"
create_topic "ride.cancelled"
create_topic "driver.location.updated"
create_topic "pool.created"
create_topic "pool.updated"
create_topic "pricing.calculated"
create_topic "payment.completed"
create_topic "payment.failed"
create_topic "wallet.transaction.created"
create_topic "safety.sos.triggered"
create_topic "fraud.detected"
create_topic "notification.send"

echo "All topics created successfully"
```

**Consumer Groups (infra/kafka/setup-consumers.sh):**
```bash
#!/bin/bash

# Create consumer groups for each service
for service in auth ride dispatch pooling gps payment wallet safety fraud pricing; do
    kafka-consumer-groups.sh --bootstrap-server localhost:9092 \
        --create \
        --group ${service}-service-group \
        --topic 'ride.*' --topic 'pool.*' --topic 'payment.*' \
        --topic 'driver.*' --topic 'safety.*' --topic 'fraud.*'
done
```

**Status:** ⏳ Not Started

---

### DAYS 5-7: Service Boundaries Documentation

**Deliverable:** Explicit boundaries for all 18 services

**Location:** `docs/service-boundaries.md`

**Template (for each service):**
```markdown
# [Service Name] Service Boundary

## Responsibilities (ONLY):
- Item 1
- Item 2
- Item 3

## NOT Responsible For:
- Other service responsibility
- Another service responsibility

## Database:
- Tables owned: table_a, table_b
- Tables accessed: (none - use gRPC)

## External APIs:
- gRPC: ServiceNameService (internal)
- REST: /service-path/* (external)

## Events Published:
- event.type.v1
- event.type.v2

## Events Consumed:
- other.service.event

## Dependencies (Services):
- auth-service (for JWT verification)
- user-service (for user profiles)

## Performance SLAs:
- Latency: <100ms (p99)
- Availability: 99.9%
- Throughput: 10,000 req/s
```

**Services to Document (All 18):**
1. auth-service
2. user-service
3. driver-service
4. ride-service
5. dispatch-service
6. pooling-service
7. pricing-service
8. gps-service
9. payment-service
10. wallet-service
11. notification-service
12. safety-service
13. fraud-service
14. analytics-service
15. subscription-service
16. smart-pickup-service
17. voice-booking-service
18. websocket-gateway

**Status:** ⏳ Not Started

---

## COMPLETION CHECKLIST (Week 2)

### Go Service Template
- [ ] cmd/service/main.go - gRPC server bootstrap
- [ ] internal/domain/service.go - business logic interface
- [ ] internal/infrastructure/ - postgres, redis, kafka
- [ ] internal/handlers/ - grpc, rest handlers
- [ ] api/proto/service.proto - gRPC definitions
- [ ] migrations/001_init.sql - database schema
- [ ] tests/unit/ - unit test examples
- [ ] Dockerfile - multi-stage build
- [ ] go.mod - all dependencies
- [ ] Makefile - 14 commands
- [ ] README.md - comprehensive guide
- [ ] Local testing - verify builds and runs

### Python FastAPI Template
- [ ] app/main.py - FastAPI server + Kafka
- [ ] app/models/ - data models
- [ ] app/services/ - business logic
- [ ] app/routes/ - endpoint definitions
- [ ] ml/ - ML pipeline structure
- [ ] migrations/ - alembic setup
- [ ] tests/ - unit + integration tests
- [ ] Dockerfile - multi-stage build
- [ ] requirements.txt - dependencies
- [ ] Makefile - build commands
- [ ] README.md - ML service guide
- [ ] Local testing - verify runs on port 8080

### Kong API Gateway
- [ ] kong.yml - service definitions (all 18 services)
- [ ] docker-compose.yml - Kong infrastructure
- [ ] JWT plugin configuration
- [ ] Rate limiting configuration
- [ ] CORS configuration
- [ ] Admin UI running on :8002
- [ ] Proxy running on :8000
- [ ] Service routing verified

### Kafka Topics
- [ ] All 15 topics created
- [ ] Consumer groups configured
- [ ] Retention policy set (7 days)
- [ ] Partition strategy defined (3 partitions)
- [ ] Schema validation ready

### Service Boundaries Documentation
- [ ] All 18 services documented
- [ ] Responsibilities explicitly defined
- [ ] No overlap identified
- [ ] Event flows mapped
- [ ] Dependency graph created
- [ ] Performance SLAs defined

---

## SUCCESS CRITERIA

**Week 2 is COMPLETE when:**

1. ✅ Go template can create new service in <5 minutes
2. ✅ Python template can create new ML service in <5 minutes
3. ✅ Kong routes all requests to correct backends
4. ✅ All 15 Kafka topics exist and accept messages
5. ✅ Service boundaries clearly defined (no ambiguity)
6. ✅ 100% specification compliance achieved
7. ✅ All templates pass local testing
8. ✅ Documentation complete and accurate

**Timeline to Week 3:**
- Monday morning: Review Week 2 deliverables
- Tuesday: Begin Auth Service (Go) with correct template
- Wednesday+: Roll out other services

---

## RISK MITIGATION

### Risk 1: Go Module Versions Conflict
**Mitigation:** Test go.mod with `go mod tidy` before committing

### Risk 2: Python Dependencies Conflict
**Mitigation:** Use Python 3.11+, test with `pip install -r requirements.txt`

### Risk 3: Kong Configuration Issues
**Mitigation:** Validate yaml with `kong config --syntax-check`

### Risk 4: Kafka Topic Creation Fails
**Mitigation:** Have backup kafka-console scripts ready

### Risk 5: Service Boundaries Unclear
**Mitigation:** Use event flow diagrams to validate

---

**This Week 2 realignment ensures 100% specification compliance and enterprise-grade architecture.**

**Ready to execute! 🚀**

---

*FamGo Platform — Week 2 Realignment Roadmap*  
*Execution Template v1.0*  
*2025-01-15*
