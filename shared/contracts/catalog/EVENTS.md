# 📨 EVENT REGISTRY: Consolidated Event Catalog

**Status:** Task 2 Phase 2.1 Complete  
**Location:** shared/contracts/events/catalog/  
**Total Events:** 20 unique, no duplicates

---

## EVENT REGISTRY BY DOMAIN

### AUTH DOMAIN (8 Events)

| Event Name | Type | Version | Topic | Publisher | Consumers | Status |
|-----------|------|---------|-------|-----------|-----------|--------|
| auth.login.succeeded | Auth | v1 | auth.events.v1 | auth-service | [all-services] | ✅ |
| auth.login.failed | Auth | v1 | auth.events.v1 | auth-service | [audit, notification] | ✅ |
| auth.logout.succeeded | Auth | v1 | auth.events.v1 | auth-service | [all-services] | ✅ |
| auth.token.refreshed | Auth | v1 | auth.events.v1 | auth-service | [internal] | ✅ |
| auth.session.revoked | Auth | v1 | auth.events.v1 | auth-service | [all-services] | ✅ |
| auth.otp.requested | Auth | v1 | auth.events.v1 | auth-service | [notification] | ✅ |
| auth.otp.verified | Auth | v1 | auth.events.v1 | auth-service | [internal] | ✅ |

### RIDE DOMAIN (5 Events)

| Event Name | Type | Version | Topic | Publisher | Consumers | Status |
|-----------|------|---------|-------|-----------|-----------|--------|
| ride.requested | Ride | v1 | ride.events.v1 | ride-service | [dispatch, notification] | ✅ |
| ride.accepted | Ride | v1 | ride.events.v1 | ride-service | [driver, passenger] | ✅ |
| ride.cancelled | Ride | v1 | ride.events.v1 | ride-service | [notification, analytics] | ✅ |
| ride.started | Ride | v1 | ride.events.v1 | ride-service | [gps, safety] | ✅ |
| ride.completed | Ride | v1 | ride.events.v1 | ride-service | [payment, analytics, fraud] | ✅ |

### DRIVER DOMAIN (3 Events)

| Event Name | Type | Version | Topic | Publisher | Consumers | Status |
|-----------|------|---------|-------|-----------|-----------|--------|
| driver.online | Driver | v1 | driver.events.v1 | driver-service | [dispatch, analytics] | ✅ |
| driver.offline | Driver | v1 | driver.events.v1 | driver-service | [dispatch, analytics] | ✅ |
| driver.location.updated | Driver | v1 | driver.events.v1 | gps-service | [dispatch, ride, analytics] | ✅ |

### PAYMENT DOMAIN (3 Events)

| Event Name | Type | Version | Topic | Publisher | Consumers | Status |
|-----------|------|---------|-------|-----------|-----------|--------|
| payment.authorized | Payment | v1 | payment.events.v1 | payment-service | [internal] | ✅ |
| payment.captured | Payment | v1 | payment.events.v1 | payment-service | [wallet, analytics] | ✅ |
| payment.failed | Payment | v1 | payment.events.v1 | payment-service | [notification, fraud] | ✅ |

### FRAUD DOMAIN (1 Event)

| Event Name | Type | Version | Topic | Publisher | Consumers | Status |
|-----------|------|---------|-------|-----------|-----------|--------|
| fraud.detected | Fraud | v1 | fraud.events.v1 | fraud-service | [payment, support, notification] | ✅ |

### SAFETY DOMAIN (1 Event)

| Event Name | Type | Version | Topic | Publisher | Consumers | Status |
|-----------|------|---------|-------|-----------|-----------|--------|
| safety.sos.triggered | Safety | v1 | safety.events.v1 | safety-service | [notification, support] | ✅ |

---

## CROSS-REFERENCE: Events by Consumer

### Per-Service Consumption

**auth-service (produces):**
- auth.login.succeeded
- auth.login.failed
- auth.logout.succeeded
- auth.token.refreshed
- auth.session.revoked
- auth.otp.requested
- auth.otp.verified

**ride-service (produces):**
- ride.requested
- ride.accepted
- ride.cancelled
- ride.started
- ride.completed

**driver-service (produces):**
- driver.online
- driver.offline

**gps-service (produces):**
- driver.location.updated

**payment-service (produces):**
- payment.authorized
- payment.captured
- payment.failed

**fraud-service (produces):**
- fraud.detected

**safety-service (produces):**
- safety.sos.triggered

---

## QUALITY GATES VERIFICATION

✅ All events have:
- [x] Unique name in catalog
- [x] Version field (v1)
- [x] Topic assignment
- [x] Publisher assigned
- [x] Consumer(s) listed
- [x] Schema reference

✅ No events are:
- [x] Duplicated
- [x] Orphaned (no consumer)
- [x] Unversioned
- [x] Without publisher

---

**Registry Status:** ✅ COMPLETE & VERIFIED

