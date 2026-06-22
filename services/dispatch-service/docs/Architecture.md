# Dispatch Service Architecture

## Layering

| Layer | Path | Responsibility |
|---|---|---|
| Domain | `internal/domain` | Aggregates, value objects, engines, ports |
| Application | `internal/application` | Use cases, saga handlers |
| Infrastructure | `internal/infrastructure` | PostgreSQL, GPS client, Kafka publisher |
| Interfaces | `interfaces/grpc`, `internal/interfaces/rest` | gRPC/REST adapters |

## Core Engines

1. **DriverCandidateEngine** — queries `gps-service`, filters online drivers, sorts by distance.
2. **MatchingEngine** — scores candidates via `MatchingService`, returns primary + proposed drivers.
3. **AssignmentEngine** — assigns first available driver, supports reassignment after rejection.
4. **TimeoutService** — enforces dispatch expiry and retry eligibility.

## State Machine

`DispatchRequest` statuses:

`pending → matching → matched → accepted → completed`

Failure paths: `rejected`, `failed`, `expired`, `cancelled`.

## Saga Integration

`internal/application/saga/dispatch_saga.go` consumes `ride.created.v1` and starts matching. Compensation cancels dispatch via `CancelDispatch`.

## Pooling Hooks

`ports.PoolingStrategyHook` is wired with `NoOpPoolingHook` until `pooling-service` integration lands in Wave 3+.

## Observability

Events carry `trace_id`, `correlation_id`, and `request_id` through the event-bus envelope. Service bootstrap uses structured Zap logging.

## External Dependencies

- PostgreSQL (`dispatch_requests`, `matching_sessions`, `match_results`)
- gps-service gRPC (`FindNearbyDrivers`)
- Kafka (dispatch domain events)

Each service owns its database schema; cross-service communication is via APIs and events only.
