# API Inventory — richxcame/ride-hailing

Source: `docs/API.md` (canonical API reference index in the repository), cross-checked against the README service table.

## Gateway & Routing Convention

- Each service exposes its own port locally (`8081`–`8093`); when fronted by Kong or Istio, ports collapse behind the gateway but the `/api/v1/...` path structure is preserved.
- Every service exposes the same operational trio: `GET /healthz`, `GET /version`, `GET /metrics` (Prometheus scrape format).

## Cross-Cutting Conventions

- **Auth**: Bearer JWT issued by the Auth service; token carries `user_id` and `role`. Three roles: `rider`, `driver`, `admin`. Service-to-service calls reuse the same middleware with a seeded admin/service account.
- **Response envelope**: All handlers return a consistent shape — a `success` boolean, a `data` payload, and an optional `meta` block (used for pagination); errors use the same envelope with an `error.code` / `error.message` pair mirroring the HTTP status.
- **Pagination**: Two competing styles exist in parallel — page-based (`page`/`per_page`, used by rides and fraud alert listings) and offset-based (`limit`/`offset`, used by wallet transactions, notifications, ride history, analytics listings).
- **Rate limiting**: Redis-backed token buckets on the Rides service by default — 120 req/min (burst 40) for authenticated callers, 60 req/min (burst 20) for anonymous callers on public endpoints. Standard `X-RateLimit-*` response headers are emitted.
- **IDs/format**: UUIDv4 identifiers, RFC 3339 UTC timestamps, floating-point currency values in a single configured currency (USD by default).

## Service-by-Service Endpoint Surface

### Auth (`:8081`, base `/api/v1/auth`)
| Method | Path | Auth | Purpose |
|---|---|---|---|
| POST | `/register` | none | Create rider/driver account |
| POST | `/login` | none | Exchange credentials for JWT |
| GET | `/profile` | bearer | Fetch current user |
| PUT | `/profile` | bearer | Update name/phone fields |

### Rides (`:8082`, base `/api/v1`)
Rider-facing: create ride, fetch one ride, list own rides (paginated), surge-info lookup by coordinates, cancel ride, rate a completed ride.
Driver-facing: list available ride requests, accept a ride, start a ride, complete a ride (accepts actual distance, computes final fare).
Ride status machine: `requested → accepted → in_progress → completed`, with `cancelled` as an alternate terminal state.

### Mobile API (`:8087`, base `/api/v1`)
A façade reusing Rides/Favorites logic for mobile clients: ride history (offset-paginated, filterable by status/date), ride receipt, rating, profile get/update, and full CRUD on favorite locations.

### Geo (`:8083`, base `/api/v1/geo`)
Driver location push (`POST /location`, driver role required), location lookup by driver ID, and a distance/ETA utility endpoint that returns `distance_km` and `eta_minutes` between two coordinate pairs.

### Payments (`:8084`, base `/api/v1`)
Wallet fetch/create-on-demand, wallet top-up (Stripe-backed or simulated if no key configured), wallet transaction history, ride payment processing, payment lookup, refund request (rider for own payments, admin for any), and an unauthenticated Stripe webhook intake endpoint. The webhook handler validates payload shape only — signature verification is explicitly noted as missing and required before production use.

### Promos (`:8089`, base `/api/v1`)
Ride type listing, fare calculation, promo code validation, referral code retrieval/application, and admin-only promo code creation.

### Notifications (`:8085`, base `/api/v1`)
User-facing: list notifications, unread count, mark-as-read, ad hoc send, scheduled send.
Ride lifecycle hooks (trusted-caller only): `/notifications/ride/{requested,accepted,started,completed,cancelled}` — these are the closest thing in the codebase to ride-lifecycle "events," though they are synchronous REST calls rather than a message-bus subscription.
Admin: bulk broadcast to a list of user IDs across a chosen channel.

### Realtime (`:8086`, base `/api/v1`)
WebSocket upgrade (`GET /ws`, JWT required), per-ride chat history, driver location lookup, connection/hub statistics (admin), and two internal broadcast endpoints (`/internal/broadcast/ride`, `/internal/broadcast/user`) used by other services to push data to connected clients. The documentation explicitly flags these two internal routes as **lacking auth middleware** and recommends network ACLs/mTLS rather than relying on application-level auth.

### Admin (`:8088`, base `/api/v1/admin`)
Dashboard stats, paginated user list/detail, suspend/activate user, list/approve/reject pending drivers, recent rides, ride statistics by date range. All routes require both authentication and the `admin` role.

### Analytics (`:8091`, base `/api/v1/analytics`)
Admin-only BI surface: dashboard snapshot, revenue, promo-code performance, ride-type usage mix, referral funnel, top drivers, geographic heat map, financial report, demand zones. Date ranges default to the trailing 30 days.

### Fraud (`:8092`, base `/api/v1/fraud`)
Admin-only: alert listing/detail/creation, investigate/resolve workflow, per-user alert history and risk profile, on-demand risk analysis, suspend/reinstate user, and two detection triggers (payment fraud, ride pattern fraud). Alert taxonomy: types `payment_fraud`, `account_fraud`, `location_fraud`, `ride_fraud`, `rating_manipulation`, `promo_abuse`; severity levels `low/medium/high/critical`.

### ML ETA (`:8093`, base `/api/v1/eta`)
Public: single and batch ETA prediction (max 100 routes per batch).
Admin: trigger async model retraining, model stats, accuracy-over-time, hyperparameter tuning.
Authenticated: historical predictions, accuracy trend, feature-importance breakdown.

### Scheduler (`:8090`)
No application API — only `healthz`/`version`/`metrics`. Operates as a polling background worker against the database and the Notifications service.

## API Surface Gaps / Risk Notes (observed in documentation, not inferred)

- Stripe webhook signature verification is explicitly called out as **not implemented**.
- The two Realtime internal broadcast endpoints are explicitly called out as **unauthenticated at the application layer**, relying on network-level controls instead.
- Two pagination conventions (page-based vs. offset-based) coexist across the API surface rather than a single standard.
