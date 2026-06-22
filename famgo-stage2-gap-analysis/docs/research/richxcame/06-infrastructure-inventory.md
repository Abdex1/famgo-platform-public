# Infrastructure Inventory — richxcame/ride-hailing

Source: `docker-compose.yml`, `docs/DEPLOYMENT.md`, `docs/observability.md`, `.env.example`, README.

## Runtime & Packaging

| Item | Detail |
|---|---|
| Language/runtime | Go 1.24+, Gin web framework |
| Containerization | Single shared `Dockerfile` parameterized by a `SERVICE_NAME` build argument; each service is built as `docker build --build-arg SERVICE_NAME=<svc>` |
| Compose variants | `docker-compose.yml` (production-shaped, full stack incl. observability and self-hosted Sentry), `docker-compose.dev.yml` (infrastructure-only: Postgres + Redis, for native `go run` development), `docker-compose.test.yml` (CI/integration tests) |
| Build tooling | `Makefile` driving build/run/migrate/test/lint/format targets |
| CI | GitHub Actions (`.github/workflows`, referenced via a `ci.yml` badge in the README); Codecov integration for coverage reporting |

## Orchestration (Kubernetes)

| Item | Detail |
|---|---|
| Manifests location | `k8s/` |
| Core resources referenced | `namespace.yaml`, `configmap.yaml`, a secrets template, and a per-service deployment/service file for each of the 13 implemented services |
| Infra manifests | `postgres.yaml`, `redis.yaml`, `nats.yaml` |
| Ingress | `ingress.yaml` |
| Scaffolding | `k8s/generate-services.sh` — script to generate new service manifests from a template |
| Service mesh (optional) | Istio configs under `k8s/istio/`: `install-istio.sh`, `gateway.yaml`, `destination-rules.yaml`, `security-policies.yaml` |

## API Gateway

| Item | Detail |
|---|---|
| Product | Kong 3.9.1 |
| Config store | Dedicated PostgreSQL instance (`kong-database`), bootstrapped via a one-shot `kong-migration` container running `kong migrations bootstrap` |
| Exposed ports | `8000`/`8443` (proxy HTTP/HTTPS), `8001`/`8444` (Admin API HTTP/HTTPS), `8002` (Kong Manager GUI) |
| Config directory | `kong/` (declarative/route configuration; not individually enumerated in this session) |
| Relationship to Istio | Both are present in the repo as alternative or layered approaches to traffic management — Kong at the edge, Istio for mesh-internal traffic — per `docs/DEPLOYMENT.md` |

## Data & Cache Infrastructure

Covered in detail in the Database Inventory; summarized here for completeness:

- PostgreSQL 15 + PostGIS (`ridehailing` database)
- Redis 7 (cache, geospatial index, Pub/Sub)
- NATS 2 with JetStream (provisioned, disabled by default — see Event Inventory)

## Observability Stack

| Tool | Role |
|---|---|
| OpenTelemetry (SDK + Collector) | Trace instrumentation and collection; collector config at `deploy/otel-collector.yml`; exposes OTLP gRPC (4317), OTLP HTTP (4318), Prometheus exporter (8889), health (13133), zpages (55679) |
| Grafana Tempo | Trace storage/query backend; config at `deploy/tempo.yml`; backed by MinIO (S3-compatible object storage) for trace blob storage |
| Prometheus | Metrics collection and alerting; config and alert rules in `monitoring/` |
| Grafana | Unified dashboarding (System Overview, Rides Service, Payments Service dashboards auto-provisioned under a "RideHailing" folder) |
| Zap | Structured JSON application logging with `request_id`/`trace_id`/`span_id` correlation |
| Sentry | Error tracking; can run as a SaaS DSN or as a full self-hosted stack (Postgres, Redis, ClickHouse, Kafka, Zookeeper, web/worker/cron containers) defined directly in `docker-compose.yml` |

### Sampling Strategy (documented)

| Environment | OTel traces | Sentry errors | Sentry traces |
|---|---|---|---|
| development | 100% | 100% | 100% |
| staging | 50% | 100% | 50% |
| production | 10% (configurable) | 100% | 10% |

### Alert Categories (in `monitoring/prometheus/alerts.yml`)

System (error rate, latency, service-down, CPU/memory, goroutine count), Business (driver availability, cancellation rate, payment failure rate, no-rides, revenue drop), Infrastructure (DB pool exhaustion, slow queries, Redis hit rate/memory, circuit breaker state, rate-limit rejection rate), Fraud (detection rate spikes, blocked-user spikes).

## Resilience Patterns

| Pattern | Configuration |
|---|---|
| Circuit breaker | Vendored `gobreaker` (`third_party/gobreaker`); global enable flag plus failure/success thresholds, timeout, and interval; a separate breaker specifically for database calls (`DB_BREAKER_*`); per-service override map supported (`CB_SERVICE_OVERRIDES`) |
| Rate limiting | Redis-backed token bucket; distinct authenticated vs. anonymous limits; configurable window and Redis key prefix |
| Timeouts | Granular timeout configuration for HTTP client calls, DB queries, Redis read/write, WebSocket connections, and per-route overrides via a JSON map (`ROUTE_TIMEOUT_OVERRIDES`) |

## Secrets Management

Pluggable provider model (`SECRETS_PROVIDER=env|vault|aws|gcp|kubernetes`), with configurable cache TTL, rotation interval, and an audit-logging toggle. JWT signing keys additionally support their own rotation schedule (`JWT_ROTATION_HOURS`, grace period, refresh interval) independent of the general secrets provider.

## Third-Party Integrations (infrastructure-adjacent)

| Integration | Purpose | Enabled by default? |
|---|---|---|
| Stripe | Payment processing | Yes (test key placeholder) |
| Firebase Cloud Messaging | Push notifications | No (`FIREBASE_ENABLED=false`) |
| Twilio | SMS notifications | Optional, no enable flag (presence of credentials gates use) |
| SMTP | Email notifications | Optional |
| Google Maps / HERE | Mapping/geocoding providers | No (`MAPS_ENABLED=false`) |
| Google Cloud Pub/Sub | Alternative event bus | No (`PUBSUB_ENABLED=false`) |
| Checkr | Driver background checks | No (`CHECKR_ENABLED=false`) |
| Onfido | Identity verification | No (`ONFIDO_ENABLED=false`) |

## Target Deployment Environment

Per the repository description, the system is "designed for Cloud Run + Cloud SQL deployment on Google Cloud," though the concrete infrastructure-as-code in the repo (Docker Compose + raw Kubernetes manifests + Kong + Istio) is platform-agnostic rather than Cloud Run-specific; no Cloud Run service YAML or Cloud Build pipeline was found in the files reviewed.

## Caveat

The `k8s/` and `kong/` directory listings could not be directly enumerated in this session (GitHub robots restrictions on directory-tree pages), so individual manifest filenames beyond those explicitly named in `docs/DEPLOYMENT.md` are not confirmed. This inventory reflects what `docker-compose.yml`, `.env.example`, and the documentation set describe.
