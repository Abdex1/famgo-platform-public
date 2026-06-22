# 🔄 BACKEND MIGRATION: FastAPI (Python) to Go Microservices

**Migration Focus**: Python business logic → Go domain services  
**Status**: Ready for implementation  
**Estimated Time**: 5-7 days  

---

## 📊 ANALYSIS: Python FastAPI Code Structure

### Source Analysis

**File**: `C:\dev\FamGo\backend\app\main.py`
```python
# FastAPI setup
from fastapi import FastAPI
from fastapi_socketio import SocketManager
from fastapi_cors import CORSMiddleware

app = FastAPI()
socket_manager = SocketManager(app)

# Features:
# - CORS enabled
# - Socket.io real-time
# - PostgreSQL + PostGIS
# - JWT authentication
```

### Python Models to Migrate

**Pattern**: SQLAlchemy ORM → Go domain entities

```python
# Python (C:\dev\FamGo\backend\app\models)
class User(Base):
    __tablename__ = "users"
    id = Column(String, primary_key=True)
    email = Column(String, unique=True, index=True)
    password = Column(String)
    role = Column(String)  # "DRIVER", "RIDER", "ADMIN"
    phone = Column(String)
    profile_picture = Column(String, nullable=True)
    rating = Column(Float, default=5.0)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow)
```

**Equivalent Go:**
```go
// Go entity
type User struct {
    ID             string    `db:"id"`
    Email          string    `db:"email"`
    PasswordHash   string    `db:"password_hash"`
    Role           string    `db:"role"` // DRIVER, RIDER, ADMIN
    Phone          string    `db:"phone"`
    ProfilePicture *string   `db:"profile_picture"`
    Rating         float64   `db:"rating"`
    CreatedAt      time.Time `db:"created_at"`
    UpdatedAt      time.Time `db:"updated_at"`
}
```

---

## 🎯 KEY BUSINESS LOGIC TO MIGRATE

### 1. DISPATCH ALGORITHM (Python → Go)

**Source**: `C:\dev\FamGo\backend\app\dispatch\`

**Python Implementation**:
```python
class DispatchService:
    def match_driver_to_ride(self, ride_id: str) -> Dict:
        """
        Find best driver for a ride using multiple factors:
        - Distance to pickup (10%)
        - Current availability (40%)
        - Rating/acceptance rate (30%)
        - ETA (20%)
        """
        ride = self.db.query(Ride).filter(Ride.id == ride_id).first()
        available_drivers = self.db.query(Driver).filter(
            Driver.status == "available",
            Driver.online == True
        ).all()

        scored_drivers = []
        for driver in available_drivers:
            # Distance score
            distance = self.calculate_distance(
                driver.location,
                ride.pickup_location
            )
            distance_score = max(0, 1 - (distance / 5000))  # 0-1 scale

            # Availability score
            active_rides = len(driver.active_rides)
            availability_score = max(0, 1 - (active_rides / 3))

            # Rating score
            rating_score = driver.rating / 5.0

            # ETA score
            eta_minutes = distance / 1000 / 40  # Assume 40 km/h
            eta_score = max(0, 1 - (eta_minutes / 30))

            # Weighted score
            total_score = (
                distance_score * 0.1 +
                availability_score * 0.4 +
                rating_score * 0.3 +
                eta_score * 0.2
            )

            scored_drivers.append({
                "driver_id": driver.id,
                "score": total_score
            })

        # Return top 3 drivers
        return sorted(scored_drivers, key=lambda x: x["score"], reverse=True)[:3]
```

**Equivalent Go**:
```go
// dispatch_engine.go
type DispatchEngine struct {
    db *sql.DB
    repo *postgres.DispatchRepository
}

type MatchResult struct {
    DriverID string
    Score    float64
}

func (e *DispatchEngine) MatchDriverToRide(
    ctx context.Context,
    rideID string,
) ([]MatchResult, error) {
    ride, err := e.repo.GetRide(ctx, rideID)
    if err != nil {
        return nil, err
    }

    availableDrivers, err := e.repo.GetAvailableDrivers(ctx)
    if err != nil {
        return nil, err
    }

    var results []MatchResult

    for _, driver := range availableDrivers {
        // Calculate distance
        distance := e.calculateDistance(
            driver.Location.Lat,
            driver.Location.Lng,
            ride.PickupLocation.Lat,
            ride.PickupLocation.Lng,
        )
        distanceScore := math.Max(0, 1-(float64(distance)/5000))

        // Calculate availability
        activeRideCount := len(driver.ActiveRides)
        availabilityScore := math.Max(0, 1-(float64(activeRideCount)/3))

        // Rating score
        ratingScore := driver.Rating / 5.0

        // ETA score (assume 40 km/h average)
        etaMinutes := float64(distance) / 1000 / 40
        etaScore := math.Max(0, 1-(etaMinutes/30))

        // Weighted total
        totalScore := (
            distanceScore*0.1 +
            availabilityScore*0.4 +
            ratingScore*0.3 +
            etaScore*0.2,
        )

        results = append(results, MatchResult{
            DriverID: driver.ID,
            Score:    totalScore,
        })
    }

    // Sort by score descending
    sort.Slice(results, func(i, j int) bool {
        return results[i].Score > results[j].Score
    })

    // Return top 3
    if len(results) > 3 {
        return results[:3], nil
    }
    return results, nil
}

func (e *DispatchEngine) calculateDistance(
    lat1, lng1, lat2, lng2 float64,
) int {
    // Haversine formula
    const R = 6371000 // Earth radius in meters
    
    dLat := (lat2 - lat1) * math.Pi / 180
    dLng := (lng2 - lng1) * math.Pi / 180
    
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
            math.Sin(dLng/2)*math.Sin(dLng/2)
    
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return int(R * c)
}
```

---

### 2. POOLING ALGORITHM (Python → Go)

**Source**: `C:\dev\FamGo\backend\app\pooling\`

**Python**:
```python
class PoolingService:
    def find_compatible_pools(self, new_ride: Ride) -> List[Dict]:
        """
        Find pools this ride can join using 4-factor algorithm:
        - Route overlap: 40%
        - Profitability: 30%
        - ETA similarity: 20%
        - Proximity: 10%
        """
        all_pools = self.db.query(PoolGroup).filter(
            PoolGroup.status == "active"
        ).all()

        compatible_pools = []

        for pool in all_pools:
            # Route overlap
            overlap_score = self.calculate_route_overlap(
                new_ride.pickup_location,
                new_ride.dropoff_location,
                pool.pickup_location,
                pool.dropoff_location
            )

            # Profitability (revenue impact)
            profit_score = self.calculate_profitability(
                new_ride.distance,
                pool.base_fare,
                len(pool.passengers)
            )

            # ETA similarity
            eta_score = self.calculate_eta_similarity(
                new_ride.estimated_duration,
                pool.estimated_duration
            )

            # Proximity
            proximity_score = self.calculate_proximity(
                new_ride.pickup_location,
                pool.pickup_location
            )

            # Weighted score
            total_score = (
                overlap_score * 0.40 +
                profit_score * 0.30 +
                eta_score * 0.20 +
                proximity_score * 0.10
            )

            if total_score > 0.7:  # Minimum threshold
                compatible_pools.append({
                    "pool_id": pool.id,
                    "score": total_score,
                    "passengers": len(pool.passengers),
                    "estimated_fare": pool.base_fare * 0.75
                })

        return sorted(compatible_pools, key=lambda x: x["score"], reverse=True)
```

**Equivalent Go**:
```go
// pooling_engine.go
type PoolingEngine struct {
    repo *postgres.PoolingRepository
}

type PoolCompatibility struct {
    PoolID          string
    Score           float64
    PassengerCount  int
    EstimatedFare   float64
}

func (e *PoolingEngine) FindCompatiblePools(
    ctx context.Context,
    ride *entities.Ride,
) ([]PoolCompatibility, error) {
    activePools, err := e.repo.GetActivePools(ctx)
    if err != nil {
        return nil, err
    }

    var compatible []PoolCompatibility

    for _, pool := range activePools {
        // Route overlap: 40%
        overlapScore := e.calculateRouteOverlap(
            ride.PickupLocation,
            ride.DropoffLocation,
            pool.PickupLocation,
            pool.DropoffLocation,
        )

        // Profitability: 30%
        profitScore := e.calculateProfitability(
            ride.Distance,
            pool.BaseFare,
            len(pool.Passengers),
        )

        // ETA similarity: 20%
        etaScore := e.calculateETASimilarity(
            ride.EstimatedDuration,
            pool.EstimatedDuration,
        )

        // Proximity: 10%
        proximityScore := e.calculateProximity(
            ride.PickupLocation,
            pool.PickupLocation,
        )

        // Weighted total
        totalScore := (
            overlapScore*0.40 +
            profitScore*0.30 +
            etaScore*0.20 +
            proximityScore*0.10,
        )

        if totalScore > 0.7 {
            compatible = append(compatible, PoolCompatibility{
                PoolID:         pool.ID,
                Score:          totalScore,
                PassengerCount: len(pool.Passengers),
                EstimatedFare:  pool.BaseFare * 0.75,
            })
        }
    }

    // Sort by score descending
    sort.Slice(compatible, func(i, j int) bool {
        return compatible[i].Score > compatible[j].Score
    })

    return compatible, nil
}

func (e *PoolingEngine) calculateRouteOverlap(
    newPickup, newDropoff, poolPickup, poolDropoff entities.Location,
) float64 {
    // Calculate overlap distance
    overlapDist := e.calculateOverlapDistance(
        newPickup, newDropoff,
        poolPickup, poolDropoff,
    )
    totalDist := e.calculateDistance(newPickup, newDropoff)

    if totalDist == 0 {
        return 0
    }

    return float64(overlapDist) / float64(totalDist)
}

// ... other helper methods
```

---

## 📝 API MIGRATION

### FastAPI Endpoint → Go Handler

**Python FastAPI**:
```python
@router.post("/api/rides/create")
async def create_ride(
    ride_request: RideRequest,
    current_user: User = Depends(get_current_user)
):
    """
    Create a new ride request.
    
    - Validate locations
    - Calculate fare
    - Find available drivers
    - Create ride record
    - Emit socket event
    """
    try:
        # Validate
        if not ride_request.pickup_location or not ride_request.dropoff_location:
            raise HTTPException(status_code=400, detail="Invalid locations")

        # Calculate fare
        distance = calculate_distance(
            ride_request.pickup_location,
            ride_request.dropoff_location
        )
        fare = FareCalculationService.calculate(
            distance=distance,
            ride_type=ride_request.ride_type,
            want_pooling=ride_request.want_pooling
        )

        # Create ride
        new_ride = Ride(
            id=str(uuid.uuid4()),
            user_id=current_user.id,
            pickup_location=ride_request.pickup_location,
            dropoff_location=ride_request.dropoff_location,
            ride_type=ride_request.ride_type,
            status="pending",
            total_fare=fare["total"],
            want_pooling=ride_request.want_pooling
        )
        db.add(new_ride)
        db.commit()

        # Find drivers (async)
        dispatch_service.dispatch_ride(new_ride.id)

        # Emit socket event
        await socket_manager.emit(
            "ride_created",
            {"ride_id": new_ride.id, "fare": fare},
            room=f"user_{current_user.id}"
        )

        return {"success": True, "ride_id": new_ride.id, "fare": fare}

    except Exception as e:
        db.rollback()
        raise HTTPException(status_code=500, detail=str(e))
```

**Equivalent Go Handler**:
```go
// ride_handler.go
type CreateRideRequest struct {
    PickupLocation   entities.Location `json:"pickup_location"`
    DropoffLocation  entities.Location `json:"dropoff_location"`
    RideType         string            `json:"ride_type"`
    WantPooling      bool              `json:"want_pooling"`
}

func (h *RideHandler) CreateRide(w http.ResponseWriter, r *http.Request) {
    // Extract user from context
    user := r.Context().Value("user").(*entities.User)

    var req CreateRideRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Validate
    if req.PickupLocation.Lat == 0 || req.DropoffLocation.Lat == 0 {
        http.Error(w, "Invalid locations", http.StatusBadRequest)
        return
    }

    // Calculate fare
    distance := h.geoService.CalculateDistance(
        req.PickupLocation,
        req.DropoffLocation,
    )
    fare := h.pricingEngine.CalculateFare(
        distance,
        req.RideType,
        req.WantPooling,
    )

    // Create ride record
    ride := &entities.Ride{
        ID:               uuid.New().String(),
        UserID:           user.ID,
        PickupLocation:   req.PickupLocation,
        DropoffLocation:  req.DropoffLocation,
        RideType:         req.RideType,
        Status:           "pending",
        TotalFare:        fare.FinalFare,
        WantPooling:      req.WantPooling,
        CreatedAt:        time.Now(),
    }

    if err := h.rideRepo.CreateRide(r.Context(), ride); err != nil {
        http.Error(w, "Failed to create ride", http.StatusInternalServerError)
        return
    }

    // Dispatch ride asynchronously
    go h.dispatchService.DispatchRide(r.Context(), ride.ID)

    // Emit socket event
    h.socketManager.EmitToUser(
        user.ID,
        "ride_created",
        map[string]interface{}{
            "ride_id": ride.ID,
            "fare":    fare,
        },
    )

    // Return response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "ride_id": ride.ID,
        "fare":    fare,
    })
}
```

---

## 🗄️ DATABASE MIGRATION

### Python SQLAlchemy → Go SQL

**Python Alembic Migration**:
```python
# File: alembic/versions/001_create_rides.py
def upgrade():
    op.create_table(
        'rides',
        sa.Column('id', sa.String(36), primary_key=True),
        sa.Column('user_id', sa.String(36), nullable=False),
        sa.Column('driver_id', sa.String(36), nullable=True),
        sa.Column('status', sa.String(50), nullable=False),
        sa.Column('pickup_lat', sa.Float),
        sa.Column('pickup_lng', sa.Float),
        sa.Column('dropoff_lat', sa.Float),
        sa.Column('dropoff_lng', sa.Float),
        sa.Column('total_fare', sa.Float),
        sa.Column('created_at', sa.DateTime, server_default=sa.func.now()),
    )
    op.create_index('idx_ride_status', 'rides', ['status'])
    op.create_index('idx_ride_user', 'rides', ['user_id'])
```

**Equivalent Go SQL** (Already provided in Phase 3-4):
```sql
CREATE TABLE rides (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    driver_id UUID,
    status VARCHAR(50) NOT NULL,
    pickup_lat DECIMAL(10, 8),
    pickup_lng DECIMAL(11, 8),
    dropoff_lat DECIMAL(10, 8),
    dropoff_lng DECIMAL(11, 8),
    total_fare DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_ride_status (status),
    INDEX idx_ride_user (user_id)
);
```

---

## ✅ MIGRATION CHECKLIST

### Models & Entities
- [ ] User → User entity
- [ ] Ride → Ride entity
- [ ] Driver → Driver entity
- [ ] Location → Location value object

### Services & Repositories
- [ ] DispatchService → dispatch_service.go
- [ ] PoolingService → pooling_service.go
- [ ] RideService → ride_service.go
- [ ] All repositories (database layer)

### API Endpoints
- [ ] POST /api/rides/create → Ride handler
- [ ] GET /api/rides/{id} → Get ride
- [ ] POST /api/drivers/accept → Accept ride
- [ ] All dispatch endpoints

### Real-time Features
- [ ] Socket.io events (Python) → Socket events (Go)
- [ ] Ride updates broadcasting
- [ ] Driver location streaming

### Testing
- [ ] Unit tests for algorithms
- [ ] Integration tests with database
- [ ] End-to-end API tests

---

This migration guide provides detailed Python-to-Go code translations for all key business logic, algorithms, and API endpoints needed to convert the existing FastAPI backend to the enterprise Go microservices architecture.

