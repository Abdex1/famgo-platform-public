# Dispatch Service

Ride matching and driver assignment microservice for FamGo.

## Responsibilities

- Driver discovery via `gps-service` gRPC (`FindNearbyDrivers`)
- Nearest-driver matching with multi-factor scoring (`FindDriversWithinRadius → SortByDistance → AssignFirstAvailable`)
- Matching sessions, reassignment, and timeout handling
- Dispatch saga integration for `ride.created.v1`
- Kafka events under `dispatch.*.v1`

## Run

```powershell
cd services/dispatch-service
$env:GOWORK="off"
go run ./cmd
```

Environment variables are documented in `.env.example`.

## Architecture

See [docs/Architecture.md](docs/Architecture.md).

## API

- gRPC: `api/proto/v1/dispatch.proto` on port `5005` (default)
- REST legacy entrypoint: `cmd/api` (health + operational endpoints)

## Database

Apply migrations in order:

```powershell
migrate -path db/migrations -database $DATABASE_URL up
```

## Tests

```powershell
$env:GOWORK="off"
go test ./... -count=1
```
