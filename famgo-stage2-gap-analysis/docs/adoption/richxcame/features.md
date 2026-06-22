# richxcame/ride-hailing — Features

Source: Stage 1 service/API/event/database/realtime/infrastructure inventories.

## Service-Level Features Worth Extracting

- **Operational trio per service**: every one of 13 implemented services exposes `GET /healthz`, `GET /version`, `GET /metrics` (Prometheus format) — a uniform, low-effort convention.
- **Unified response envelope**: `success` boolean + `data` payload + optional `meta` (pagination), with errors mirrored as `error.code`/`error.message`.
- **Role-based JWT auth** with three roles (`rider`, `driver`, `admin`), service-to-service calls reusing the same middleware via a seeded admin/service account.
- **Fraud alert taxonomy**: types `payment_fraud`, `account_fraud`, `location_fraud`, `ride_fraud`, `rating_manipulation`, `promo_abuse`; severities `low/medium/high/critical`; investigate/resolve workflow with per-user risk profile and on-demand analysis triggers.
- **Ride lifecycle notification hooks**: a named, complete set of trigger points (`requested`, `accepted`, `started`, `completed`, `cancelled`) that any ride-state-machine-driven notification system needs to cover, regardless of transport mechanism.
- **Wallet model**: fetch/create-on-demand wallet, top-up (provider-backed or simulated when no key is configured — useful for local/dev environments), transaction history, refund request scoped to rider-for-own vs. admin-for-any.
- **Promo/referral domain split**: ride type listing, fare calculation, promo code validation, referral code retrieval/application, admin-only promo creation — a clean separation of "pricing inputs" from "pricing computation."
- **Analytics surface breadth**: revenue, promo-code performance, ride-type usage mix, referral funnel, top drivers, geographic heat map, financial report, demand zones — a useful checklist of what a BI surface for a ride-hailing platform typically needs to cover, even though FamGo's own `analytics-service` is currently empty.
- **Scheduler as a pure background worker**: no application API surface at all beyond the operational trio — a clean example of a worker-only service that doesn't need to be a "real" API service.
- **Seed data tiers**: light (11 users/9 rides/5 payments) / medium (50 users/200 rides) / heavy (1,000 users/5,000 full-lifecycle rides) — practical scale points for dev/test/load environments.
- **Database operational tooling**: scripted backup (`--compress`, `--encrypt`, `--storage s3|gcs`, `--retention`), scripted restore (latest/file/remote/timestamp/new-database/validate-only modes), automated backup health checks with Slack/email alerting, a documented PITR runbook.
- **Pluggable secrets provider**: `SECRETS_PROVIDER=env|vault|aws|gcp|kubernetes` with cache TTL, rotation interval, and an audit-logging toggle — independent of JWT key rotation, which has its own schedule.
- **Resilience configuration granularity**: a circuit breaker with a global enable flag, separate DB-specific breaker, and a per-service override map (`CB_SERVICE_OVERRIDES`); granular timeout configuration per HTTP client call, DB query, Redis op, WebSocket connection, and per-route overrides via a JSON map.
- **Environment-scaled observability sampling**: a concrete dev/staging/production OTel-trace and Sentry sampling table (100%/50%/10% traces; 100% errors at every tier).

## Explicit Non-Features (documented gaps in richxcame itself — useful as a "don't repeat this" list)

- Stripe webhook **signature verification is not implemented** — payload shape only.
- Two internal Realtime broadcast routes have **no application-layer auth**, relying on network ACLs/mTLS instead.
- **Two competing pagination conventions** (page-based and offset-based) coexist with no single standard.
- NATS is **provisioned but disabled by default**, with no documented application wiring.
- The Negotiation service is **documented but never shipped**.
