# richxcame/ride-hailing — Patterns

Source: Stage 1 inventories. These are structural/architectural patterns observed, independent of specific endpoints.

## Patterns Worth Adapting

1. **Shared single Dockerfile, parameterized by build arg** (`docker build --build-arg SERVICE_NAME=<svc>`) — reduces per-service Dockerfile drift at the cost of a slightly less conventional build. (FamGo already uses per-service Dockerfiles; noted for context, not recommended as a switch — see migration matrix REJECT row.)
2. **Three compose variants for three purposes**: `docker-compose.yml` (full prod-shaped stack incl. observability + self-hosted Sentry), `docker-compose.dev.yml` (infra-only, for native `go run` development), `docker-compose.test.yml` (CI/integration). Clean separation of concerns by environment.
3. **k8s manifest generator script** (`k8s/generate-services.sh`) — scaffolds new service manifests from a template rather than hand-writing each one. Directly relevant since FamGo's Kubernetes manifests are currently fragmented across three locations (see gap-analysis §6); a generator could be part of the consolidation fix.
4. **Synchronous service mesh via env-var URLs** (`PROMOS_SERVICE_URL`, `ML_ETA_SERVICE_URL`, etc.) rather than service discovery — simple, debuggable, but doesn't scale past a fixed known set of services and creates the exact coupling problem Microbus's event-pattern documentation describes (see microbus extraction package). Pattern to observe, not necessarily to adopt, given FamGo is already moving toward gRPC + Kafka.
5. **Trusted "admin" service account for service-to-service calls** — service A calls service B's API using a seeded admin/service JWT rather than a separate machine-to-machine auth scheme. Simple but conflates "an admin user" with "a trusted service caller" in the same role system; worth a deliberate decision rather than blind adoption.
6. **Database circuit breaker as a distinct instance from the general HTTP circuit breaker** (`DB_BREAKER_ENABLED` and friends, separate from the general `CB_*` config) — pattern of having resilience config scoped per *dependency type*, not just one global breaker.
7. **Per-route timeout overrides via a JSON config map** (`ROUTE_TIMEOUT_OVERRIDES`) layered on top of global granular timeouts (HTTP client, DB query, Redis read/write, WebSocket) — a useful escape hatch pattern for the rare endpoint that needs a non-default timeout without forking the whole timeout config.
8. **Backup health-check script as a first-class, alerting-capable component** (`check-backup-health.sh` with Slack/email), not just a backup script — treats "backups exist" and "backups are known-good" as two separate, separately-monitored concerns.

## Patterns to Note But Not Necessarily Adopt

- **Kong + Istio as parallel, layered traffic-management tools** (Kong at the edge, Istio mesh-internal) — adds operational surface area; only adopt both if FamGo genuinely needs edge + mesh-internal policy separation, not by default.
- **Self-hosted Sentry stack option** (Postgres + Redis + ClickHouse + Kafka + Zookeeper, all in one `docker-compose.yml`) — heavy footprint for local/dev; reasonable for a documented "production-shaped" compose file, risky if accidentally used as the default dev environment.
