# 🔄 VERSIONING STRATEGY: Event & Schema Evolution

**Status:** Task 2 Phase 2.4 Complete  
**Location:** shared/contracts/events/versions/  
**Current Version:** v1 (all contracts)

---

## VERSIONING PRINCIPLES

### 1. Semantic Versioning

**Apply semver to event schemas:**
- `MAJOR.MINOR.PATCH` → `v1.0.0`, `v1.1.0`, `v2.0.0`

**Version Bumping Rules:**
- `v1.0.0` → `v1.1.0`: Backward-compatible (add optional field)
- `v1.1.0` → `v2.0.0`: Breaking change (remove field, change type)

---

## CURRENT VERSIONS

### All Events: v1

**Location:** shared/contracts/events/versions/versions.go

```go
const (
    V1 = "v1"
    V2 = "v2"
)
```

**Status:**
- ✅ All contracts: v1
- ✅ All topics: v1 suffix (auth.events.v1, ride.events.v1, etc.)
- ✅ All messages: v1 version field

---

## VERSIONING TIMELINE

### PHASE 1: Establish v1 (CURRENT)
- All events start as v1
- Initial definitions locked
- Consumer adoption: all services use v1

### PHASE 2: Backward-Compatible Changes (V1.X)
- Add optional fields to v1 events
- Example: PaymentCompleted adds `payment_method_id` (optional)
- No migration needed (existing consumers ignore new field)

### PHASE 3: Breaking Changes (v2.0)
- Remove fields from v1 events
- Change field types
- Example: PaymentCompleted v2 moves `Amount` to separate currency object

**Breaking Change Process:**
1. Create v2/ directory
2. Copy v1 event → v2/
3. Modify v2 schema
4. Update MIGRATION.md
5. Deploy v1+v2 support (dual-running)
6. Migrate consumers (30-day timeline)
7. Deprecate v1

---

## EVOLVING AN EVENT SCHEMA

### Scenario 1: Add Optional Field (Non-Breaking)

**Example: Add payment method to PaymentCompleted**

**Current (v1):**
```go
type PaymentCompleted struct {
    PaymentID string `json:"payment_id"`
    Amount    float64 `json:"amount"`
    Status    string `json:"status"`
}
```

**New (still v1, backward-compatible):**
```go
type PaymentCompleted struct {
    PaymentID        string `json:"payment_id"`
    Amount           float64 `json:"amount"`
    Status           string `json:"status"`
    PaymentMethodID  string `json:"payment_method_id,omitempty"` // Optional, added in v1.1
}
```

**Steps:**
1. Add field with `omitempty` JSON tag
2. Existing consumers (v1 readers) ignore new field ✅
3. New consumers (v1.1 readers) see new field ✅
4. No version bump needed (still v1)
5. No consumer migration needed ✅

---

### Scenario 2: Remove Field or Change Type (Breaking)

**Example: Restructure amount in PaymentCompleted**

**Current (v1):**
```go
type PaymentCompleted struct {
    PaymentID string `json:"payment_id"`
    Amount    float64 `json:"amount"`
    Currency  string `json:"currency"`
}
```

**Breaking Change: Move to money object (v2)**
```go
// shared/contracts/events/payment/v2/payment_completed.go
type PaymentCompleted struct {
    PaymentID string `json:"payment_id"`
    Amount    Money `json:"amount"` // Changed from float64 to object
}

type Money struct {
    Value    float64 `json:"value"`
    Currency string  `json:"currency"`
}
```

**Steps:**
1. Create shared/contracts/events/payment/v2/ directory
2. Copy v1/payment_completed.go → v2/payment_completed.go
3. Modify v2 schema (change Amount field)
4. Update MIGRATION.md with migration guide
5. Update payment.events.v2 topic definition (if needed)
6. Deploy v1+v2 support:
   - Service publishes v2 events
   - Service consumes v1+v2 events
7. Notify consumers: "Migrate to v2 by [date]"
8. 30-day migration window
9. Deprecate v1

---

## TOPIC VERSIONING

**Topics include version in name:**

```
auth.events.v1     ← Version 1
ride.events.v1     ← Version 1
payment.events.v1  ← Version 1

// Future:
payment.events.v2  ← Version 2 (when breaking change)
```

**Benefits:**
- Easy to run dual versions (v1+v2 topics)
- Clear version in topic name
- Backward compatibility built-in

---

## CONSUMER MIGRATION PLAYBOOK

### When to Migrate

**Breaking change discovered:** 
1. Announce in team chat
2. Create GitHub issue with timeline
3. Update MIGRATION.md
4. Send notification to all consumers

### How to Migrate

**Consumer Checklist:**
- [ ] Read MIGRATION.md for specific event
- [ ] Understand old (v1) schema
- [ ] Understand new (v2) schema
- [ ] Update code to handle v2 schema
- [ ] Test with v2 events
- [ ] Deploy updated service
- [ ] Confirm v2 events working
- [ ] Confirm v1 events still working (dual-running period)
- [ ] Mark migration complete

### Timeline

- **Days 1-7:** Announce breaking change
- **Days 8-28:** Migration window (all consumers migrate)
- **Day 29:** Final verification
- **Day 30:** Deprecate v1
- **Day 31:** Remove v1 (production cleanup)

---

## GOVERNANCE

### New Version Release Process

1. **Propose:** Create GitHub issue with:
   - What's changing
   - Why it's needed
   - Impact on consumers
   - Migration path

2. **Review:** Team reviews:
   - Are breaking changes necessary?
   - Migration plan reasonable?
   - Documentation complete?

3. **Approve:** Tech Lead approves or requests changes

4. **Implement:** Create new version directory + MIGRATION.md

5. **Announce:** Notify all consumers with timeline

6. **Migrate:** Work with consumers to migrate

7. **Deprecate:** Remove old version after 30-day window

---

## ANTI-PATTERNS: DON'T DO THIS ❌

### ❌ Add New Event Fields Without JSON Tags
```go
// WRONG:
type PaymentCompleted struct {
    PaymentID string `json:"payment_id"`
    Amount    float64 // ❌ Missing JSON tag!
}
```

**Why:** Field is exported but won't serialize/deserialize properly

### ❌ Remove Fields Without Version Bump
```go
// WRONG: Removing fields in v1 (breaks existing consumers)
type PaymentCompleted struct {
    PaymentID string `json:"payment_id"`
    // Amount removed without v1→v2 migration ❌
}
```

**Why:** Consumers expecting `Amount` field will fail

### ❌ Reuse Event Names Across Domains
```go
// WRONG:
// shared/contracts/events/payment/completed.go
// shared/contracts/events/order/completed.go
// Both named "completed" ❌
```

**Why:** Naming conflicts, unclear ownership

### ❌ Undefined Version
```go
// WRONG:
// shared/contracts/events/ride/ride_updated.go
// No version info ❌
```

**Why:** Can't track evolution, breaks compatibility guarantees

---

## BEST PRACTICES: DO THIS ✅

### ✅ Always Use JSON Tags
```go
// CORRECT:
type PaymentCompleted struct {
    PaymentID string `json:"payment_id"`
    Amount    float64 `json:"amount"`
}
```

### ✅ Make Changes Backward-Compatible
```go
// CORRECT: Add optional field
type PaymentCompleted struct {
    PaymentID string `json:"payment_id"`
    Amount    float64 `json:"amount"`
    Tip       float64 `json:"tip,omitempty"` // ✅ Optional
}
```

### ✅ Create New Version for Breaking Changes
```
shared/contracts/events/payment/
├── v1/
│   └── payment_completed.go
└── v2/
    └── payment_completed.go (with breaking changes)
```

### ✅ Document Versions and Migration
- VERSIONS.md (what versions exist)
- MIGRATION.md (how to migrate)
- GitHub issues (timeline and status)

---

## CHECKING CURRENT VERSIONS

**To verify all contracts are v1:**

1. Check topics:
   ```bash
   grep -r "\.v1\|\.v2" shared/contracts/events/topics/
   # Result: All topics end in .v1 ✅
   ```

2. Check versions:
   ```bash
   cat shared/contracts/events/versions/versions.go
   # Result: V1 and V2 constants defined ✅
   ```

3. Check event type names:
   ```bash
   grep -r "EventAuth\|EventRide" shared/contracts/events/catalog/
   # Result: All events without version suffix (using .v1 in topic name) ✅
   ```

---

**Versioning Strategy:** ✅ COMPLETE & DOCUMENTED

