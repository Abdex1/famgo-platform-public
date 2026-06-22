# richxcame/ride-hailing — Contracts

Source: `docs/API.md` as inventoried in Stage 1. richxcame has **no event/message contracts** (NATS is disabled by default, no schemas documented) — everything below is an HTTP API contract.

## Cross-Cutting Contract Conventions

| Convention | Detail |
|---|---|
| Auth | Bearer JWT; claims carry `user_id` and `role` (`rider`/`driver`/`admin`) |
| Response shape | `{ success: bool, data: <payload>, meta?: {...} }`; errors use the same envelope with `error.code` / `error.message` |
| Pagination | **Not unified** — page-based (`page`/`per_page`) on rides and fraud alerts; offset-based (`limit`/`offset`) on wallet transactions, notifications, ride history, analytics |
| IDs | UUIDv4 |
| Timestamps | RFC 3339 UTC |
| Currency | Floating-point, single configured currency (USD by default) |
| Rate limiting | Redis-backed token bucket; 120 req/min (burst 40) authenticated, 60 req/min (burst 20) anonymous; `X-RateLimit-*` headers |

## Service Contract Summaries (method/path/auth only — see Stage 1 `02-api-inventory.md` for full detail)

- **Auth** (`/api/v1/auth`): `POST /register` (none), `POST /login` (none), `GET /profile` (bearer), `PUT /profile` (bearer)
- **Rides** (`/api/v1`): rider create/fetch/list/cancel/rate; driver list/accept/start/complete. State machine: `requested → accepted → in_progress → completed`, alternate terminal `cancelled`.
- **Geo** (`/api/v1/geo`): `POST /location` (driver role), `GET /drivers/:id/location`, `POST /distance` → `{distance_km, eta_minutes}`
- **Payments** (`/api/v1`): wallet fetch/top-up/history, ride payment processing, refund (rider-own / admin-any), unauthenticated `POST /webhooks/stripe` (signature verification **not implemented** — do not treat as a safe contract to copy as-is)
- **Notifications** (`/api/v1`): list/unread-count/mark-read/send/scheduled-send; trusted-caller-only ride hooks `/notifications/ride/{requested,accepted,started,completed,cancelled}`; admin bulk broadcast
- **Realtime** (`/api/v1`): `GET /ws` (bearer), `GET /rides/:ride_id/chat`, `GET /drivers/:driver_id/location`, `GET /stats` (admin); **unauthenticated** internal routes `POST /internal/broadcast/ride`, `POST /internal/broadcast/user`
- **Admin** (`/api/v1/admin`): dashboard stats, user list/detail/suspend/activate, driver approve/reject, recent rides, ride stats — all admin-role-gated
- **Analytics** (`/api/v1/analytics`): dashboard, revenue, promo performance, ride-type mix, referral funnel, top drivers, heat map, financial report, demand zones — admin-only, 30-day default window
- **Fraud** (`/api/v1/fraud`): alert list/detail/create, investigate/resolve, per-user history/risk, on-demand analysis, suspend/reinstate, two detection triggers — admin-only
- **ML ETA** (`/api/v1/eta`): public single/batch prediction (max 100/batch); admin retrain/stats/accuracy/tuning; authenticated historical predictions/trend/feature-importance
- **Mobile** (`/api/v1`): ride history (offset-paginated), receipt, rating, profile, favorite-location CRUD — façade over Rides/Favorites

## Contract Risks Flagged in Source Documentation (do not silently inherit)

- Stripe webhook: payload-shape validation only, **no signature verification**.
- Internal Realtime broadcast routes: **no application-layer auth middleware**.
