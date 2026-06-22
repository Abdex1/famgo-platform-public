# Service Inventory — richxcame/ride-hailing

Source: https://github.com/richxcame/ride-hailing (branch: `main`)

## Summary

The platform is a Go-based microservices backend for a ride-hailing product (passenger, driver, and admin surfaces). The README advertises **14 services**; only **13 have a corresponding entry point under `cmd/` and a service definition in `docker-compose.yml`**. The 14th (Negotiation, fare negotiation) is documented in the README's service table but has no observable `cmd/`, `internal/`, Dockerfile build target, or compose entry — it appears to be a planned/aspirational service rather than a shipped one as of the inspected commit. This is flagged here as a documentation/implementation discrepancy rather than resolved or removed.

## Implemented Services (present in `cmd/`, `internal/`, and `docker-compose.yml`)

| Service | Port | Purpose | Primary Dependencies |
|---|---|---|---|
| Auth | 8081 | JWT authentication, registration/login, profile management, role-based access control (rider/driver/admin), JWT key rotation support | PostgreSQL, Redis |
| Rides | 8082 | Core ride lifecycle (request → accept → start → complete/cancel), scheduled rides, surge pricing lookup, rating | PostgreSQL, Redis, calls Promos service, calls ML ETA service (optional) |
| Geo | 8083 | Driver location ingestion and lookup, distance/ETA helper utility, geospatial matching | Redis (GeoSpatial index) |
| Payments | 8084 | Wallets, ride payment capture, refunds, Stripe integration and webhook intake | PostgreSQL, Redis, Stripe |
| Notifications | 8085 | Multi-channel notification delivery (push/SMS/email), ride lifecycle notification hooks, admin bulk broadcast | PostgreSQL, Redis, Firebase, Twilio, SMTP |
| Realtime | 8086 | WebSocket gateway, in-ride chat history, driver location relay, internal broadcast fan-out | PostgreSQL, Redis |
| Mobile | 8087 | Mobile-optimized API façade (ride history, receipts, favorites, profile) — reuses Rides/Favorites logic | PostgreSQL |
| Admin | 8088 | Back-office dashboard, user/driver governance (suspend/activate/approve/reject), ride statistics | PostgreSQL |
| Promos | 8089 | Ride types, fare calculation, promo code validation/management, referral program | PostgreSQL |
| Scheduler | 8090 | Background worker — polls for scheduled rides and time-based tasks; exposes only health/version/metrics, no application API | PostgreSQL, calls Notifications service URL |
| Analytics | 8091 | Business intelligence endpoints (revenue, promo performance, ride-type mix, referral funnel, heat maps, demand zones) | PostgreSQL |
| Fraud | 8092 | Fraud alerting, risk scoring, investigation/resolution workflow, automated payment/ride pattern detection | PostgreSQL |
| ML ETA | 8093 | ML-based ETA prediction, model training/tuning, accuracy analytics | PostgreSQL, Redis |

## Documented but Not Found in Implementation

| Service | Port (per README) | Purpose (per README) | Status |
|---|---|---|---|
| Negotiation | 8094 | Fare negotiation between rider and driver | Listed in README service table only. No `cmd/negotiation`, no compose service, no k8s manifest reference found during this review. |

## Supporting / Non-Domain Components

| Component | Role |
|---|---|
| `pkg/` | Shared libraries used across all services: common response envelope, middleware (auth, rate limiting, request logging, correlation IDs, CORS, security headers), config loading, resilience (circuit breaker wiring), logger (Zap), tracing (OTel) |
| `third_party/gobreaker` | Vendored copy of the `gobreaker` circuit-breaker library |
| Kong | API gateway sitting in front of the service mesh (see Infrastructure Inventory) |
| Istio (optional) | Service-mesh layer, configured separately from Kong |
| Scheduler | Technically a service process, but functions as a background worker rather than an API-serving microservice — listed above for completeness |

## Inter-Service Call Graph (as documented)

- **Rides → Promos**: optional, for fare/promo calculation
- **Rides → ML ETA**: optional, for ETA prediction
- **Rides/Payments/Admin/etc. → Notifications**: ride lifecycle notification hooks (requested/accepted/started/completed/cancelled) are invoked by other services using a trusted "admin" service account
- **Scheduler → Notifications**: polls DB and calls the notifications service URL for due scheduled rides/tasks
- **Realtime**: receives internal, network-restricted broadcast calls from other services to push events to connected WebSocket clients

## Notes / Limitations of This Inventory

- GitHub's robots policy blocked direct retrieval of the `internal/`, `cmd/`, `k8s/`, and `kong/` directory tree listings during this session, so service boundaries above are derived from the README, `docs/API.md`, `docs/QUICKSTART.md`, `docs/DEPLOYMENT.md`, and `docker-compose.yml` rather than from enumerating every source file.
- Service-to-service URLs (`PROMOS_SERVICE_URL`, `ML_ETA_SERVICE_URL`, `GEO_SERVICE_URL`, `REALTIME_SERVICE_URL`, `NOTIFICATIONS_SERVICE_URL`) are configured via environment variables, confirming a synchronous HTTP-based service mesh rather than a fully event-driven inter-service contract (see Event Inventory for the asynchronous surface that does exist).
