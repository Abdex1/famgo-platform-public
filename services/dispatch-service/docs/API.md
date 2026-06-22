# Dispatch API

## gRPC (`dispatch.v1.DispatchService`)

| RPC | Description |
|---|---|
| `MatchRide` | Start matching for a ride |
| `GetMatches` | Fetch dispatch request status |
| `AcceptMatch` | Driver accepts proposed match |
| `RejectMatch` | Driver rejects; optional retry/reassign |
| `CancelDispatch` | Cancel active dispatch |
| `GetDispatchStats` | Aggregate matching metrics |

Proto: `api/proto/v1/dispatch.proto`

## REST (legacy bootstrap)

`cmd/api` exposes operational endpoints:

- `GET /v1/health`
- `POST /v1/dispatch/match`
- `POST /v1/dispatch/assign`
- `POST /v1/dispatch/cancel`
- `GET /v1/dispatch/status`
- `GET /v1/dispatch/metrics`

Production clients should prefer gRPC.

## Auth

Requests are validated at the API gateway using auth-service JWTs. Internal service-to-service calls require service credentials at the gateway policy layer.

## Error Model

gRPC status codes:

- `InvalidArgument` — missing ride/dispatch identifiers
- `NotFound` — dispatch request not found
- `Internal` — matching/assignment persistence failures
