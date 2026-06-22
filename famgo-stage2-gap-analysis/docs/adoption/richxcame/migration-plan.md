# richxcame/ride-hailing — Migration Plan

Source: derived from `features.md`, `contracts.md`, `patterns.md`, `anti-patterns.md` in this folder, cross-referenced with `docs/adoption/gap-analysis.md` and `docs/adoption/migration-matrix.md`.

## Principle

Never import richxcame code directly. Every item below is a **re-implementation against FamGo's own stack** (Go DDD layout, gRPC inter-service calls, Kafka events, existing `platform/`/`packages/` structure), using richxcame only as a checklist or cautionary reference.

## Phase 1 — Low-risk, low-effort conventions (adopt as-is in spirit)

| Step | Action | Target | Depends on |
|---|---|---|---|
| 1.1 | Add `/healthz`, `/version`, `/metrics` to every FamGo service that doesn't already expose them | `platform/runtime/health/` | none |
| 1.2 | Verify FamGo's response envelope is actually consistent across services (cannot be confirmed from structure alone — read code) | `standards/api/`, `packages/api-client/` | code review |
| 1.3 | Pick one pagination convention and document it | `standards/api/` | 1.2 |
| 1.4 | Add the dev/staging/prod OTel sampling table to FamGo's telemetry config | `shared/telemetry/`, `infra/monitoring/otel-collector/` | none |
| 1.5 | Add `JWT_ROTATION_HOURS`-equivalent config if not already present | `platform/security/jwt/rotation/` | code review |

## Phase 2 — Medium-effort capability gaps

| Step | Action | Target | Depends on |
|---|---|---|---|
| 2.1 | Populate `platform/resilience/ratelimit/` with a Redis-token-bucket limiter modeled on richxcame's auth/anonymous tiering | `platform/resilience/ratelimit/`, `platform/cache/redis/` | none |
| 2.2 | Build out `services/notification-service` to cover the same lifecycle-hook trigger points richxcame documents (`requested/accepted/started/completed/cancelled`) — **as Kafka event consumers, not synchronous REST hooks** | `services/notification-service/`, `packages/event-bus/contracts/ride/`, `packages/event-bus/contracts/notification/` | `ride-service` emits the corresponding events first |
| 2.3 | Build out `services/analytics-service` (currently empty) using richxcame's endpoint list as a feature checklist, not a literal spec | `services/analytics-service/` | data availability from source services |
| 2.4 | Write database backup/restore/health-check scripts for FamGo's actual storage targets, modeled on richxcame's documented capabilities (compression, encryption, S3/GCS, retention, alerting) | `scripts/`, `database/backups/` | resolve `database/migrations/` consolidation first (see gap-analysis §4) |
| 2.5 | Write seed-data scripts at light/medium/heavy tiers, scaled to FamGo's actual domain entities (including pooling, safety — which richxcame has no equivalent for) | `database/seeds/` | schema consolidation |
| 2.6 | Decide and document the pluggable secrets provider abstraction, consolidating FamGo's currently-split `security/vault/`, `security/secrets/`, `platform/config/vault/`, `packages/vault-sdk/` | `platform/config/vault/`, `packages/vault-sdk/` | none |

## Phase 3 — Audit items derived from richxcame's documented anti-patterns (do first, not last, where flagged High in gap-analysis)

| Step | Action | Target | Why urgent |
|---|---|---|---|
| 3.1 | Confirm whether FamGo's internal/broadcast-style endpoints (`services/websocket-gateway/`, `gateway/`) have application-layer auth, not just network ACLs | `services/websocket-gateway/`, `gateway/policies/` | richxcame's own documented failure on the identical mechanism |
| 3.2 | Confirm payment webhook signature verification exists in `wallet-service`/`packages/payment-sdk/` | `services/wallet-service/`, `packages/payment-sdk/` | richxcame's own documented failure on the identical mechanism; payment fraud vector |
| 3.3 | Do not conflate "admin" and "trusted service caller" identities if FamGo's gRPC inter-service calls use any shared/seeded credential pattern — verify each service-to-service call uses a properly scoped identity | `platform/security/jwt/` | richxcame anti-pattern #6 |

## Explicitly Out of Scope for This Migration

- richxcame's NATS bus — FamGo is committed to Kafka; do not introduce NATS.
- richxcame's Negotiation service — never shipped, nothing to migrate.
- richxcame's shared single-Dockerfile build pattern — FamGo's per-service Dockerfiles stay as-is.
