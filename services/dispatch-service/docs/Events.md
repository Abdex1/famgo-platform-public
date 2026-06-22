# Dispatch Events

Event naming follows `domain.entity.action.v1`.

## Published

| Event | Topic | When |
|---|---|---|
| MatchingStarted | `dispatch.matching.started.v1` | Dispatch request persisted, matching begins |
| DriverMatched | `dispatch.driver.matched.v1` | Primary driver selected and proposed list stored |
| DriverAssigned | `dispatch.driver.assigned.v1` | Driver accepts match |
| MatchingFailed | `dispatch.matching.failed.v1` | No eligible drivers or assignment failure |
| MatchingExpired | `dispatch.matching.expired.v1` | Dispatch request timed out |

Contracts: `packages/event-bus/contracts/dispatch/events.go`

## Consumed

| Event | Topic | Handler |
|---|---|---|
| RideCreated | `ride.created.v1` | `DispatchSagaHandler.HandleRideCreated` |

## Envelope

All events use `packages/event-bus/envelope.EventEnvelope` with:

- `trace_id`
- `correlation_id`
- `request_id`
- `partition_key` (ride ID)

## Registry

Topics registered in `packages/event-bus/registry/topics.go` and `packages/event-bus/topics/dispatch.go`.
