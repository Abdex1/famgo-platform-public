# Dispatch Service Runbook

## Health Checks

1. Verify PostgreSQL connectivity (`DB_*` env vars).
2. Verify gps-service gRPC endpoint (`GPS_SERVICE_URL`, default `localhost:5002`).
3. Confirm gRPC listener on `GRPC_PORT` (default `5005`).

## Common Failures

### No drivers matched

- Confirm drivers are online in gps-service geo index.
- Increase `SEARCH_RADIUS_KM` or `MAX_SEARCH_RADIUS_KM`.
- Check `MIN_ACCEPTANCE_RATE_PERCENT` and `MIN_RATING` thresholds.

### Matching timeouts

- Review `MATCH_EXPIRY_SECONDS` and `MATCH_REQUEST_TTL`.
- Inspect `matching_sessions` for stuck `active` rows past `expires_at`.

### Build failures in workspace mode

Run with isolated module resolution:

```powershell
$env:GOWORK="off"
go test ./... -count=1
```

Root `go.work` currently references services with invalid `famgo/shared` module paths.

## Migrations

```powershell
migrate -path db/migrations -database $DATABASE_URL up
```

Rollback one step:

```powershell
migrate -path db/migrations -database $DATABASE_URL down 1
```

## Kafka Operations

Ensure dispatch topics exist:

- `dispatch.matching.started.v1`
- `dispatch.driver.matched.v1`
- `dispatch.driver.assigned.v1`
- `dispatch.matching.failed.v1`
- `dispatch.matching.expired.v1`

Monitor consumer lag for ride saga subscribers.

## Escalation

1. Check gps-service availability and Redis geo index population.
2. Check dispatch DB row counts by status.
3. Replay `ride.created.v1` only through controlled saga tooling (never mutate dispatch state manually in production).
