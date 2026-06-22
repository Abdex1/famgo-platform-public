# 📋 EVENT CATALOG: Complete Event Registry

**Status:** Comprehensive audit of `shared/contracts/events/`  
**Created:** Day 1 of Weeks 3-4 Audit Phase  
**Repository:** github.com/Abdex1/FamGo-platform  
**Purpose:** Single source of truth for all event types, topics, and governance

---

## 🎯 EXECUTIVE SUMMARY

The FamGo platform uses a **centralized, versioned event system** for all inter-service communication.

**Critical Rule:** NO service may define local events. ALL events originate from `shared/contracts/events/`.

**Key Components:**
- ✅ Centralized event catalog
- ✅ Versioning strategy
- ✅ Dead-letter queue (DLQ) handling
- ✅ Idempotency management
- ✅ Retry policies
- ✅ Event envelopes (wrappers)
- ✅ Domain-specific event schemas

---

## 📁 DIRECTORY STRUCTURE ANALYSIS

### Location: `shared/contracts/events/`

```
shared/contracts/events/
├── catalog/              # Central event type registry
├── common/               # Shared event metadata
├── dlq/                  # Dead-letter queue policies
├── driver/               # Driver domain events
├── envelopes/            # Event wrapper structure
├── idempotency/          # Idempotency key management
├── payment/              # Payment domain events
├── policies/             # Event handling policies
├── rating/               # Rating domain events
├── retry/                # Retry policy definitions
├── ride/                 # Ride domain events
├── schemas/              # JSON schemas for events
├── topics/               # Kafka topic definitions
├── trip/                 # Trip domain events
└── versions/             # Event versioning strategies
```

---

## 🔍 LAYER-BY-LAYER AUDIT

### Layer 1: Event Catalog (`catalog/`)

**Purpose:** Central registry of all event types in the system

**Event Categories:**

#### 1a. User/Authentication Events
- `user.created` - New user registered
- `user.verified` - User email/phone verified
- `user.updated` - User profile updated
- `user.deleted` - User account deleted
- `login.success` - User logged in successfully
- `login.failed` - Login failed (too many attempts, wrong credentials)
- `logout` - User logged out
- `password.reset.requested` - Password reset initiated
- `password.reset.completed` - Password reset finished
- `otp.sent` - OTP sent to user
- `otp.verified` - OTP verification successful
- `device.registered` - Device registered for push notifications
- `device.unregistered` - Device removed

#### 1b. Driver Events
- `driver.registered` - New driver onboarded
- `driver.approved` - Driver KYC approved
- `driver.rejected` - Driver KYC rejected
- `driver.verified` - Driver documents verified
- `driver.status.changed` - Driver status changed (active, inactive, suspended)
- `driver.location.updated` - Driver location updated
- `driver.went.online` - Driver came online
- `driver.went.offline` - Driver went offline
- `driver.rating.received` - Driver received rating
- `driver.rating.updated` - Driver average rating updated
- `driver.document.uploaded` - Driver uploaded required document
- `driver.document.expired` - Driver document expired (license, insurance, etc.)

#### 1c. Rider Events
- `rider.registered` - New rider account created
- `rider.profile.updated` - Rider profile changed
- `rider.address.added` - Rider added home/work address
- `rider.address.removed` - Rider deleted address
- `rider.rating.received` - Rider received rating
- `rider.payment.method.added` - New payment method added
- `rider.payment.method.removed` - Payment method deleted
- `rider.favorite.created` - Rider marked location as favorite

#### 1d. Ride Events
- `ride.requested` - Rider requested a ride (main event)
- `ride.assigned` - Driver assigned to ride
- `ride.driver.arrived` - Driver arrived at pickup
- `ride.started` - Ride has started (driver at passenger location)
- `ride.passenger.updated` - Passenger count/details changed mid-ride
- `ride.route.updated` - Pickup/dropoff location changed
- `ride.completed` - Ride successfully completed
- `ride.cancelled` - Ride was cancelled
- `ride.cancelled.by.driver` - Driver cancelled ride
- `ride.cancelled.by.rider` - Rider cancelled ride
- `ride.no.show` - Rider didn't show up
- `ride.redirected` - Ride redirected to different driver
- `ride.shared` - Ride shared/pooled with another rider
- `ride.unpooled` - Pooled ride split into separate rides
- `ride.disputed` - Dispute created on ride

#### 1e. Dispatch Events
- `driver.search.initiated` - System started searching for drivers
- `driver.searched` - System found drivers
- `driver.offered` - Driver offered the ride
- `driver.accepted` - Driver accepted the ride
- `driver.rejected` - Driver rejected the ride
- `driver.no.response` - Driver didn't respond to offer (timeout)
- `driver.offer.cancelled` - Offer cancelled (ride already assigned)
- `driver.assignment.confirmed` - Assignment confirmed
- `dispatch.failed` - Dispatch process failed

#### 1f. Pricing Events
- `price.calculated` - Fare calculated
- `price.surging` - Surge pricing applied
- `surge.multiplier.updated` - Surge multiplier changed
- `promotion.applied` - Promotion/discount applied
- `promotion.removed` - Promotion removed
- `pool.discount.applied` - Pooling discount applied
- `price.estimate.requested` - Rider requested fare estimate
- `price.estimate.provided` - Fare estimate sent to rider
- `price.adjustment` - Fare adjusted (due to route change, etc.)

#### 1g. Payment Events
- `payment.initiated` - Payment process started
- `payment.pending` - Payment is processing
- `payment.succeeded` - Payment completed successfully
- `payment.failed` - Payment failed
- `payment.retried` - Retry payment initiated
- `payment.method.changed` - Payment method changed mid-ride
- `payment.refund.requested` - Refund requested
- `payment.refund.initiated` - Refund processing started
- `payment.refund.completed` - Refund completed
- `payment.dispute.raised` - Payment dispute raised

#### 1h. Wallet Events
- `wallet.created` - New wallet created for user
- `wallet.credited` - Money added to wallet
- `wallet.debited` - Money deducted from wallet
- `wallet.balance.updated` - Wallet balance changed
- `wallet.frozen` - Wallet frozen (due to fraud, etc.)
- `wallet.unfrozen` - Wallet unfrozen
- `wallet.threshold.exceeded` - Wallet balance exceeded limit

#### 1i. Trip/Route Events
- `trip.started` - Trip tracking started
- `trip.location.updated` - Trip location updated
- `trip.route.changed` - Route recalculated
- `trip.eta.updated` - ETA recalculated
- `trip.completed` - Trip tracking completed
- `geofence.entered` - Driver entered geofence
- `geofence.exited` - Driver exited geofence
- `route.deviation` - Driver deviated from planned route

#### 1j. Safety Events
- `sos.triggered` - SOS/emergency button pressed
- `sos.cancelled` - SOS cancelled
- `sos.emergency.resolved` - Emergency resolved
- `sos.alert.sent` - Alert sent to emergency services
- `trip.shared` - Trip shared with trusted contact
- `trip.share.cancelled` - Trip share cancelled
- `unsafe.behavior.reported` - Unsafe behavior reported
- `safety.rating.updated` - Safety score updated

#### 1k. Fraud Events
- `fraud.check.triggered` - Fraud check initiated
- `fraud.suspicious.activity` - Suspicious activity detected
- `fraud.account.flagged` - Account flagged for review
- `fraud.case.created` - Fraud investigation case created
- `fraud.case.resolved` - Fraud case resolved
- `fraud.refund.issued` - Fraud refund issued

#### 1l. Rating/Review Events
- `rating.submitted` - Rating submitted by user
- `rating.confirmed` - Rating confirmed
- `rating.disputed` - Rating dispute raised
- `rating.removed` - Rating removed (by admin or user)
- `driver.rating.avg.updated` - Driver average rating recalculated
- `rider.rating.avg.updated` - Rider average rating recalculated

#### 1m. Notification Events
- `notification.sent` - Notification sent to user
- `notification.delivered` - Notification delivered
- `notification.read` - Notification read by user
- `push.notification.sent` - Push notification sent
- `sms.sent` - SMS sent
- `email.sent` - Email sent

#### 1n. Analytics Events
- `analytics.recorded` - Analytics event recorded
- `metric.captured` - System metric captured
- `error.tracked` - Error tracked for analytics
- `user.behavior.tracked` - User behavior tracked

---

### Layer 2: Common Event Fields (`common/`)

**Purpose:** Shared metadata that EVERY event must include

#### Required Event Metadata:
```go
EventID       string    // Unique event identifier (UUID)
EventType     string    // Type from catalog
EventVersion  int       // Version (v1, v2, etc.)
Timestamp     time.Time // When event occurred
AggregateID   string    // Primary entity ID (ride_id, driver_id, etc.)
```

#### Tracing/Correlation Fields:
```go
CorrelationID string    // Request correlation ID (trace all related events)
CausationID   string    // ID of event that caused this event
TraceID       string    // OpenTelemetry trace ID
SpanID        string    // OpenTelemetry span ID
```

#### Source Fields:
```go
Source        string    // Which service created event (auth-service, ride-service)
Environment   string    // Environment (development, staging, production)
```

---

### Layer 3: Event Envelopes (`envelopes/`)

**Purpose:** Standard wrapper for all events

```go
type EventEnvelope struct {
    // Event metadata (common/)
    EventID       string
    EventType     string
    EventVersion  int
    Timestamp     time.Time
    AggregateID   string
    
    // Tracing
    CorrelationID string
    CausationID   string
    TraceID       string
    SpanID        string
    
    // Source
    Source        string
    Environment   string
    
    // Payload
    Data          json.RawMessage  // Event-specific data
    
    // Reliability
    RetryCount    int
    DeadLettered  bool
}
```

---

### Layer 4: Versioning Strategy (`versions/`)

**Purpose:** Manage event schema evolution

#### Version Strategy:
- v1: Initial version (stable)
- v2: Added fields (backward compatible)
- v3: Schema breaking change (requires migration)

#### Rules:
- Old version consumers must continue to work
- New fields should have defaults
- Deprecated fields kept for compatibility
- Clear migration path documented

**Example:**
```go
// v1: Initial RideRequested
type RideRequestedV1 struct {
    RideID    string
    RiderID   string
    Pickup    Location
    Dropoff   Location
}

// v2: Added passenger count
type RideRequestedV2 struct {
    RideID         string
    RiderID        string
    PassengerCount int    // NEW - defaults to 1
    Pickup         Location
    Dropoff        Location
}
```

---

### Layer 5: Idempotency Management (`idempotency/`)

**Purpose:** Guarantee event handlers run exactly once

#### Idempotency Key Generation:
```go
// Format: event_type:aggregate_id:sequence_number
"ride.requested:ride_123:1"
"ride.started:ride_123:2"
```

#### Rules:
- Each event has unique idempotency key
- Same key = same event (duplicate detection)
- Handlers store processed keys
- Replayed events rejected if already processed

---

### Layer 6: Retry Policies (`retry/`)

**Purpose:** Define when and how to retry events

#### Retry Strategies:
1. **Exponential Backoff**
   - Attempt 1: 0 seconds
   - Attempt 2: 1 second
   - Attempt 3: 4 seconds
   - Attempt 4: 16 seconds
   - Max: 5 attempts

2. **Linear Backoff**
   - Attempt 1: 0 seconds
   - Attempt 2: 5 seconds
   - Attempt 3: 10 seconds
   - Max: 3 attempts

3. **No Retry**
   - For events that must not be retried

#### Configuration:
- Max attempts: (typically 5)
- Timeout: (typically 30s per attempt)
- Backoff type: (exponential, linear)

---

### Layer 7: Dead-Letter Queue (`dlq/`)

**Purpose:** Handle events that can't be processed

#### When Event Goes to DLQ:
- Max retries exceeded
- Unrecoverable error (bad data, schema mismatch)
- No handler found for event type

#### DLQ Processing:
1. Event logged with full context
2. Alert sent to operations team
3. Retry manually after fix
4. Analyzed for patterns

#### Critical Events (never DLQ):
- payment.succeeded ← MUST process (refund if fails)
- ride.completed ← MUST process
- wallet.debited ← MUST process

---

### Layer 8: Topics (Kafka) (`topics/`)

**Purpose:** Define Kafka topics and configurations

#### Topic Naming Convention:
```
{domain}.{entity}.{action}

Examples:
ride.lifecycle            ← All ride events
driver.location           ← Driver location updates
payment.transactions      ← Payment events
dispatch.matching         ← Dispatch events
```

#### Topic Configuration:
```go
Topic:           "ride.lifecycle"
Partitions:      10              // Scale by demand
Replication:     3               // High availability
Retention:       7 days          // Keep data 7 days
CompressType:    "snappy"        // Compression
```

---

### Layer 9: Policies (`policies/`)

**Purpose:** Define event handling policies per domain

#### Example: Ride Event Policies
- `ride.requested` → Handle within 30 seconds
- `ride.completed` → Handle within 5 minutes (critical for payment)
- `ride.cancelled` → Handle within 10 seconds

#### Policy Components:
- Max processing time
- Retry strategy
- Alert thresholds
- Required acknowledgment

---

### Layer 10: Domain-Specific Schemas (`schemas/`)

**Purpose:** JSON schemas for validation

#### Structure:
```
schemas/
├── driver/
│   ├── driver.registered.json
│   ├── driver.location.updated.json
│   └── ...
├── ride/
│   ├── ride.requested.json
│   ├── ride.completed.json
│   └── ...
├── payment/
│   ├── payment.succeeded.json
│   └── ...
└── ...
```

#### Validation:
- Events validated against schemas before publishing
- Validation errors prevent event publication
- Schema mismatch treated as critical error

---

## ✅ EVENT CATALOG SUMMARY

### Total Event Types: 100+

**By Domain:**
- User/Auth: 13 events
- Driver: 12 events
- Rider: 8 events
- Ride: 15 events
- Dispatch: 9 events
- Pricing: 8 events
- Payment: 10 events
- Wallet: 7 events
- Trip/Route: 8 events
- Safety: 8 events
- Fraud: 6 events
- Rating: 6 events
- Notifications: 4 events
- Analytics: 3 events

### Topics: 15+
- ride.lifecycle
- driver.lifecycle
- driver.location
- payment.transactions
- dispatch.matching
- pricing.calculations
- wallet.ledger
- trip.tracking
- safety.alerts
- fraud.detection
- rating.submissions
- notifications.channel
- analytics.events

### Critical Infrastructure:
- ✅ Versioning system (v1, v2, v3)
- ✅ Idempotency guarantees (exact-once processing)
- ✅ Retry mechanisms (exponential backoff)
- ✅ Dead-letter queues (unreliable event handling)
- ✅ Schema validation (prevent bad data)
- ✅ Tracing/correlation (debug complex flows)

---

## 🚨 VIOLATIONS TO PREVENT

❌ **DO NOT:**
- Define events in services/ (must use shared/contracts/events/)
- Publish events directly to Kafka (must use event-bus)
- Create new event types without catalog entry
- Skip idempotency keys
- Ignore DLQ handling
- Modify events after publishing (versioning required)

✅ **MUST:**
- Use `shared/contracts/events/catalog/` for all events
- Use `packages/event-bus/` for publishing
- Include all common metadata
- Test idempotency
- Handle retries
- Monitor DLQ

---

## 📊 NEXT STEPS

1. **Days 1-2:** Complete package and contract audits
2. **Day 3-4:** Document auth-service as reference
3. **Days 5-9:** Build services using these events
4. **Days 8-9:** Wire services through event workflows
5. **Days 9-10:** Verify all events flowing end-to-end

---

**EVENT CATALOG AUDIT COMPLETE** ✅

Repository: github.com/Abdex1/FamGo-platform  
All 100+ events documented and categorized.  
Ready for service implementation phase.

