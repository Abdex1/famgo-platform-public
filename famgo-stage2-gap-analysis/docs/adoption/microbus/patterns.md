# Microbus (microbus-io/fabric) — Patterns

Source: Stage 1 event-driven architecture notes, adapted for FamGo's Kafka-based stack.

## Pattern 1 — Pub/Sub Over Direct Calls (decoupling)

Microbus's own framing example is the cleanest available argument for finishing FamGo's `packages/event-bus/contracts/` buildout (currently only 2 of 7 domain folders populated, per gap-analysis §3): a service that directly calls every interested downstream service accumulates a cyclic, ever-growing dependency graph as features are added. Publishing an event and letting consumers subscribe independently keeps that graph a DAG. **Adapt as**: every FamGo cross-service ride-lifecycle transition (request/match/accept/start/complete/cancel) should be a published Kafka event with N independent consumers (notification-service, analytics-service, fraud-service, etc.), not a fan-out of direct gRPC calls from `ride-service`.

## Pattern 2 — Capability-Derived ACLs Instead of Hand-Maintained Lists

`gencreds` scans actual code to generate scoped credentials. **Adapt as**: a CI-time tool that scans each FamGo service's `packages/event-bus/contracts/` usage (which topics it publishes/subscribes to) and generates Kafka ACL bindings (or SASL scopes) automatically, rather than a manually maintained topic-permission spreadsheet. This directly strengthens `packages/event-bus/governance/ownership.go`, which currently exists as a file but its enforcement mechanism is unverified from structure alone.

## Pattern 3 — Adapted Topic Naming Convention

Microbus's `<plane>.<trust>.<port>.<src>.<dest>.<id_or_locality>.<method>.<path...>` doesn't map onto Kafka, but the *dimensions* it encodes are worth preserving in a Kafka-appropriate scheme, e.g.:

```
<domain>.<entity>.<event-name>.v<version>
```
example: `ride.lifecycle.ride_created.v1`, `dispatch.match.driver_accepted.v1`

This gives FamGo's currently single-domain-populated `packages/event-bus/topics/auth.go` a template to replicate across `ride`, `dispatch`, `driver`, `payment`, `safety`, `notification`.

## Pattern 4 — Typed Workflow Transition Graph for Sagas

The Foreman model (conditional/switch/fan-out/fan-in/error/timeout transitions, durable per-step state) is the strongest available reference for populating FamGo's empty `platform/saga/` directory and for deciding what to do with `ride-service/internal/application/saga.go`'s existing inline logic. **Adapt as**: define a small Go type for a workflow step (input state, output state, transition rules), persist each step's state to Postgres (FamGo already has `platform/database/postgres/`), and migrate the ride-lifecycle saga (request → match → accept → start → complete, with cancel/refund as the error-transition path) onto this shared engine instead of leaving saga logic embedded in one service.

## Pattern 5 — Bounded Retry as an Explicit, Caller-Supplied Primitive

`flow.Retry(maxAttempts, initialDelay, multiplier, maxDelay)`, with the retryable condition supplied by the calling code, not inferred by the framework — and a hard rule against self-targeted retry loops with no backoff budget. **Adapt as**: a small Go function in `packages/kafka-sdk/internal/retries/` (currently empty) with the same signature shape, wired to the already-populated `packages/kafka-sdk/consumer/retry_handler.go`.

## Pattern 6 — Type Aliasing for Cross-Service Contract Reuse

When one service's contract needs a type another service owns, alias it rather than redefining it. **Adapt as**: a lint rule or code-review checklist item for FamGo's `packages/event-bus/contracts/*` folders — if `dispatch`'s contract needs a `Ride` status enum, it should reference `ride`'s definition, not redeclare its own copy that can drift.

## Pattern Explicitly Flagged as Needing Modification Before Adoption

- **Parked-backlog DLQ analog with no auto-expiry.** Microbus's own documentation calls this out as a deliberate scope boundary, not a complete solution. FamGo's `packages/kafka-sdk/internal/dlq/` (currently empty) and `standards/events/dlq/` (currently empty) should define an explicit expiry/poison-message threshold from day one — this is the one pattern Microbus itself says is incomplete.
