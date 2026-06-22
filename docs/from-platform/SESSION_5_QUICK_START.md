# SESSION 5 QUICK START - Dispatch Service Build

## ⏱️ TARGET: 3-4 Hours (18 files)

## 🎯 WHAT YOU'RE BUILDING

**Dispatch Service**: Matches riders with drivers using multi-factor scoring algorithm.

**Core Function**: When a rider requests a ride, find the best available drivers based on:
- Proximity (40% weight) - Closer is better
- Acceptance Rate (30% weight) - More reliable drivers first
- Rating (20% weight) - Higher rated drivers
- Availability (10% weight) - Currently online status

---

## 📋 FILES TO CREATE (18 total)

### 1. Configuration & Setup (2 files) - 15 minutes
```
dispatch-service/
├── go.mod                          [Dependencies]
└── internal/config/config.go       [50+ dispatch parameters]
```

### 2. Domain Layer (3 files) - 45 minutes
```
domain/
├── entities/dispatch_request.go    [Matching state machine]
├── valueobjects/match_score.go    [Scoring calculations]
└── services/matching_service.go   [Multi-factor algorithm]
```

### 3. Infrastructure (2 files) - 30 minutes
```
infrastructure/repositories/
├── dispatch_repository.go         [PostgreSQL CRUD]
└── matching_repository.go         [Scoring queries]
```

### 4. Application (1 file) - 30 minutes
```
application/usecases/
└── dispatch_usecases.go           [5 use cases: MatchRide, GetMatches, Accept, Reject, GetStats]
```

### 5. Interface (2 files) - 30 minutes
```
interfaces/grpc/
├── dispatch.proto                 [6 gRPC endpoints]
└── dispatch_handler.go            [gRPC implementation]
```

### 6. Bootstrap & Deployment (3 files) - 30 minutes
```
├── cmd/main.go                    [Bootstrap, DI, server]
├── Dockerfile                     [Multi-stage build]
└── .env.example                   [Configuration template]
```

### 7. Tests (2 files) - 30 minutes
```
├── domain/services/matching_service_test.go   [Scoring algorithm tests]
└── domain/entities/dispatch_request_test.go   [State machine tests]
```

### 8. Additional (2 files)
```
├── internal/domain/services/kafka_publisher.go [Event publishing]
└── README.md                      [Service documentation]
```

---

## 🔑 KEY COMPONENTS EXPLAINED

### MatchingService - Multi-Factor Scoring Algorithm

```go
// Pseudocode for scoring
for each driver in nearby_drivers {
    proximity_score = (1 - distance/max_distance) * 100 * 0.40
    acceptance_score = acceptance_rate * 100 * 0.30
    rating_score = (rating / 5.0) * 100 * 0.20
    availability_score = (is_online ? 100 : 0) * 0.10
    
    total_score = proximity_score + acceptance_score + rating_score + availability_score
    
    matches = append(matches, {driver, total_score})
}

sort matches by total_score DESC
return top N matches
```

### DispatchRequest Entity - State Machine

```
States: PENDING → MATCHED → ACCEPTED → COMPLETED/CANCELLED

Methods needed:
- CreateRequest(riderId, location)
- FindMatches(gpsService, rideService) -> []Driver
- AcceptMatch(driverId)
- RejectMatch(driverId, reason)
- Complete()
- Cancel()
```

### Key Inputs (from other services)
- GPS Service: `FindNearbyDrivers(latitude, longitude, radius)`
- Ride Service: `GetRideDetails(rideId)`

### Key Outputs (Kafka Events)
- `dispatch.ride.matched` - When driver matched
- `dispatch.match.accepted` - When driver accepts
- `dispatch.match.rejected` - When driver rejects
- `dispatch.match.expired` - When match times out

---

## 📊 DEPENDENCIES & DATA FLOW

```
┌─────────────┐
│   Rider     │ Creates ride request
└──────┬──────┘
       │ RideService: CreateRide
       ▼
┌─────────────────────┐
│ Ride Service        │ Creates ride entity
└──────┬──────────────┘
       │ Publishes: ride.requested
       ▼
┌──────────────────────────┐
│ Dispatch Service (NEW!)  │ Matches rider with drivers
│ 1. Query GPS service     │ - GetNearby drivers
│ 2. Score drivers         │ - Multi-factor algorithm
│ 3. Propose matches       │ - Top 3-5 drivers
│ 4. Handle acceptance     │ - Driver accepts/rejects
└──────┬───────────────────┘
       │ Publishes: dispatch.ride.matched
       ▼
┌──────────────────┐
│ Ride Updated     │ With driver_id
└──────────────────┘
       │
       ▼
Payment Service (next session)
```

---

## 💡 IMPLEMENTATION TIPS

### 1. Scoring Algorithm
- **Keep it fast**: Sub-second responses required
- **Configurable weights**: Make 40/30/20/10 adjustable via config
- **Handle edge cases**: No online drivers, all drivers too far away
- **Tie-breaking**: Use rating when scores equal

### 2. State Machine
- **Valid transitions only**: PENDING → MATCHED, MATCHED → ACCEPTED/REJECTED, etc.
- **Timestamp tracking**: When each state changed
- **Expiration handling**: Matches expire after 60 seconds

### 3. Database Queries
- **Index on location**: For spatial queries
- **Index on driver_id + status**: For matching queries
- **Cached results**: Use Redis for scoring to avoid DB hits

### 4. gRPC Endpoints
```protobuf
service DispatchService {
  rpc MatchRide(MatchRideRequest) returns (MatchRideResponse);           // Find drivers
  rpc GetMatchedDrivers(GetMatchesRequest) returns (GetMatchesResponse); // List options
  rpc AcceptMatch(AcceptMatchRequest) returns (Empty);                   // Driver accepts
  rpc RejectMatch(RejectMatchRequest) returns (Empty);                   // Driver rejects
  rpc CancelRequest(CancelRequestRequest) returns (Empty);               // Rider cancels
  rpc GetDispatchStats(Empty) returns (DispatchStatsResponse);           // Metrics
}
```

---

## 🧪 TEST CASES TO CREATE

### Scoring Algorithm Tests
- [ ] Proximity scoring (closer driver scores higher)
- [ ] Acceptance rate effect (higher % scores higher)
- [ ] Rating effect (5-star scores higher)
- [ ] Availability effect (online scores higher)
- [ ] Weighted combination (all factors together)
- [ ] Edge case: No drivers available
- [ ] Edge case: All drivers equally good

### State Machine Tests
- [ ] PENDING → MATCHED transition
- [ ] MATCHED → ACCEPTED transition
- [ ] MATCHED → REJECTED transition
- [ ] ACCEPTED → COMPLETED transition
- [ ] Rejection resets to PENDING
- [ ] Expiration handling

---

## 🚀 QUICK BUILD CHECKLIST

- [ ] Create go.mod (copy from ride-service, update module name)
- [ ] Create config.go (use ride-service as template)
- [ ] Create DispatchRequest entity with state machine
- [ ] Create MatchingService with scoring algorithm
- [ ] Create repositories (dispatch_request + matching queries)
- [ ] Create use cases (5 orchestration flows)
- [ ] Create proto definitions (6 endpoints)
- [ ] Create gRPC handler (implement all endpoints)
- [ ] Create bootstrap (main.go with DI)
- [ ] Create Dockerfile (multi-stage, same pattern)
- [ ] Create test files (state machine + scoring algorithm)
- [ ] Test build: `docker build -t famgo/dispatch-service:latest .`
- [ ] Verify integration with GPS + Ride services

---

## 📖 REFERENCE PATTERNS

### From GPS Service (Location-based)
- Use Redis for caching driver lists
- Batch operations for performance
- Use gRPC for service-to-service calls

### From Ride Service (Entity with state machine)
- Similar entity lifecycle management
- Use same repository patterns
- Similar gRPC handler structure

### New Pattern - Multi-factor Scoring
- Weighted algorithm with configurable weights
- Service composition (GPS + Ride)
- Complex business logic in domain service

---

## ⏰ TIME ALLOCATION

| Task | Time | Status |
|------|------|--------|
| Setup (go.mod, config) | 15m | ⏳ |
| Domain Layer | 45m | ⏳ |
| Infrastructure | 30m | ⏳ |
| Application (Use Cases) | 30m | ⏳ |
| Interface (gRPC) | 30m | ⏳ |
| Bootstrap + Docker | 30m | ⏳ |
| Tests | 30m | ⏳ |
| Integration Testing | 30m | ⏳ |
| **TOTAL** | **3.5 hours** | ⏳ |

---

## ✅ SUCCESS CRITERIA

- [ ] All 18 files created
- [ ] Docker image builds successfully
- [ ] gRPC endpoints accessible
- [ ] Scoring algorithm tested
- [ ] State machine working
- [ ] Integrates with GPS service (calls GetNearbyDrivers)
- [ ] Integrates with Ride service (calls GetRideDetails)
- [ ] Publishes Kafka events
- [ ] 80%+ test coverage
- [ ] All tests passing

---

## 🔗 USEFUL COMMANDS

```bash
# Generate gRPC code from proto
protoc --go_out=. --go-grpc_out=. proto/dispatch.proto

# Build Docker image
docker build -t famgo/dispatch-service:latest .

# Run tests
go test ./...

# Run with debugging
go run cmd/main.go --debug

# Check code coverage
go test -cover ./...
```

---

## 📝 TEMPLATE REMINDERS

✓ Use same 7-layer DDD pattern (proven in GPS + Ride)
✓ Use same repository pattern (PostgreSQL + scanning)
✓ Use same use case pattern (input/output DTOs)
✓ Use same gRPC handler pattern (request mapping)
✓ Use same bootstrap pattern (DI, pooling, graceful shutdown)
✓ Use same test patterns (unit + integration)

**No reinventing - copy patterns from GPS/Ride, just change the business logic!**

---

## 🎯 WHEN YOU'RE DONE

1. Update docker-compose.yml to include dispatch-service
2. Test all three services together (auth, gps, ride, dispatch)
3. Verify Kafka events flow
4. Create SESSION_5_DISPATCH_DELIVERY.md (same format as GPS/Ride)

---

**Ready? Let's build the Dispatch Service!** 🚀
