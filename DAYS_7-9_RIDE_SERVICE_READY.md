# 🚀 DAYS 7-9: RIDE SERVICE EXECUTION READY

**Status:** User Service 100% Complete ✅  
**Next:** Build Ride Service (Days 7-9, 12 hours)  
**Overall Phase:** 81% Complete (65 of 80 hours)  

---

## 🎯 RIDE SERVICE: BUILD PATTERN

**Copy User Service structure exactly:**
```
services/ride-service/
├── internal/domain/          ← Replace User entities with Ride + state machine
├── internal/application/     ← Copy command/query pattern
├── internal/infrastructure/  ← Copy repo pattern
├── internal/transport/       ← Copy HTTP handlers pattern
├── internal/bootstrap/       ← Copy DI container
├── internal/config/          ← Copy config loader
├── cmd/main.go               ← Copy main entry point
├── db/migrations/            ← Create ride schema
├── deployments/              ← Copy K8s manifests
├── Dockerfile                ← Copy Dockerfile (port 5005)
└── tests/unit/               ← Copy test pattern
```

---

## 🏗️ DOMAIN LAYER: STATE MACHINE

**Create entities with lifecycle states:**

```go
type RideStatus string

const (
    RideStatusRequested      RideStatus = "REQUESTED"      // User creates ride
    RideStatusSearching      RideStatus = "SEARCHING"      // Dispatch searches
    RideStatusAssigned       RideStatus = "ASSIGNED"       // Driver found
    RideStatusDriverArriving RideStatus = "DRIVER_ARRIVING" // Driver en route
    RideStatusStarted        RideStatus = "STARTED"        // Pickup complete
    RideStatusCompleted      RideStatus = "COMPLETED"      // Dropoff complete
    RideStatusCancelled      RideStatus = "CANCELLED"      // Cancelled
)

type Ride struct {
    ID              string    // UUID
    PassengerID     string    // Foreign key
    DriverID        string    // Assigned later
    PickupLat       float64
    PickupLon       float64
    DropoffLat      float64
    DropoffLon      float64
    Status          RideStatus
    EstimatedFare   float32
    ActualFare      float32
    PickupTime      time.Time
    DropoffTime     time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// State transitions with validation
func (r *Ride) CanTransitionTo(newStatus RideStatus) bool {
    // Validate allowed transitions
    allowed := map[RideStatus][]RideStatus{
        RideStatusRequested: {RideStatusSearching, RideStatusCancelled},
        RideStatusSearching: {RideStatusAssigned, RideStatusCancelled},
        RideStatusAssigned: {RideStatusDriverArriving, RideStatusCancelled},
        RideStatusDriverArriving: {RideStatusStarted, RideStatusCancelled},
        RideStatusStarted: {RideStatusCompleted, RideStatusCancelled},
    }
    
    for _, s := range allowed[r.Status] {
        if s == newStatus {
            return true
        }
    }
    return false
}
```

---

## 📋 APPLICATION LAYER: LIFECYCLE COMMANDS

**Create commands for state transitions:**

```go
// Command: Create Ride (REQUESTED)
type CreateRideCommand struct {
    PassengerID   string
    PickupLat     float64
    PickupLon     float64
    DropoffLat    float64
    DropoffLon    float64
}
// Handler publishes: ride.requested

// Command: Assign Driver (ASSIGNED)
type AssignDriverCommand struct {
    RideID  string
    DriverID string
}
// Handler publishes: ride.assigned

// Command: Start Ride (STARTED)
type StartRideCommand struct {
    RideID string
}
// Handler publishes: ride.started

// Command: Complete Ride (COMPLETED)
type CompleteRideCommand struct {
    RideID    string
    ActualFare float32
}
// Handler publishes: ride.completed

// Command: Cancel Ride (CANCELLED)
type CancelRideCommand struct {
    RideID string
    Reason string
}
// Handler publishes: ride.cancelled
```

---

## 🔄 EVENTS: INTEGRATE WITH EVENT-DRIVEN WORKFLOW

**Use from shared/contracts/events:**

```go
// Domain publishes these events via platform/event-bus
- ride.requested      (passengerId, pickupLat, pickupLon, dropoffLat, dropoffLon)
- ride.assigned       (rideId, driverId)
- ride.started        (rideId, pickupTime)
- ride.completed      (rideId, dropoffTime, actualFare)
- ride.cancelled      (rideId, reason)

// Domain consumes these events
- dispatch.driver.assigned  (rideId, driverId)
- pricing.fare.calculated   (rideId, estimatedFare)
- payment.authorized        (rideId, amount)
```

---

## 🗄️ DATABASE SCHEMA

```sql
CREATE TABLE rides (
    id UUID PRIMARY KEY,
    passenger_id UUID NOT NULL REFERENCES users(id),
    driver_id UUID REFERENCES users(id),
    pickup_lat DECIMAL(10,8),
    pickup_lon DECIMAL(11,8),
    dropoff_lat DECIMAL(10,8),
    dropoff_lon DECIMAL(11,8),
    status VARCHAR(50) DEFAULT 'REQUESTED',
    estimated_fare DECIMAL(10,2),
    actual_fare DECIMAL(10,2),
    pickup_time TIMESTAMP,
    dropoff_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ride_status_history (
    id UUID PRIMARY KEY,
    ride_id UUID NOT NULL REFERENCES rides(id),
    old_status VARCHAR(50),
    new_status VARCHAR(50),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 📊 DAYS 7-9 TIMELINE

**Day 7 (4 hours):**
- Domain layer: Ride entity + state machine + service (2 hrs)
- Application layer: Commands + queries (2 hrs)

**Day 8 (4 hours):**
- Infrastructure: Repos + cache (2 hrs)
- Transport: HTTP handlers (2 hrs)

**Day 9 (4 hours):**
- Bootstrap + config (1 hr)
- Database + deployment (2 hrs)
- Tests (1 hr)

---

## ✅ READY TO START

All User Service files available as reference:
- Domain pattern: `services/user-service/internal/domain/`
- Application pattern: `services/user-service/internal/application/`
- Infrastructure pattern: `services/user-service/internal/infrastructure/`
- Transport pattern: `services/user-service/internal/transport/`
- Deployment pattern: `services/user-service/deployments/`
- Dockerfile template: `services/user-service/Dockerfile`

**Copy, adjust for Ride domain, and build.**

