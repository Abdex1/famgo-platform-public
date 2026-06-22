# PHASE 2 - EXECUTION GUIDE: USER & DRIVER SERVICES

**Status**: ✅ Ready to execute
**Duration**: 2-3 weeks
**Services to build**: 3 (User, Driver, Notification)
**Database tables to create**: 8+ new tables
**API endpoints**: 25+ total

---

## WEEKLY BREAKDOWN

### Week 1: Foundation & User Service

#### Day 1-2: Database Setup
```bash
# Create migration files
database/migrations/003_user_service_tables.sql
database/migrations/004_driver_service_tables.sql
database/migrations/005_notification_service_tables.sql

# Run migrations
psql -h localhost -U famgo -d famgo < database/migrations/003_user_service_tables.sql
psql -h localhost -U famgo -d famgo < database/migrations/004_user_service_tables.sql
psql -h localhost -U famgo -d famgo < database/migrations/005_notification_service_tables.sql
```

#### Day 3-4: User Service Development
```bash
# Setup structure
mkdir -p services/user-service/{cmd/api,internal/{domain/entities,infrastructure/postgres,interfaces/rest/{handlers,routes}}}

# Build
cd services/user-service
go mod download
go build -o bin/user-service cmd/api/main.go

# Test
./bin/user-service
# Verify health: curl http://localhost:3001/v1/health
```

#### Day 5: User Service Testing
```bash
# Test endpoints
curl -X GET http://localhost:3001/v1/users/{id}/profile
curl -X PUT http://localhost:3001/v1/users/{id}/profile -d '{...}'
curl -X GET http://localhost:3001/v1/users/{id}/preferences
```

---

### Week 2: Driver Service & Event Bus

#### Day 1-2: Driver Service Development
```bash
# Setup structure
mkdir -p services/driver-service/{cmd/api,internal/{domain/entities,infrastructure/postgres,interfaces/rest/{handlers,routes}}}

# Build
cd services/driver-service
go mod download
go build -o bin/driver-service cmd/api/main.go

# Test
./bin/driver-service
# Verify health: curl http://localhost:3002/v1/health
```

#### Day 3: Event Bus Implementation
```bash
# Create shared event bus
shared/event-bus/
  ├── producer.go  (13 event types)
  ├── consumer.go  (handlers for each service)
  └── models.go    (event structures)

# Integration test
# Publish event from user service
# Verify notification service consumed it
```

#### Day 4-5: Notification Service
```bash
# Setup structure
mkdir -p services/notification-service/{cmd/api,internal/{domain/entities,infrastructure/{twilio,firebase}}}

# Build
cd services/notification-service
go mod download
go build -o bin/notification-service cmd/api/main.go

# Mock providers for testing
# Verify SMS/Push event consumption
```

---

### Week 3: Integration & Hardening

#### Day 1-2: End-to-End Testing
```bash
# Test complete flow:
1. User registers (Auth Service)
2. User event published
3. User Service creates profile
4. Notification Service sends welcome SMS
5. All data persisted correctly
```

#### Day 3: Performance Tuning
```bash
# Database indexing verification
# Query performance testing
# Load testing with multiple concurrent requests
# Memory profiling
```

#### Day 4-5: Documentation & Deployment
```bash
# Unit tests for each service
# Integration tests for event flow
# Docker Compose updates for new services
# Deployment documentation
```

---

## DAILY TASKS CHECKLIST

### Task Template
```
Day X - [Service Name] Development
├── [ ] Code skeleton created
├── [ ] Entities implemented
├── [ ] Repositories implemented
├── [ ] HTTP handlers implemented
├── [ ] Routes registered
├── [ ] Service builds without errors
├── [ ] Service starts on correct port
├── [ ] Health endpoint responds
├── [ ] Database migrations run
├── [ ] Endpoints tested
├── [ ] Event publishing working
└── [ ] Documentation updated
```

---

## COMMAND REFERENCE

```bash
# Setup new service
mkdir -p services/{service-name}/{cmd/api,internal/{domain/entities,infrastructure/postgres,interfaces/rest/{handlers,routes}}}

# Initialize Go module
cd services/{service-name}
go mod init github.com/FamGo/platform/services/{service-name}

# Download dependencies
go mod download

# Build service
go build -o bin/{service-name} cmd/api/main.go

# Run service
./bin/{service-name}

# Test endpoint
curl -X GET http://localhost:port/v1/health

# View logs
docker logs {service-name}

# Database migration
psql -h localhost -U famgo -d famgo < database/migrations/XXX.sql

# Kafka topic creation
kafka-topics --create --topic {topic-name} --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1

# Redis setup
redis-cli < setup_redis.sh
```

---

## TESTING PROCEDURES

### Unit Tests
```go
// For each repository
func TestUserProfileRepositoryCreate(t *testing.T) {}
func TestUserProfileRepositoryGetByUserID(t *testing.T) {}
func TestUserProfileRepositoryUpdate(t *testing.T) {}

// For each handler
func TestGetProfileHandler(t *testing.T) {}
func TestUpdateProfileHandler(t *testing.T) {}
```

### Integration Tests
```bash
# Test complete flow
1. Start all services
2. Create user via Auth Service
3. Verify user event published
4. Check User Service received event
5. Verify profile created in database
6. Check notification queued
```

### Load Testing
```bash
# Using Apache Bench or k6
ab -n 1000 -c 10 http://localhost:3001/v1/health
```

---

## TROUBLESHOOTING

### User Service won't start
```bash
# Check port 3001 is free
lsof -i :3001

# Check database connection
psql -h localhost -U famgo -d famgo -c "SELECT 1"

# Check Kafka connectivity
kafka-topics --list --bootstrap-server localhost:9092
```

### Driver Service database errors
```bash
# Check tables exist
psql -h localhost -U famgo -d famgo -c "\dt driver*"

# Verify indexes
psql -h localhost -U famgo -d famgo -c "\di driver*"
```

### Event not publishing
```bash
# Check Kafka topics
kafka-topics --list --bootstrap-server localhost:9092

# Monitor Kafka messages
kafka-console-consumer --bootstrap-server localhost:9092 --topic driver.registered --from-beginning
```

### Notification not sending
```bash
# Check Twilio credentials
echo $TWILIO_ACCOUNT_SID

# Check Firebase config
ls -la /path/to/firebase-credentials.json

# Monitor notification queue
curl -X GET http://localhost:3003/v1/notifications/pending
```

---

## SUCCESS CRITERIA

✅ All services building without errors
✅ All services starting on correct ports
✅ All health endpoints responding
✅ Database migrations applied
✅ All endpoints returning 200 OK
✅ Event publishing working
✅ Event consuming working
✅ Notifications queued correctly
✅ User profiles persisting
✅ Driver profiles persisting
✅ All integration tests passing

---

## AFTER PHASE 2

Expected state:
- 3 new microservices running
- 25+ new API endpoints
- Event-driven communication established
- Notification infrastructure ready
- 8+ new database tables
- Complete user and driver workflows

Ready for Phase 3:
- Ride Service
- Dispatch Service
- GPS Service
- WebSocket Gateway

