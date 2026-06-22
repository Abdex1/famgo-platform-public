# Microbus (microbus-io/fabric) — Features

Source: Stage 1 event-driven architecture notes (`event-driven-architecture-notes.md`).

## Framing Note

Microbus is a Go microservice fabric built for **agentic (LLM-driven) workflows**, not a ride-hailing backend. Every feature below is extracted for its *mechanism*, to be adapted to FamGo's domain and stack (Kafka, not NATS) — not because Microbus and FamGo solve the same business problem.

## Event-Architecture Features

- **First-class pub/sub events** that decouple producers from consumers — the documented motivating example (a `DeleteUser` handler directly calling `filestoreapi`, `creditcardapi`, `groupmanagerapi`, etc., creating a cyclic, ever-growing dependency graph) is a direct, named illustration of the exact problem FamGo's `packages/event-bus/` is structurally trying to solve.
- **Dedicated event port** separate from request traffic, enabling port-based ACLs (publish-only for the source, subscribe-only for everyone else) — the *authorization-boundary* idea is portable even though the literal port mechanic is NATS-specific.
- **Dynamic, decoupled consumers** — new event subscribers can be deployed without releasing a new version of the producing service, keeping the service dependency graph a DAG.
- **Structured, parseable topic/subject naming** with explicit dimensions for tenancy isolation (plane), trust tier, source, destination, and routing locality.
- **Mechanically-derived ACLs from source code** (`gencreds`) — a tool scans each microservice's actual publish/subscribe call patterns and emits a capability-scoped credentials file, rather than a hand-maintained permission list that can drift from reality.
- **Locality-aware routing** — nested subject slots that let callers narrow toward the most specific live region/zone.
- **Contract-as-code via `api` sub-packages** generating typed client stubs (four shapes: unicast, multicast client, multicast trigger, hook), plus a `manifest.yaml` kept in lockstep with code by codegen, plus auto-published OpenAPI per endpoint.
- **Type aliasing across service boundaries** to avoid duplicating shared types — one service owns a type, others alias it.
- **Orchestrated workflow graphs** (Foreman core service): typed transitions — unconditional, conditional, switch, fan-out, dynamic fan-out (`forEach`), fan-in, error, timeout — executed as a directed graph of task endpoints.
- **Compensation via designated error-transition handler tasks**, with the serialized error placed in that handler's state under `onErr`.
- **Durable, crash-resumable workflow state** — every step's input state and resulting changes are persisted to SQL, creating a full execution history.
- **`goto`-based imperative loop/retry control flow**, distinct from state-driven conditional/switch branching — useful for things like manual-review cycles.
- **Bounded, explicit retry primitive**: `flow.Retry(maxAttempts, initialDelay, multiplier, maxDelay)`, with the retry condition supplied by the calling task, not inferred by the framework.
- **Ack-or-fail-fast transport layer**: a short ack window (250ms default) before the full request timeout (20s default), giving fast failure detection independent of retry logic.
- **Adaptive dispatch-layer controllers**: a rate-limit valve (on `429`, paces down and recovers along a TCP-CUBIC-shaped curve) and a circuit breaker (on `404` ack-timeout/`503`/`529`, parks backlog and probes on an exponential schedule capped at 1 minute).
- **"Parked" backlog as the DLQ-equivalent** — work is held against its original workflow step, released as a "rolling wave" on recovery rather than all at once.

## Explicitly Confirmed Absent in Microbus (do not assume Microbus has these — useful for calibrating effort)

- No "saga" terminology — concept exists, name doesn't.
- No transactional-outbox-for-general-publish pattern — only workflow-step durability.
- No "dead letter queue" terminology or component — only the parked-backlog analog, which has **no auto-expiry**.
