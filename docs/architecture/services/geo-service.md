# Geo Service (gps-service)

## Responsibilities

- Ingest and serve driver locations in real time.
- Maintain geospatial indexes for nearby-driver search (Redis Geo or PostGIS).
- Compute distance and ETA helpers for dispatch and ride domains.
- Publish location updates for surge/pricing and realtime fan-out.

## Database / Index Stores

| Store | Purpose |
|---|---|
| Redis `drivers:geo:index` | Live driver positions keyed by geohash |
| `driver_location_history` (PostgreSQL, optional) | Archival trail for analytics and fraud |
| In-memory ETA cache | Short-lived route/distance estimates |

## Publishes

| Event | Topic | Consumers |
|---|---|---|
| Driver location updated | `geo.driver_location_updated.v1` | dispatch-service, pricing-service, realtime gateway |
| Driver went offline | `driver.status_changed.v1` | dispatch-service, notification-service |

## Consumes

| Event | Source | Action |
|---|---|---|
| `driver.status_changed.v1` | driver-service | Add/remove driver from geo index |
| `driver.onboarded.v1` | driver-service | Seed initial driver profile in geo layer |

## External Dependencies

- **Redis** — primary geospatial index and pub/sub for live updates
- **driver-service** — driver eligibility and vehicle metadata

## Notes

FamGo maps this domain to `gps-service`. Dispatch-service calls `FindNearbyDrivers` gRPC during matching; event stream decouples surge and analytics consumers.
