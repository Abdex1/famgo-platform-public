# PHASE 1 EXECUTION GUIDE - STEP BY STEP

## Prerequisites
- Docker Desktop installed
- Docker Compose installed
- Go 1.21+
- PostgreSQL client tools (psql)
- Redis CLI
- Kafka CLI tools

---

## STEP 1: Start Docker Infrastructure (5 minutes)

```bash
# Navigate to project
cd C:\dev\FamGo-platform

# Start all services (Docker Compose)
docker-compose -f infra/docker/docker-compose.yml up -d

# Verify all services are running
docker ps
# Should show 15 services

# Check service health
docker-compose -f infra/docker/docker-compose.yml ps
```

**Expected Output:**
- PostgreSQL: healthy
- Redis: healthy
- Kafka: healthy
- Kong: healthy
- Grafana: running
- Jaeger: running

---

## STEP 2: Initialize PostgreSQL Database (3 minutes)

```bash
# Connect to PostgreSQL
psql -h localhost -U famgo -d famgo

# Run migration 1 (in psql prompt)
\i database/migrations/001_initial_schema.sql

# Run migration 2
\i database/migrations/002_advanced_indexes_procedures.sql

# Verify tables
\dt

# Exit psql
\q
```
 
**Expected:**
- 10 tables created
- Extensions enabled (PostGIS, pgvector)
- Indexes created
- Functions/procedures available

---

## STEP 3: Setup Redis (2 minutes)

```bash
# Run Redis setup script
bash infra/docker/scripts/setup_redis.sh

# Verify Redis GEO index
redis-cli
GEOADD drivers:geo 9.0320 8.9868 "driver:test-1"
GEOPOS drivers:geo driver:test-1
exit
```

**Expected:**
- GEO index initialized
- Session storage ready
- Rate limiting keys set

---

## STEP 4: Setup Kafka Topics (3 minutes)

```bash
# Connect to Kafka container
docker exec -it <kafka-container-id> bash

# Create topics
kafka-topics --create --topic auth.user.registered --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
kafka-topics --create --topic ride.created --bootstrap-server localhost:9092 --partitions 5 --replication-factor 1
kafka-topics --create --topic driver.location.updated --bootstrap-server localhost:9092 --partitions 10 --replication-factor 1
kafka-topics --create --topic payment.completed --bootstrap-server localhost:9092 --partitions 5 --replication-factor 1

# List topics
kafka-topics --list --bootstrap-server localhost:9092

# Exit container
exit
```

**Expected:**
- 15+ topics created
- Topics have correct partition counts
- Replication factor set

---

## STEP 5: Configure Kong Gateway (2 minutes)

```bash
# Check Kong status
curl -s http://localhost:8001/status | jq

# Add auth-service upstream
curl -X POST http://localhost:8001/upstreams \
  -d "name=auth-service" \
  -d "algorithm=round_robin" \
  -d "slots=10"

# Add route
curl -X POST http://localhost:8001/services/auth-service/routes \
  -d "paths[]=/v1/auth" \
  -d "name=auth-routes"

# Verify configuration
curl -s http://localhost:8001/services | jq
```

**Expected:**
- Services registered in Kong
- Routes configured
- Gateway responding to requests

---

## STEP 6: Build & Start Auth Service (5 minutes)

```bash
# Navigate to auth service
cd services/auth-service

# Install Go dependencies
go mod download

# Build binary
go build -o bin/auth-service cmd/api/main.go

# Start service
./bin/auth-service

# In another terminal, verify it's running
curl http://localhost:3000/v1/health
```

**Expected:**
- Service starts without errors
- Health endpoint responds
- Database connection successful

---

## STEP 7: Test Auth Service Endpoints (3 minutes)

```bash
# Test registration
curl -X POST http://localhost:3000/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "phone": "+251911234567",
    "password": "Test@1234",
    "password_confirm": "Test@1234",
    "first_name": "John",
    "last_name": "Doe",
    "role": "rider"
  }'

# Test login
curl -X POST http://localhost:3000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test@1234"
  }'

# Test health
curl http://localhost:3000/v1/health
```

**Expected:**
- Registration: 201 Created
- Login: 200 OK with tokens
- Health: 200 OK

---

## STEP 8: Run Integration Tests (2 minutes)

```bash
# Navigate back to project root
cd C:\dev\FamGo-platform

# Run integration tests
go test ./services/auth-service/internal/tests/...

# Or run specific test
go test -v ./services/auth-service/internal/tests/ -run TestDatabaseConnectivity
```

**Expected:**
- All tests pass (6/6)
- Database connectivity: PASS
- Redis connectivity: PASS
- Kong Gateway: PASS
- Auth Service: PASS

---

## STEP 9: Verify Data Flow (5 minutes)

```bash
# 1. Check PostgreSQL data
psql -h localhost -U famgo -d famgo -c "SELECT COUNT(*) FROM users;"

# 2. Check Redis keys
redis-cli KEYS "*"

# 3. Verify Kafka messages
docker exec -it <kafka-container-id> kafka-console-consumer --bootstrap-server localhost:9092 --topic auth.user.registered --from-beginning

# 4. Check Kong access
curl -v http://localhost:8000/v1/auth/me

# 5. Monitor in Grafana
# Open http://localhost:3000
# Login: admin/admin
# Check Prometheus datasource
```

**Expected:**
- Data persisted in PostgreSQL
- Keys stored in Redis
- Events flowing through Kafka
- Kong routing working
- Grafana dashboards accessible

---

## STEP 10: Health Check Dashboard (1 minute)

```bash
# Access monitoring
echo "Grafana:        http://localhost:3000"
echo "Jaeger:         http://localhost:16686"
echo "Kong Admin:     http://localhost:8001/status"
echo "Prometheus:     http://localhost:9090"
echo "Auth Service:   http://localhost:3000/v1/health"
```

---

## Troubleshooting

### PostgreSQL Connection Failed
```bash
# Check if postgres is running
docker ps | grep postgres

# View postgres logs
docker logs <postgres-container-id>

# Verify credentials in docker-compose.yml
```

### Redis Connection Failed
```bash
# Check redis
redis-cli ping

# View redis logs
docker logs <redis-container-id>
```

### Kafka Topics Not Found
```bash
# List topics
kafka-topics --list --bootstrap-server localhost:9092

# Create manually if needed
kafka-topics --create --topic auth.user.registered --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
```

### Auth Service Won't Start
```bash
# Check Go version
go version  # Should be 1.21+

# Check port 3000
netstat -an | grep 3000

# View service logs
tail -f /path/to/logs
```

---

## Success Criteria ✓

- [x] All 15 Docker services running
- [x] PostgreSQL tables created (10 tables)
- [x] Redis GEO index initialized
- [x] Kafka topics created (15 topics)
- [x] Kong gateway configured
- [x] Auth service running on port 3000
- [x] Auth endpoints responding
- [x] Integration tests passing
- [x] Data flowing through system
- [x] Monitoring dashboards accessible

---

## What's Next (Phase 2)

After Phase 1 is complete:

1. **User Service** (same structure as Auth Service)
2. **Driver Service** (with vehicle management)
3. **Notification Service** (SMS/Push)
4. **More microservices** (Ride, Dispatch, Payment, etc.)

Each service follows the same pattern:
- Go skeleton with `go.mod`
- Domain entities (GORM models)
- Repositories (database layer)
- Services (business logic)
- HTTP handlers (REST API)
- Kafka producers/consumers

---

## Common Commands

```bash
# View logs
docker logs -f <service-name>

# Enter service
docker exec -it <service-name> bash

# Stop all services
docker-compose -f infra/docker/docker-compose.yml down

# Clean up volumes
docker-compose -f infra/docker/docker-compose.yml down -v

# Rebuild images
docker-compose -f infra/docker/docker-compose.yml build --no-cache

# View Docker stats
docker stats
```

---

**Phase 1 Execution Time: ~30-45 minutes total**

**Status**: Ready for Phase 1 execution
