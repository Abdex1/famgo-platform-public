# PHASE 4: POOLING SERVICE - PRODUCTION EXECUTION GUIDE

**Status**: ✅ COMPLETE DELIVERY  
**Duration**: 2 weeks  
**Services**: 1 (Pooling Service)  
**Endpoints**: 5 REST  
**Database Tables**: 5 new tables  
**Code Files**: 6 Go files (~25 KB)

---

## 🎯 PHASE 4 OVERVIEW

### What Phase 4 Delivers
✅ **Pooling Service** (Port 3013) - Intelligent ride pooling  
✅ **Pool Matching Algorithm** - Route compatibility scoring  
✅ **Route Optimization** - Nearest neighbor algorithm (ready for TSP solver)  
✅ **Constraint Validation** - Detour/wait time enforcement  
✅ **Pool Metrics** - Performance tracking & analytics  
✅ **Database Schema** - 5 tables with triggers & views  

### Core Algorithm
```
Pool Compatibility Score = 
  (route_overlap * 0.4) +
  (profitability * 0.3) +
  (eta_similarity * 0.2) +
  (pickup_proximity * 0.1)

Constraints:
  - Max pool size: 3
  - Min route overlap: 70%
  - Max detour: 10 minutes
  - Max extra wait: 5 minutes
```

---

## 📁 FILES CREATED (Phase 4)

### Code Files
```
services/pooling-service/
├── go.mod                              (702 bytes)
├── cmd/api/main.go                     (1.8 KB)
├── internal/domain/entities/pool.go    (5.0 KB) - 5 entities
├── internal/domain/services/
│   └── pooling_engine.go               (9.7 KB) - Matching algorithm
├── internal/infrastructure/postgres/
│   └── pool_repository.go              (8.3 KB) - Database operations
└── internal/interfaces/rest/
    └── pooling_handler.go              (7.0 KB) - HTTP handlers (5 endpoints)
```

### Database
```
database/migrations/
└── 004_phase4_pooling_service.sql     (10.5 KB) - 5 tables + triggers + views
```

---

## 🔌 API ENDPOINTS (5 Total)

### POOLING SERVICE (Port 3013)

**1. Find Pool Matches** - POST `/v1/pooling/find-matches`
```json
Request:
{
  "ride_id": "ride_001",
  "driver_id": "driver_123",
  "pickup_lat": 9.0320,
  "pickup_lng": 38.7469,
  "dropoff_lat": 9.0456,
  "dropoff_lng": 38.7618,
  "pickup_address": "Addis Ababa",
  "dropoff_address": "Bole",
  "estimated_distance_meters": 5500,
  "estimated_duration_seconds": 900,
  "estimated_fare": 150.00,
  "female_only": false,
  "max_detour_minutes": 10,
  "max_wait_minutes": 5,
  "min_route_overlap": 0.70
}

Response (200 OK):
{
  "ride_id": "ride_001",
  "total": 3,
  "candidates": [
    {
      "ride_id": "ride_002",
      "pickup_lat": 9.0325,
      "pickup_lng": 38.7475,
      "dropoff_lat": 9.0460,
      "dropoff_lng": 38.7620,
      "route_overlap": 0.82,
      "detour_minutes": 5,
      "wait_minutes": 2,
      "compatibility_score": 0.78,
      "savings_percentage": 0.25,
      "rank": 1
    }
  ],
  "timestamp": "2026-06-10T10:30:00Z"
}
```

**2. Create Pool** - POST `/v1/pooling/pools`
```json
Request:
{
  "driver_id": "driver_123",
  "ride_ids": ["ride_001", "ride_002", "ride_003"],
  "max_size": 3
}

Response (201 Created):
{
  "pool_id": "pool_xyz789",
  "driver_id": "driver_123",
  "rides": ["ride_001", "ride_002", "ride_003"],
  "size": 3,
  "status": "FORMING",
  "created_at": "2026-06-10T10:30:15Z"
}
```

**3. Activate Pool** - POST `/v1/pooling/pools/{poolID}/activate`
```json
Response (200 OK):
{
  "pool_id": "pool_xyz789",
  "status": "ACTIVE",
  "rides": ["ride_001", "ride_002", "ride_003"]
}
```

**4. Complete Pool** - POST `/v1/pooling/pools/{poolID}/complete`
```json
Response (200 OK):
{
  "pool_id": "pool_xyz789",
  "status": "COMPLETED",
  "completed_at": "2026-06-10T10:45:30Z"
}
```

**5. Get Pool Statistics** - GET `/v1/pooling/statistics`
```json
Response (200 OK):
{
  "total_pools": 245,
  "active_pools": 12,
  "completed_pools": 220,
  "avg_pool_size": 2.4,
  "avg_compatibility_score": 0.72,
  "total_estimated_profit": 28500.00
}
```

**6. Health Check** - GET `/v1/health`
```json
Response (200 OK):
{
  "status": "healthy",
  "service": "pooling-service",
  "timestamp": "2026-06-10T10:30:00Z"
}
```

---

## 🗄️ DATABASE SCHEMA

### 5 NEW TABLES

**1. pool_groups** (Active ride pools)
- id (UUID PK)
- driver_id (FK)
- status (FORMING, ACTIVE, COMPLETED, CANCELLED)
- max_size (2-3)
- current_size
- ride_ids (JSON array)
- compatibility_score
- estimated_profit
- created_at, completed_at, updated_at

**2. pool_requests** (Rides eligible for pooling)
- id (UUID PK)
- ride_id (unique)
- pickup_lat, pickup_lng, dropoff_lat, dropoff_lng
- estimated_distance_meters, estimated_duration_seconds
- estimated_fare
- female_only, max_detour_minutes, max_wait_minutes, min_route_overlap
- status (PENDING, POOLED, SOLO, EXPIRED)
- expires_at

**3. pool_routes** (Optimized delivery routes)
- id (UUID PK)
- pool_id (FK)
- pickup_sequence (JSON)
- dropoff_sequence (JSON)
- total_distance_meters
- waypoints (GeoJSON)

**4. pool_compatibility_matrix** (Audit trail)
- Stores all compatibility scores between ride pairs
- Used for analytics and optimization

**5. pool_metrics** (Time-series metrics)
- Tracks pools created, activated, completed
- Records earnings and satisfaction scores

---

## 🔧 SETUP INSTRUCTIONS

### Prerequisites
```bash
go version      # 1.21+
psql --version  # 14+
```

### 1. Create Pooling Service Directory
```bash
cd C:\dev\FamGo-platform\services
mkdir -p pooling-service/{cmd/api,internal/{domain/entities,domain/services,infrastructure/postgres,interfaces/rest}}
```

### 2. Initialize Go Module
```bash
cd services/pooling-service
go mod init github.com/FamGo/platform/services/pooling-service
go mod download
```

### 3. Apply Database Migration
```bash
psql -U famgo_user -d famgo_platform
-- Apply Phase 4 migration
\i database/migrations/004_phase4_pooling_service.sql

-- Verify tables created
\dt pool_groups
\dt pool_requests
\dt pool_routes
\dv v_active_pools
```

### 4. Build Service
```bash
cd services/pooling-service
go build -o bin/pooling-service cmd/api/main.go
```

### 5. Run Service
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=famgo_user
export DB_PASSWORD=secure_password
export DB_NAME=famgo_platform
export SERVICE_PORT=3013

./bin/pooling-service
# Output: 🚀 Pooling Service starting on port 3013
```

---

## ✅ VERIFICATION CHECKLIST

### Service Health
- [ ] `curl http://localhost:3013/v1/health` returns 200
- [ ] Service logs show "Connected to PostgreSQL"

### Database
- [ ] All 5 tables created: `SELECT table_name FROM information_schema.tables WHERE table_schema='public'`
- [ ] Triggers created: `SELECT trigger_name FROM information_schema.triggers`
- [ ] Views accessible: `SELECT * FROM v_active_pools LIMIT 1`

### API Endpoints
- [ ] Find matches: `POST /v1/pooling/find-matches` with test data
- [ ] Create pool: `POST /v1/pooling/pools` with 2-3 ride IDs
- [ ] Activate pool: `POST /v1/pooling/pools/{poolID}/activate`
- [ ] Complete pool: `POST /v1/pooling/pools/{poolID}/complete`
- [ ] Statistics: `GET /v1/pooling/statistics`

### Algorithm
- [ ] Route overlap calculation returns 0-1 score
- [ ] Compatibility calculation weighs all 4 factors
- [ ] Constraint validation rejects invalid combinations
- [ ] Pool creation transitions through states correctly

---

## 🧪 INTEGRATION TEST FLOW

### Complete Pooling Workflow

```bash
# 1. Request pool for 2+ rides
curl -X POST http://localhost:3013/v1/pooling/find-matches \
  -H "Content-Type: application/json" \
  -d '{
    "ride_id": "ride_001",
    "driver_id": "driver_001",
    "pickup_lat": 9.0320,
    "pickup_lng": 38.7469,
    "dropoff_lat": 9.0456,
    "dropoff_lng": 38.7618,
    "pickup_address": "Addis Ababa",
    "dropoff_address": "Bole",
    "estimated_distance_meters": 5500,
    "estimated_duration_seconds": 900,
    "estimated_fare": 150.00,
    "female_only": false,
    "max_detour_minutes": 10,
    "max_wait_minutes": 5,
    "min_route_overlap": 0.70
  }'
# Returns candidates with compatibility scores

# 2. Create pool with matching rides
curl -X POST http://localhost:3013/v1/pooling/pools \
  -H "Content-Type: application/json" \
  -d '{
    "driver_id": "driver_001",
    "ride_ids": ["ride_001", "ride_002"],
    "max_size": 2
  }'
# Returns pool_id (e.g., pool_xyz)

# 3. Activate pool (transitions to ACTIVE)
curl -X POST http://localhost:3013/v1/pooling/pools/pool_xyz/activate

# 4. Complete pool (transitions to COMPLETED)
curl -X POST http://localhost:3013/v1/pooling/pools/pool_xyz/complete

# 5. Get statistics
curl http://localhost:3013/v1/pooling/statistics
# Shows total_pools, active_pools, avg_compatibility_score, total_profit
```

---

## 🔄 KAFKA INTEGRATION

### Topics Published by Pooling Service
```yaml
pool.created              # New pool group formed
pool.matching.started     # Compatibility matching begun
pool.matched              # Compatible rides found
pool.activated            # Pool transitioned to ACTIVE
pool.completed            # Pool successfully completed
pool.cancelled            # Pool cancelled
```

### Consumers
- **Analytics Service**: Aggregates pool metrics
- **Pricing Service**: Calculates pool discounts
- **Driver Service**: Notifies driver of pool assignments

---

## 📊 ALGORITHM DETAILS

### Route Overlap Calculation
```
Uses bounding box approach:
1. Calculate bounding box for each route
2. Find overlap area
3. Score = overlap_area / total_area covered
```

### Compatibility Scoring
```
Final Score = (overlap * 0.4) + (profit * 0.3) + (eta * 0.2) + (proximity * 0.1)

- Route Overlap: 40% weight (most important - must match route)
- Profitability: 30% weight (higher fares better)
- ETA Similarity: 20% weight (similar duration)
- Pickup Proximity: 10% weight (close pickups easier)

Viable if: score > 0.5 AND route_overlap >= 0.70
```

### Constraint Validation
```
Valid if all true:
1. Detour <= 10 minutes
2. Wait time <= 5 minutes
3. Route overlap >= 70%
4. Female-only flags match (if applicable)
```

---

## 🚀 PERFORMANCE OPTIMIZATION

### Current Implementation
- Nearest neighbor for route optimization (O(n²))
- Single-threaded matching

### Future Enhancements
- [ ] Traveling Salesman Problem (TSP) solver for optimal routes
- [ ] Parallel candidate evaluation
- [ ] Redis caching for frequently-matched routes
- [ ] ML model to predict pool success rate

---

## 📈 METRICS & MONITORING

### Key Metrics to Track
1. **Pool Formation Rate**: Pools created per hour
2. **Match Success Rate**: % of pools that complete
3. **Avg Pool Size**: Average rides per pool (target: 2.3)
4. **Avg Compatibility Score**: Higher is better (target: >0.75)
5. **Total Savings**: Rider cost reduction (target: 25-30%)
6. **Platform Profit**: Revenue from pooling (target: 15-20% of ride fare)

### Dashboards
- Real-time active pools
- Historical pool completion rate
- Compatibility score distribution
- Profit per pool over time

---

## 🏁 NEXT PHASE

### Phase 5: Pricing Service (2 weeks)
- Fare calculation formula
- Surge multiplier logic
- Discount application (including pool discounts)
- Pricing history audit

**Phase 5 depends on Phase 4** ✅ Complete for:
- Pool-specific fare logic
- Pool discount percentages
- Combined pricing (ride + pool optimization fee)

---

**Status**: ✅ **PHASE 4 COMPLETE & READY FOR EXECUTION**

**Files**: 6 Go files (~25 KB code) + database migration  
**Services**: 1 operational (Pooling Service)  
**Endpoints**: 5 REST + statistics  
**Database**: 5 tables + triggers + views  
**Next Phase**: Phase 5 - Pricing Service (2 weeks)

