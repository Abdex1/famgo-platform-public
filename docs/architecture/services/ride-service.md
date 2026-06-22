# Ride Service

## Responsibilities

- Owns the ride lifecycle from request through completion or cancellation.
- Persists ride state, pickup/dropoff coordinates, rider identity, and fare metadata.
- Publishes lifecycle events that downstream domains (dispatch, payment, notification) react to.
- Exposes APIs for riders and drivers to create, track, and transition rides.

## Database Tables (proposed)

| Table | Purpose |
|---|---|
| `rides` | Core ride record: status, rider_id, pickup/dropoff, vehicle type, timestamps |
| `ride_events` | Append-only audit of state transitions |
| `ride_ratings` | Post-trip ratings linked to completed rides |
| `scheduled_rides` | Future-dated ride requests polled by scheduler workers |

## Publishes

| Event | Topic | Consumers |
|---|---|---|
| Ride created | `ride.created.v1` | dispatch-service, pricing-service, fraud-service |
| Ride assigned | `ride.assigned.v1` | notification-service, realtime gateway |
| Ride started | `ride.started.v1` | notification-service, payment-service |
| Ride completed | `ride.completed.v1` | payment-service, wallet-service, analytics |
| Ride cancelled | `ride.cancelled.v1` | dispatch-service, payment-service, notification-service |

## Consumes

| Event | Source | Action |
|---|---|---|
| `dispatch.driver.assigned.v1` | dispatch-service | Transition ride to assigned; notify rider |
| `payment.completed.v1` | payment-service | Mark fare settled on completion flows |
| `payment.failed.v1` | payment-service | Hold or cancel ride when capture fails |

## External Dependencies

- **gps-service** — pickup ETA, nearby driver context
- **pricing-service** — fare estimate at request time

## Notes

Reference implementation for event-driven ride orchestration in FamGo. Saga hooks live under `internal/application/`.
