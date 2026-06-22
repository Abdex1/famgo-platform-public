# Event Inventory — richxcame/ride-hailing

Source: `.env.example`, `docker-compose.yml`, `docs/API.md`, `docs/DEPLOYMENT.md`.

## Headline Finding

This platform is predominantly **synchronous and HTTP-driven** between services. A message bus (NATS with JetStream) is provisioned in infrastructure but is **disabled by default** (`NATS_ENABLED=false` in `.env.example`), and no event payload contracts, subject/topic names, consumer groups, sagas, outbox tables, retry policies, or dead-letter queues are documented anywhere in the repository's docs. What follows is an inventory of every asynchronous or event-like mechanism that *is* observable, plus an explicit list of what is provisioned but not (yet) wired into application logic.

## Provisioned Message Bus (infrastructure only)

| Item | Value / Observation |
|---|---|
| Broker | NATS 2.x with JetStream (`nats:2-alpine`, started with `--js` flag) |
| Client URL | `NATS_URL=nats://localhost:4222` |
| Stream name | `NATS_STREAM_NAME=ridehailing` |
| Enabled flag | `NATS_ENABLED=false` by default in `.env.example`; the production checklist in `docs/DEPLOYMENT.md` explicitly instructs operators to "Enable NATS (`NATS_ENABLED=true`) with persistent JetStream storage" before go-live, implying it is **not used by default in the documented deployment path** |
| Monitoring | NATS HTTP monitoring port `8222` exposed in `docker-compose.yml` |
| Application wiring | No NATS publisher/subscriber code, subject names, or message schemas appear in any of the docs reviewed (API reference, deployment guide, database operations guide, observability guide) |

## Optional Secondary Bus (also provisioned, also unused by default)

| Item | Value / Observation |
|---|---|
| Provider | Google Cloud Pub/Sub |
| Config | `PUBSUB_PROJECT_ID`, `PUBSUB_ENABLED=false` |
| Application wiring | Not referenced elsewhere in the documentation reviewed |

## Event-Like Mechanisms That Are Actually Implemented

These are not "events" in the pub/sub sense — they are synchronous REST calls that play the *role* events would play in a more decoupled design:

| Mechanism | How it works | Trigger |
|---|---|---|
| Ride lifecycle notification hooks | `POST /api/v1/notifications/ride/{requested,accepted,started,completed,cancelled}` on the Notifications service | Called directly (HTTP) by whichever service drives the corresponding ride-state transition, using a trusted service-account JWT |
| Realtime broadcast fan-out | `POST /internal/broadcast/ride` and `POST /internal/broadcast/user` on the Realtime service push a JSON payload to all WebSocket clients subscribed to a ride or user | Called directly (HTTP) by other services when they want connected clients (rider/driver apps) to receive a live update |
| Stripe webhook intake | `POST /api/v1/webhooks/stripe` on the Payments service | Inbound event from Stripe (e.g., `payment_intent.succeeded`); payload shape is validated but signature verification is explicitly noted as not yet implemented |
| Admin bulk notification broadcast | `POST /api/v1/admin/notifications/bulk` | Synchronous fan-out to a list of user IDs, queued for delivery across a chosen channel |
| Scheduler polling | Background worker polls the database directly for due scheduled rides/tasks, then calls the Notifications service | Time-based polling, not event-driven |

## Pub/Sub-Adjacent Infrastructure (data plane, not message bus)

| Mechanism | Purpose |
|---|---|
| Redis Pub/Sub | Listed in the README tech stack ("Redis 7 (GeoSpatial, Pub/Sub)") as a capability of the Redis layer, used in conjunction with the Geo and Realtime services for live location data. No specific channel names or message schemas are documented. |
| Redis GeoSpatial index (`drivers:geo:index`) | Referenced in the README as the mechanism for nearby-driver search; this is a query/index pattern, not an event stream. |

## Explicitly Not Found / Not Documented

- **Event contracts/schemas** — no documented message envelope or versioning scheme for any bus.
- **Topics/subjects** — no NATS subject naming convention or Pub/Sub topic list beyond the single configured stream name (`ridehailing`).
- **Sagas** — no orchestration or choreography pattern documented for multi-service ride transactions (e.g., ride accept → payment hold → driver dispatch). The ride lifecycle instead relies on direct, synchronous service-to-service calls and shared database state.
- **Outbox pattern** — no transactional outbox table or relay process referenced in the database operations guide or migrations description.
- **Retry policy for async messages** — none documented (the only documented retry/resilience mechanism is the HTTP-level circuit breaker, covered in the Infrastructure Inventory).
- **Dead-letter queues** — not referenced anywhere in the reviewed documentation.

## Caveat

This inventory reflects what is described in the repository's documentation set and infrastructure definitions (`.env.example`, `docker-compose.yml`, `docs/`). GitHub's robots policy blocked direct retrieval of the `internal/` and `pkg/` source-tree listings during this session, so it is possible that NATS publisher/consumer code exists in source files that were not individually inspected. The default-disabled flags and the absence of any event-contract documentation are nonetheless strong signals that event-driven messaging is, at most, a partially built or aspirational capability rather than a load-bearing part of the current architecture.
