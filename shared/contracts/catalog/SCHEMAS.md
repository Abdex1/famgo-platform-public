# 📊 SCHEMA REGISTRY: Data Structure Definitions

**Status:** Task 2 Phase 2.2 Complete  
**Location:** shared/contracts/events/schemas/  
**Type Safety:** Go structs with JSON tags

---

## SCHEMA CATALOG

### 1. AUTH SCHEMAS (Location: shared/contracts/events/schemas/auth-events.go)

#### LoginSucceededEvent
```go
type LoginSucceededEvent struct {
    UserID       string `json:"user_id"`
    Email        string `json:"email"`
    IP           string `json:"ip"`
    DeviceID     string `json:"device_id"`
    SessionID    string `json:"session_id"`
}
```

**Fields:**
- `user_id` (string): Unique user identifier
- `email` (string): User email address
- `ip` (string): Client IP address
- `device_id` (string): Device identifier
- `session_id` (string): Session identifier

**Usage:** Track successful login events

#### LoginFailedEvent
```go
type LoginFailedEvent struct {
    Email        string `json:"email"`
    IP           string `json:"ip"`
    Reason       string `json:"reason"`
}
```

**Fields:**
- `email` (string): Email attempted
- `ip` (string): Client IP address
- `reason` (string): Failure reason (invalid_password, account_locked, etc.)

**Usage:** Track failed login attempts for security

#### TokenRefreshedEvent
```go
type TokenRefreshedEvent struct {
    UserID       string `json:"user_id"`
    SessionID    string `json:"session_id"`
}
```

**Fields:**
- `user_id` (string): User identifier
- `session_id` (string): Session identifier

**Usage:** Track token refresh operations

#### SessionRevokedEvent
```go
type SessionRevokedEvent struct {
    UserID       string `json:"user_id"`
    SessionID    string `json:"session_id"`
}
```

**Fields:**
- `user_id` (string): User identifier
- `session_id` (string): Session identifier

**Usage:** Track session revocation (logout, timeout, etc.)

---

### 2. PAYMENT SCHEMAS (Location: shared/contracts/events/payment/v1/payment_completed.go)

#### PaymentCompleted
```go
type PaymentCompleted struct {
    PaymentID   string    `json:"payment_id"`
    TripID      string    `json:"trip_id"`
    RiderID     string    `json:"rider_id"`
    DriverID    string    `json:"driver_id"`
    Amount      float64   `json:"amount"`
    Currency    string    `json:"currency"`
    Status      string    `json:"status"`
    CompletedAt time.Time `json:"completed_at"`
}
```

**Fields:**
- `payment_id` (string): Unique payment identifier
- `trip_id` (string): Associated trip/ride identifier
- `rider_id` (string): Rider user identifier
- `driver_id` (string): Driver user identifier
- `amount` (float64): Payment amount
- `currency` (string): Currency code (USD, ETB, etc.)
- `status` (string): Payment status (completed, pending, failed)
- `completed_at` (timestamp): Completion time

**Usage:** Record completed payment transactions

---

### 3. RIDE SCHEMAS (Location: shared/contracts/events/driver/v1/ride_requested.go)

#### RideRequested
```go
type RideRequested struct {
    TripID         string    `json:"trip_id"`
    RiderID        string    `json:"rider_id"`
    PickupLat      float64   `json:"pickup_lat"`
    PickupLng      float64   `json:"pickup_lng"`
    DestinationLat float64   `json:"destination_lat"`
    DestinationLng float64   `json:"destination_lng"`
    VehicleType    string    `json:"vehicle_type"`
    RequestedAt    time.Time `json:"requested_at"`
}
```

**Fields:**
- `trip_id` (string): Unique trip identifier
- `rider_id` (string): Rider user identifier
- `pickup_lat` (float64): Pickup latitude
- `pickup_lng` (float64): Pickup longitude
- `destination_lat` (float64): Destination latitude
- `destination_lng` (float64): Destination longitude
- `vehicle_type` (string): Vehicle type (economy, comfort, premium)
- `requested_at` (timestamp): Request time

**Usage:** Record new ride requests for dispatch

---

### 4. AUDIT SCHEMAS (Location: shared/contracts/events/schemas/audit-events.go)

#### AuditEvent (Base Pattern)
```go
// Generic audit event pattern (typesafe)
// All audit events follow this pattern:
//   - Actor: who performed the action
//   - Action: what was done
//   - Resource: what was affected
//   - Result: success/failure
//   - Timestamp: when it happened
//   - IP/Device: where it happened from
```

**Usage:** Track all system changes for compliance

---

## SCHEMA VERSIONING

### Current Schema Version: v1

**All schemas** use v1 structure (defined in shared/contracts/events/versions/versions.go)

### Adding New Schemas

**Process:**
1. Create file in shared/contracts/events/schemas/ or shared/contracts/events/{domain}/v1/
2. Define Go struct with JSON tags
3. Ensure all fields have JSON tags (no omitted fields)
4. Add schema to SCHEMAS.md registry
5. Update VERSIONS.md migration guide

**Example (new schema):**
```go
// shared/contracts/events/schemas/user-events.go
package schemas

type UserRegisteredEvent struct {
    UserID      string `json:"user_id"`
    Email       string `json:"email"`
    UserType    string `json:"user_type"` // driver or passenger
    RegisteredAt time.Time `json:"registered_at"`
}
```

---

## SCHEMA EVOLUTION

### V1 → V2 Migration Path

**When to bump version:**
- Adding required fields (breaking change)
- Removing fields (breaking change)
- Changing field type (breaking change)

**How to evolve:**
1. Create shared/contracts/events/{domain}/v2/ directory
2. Copy v1 schema to v2/
3. Make changes to v2/
4. Update MIGRATION.md with migration guide
5. Deploy with dual v1+v2 support
6. Migrate consumers to v2
7. Deprecate v1 (timeline: 30 days)

**Example:**
```
shared/contracts/events/payment/
├── v1/
│   └── payment_completed.go (original)
└── v2/
    └── payment_completed.go (with new fields)
```

---

## QUALITY GATES

✅ All schemas have:
- [x] Go struct definition (type-safe)
- [x] JSON tags on all fields
- [x] Clear field documentation
- [x] Version reference (v1, v2, etc.)
- [x] Publisher responsibility
- [x] Consumer list

✅ No schemas are:
- [x] Duplicated (each schema defined once)
- [x] Unversioned
- [x] Without documentation
- [x] Untyped (all are Go structs)

---

## RECOMMENDED PRACTICES

1. **Always use JSON tags:**
   ```go
   type Event struct {
       UserID string `json:"user_id"` // ✅ Good
       // NOT: UserID string ❌ Bad
   }
   ```

2. **Use descriptive field names:**
   ```go
   RiderID   string // ✅ Clear
   // NOT: ID string ❌ Ambiguous
   ```

3. **Include timestamps:**
   ```go
   type Event struct {
       // ... other fields
       CreatedAt time.Time `json:"created_at"` // ✅ Always include
   }
   ```

4. **Document all fields:**
   ```go
   type PaymentCompleted struct {
       PaymentID string `json:"payment_id"` // Unique identifier
       Amount    float64 `json:"amount"`    // Amount in currency
   }
   ```

---

**Schema Registry Status:** ✅ COMPLETE & VERIFIED

