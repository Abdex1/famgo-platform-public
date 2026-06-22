# FamGo Target Architecture — Stage 5

**Status: PROPOSED / DESIGN, not inventory.** Everything in this document is a target to build toward, not a description of code that exists. Where a claim is grounded in something confirmed earlier (Stage 1 source docs, or the FamGo structure listing), that's stated explicitly. Where it's a new design decision being made in this document, that's stated too. Don't read any event name, saga name, or domain boundary below as "this is what FamGo's code currently does" — cross-check `docs/adoption/gap-analysis.md` for what's actually been observed.

**Inputs:** `docs/adoption/gap-analysis.md`, `docs/adoption/migration-matrix.md`, `docs/adoption/{richxcame,autofleet,microbus}/`, the FamGo structure listing, and the Stage 5 brief that introduced Deliverables 1–7 below.

---

## Deliverable 1 — Service Ownership Registry

Confirmed FamGo services (from the structure listing) are registered below. Six service directories exist but contain **zero implementation files** (`analytics-service`, `api-gateway`, `smart-pickup-service`, `subscription-service`, `voice-booking-service`, `websocket-gateway`) — these get registry entries too, marked `[STUB]`, because the registry's job is to declare ownership and contract *before* code exists, not just to catalog what's already there.

Event names (`Publishes`/`Consumes`) are **proposed**, built from the `domain.entity.action.version` convention adopted in Deliverable 5, anchored to the two contracts FamGo has actually started (`auth/login_succeeded.go`, `ride/ride_created.go`) and to richxcame's documented lifecycle hook points where a FamGo equivalent doesn't yet exist to check against.

### auth-service
```
Owner:          User Domain (Auth sub-domain)
Database:       (shared user schema — see Deliverable 7 decision on schema-per-service)
Publishes:      user.login_succeeded.v1, user.login_failed.v1, user.registered.v1
Consumes:       (none — Auth is upstream of most flows)
External APIs:  none
Sagas:          DriverOnboardingSaga (participant)
Confirmed:      packages/event-bus/contracts/auth/login_succeeded.go exists with content
```

### user-service
```
Owner:          User Domain
Database:       users_db (proposed; currently inferred shared with auth — verify)
Publishes:      user.profile_updated.v1, user.deactivated.v1
Consumes:       user.registered.v1
External APIs:  none
Sagas:          DriverOnboardingSaga (participant)
```

### driver-service
```
Owner:          Driver Domain
Database:       drivers_db (proposed)
Publishes:      driver.onboarded.v1, driver.approved.v1, driver.rejected.v1, driver.status_changed.v1
Consumes:       user.registered.v1, safety.background_check_completed.v1
External APIs:  gps-service
Sagas:          DriverOnboardingSaga (owner)
```

### ride-service
```
Owner:          Ride Domain
Database:       rides_db (proposed)
Publishes:      ride.requested.v1, ride.assigned.v1, ride.started.v1, ride.completed.v1, ride.cancelled.v1
Consumes:       dispatch.driver_assigned.v1, payment.completed.v1, payment.failed.v1
External APIs:  geo-service (gps-service), pricing-service
Sagas:          RideBookingSaga (owner), RideCompletionSaga (owner)
Confirmed:      packages/event-bus/contracts/ride/ride_created.go exists with content;
                internal/application/saga.go, events.go, event_subscribers.go already present —
                this is FamGo's most architecturally mature service and should be the
                reference implementation other services are brought up to match
```

### dispatch-service
```
Owner:          Dispatch Domain
Database:       dispatch_db (proposed)
Publishes:      dispatch.matching_started.v1, dispatch.driver_assigned.v1, dispatch.matching_failed.v1
Consumes:       ride.requested.v1, ride.cancelled.v1
External APIs:  geo-service (gps-service), driver-service
Sagas:          RideBookingSaga (participant)
Confirmed:      matching_algorithm.go, matching_service.go, match_score.go present —
                domain logic exists; event publish/consume wiring not verifiable from structure
```

### gps-service (Geo Domain target — see naming note in Deliverable 2)
```
Owner:          Geo Domain
Database:       (Redis-backed geospatial index; no relational schema expected)
Publishes:      geo.driver_location_updated.v1
Consumes:       driver.status_changed.v1
External APIs:  none
Sagas:          none
Confirmed:      platform/cache/geo/driver_index.go present
```

### wallet-service (covers both Wallet and Payment domains — see Deliverable 2 note)
```
Owner:          Payment Domain (Wallet sub-domain)
Database:       wallets_db (proposed)
Publishes:      payment.authorized.v1, payment.completed.v1, payment.failed.v1, wallet.debited.v1, wallet.credited.v1
Consumes:       ride.completed.v1, subscription.renewal_due.v1
External APIs:  packages/payment-sdk (external payment provider)
Sagas:          RideCompletionSaga (participant), SubscriptionRenewalSaga (participant)
```

### pricing-service
```
Owner:          Pricing Domain (new — not in the original 12-domain list; see Deliverable 2)
Database:       pricing_db (proposed)
Publishes:      pricing.fare_calculated.v1, pricing.surge_updated.v1
Consumes:       ride.requested.v1, geo.driver_location_updated.v1 (for surge/demand signal)
External APIs:  none
Sagas:          RideBookingSaga (participant, fare calc step)
```

### pooling-service
```
Owner:          Pooling Domain
Database:       pooling_db (proposed)
Publishes:      pooling.match_found.v1, pooling.group_formed.v1
Consumes:       ride.requested.v1, dispatch.driver_assigned.v1
External APIs:  none
Sagas:          none yet defined — candidate for a future PoolingMatchSaga
Confirmed:      cmd/, internal/, tests/ present; no api/, interfaces/, proto/ found — see gap-analysis §1
```

### safety-service
```
Owner:          Safety Domain
Database:       safety_db (proposed)
Publishes:      safety.sos_triggered.v1, safety.background_check_completed.v1, safety.incident_reported.v1
Consumes:       ride.started.v1, ride.cancelled.v1
External APIs:  none
Sagas:          DriverOnboardingSaga (participant)
```

### fraud-service
```
Owner:          Safety Domain (Fraud sub-domain — folded in per the richxcame mapping in Deliverable 3)
Database:       fraud_db (proposed)
Publishes:      safety.fraud_alert_raised.v1, safety.user_suspended.v1
Consumes:       payment.authorized.v1, payment.failed.v1, ride.completed.v1
External APIs:  none
Sagas:          none yet defined
```

### notification-service
```
Owner:          Notification Domain
Database:       (likely none owned — delivery-tracking only; proposed notifications_db if read receipts needed)
Publishes:      notification.sent.v1, notification.delivery_failed.v1
Consumes:       ride.requested.v1, ride.assigned.v1, ride.started.v1, ride.completed.v1, ride.cancelled.v1,
                payment.completed.v1, payment.failed.v1, safety.sos_triggered.v1
External APIs:  none (push/SMS/email providers live behind this service)
Sagas:          RideBookingSaga, RideCompletionSaga, DriverOnboardingSaga, SubscriptionRenewalSaga
                (participant in all four — Notification is a broad fan-in consumer by design)
Confirmed:      Currently only 2 files exist (cmd/api/main.go, domain/entities/notification.go) —
                this registry entry is the build target, not the current state; see gap-analysis §1
```

### analytics-service `[STUB]`
```
Owner:          Analytics Domain
Database:       analytics_db (proposed, likely read-model/OLAP-shaped, not transactional)
Publishes:      (none expected — Analytics is a terminal consumer)
Consumes:       ride.*, payment.*, dispatch.*, safety.fraud_alert_raised.v1 (broad fan-in, all domains)
External APIs:  none
Sagas:          none
Confirmed:      Directory exists, zero files — full build-out required
```

### subscription-service `[STUB]`
```
Owner:          Subscription Domain
Database:       subscriptions_db (proposed)
Publishes:      subscription.created.v1, subscription.renewal_due.v1, subscription.cancelled.v1
Consumes:       wallet.debited.v1, payment.failed.v1
External APIs:  none
Sagas:          SubscriptionRenewalSaga (owner)
Confirmed:      Directory exists, zero files — full build-out required
```

### smart-pickup-service `[STUB]`
```
Owner:          Dispatch Domain (Smart Pickup sub-domain) — or Geo Domain; needs an explicit decision,
                flagged here rather than silently assumed
Database:       none expected (likely a stateless optimization layer over gps-service + dispatch-service)
Publishes:      dispatch.pickup_point_suggested.v1 (proposed)
Consumes:       geo.driver_location_updated.v1, ride.requested.v1
External APIs:  gps-service
Sagas:          none
Confirmed:      Directory exists, zero files — full build-out required; domain ownership undecided
```

### voice-booking-service `[STUB]`
```
Owner:          Ride Domain (alternate booking channel, not a new domain)
Database:       none expected (translates voice input into a standard ride.requested.v1)
Publishes:      (none directly — should call ride-service / publish ride.requested.v1 via the same
                contract every other booking channel uses, not a parallel one)
Consumes:       none
External APIs:  ride-service
Sagas:          RideBookingSaga (entry point only, not a participant in saga logic itself)
Confirmed:      Directory exists, zero files — full build-out required
```

### websocket-gateway `[STUB]`
```
Owner:          Cross-cutting (Realtime) — not a domain owner, a delivery channel for Dispatch/Ride/
                Notification domain events to connected clients
Database:       none (in-memory connection hub + Redis pub/sub, per richxcame's pattern — see
                docs/adoption/richxcame/patterns.md)
Publishes:      none (terminal fan-out to WebSocket clients, not back onto the event bus)
Consumes:       ride.*, dispatch.driver_assigned.v1, notification.sent.v1
External APIs:  none
Sagas:          none
Confirmed:      Directory exists, zero files. Currently fragmented across
                services/ride-service/internal/transport/websocket.go, gateway/websocket/, and
                packages/websocket-sdk/ — gap-analysis §5 priority item; this stub should become
                the single owner of all WebSocket fan-out
```

### api-gateway `[STUB]`
```
Owner:          Cross-cutting (Platform) — not a domain owner
Database:       none
Publishes:      none
Consumes:       none (routes/auth-checks requests, doesn't participate in the event bus)
External APIs:  all services (as a reverse proxy)
Sagas:          none
Confirmed:      Directory exists, zero files. Gateway concerns currently live in gateway/ (Kong config,
                JWT proxy, handlers.go, middleware.go) and infra/kong/ — see gap-analysis §6
                duplication note. Recommend retiring this empty stub in favor of gateway/, rather
                than building a third location.
```

---

## Deliverable 2 — Canonical Domain Model

The brief's 12 domains are adopted with two amendments, stated explicitly:

- **Pricing** is added as its own domain — FamGo's structure has a dedicated `pricing-service` that richxcame splits out from promos/fares too; folding it into Ride or Dispatch would hide a real bounded context.
- **Fraud** is folded into **Safety** rather than kept separate, matching the richxcame-to-FamGo mapping already used in Deliverable 3 below (`Fraud → Safety`) — keeping one mapping consistent rather than having Deliverable 2 and Deliverable 3 disagree.

| Domain | Aggregates | Key Entities | Value Objects | Core Events | Owning Service(s) |
|---|---|---|---|---|---|
| **User** | User Account | User, Profile | Email, Phone, Role | `user.registered.v1`, `user.login_succeeded.v1`, `user.profile_updated.v1` | `auth-service`, `user-service` |
| **Driver** | Driver Profile | Driver, Vehicle, DocumentSet | LicenseNumber, VehiclePlate, ApprovalStatus | `driver.onboarded.v1`, `driver.approved.v1`, `driver.status_changed.v1` | `driver-service` |
| **Ride** | Ride | Ride, RideRating, FavoriteLocation | Coordinates, RideStatus, FareAmount | `ride.requested.v1`, `ride.assigned.v1`, `ride.started.v1`, `ride.completed.v1`, `ride.cancelled.v1` | `ride-service`, `voice-booking-service` (channel only) |
| **Dispatch** | MatchRequest | MatchRequest, MatchCandidate, MatchScore | MatchScore (value), DistanceKm | `dispatch.matching_started.v1`, `dispatch.driver_assigned.v1`, `dispatch.matching_failed.v1`, `dispatch.pickup_point_suggested.v1` | `dispatch-service`, `smart-pickup-service` |
| **Pricing** | FareQuote | FareQuote, SurgeZone | FareAmount, SurgeMultiplier | `pricing.fare_calculated.v1`, `pricing.surge_updated.v1` | `pricing-service` |
| **Payment** | Payment, Wallet | Payment, WalletTransaction | Money, PaymentStatus | `payment.authorized.v1`, `payment.completed.v1`, `payment.failed.v1`, `wallet.debited.v1`, `wallet.credited.v1` | `wallet-service` |
| **Safety** | SafetyIncident, FraudAlert | SafetyIncident, FraudAlert, RiskProfile | Severity, AlertType (richxcame's 6-type taxonomy — see `docs/adoption/richxcame/features.md`) | `safety.sos_triggered.v1`, `safety.fraud_alert_raised.v1`, `safety.background_check_completed.v1` | `safety-service`, `fraud-service` |
| **Notification** | NotificationJob | NotificationJob, DeliveryReceipt | Channel (push/SMS/email), DeliveryStatus | `notification.sent.v1`, `notification.delivery_failed.v1` | `notification-service` |
| **Analytics** | (read models only — no transactional aggregates by design) | RevenueSnapshot, DemandZone | DateRange | (none published — terminal consumer) | `analytics-service` |
| **Geo** | DriverLocation | DriverLocation | Coordinates, ETA | `geo.driver_location_updated.v1` | `gps-service` |
| **Subscription** | Subscription | Subscription, BillingCycle | RenewalDate, PlanTier | `subscription.created.v1`, `subscription.renewal_due.v1`, `subscription.cancelled.v1` | `subscription-service` |
| **Pooling** | PoolGroup | PoolGroup, PoolMatch | GroupCapacity | `pooling.match_found.v1`, `pooling.group_formed.v1` | `pooling-service` |

**Cross-cutting, not domains:** Realtime delivery (`websocket-gateway`) and the API edge (`api-gateway`/`gateway/`) are deliberately excluded from this table — they carry events and requests for the domains above but don't own a bounded context of their own. Listing them as "domains" would blur ownership; they're registered in Deliverable 1 as cross-cutting instead.

---

## Deliverable 3 — Repository-to-Domain Mapping

### Richxcame → FamGo Domain

| richxcame Component | FamGo Domain | Decision | Grounding |
|---|---|---|---|
| Auth | User | ADAPT | Stage 1 service inventory; richxcame's Auth is HTTP-only, FamGo's auth-service already has gRPC + DDD layering — adapt the *capability list* (JWT, role-based access, key rotation), not the transport |
| Ride | Ride | ADAPT | richxcame's Rides is synchronous-call-driven; FamGo's ride-service already has saga/event files richxcame's doesn't — adapt the *lifecycle state machine* (`requested→accepted→in_progress→completed`, alt `cancelled`), not the synchronous mechanism |
| Geo | Geo | ADAPT | richxcame's Redis GeoSpatial pattern maps directly onto FamGo's `platform/cache/geo/driver_index.go` |
| Realtime | Dispatch* | ADAPT | *Mapped to Dispatch per the brief's table, but per Deliverable 1/2 above, FamGo's realtime concern is cross-cutting (`websocket-gateway`), not a Dispatch-domain aggregate — flagging this disagreement rather than silently following the brief; recommend treating Realtime as infrastructure serving Dispatch/Ride/Notification, not as owned by Dispatch |
| Payment | Payment | ADAPT | richxcame's Stripe integration pattern is reusable; its **missing webhook signature verification is an anti-pattern, not a feature** — see `docs/adoption/richxcame/anti-patterns.md` |
| Notification | Notification | ADAPT | richxcame's lifecycle-hook trigger list is the most useful single artifact to carry over (see Deliverable 1, notification-service consumes list) |
| Analytics | Analytics | ADAPT | richxcame's endpoint breadth (revenue, promo performance, heat maps, demand zones) is a feature checklist for FamGo's currently-empty analytics-service |
| Fraud | Safety | ADAPT | Folded per Deliverable 2's domain model amendment; richxcame's 6-type/4-severity taxonomy is directly reusable as a starting enum |

### Autofleet → FamGo Domain

| Autofleet Component | FamGo Domain | Decision | Grounding |
|---|---|---|---|
| Ride Lifecycle | Ride | ADOPT WORKFLOW ONLY | Per Stage 1, Autofleet's repo shows only the rider-facing screen sequence and the ASAP/scheduled request distinction — there is no backend logic to adopt beyond this workflow shape; see `docs/adoption/autofleet/features.md` |
| Driver Lifecycle | Driver | ADOPT WORKFLOW ONLY | Not actually present in the Autofleet repo (driver app is out-of-repo per Stage 1) — included in the brief's table but there is **no source material to back this row**; flagged honestly rather than silently treated as grounded |
| Support Flow | Safety | ADOPT WORKFLOW ONLY | Autofleet's own support pages are empty stubs (per Stage 1) — there is **no workflow to adopt here either**; this row should be read as "no source available," consistent with `docs/adoption/autofleet/anti-patterns.md` |

### Microbus → FamGo Domain

| Microbus Component | FamGo Domain | Decision | Grounding |
|---|---|---|---|
| Contracts | Contracts (cross-cutting, `packages/event-bus/contracts/`) | ADOPT (concept) | Type-aliasing discipline and `manifest.yaml`-style drift prevention — see `docs/adoption/microbus/patterns.md` Pattern 6 |
| Sagas | Event Bus (cross-cutting, `platform/saga/`) | ADOPT (concept) | Typed-transition workflow graph model — see `docs/adoption/microbus/patterns.md` Pattern 4; this is the primary design input for Deliverable 6 below |
| Outbox | Shared Infrastructure (`platform/outbox/`) | ADOPT (concept, exceeding the source) | Microbus itself has no general-purpose outbox (only workflow-step durability) — FamGo's outbox work goes further than Microbus's own implementation; see `docs/adoption/microbus/migration-plan.md` Phase 3 |
| DLQ | Messaging (`packages/kafka-sdk/internal/dlq/`) | ADOPT (concept, with a fix) | Microbus's parked-backlog has no auto-expiry — FamGo's DLQ must add one; see `docs/adoption/microbus/anti-patterns.md` #2 |

### OpenRide → FamGo Domain

**No Stage 1 research exists for OpenRide.** The mapping below is carried over from the Stage 5 brief as a *placeholder structure only* — it cannot be marked ADOPT/ADAPT/REJECT with any grounding, and is included here so the table shape is ready the moment Stage 1 research for OpenRide is actually run.

| OpenRide Component | FamGo Domain | Decision |
|---|---|---|
| Rider UX | Mobile (`apps/mobile`, `apps/rider-web`) | **UNGROUNDED — do not act on this row until Stage 1 research exists** |
| Driver UX | Mobile (`apps/driver-web`) | **UNGROUNDED — do not act on this row until Stage 1 research exists** |
| Admin UX | Admin (`apps/admin-dashboard`) | **UNGROUNDED — do not act on this row until Stage 1 research exists** |

### RidePy → FamGo Domain

**No Stage 1 research exists for RidePy either.** Same caveat as OpenRide above.

| RidePy Component | FamGo Domain | Decision |
|---|---|---|
| Pooling | Pooling (`pooling-service`) | **UNGROUNDED — do not act on this row until Stage 1 research exists** |
| Matching | Dispatch (`dispatch-service`) | **UNGROUNDED — do not act on this row until Stage 1 research exists** |
| Route Optimization | Geo (`gps-service`) | **UNGROUNDED — do not act on this row until Stage 1 research exists** |

---

## Deliverable 4 — Service Adoption Matrix

Influence ratings (`High`/`Medium`/`Low`/`None`) reflect **how much usable material each source repo offers for that FamGo service**, combined with the priority levels already established in `docs/adoption/gap-analysis.md` §9. OpenRide and RidePy are rated `None` everywhere except where the gap-analysis already flagged a structural reason to expect future relevance — rated `None` rather than guessed, since there's no source material to rate them by.

| FamGo Service | Richxcame | Autofleet | Microbus | OpenRide | RidePy | Priority | Grounding |
|---|---|---|---|---|---|---|---|
| auth-service | High | None | Low | None | None | P1 | richxcame Auth is the closest comparable; gap-analysis rates this service Implemented already |
| user-service | Medium | None | Low | None | None | P2 | richxcame's User model referenced across services |
| driver-service | Medium | Low | Low | None | None | P2 | richxcame has no standalone driver-management service distinct from Admin; Autofleet's driver app is out-of-repo |
| ride-service | High | Medium | High | None | Low | **P0** | richxcame's full lifecycle docs + Microbus's saga model are the two strongest inputs; gap-analysis already flags ride-service as FamGo's most mature service and the saga reference implementation |
| dispatch-service | Low | None | Medium | None | None | P1 | richxcame has no standalone matching service to compare against; Microbus's workflow-graph model is relevant to dispatch's internal state machine |
| gps-service | High | Low | None | None | None | P2 | richxcame's Geo service is a close structural match |
| pricing-service | Medium | None | None | None | None | P2 | richxcame's Promos (fare calc half) is the only comparable source |
| pooling-service | None | None | None | None | None | P2 | No source repo documents a comparable capability — confirmed in gap-analysis §1: "richxcame has no pooling/shared-ride concept at all" |
| wallet-service | High | Low | Low | None | None | **P0** | richxcame's Payments + documented webhook-signature gap make this a high-grounding, high-risk priority — see gap-analysis §2 |
| safety-service | Medium | None | None | None | None | P2 | richxcame has no Safety-equivalent service (only Fraud); low transferable material |
| fraud-service | High | None | Low | None | None | P1 | richxcame's Fraud taxonomy is directly reusable |
| notification-service | High | None | Medium | None | None | **P0** | gap-analysis §1 already flags notification-service as High priority (currently 2 files); richxcame's lifecycle-hook list plus Microbus's pub/sub-over-direct-call pattern are both directly applicable |
| analytics-service | High | None | Low | None | None | P1 | richxcame's endpoint breadth is the primary usable reference; service is currently empty |
| websocket-gateway | High | None | Medium | None | None | **P0** | gap-analysis §5 already flags this High priority; richxcame's standalone Realtime service is the closest structural match |
| subscription-service | None | None | Low | None | None | P2 | No source repo documents subscriptions; Microbus's saga model is relevant to the renewal workflow only |
| smart-pickup-service | None | None | None | None | None | P3 | No comparable capability in any source repo |
| voice-booking-service | None | None | None | None | None | P3 | No comparable capability in any source repo |
| api-gateway | High | None | Low | None | None | P2 | richxcame's Kong setup is directly comparable; FamGo already has gateway/ + infra/kong/ partially built |

---

## Deliverable 5 — Event Contract Catalog

**Convention:** `domain.entity.action.version` (lowercase, snake_case actions). This matches the one populated example already in FamGo's tree (`ride/ride_created.go` — note: the *file* uses `ride_created`, the catalog below standardizes the wire-format event name as `ride.requested.v1` for the request-time event specifically, to avoid the ambiguity between "a Ride row was created in the DB" and "a rider requested a ride," which are different moments. This is a deliberate naming refinement, flagged so it doesn't read as a contradiction.)

| Event | Producer | Consumers (proposed) |
|---|---|---|
| `user.registered.v1` | auth-service | user-service, driver-service |
| `user.login_succeeded.v1` | auth-service | (audit/analytics only) |
| `user.login_failed.v1` | auth-service | safety-service (fraud signal) |
| `user.profile_updated.v1` | user-service | — |
| `driver.onboarded.v1` | driver-service | notification-service |
| `driver.approved.v1` | driver-service | notification-service, dispatch-service |
| `driver.rejected.v1` | driver-service | notification-service |
| `driver.status_changed.v1` | driver-service | gps-service, dispatch-service |
| `ride.requested.v1` | ride-service | dispatch-service, pricing-service, notification-service |
| `ride.assigned.v1` | ride-service | notification-service, websocket-gateway |
| `ride.started.v1` | ride-service | safety-service, notification-service, websocket-gateway |
| `ride.completed.v1` | ride-service | wallet-service, fraud-service, notification-service, analytics-service |
| `ride.cancelled.v1` | ride-service | dispatch-service, notification-service, analytics-service |
| `dispatch.matching_started.v1` | dispatch-service | websocket-gateway |
| `dispatch.driver_assigned.v1` | dispatch-service | ride-service, websocket-gateway |
| `dispatch.matching_failed.v1` | dispatch-service | ride-service, notification-service |
| `dispatch.pickup_point_suggested.v1` | smart-pickup-service | dispatch-service |
| `pricing.fare_calculated.v1` | pricing-service | ride-service |
| `pricing.surge_updated.v1` | pricing-service | (dashboard/analytics consumers) |
| `payment.authorized.v1` | wallet-service | ride-service, fraud-service |
| `payment.completed.v1` | wallet-service | ride-service, notification-service |
| `payment.failed.v1` | wallet-service | ride-service, notification-service, fraud-service |
| `wallet.debited.v1` | wallet-service | subscription-service, analytics-service |
| `wallet.credited.v1` | wallet-service | analytics-service |
| `safety.sos_triggered.v1` | safety-service | notification-service, websocket-gateway |
| `safety.background_check_completed.v1` | safety-service | driver-service |
| `safety.incident_reported.v1` | safety-service | notification-service |
| `safety.fraud_alert_raised.v1` | fraud-service | notification-service, analytics-service |
| `safety.user_suspended.v1` | fraud-service | auth-service, notification-service |
| `notification.sent.v1` | notification-service | analytics-service |
| `notification.delivery_failed.v1` | notification-service | (retry/alerting consumers) |
| `geo.driver_location_updated.v1` | gps-service | dispatch-service, pricing-service, websocket-gateway |
| `subscription.created.v1` | subscription-service | notification-service |
| `subscription.renewal_due.v1` | subscription-service | wallet-service |
| `subscription.cancelled.v1` | subscription-service | notification-service |
| `pooling.match_found.v1` | pooling-service | ride-service |
| `pooling.group_formed.v1` | pooling-service | notification-service |

This becomes the working content for `docs/adoption/contracts-catalog.md` (see file of that name, included alongside this document).

---

## Deliverable 6 — Saga Catalog

Saga *shape* (steps, compensation) is proposed design; the *engine* it should run on is the typed-transition workflow-graph model from `docs/adoption/microbus/patterns.md` Pattern 4, targeting FamGo's currently-empty `platform/saga/`.

### RideBookingSaga (owner: ride-service)
```
Steps:
  1. ride.requested.v1 emitted (ride-service)
  2. pricing.fare_calculated.v1 (pricing-service) — informational, not blocking
  3. dispatch.driver_assigned.v1 (dispatch-service) — on success, continue;
     on dispatch.matching_failed.v1, go to compensation
  4. notification.sent.v1 to rider + driver (notification-service)

Compensation (error transition):
  - On matching_failed: ride.cancelled.v1, notify rider, release any pricing hold
  - On payment pre-auth failure (if pre-auth is added later): ride.cancelled.v1, notify rider

Participants: ride-service, dispatch-service, pricing-service, notification-service
```

### RideCompletionSaga (owner: ride-service)
```
Steps:
  1. ride.completed.v1 emitted (ride-service)
  2. payment.authorized.v1 → payment.completed.v1 (wallet-service)
  3. wallet.debited.v1 (rider) / wallet.credited.v1 (driver) (wallet-service)
  4. notification.sent.v1 — receipt (notification-service)

Compensation (error transition):
  - On payment.failed.v1: flag ride for manual payment retry, notify rider,
    do NOT reverse ride.completed.v1 (the ride happened; only payment failed) —
    this is a deliberate FamGo-specific decision, since reversing a completed ride
    to "undo" a payment failure would be a worse UX than retrying payment async

Participants: ride-service, wallet-service, notification-service
```

### DriverOnboardingSaga (owner: driver-service)
```
Steps:
  1. user.registered.v1 (driver applicant) (auth-service/user-service)
  2. driver.onboarded.v1 — document submission (driver-service)
  3. safety.background_check_completed.v1 (safety-service) — on pass, continue;
     on fail, go to compensation
  4. driver.approved.v1 (driver-service)
  5. notification.sent.v1 (notification-service)

Compensation (error transition):
  - On background check failure: driver.rejected.v1, notify applicant with reason
    (subject to FamGo's own policy on what reason detail is disclosable)

Participants: user-service, driver-service, safety-service, notification-service
```

### SubscriptionRenewalSaga (owner: subscription-service)
```
Steps:
  1. subscription.renewal_due.v1 (subscription-service)
  2. wallet.debited.v1 (wallet-service) — on success, continue;
     on payment.failed.v1, go to compensation
  3. subscription stays active; notification.sent.v1 — renewal confirmation

Compensation (error transition):
  - On payment.failed.v1: subscription.cancelled.v1 or grace-period flag
    (FamGo policy decision needed: hard-cancel vs. grace period — not resolved here),
    notification.sent.v1 — payment failure / action needed

Participants: subscription-service, wallet-service, notification-service
```

**Note on saga count:** The brief lists exactly these four. Two more are visible as candidates from Deliverable 1/2 above but are **not yet specified** here, to avoid inventing saga detail beyond what's been asked for: a `PoolingMatchSaga` (pooling-service + dispatch-service + ride-service) and a `FraudInvestigationSaga` (fraud-service + safety-service + auth-service, for the suspend/reinstate workflow richxcame documents). Flagged for Stage 9 (Saga Catalog, per the 12-stage roadmap) rather than defined here.

---

## Deliverable 7 — Infrastructure Adoption Matrix

| Capability | Source | Decision | Grounding |
|---|---|---|---|
| Kafka | FamGo | Keep | `packages/kafka-sdk/`, `infra/kafka/` already committed to; richxcame's NATS is provisioned-but-disabled and not a stronger reference — see gap-analysis §3 |
| Kong | FamGo | Keep (but de-duplicate) | FamGo has Kong in **two** locations (`gateway/kong/`, `infra/kong/`) — "keep Kong," not "keep both copies"; see gap-analysis §6 |
| Outbox | Microbus | Adopt (concept only) | `platform/outbox/` is currently empty; Microbus's own outbox is narrower than what FamGo needs (workflow-step durability only) — adopt the *intent*, design a general-purpose one; see `docs/adoption/microbus/migration-plan.md` Phase 3 |
| DLQ | Microbus | Adopt (with a fix) | `packages/kafka-sdk/internal/dlq/` is currently empty; must add the expiry/poison-message threshold Microbus's own design admits it lacks |
| Saga Coordinator | Microbus | Adopt (concept only) | `platform/saga/` is currently empty; typed-transition-graph model from Microbus's Foreman is the design reference, re-implemented in Go against Postgres, not NATS |
| Realtime Architecture | Richxcame | Adapt | richxcame's standalone Realtime service (hub stats, chat history, internal broadcast) is the closest structural match for FamGo's currently-empty `websocket-gateway`; **richxcame's own documented unauthenticated-internal-route anti-pattern must NOT be carried over** — see `docs/adoption/richxcame/anti-patterns.md` #1 |
| Pooling Algorithms | RidePy | Later | **No Stage 1 research exists for RidePy** — "Later" here means "blocked on Stage 1 research being run for RidePy," not "deprioritized despite having material available" |
| Resilience (circuit breaker, retry, timeout, bulkhead) | Richxcame + Microbus (blended) | Adapt | richxcame's per-service override map + DB-specific breaker is the stronger *working* reference (real running system); Microbus's bounded-retry primitive and exponential probe schedule are the stronger *design* reference for FamGo's empty `platform/resilience/{fallback,hedging,ratelimit}/` |
| Secrets management | Richxcame | Adapt | Pluggable-provider model; consolidates FamGo's currently-duplicated `security/vault/`, `security/secrets/`, `platform/config/vault/`, `packages/vault-sdk/` — see gap-analysis §6 |
| Observability sampling strategy | Richxcame | Adopt | Concrete dev/staging/prod trace-sampling table is low-risk and directly portable |

---

## What This Document Deliberately Does Not Do

Per the Stage 5 brief, this document stops short of regenerating `migration-matrix.xlsx`, `contracts-catalog.md` (full standalone file), `service-adoption-matrix.md` (full standalone file), `domain-model-mapping.md`, `saga-catalog.md`, `realtime-adoption-plan.md`, `observability-adoption-plan.md`, `security-adoption-plan.md`, and `repository-adoption-packages/`. Those are Stage 12 (Final Migration Matrix) and the supporting Stage 6–11 artifacts, built **on top of** this document, not duplicated inside it. The two standalone files explicitly named in the brief as immediate Deliverable-5/6 outputs (`contracts-catalog.md`, drawn from Deliverable 5; a dedicated `saga-catalog.md`, drawn from Deliverable 6) are included alongside this document so they exist as independently linkable files, per the brief's stated target paths.
