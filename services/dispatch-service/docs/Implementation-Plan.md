# Dispatch Service — Implementation Plan

## Phase 1 — Core (complete)

- [x] Domain entities and matching engines
- [x] PostgreSQL repositories and migrations (`db/migrations`)
- [x] gRPC API (`api/proto/v1/dispatch.proto`)
- [x] CQRS application layer (commands, queries, handlers, use cases)
- [x] Kafka publisher and `ride.created.v1` consumer saga
- [x] GPS discovery client integration

## Phase 2 — Platform integration (complete)

- [x] Auth-service JWT validation and RBAC scopes (`dispatch:read`, `dispatch:write`)
- [x] Redis-backed rate limiting via `redis-platform`
- [x] OpenTelemetry tracing, HTTP/gRPC middleware, correlation IDs
- [x] REST API with security middleware on `/v1/dispatch/*`
- [x] K8s deployment manifest (`deploy/k8s/dispatch-service.yaml`)
- [x] OpenAPI spec (`api/openapi/v1/dispatch.openapi.yaml`)

## Phase 3 — Hardening (next)

- [ ] Contract tests against ride-service and gps-service stubs
- [ ] Integration tests with Testcontainers (Postgres, Redis, Kafka)
- [ ] Load testing for concurrent match requests
- [ ] Circuit breakers for external gRPC clients
- [ ] Dead-letter queue for failed Kafka messages

## Phase 4 — Production rollout

- [ ] Deploy to staging with feature flag `KAFKA_ENABLED=true`
- [ ] Validate end-to-end saga: ride.created → match → dispatch.matched.v1
- [ ] Enable autoscaling on CPU and Kafka consumer lag
- [ ] Runbook validation and on-call handoff

## Success criteria

| Metric | Target |
|--------|--------|
| Unit test coverage | ≥ 80% on domain + application layers |
| P95 match latency | < 2s with GPS service available |
| Kafka consumer lag | < 5s under normal load |
| Error budget | 99.9% availability |

## Dependencies

| Service | Protocol | Purpose |
|---------|----------|---------|
| auth-service | JWT | Token validation |
| gps-service | gRPC | Nearby driver discovery |
| ride-service | Kafka | `ride.created.v1` trigger |
| redis-platform | TCP | Rate limiting |
| telemetry | OTLP | Traces and logs |
