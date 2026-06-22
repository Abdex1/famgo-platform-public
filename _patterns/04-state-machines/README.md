# 🔄 STATE MACHINE PATTERNS
## Extracted from uber-master

**Status:** Pattern 4/8

---

## Driver State Machine

```go
type DriverStatus string

const (
    DriverPending    DriverStatus = "pending"
    DriverApproved   DriverStatus = "approved"
    DriverActive     DriverStatus = "active"
    DriverSuspended  DriverStatus = "suspended"
)

var ValidTransitions = map[DriverStatus][]DriverStatus{
    DriverPending:   {DriverApproved},
    DriverApproved:  {DriverActive},
    DriverActive:    {DriverSuspended, DriverPending},
    DriverSuspended: {DriverActive},
}

func (s DriverStatus) CanTransitionTo(next DriverStatus) bool {
    valid := ValidTransitions[s]
    for _, v := range valid {
        if v == next {
            return true
        }
    }
    return false
}
```

## Trip State Machine

```go
type TripStatus string

const (
    TripRequested      TripStatus = "requested"
    TripDriverAssigned TripStatus = "driver_assigned"
    TripStarted        TripStatus = "started"
    TripCompleted      TripStatus = "completed"
    TripCancelled      TripStatus = "cancelled"
)

var ValidTripTransitions = map[TripStatus][]TripStatus{
    TripRequested:      {TripDriverAssigned, TripCancelled},
    TripDriverAssigned: {TripStarted, TripCancelled},
    TripStarted:        {TripCompleted},
    TripCompleted:      {},
    TripCancelled:      {},
}
```

**Pattern 4 Status:** READY FOR USE

---
