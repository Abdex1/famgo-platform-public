# Realtime Inventory — richxcame/ride-hailing

Source: `docs/API.md` (Realtime and Geo service sections), README tech stack, `.env.example`, `docker-compose.yml`.

## Realtime Service (`:8086`)

| Capability | Detail |
|---|---|
| WebSocket upgrade | `GET /api/v1/ws` — requires a Bearer JWT; the token's `user_id` and `role` claims are used to track the connection inside an in-memory "hub" |
| Client contract | Referenced as `pkg/websocket` in `docs/API.md`; structured JSON payloads sent/received over the socket (exact message schema not enumerated in the docs reviewed) |
| Chat history | `GET /api/v1/rides/:ride_id/chat` — returns stored in-ride chat messages for participants of that ride |
| Driver location relay | `GET /api/v1/drivers/:driver_id/location` — fetches latest coordinates from the Redis-backed store |
| Connection stats | `GET /api/v1/stats` — admin-only; connection counts and hub statistics |
| Internal broadcast (ride-scoped) | `POST /api/v1/internal/broadcast/ride` — pushes an arbitrary payload to every client subscribed to a given ride |
| Internal broadcast (user-scoped) | `POST /api/v1/internal/broadcast/user` — pushes a typed payload (e.g., `"type": "notification"`) to every connection belonging to a given user |
| Connection timeout | `WS_CONNECTION_TIMEOUT` env var, default 60 seconds |

### Documented Security Gap

The API reference explicitly states that the two internal broadcast routes **do not have authentication middleware attached**, and recommends deploying them behind mTLS or network ACLs rather than exposing them outside the service mesh. This is a stated fact from the project's own documentation, not an inferred risk.

## Geo Service (`:8083`) — Realtime-Adjacent

| Capability | Detail |
|---|---|
| Location ingestion | `POST /api/v1/geo/location` — driver-role-only; updates the driver's last known position |
| Location lookup | `GET /api/v1/geo/drivers/:id/location` — returns latitude/longitude/`updated_at` for a given driver |
| Distance/ETA utility | `POST /api/v1/geo/distance` — computes straight-line distance and a simple ETA between two coordinate pairs |
| Storage | Redis, using a geospatial index referenced in the README as `drivers:geo:index` |

## Underlying Real-Time Transport

| Mechanism | Role |
|---|---|
| WebSocket (Realtime service) | Client-facing push channel for ride status updates, chat, and notifications |
| Redis Pub/Sub | Described in the README tech stack as a Redis capability used alongside the GeoSpatial index; exact channel/topic names are not documented |
| Redis GeoSpatial commands | Backing store for "nearby driver" queries — a request/response geospatial query pattern rather than a streaming one |

## What Is Not Real-Time (for contrast)

- Ride lifecycle notification hooks (requested/accepted/started/completed/cancelled) are synchronous REST calls into the Notifications service, not a streaming or event-bus mechanism (see Event Inventory).
- Scheduler operates on a polling loop against the database, not a push/subscribe model.

## Caveat

The Hub implementation, exact WebSocket message schema, and any reconnection/heartbeat behavior live in source code (`internal/realtime/`, `pkg/websocket/`) that could not be directly enumerated in this session due to GitHub robots restrictions on directory-tree pages. This inventory reflects what the project's own API reference documents about the realtime surface.
