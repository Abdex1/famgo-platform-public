# Microbus (microbus-io/fabric) — Contracts

Source: Stage 1 event-driven architecture notes.

## Contract Model

| Element | Mechanism |
|---|---|
| Contract unit | A microservice's `api` sub-package (e.g., `calculator/calculatorapi`) — typed request/response structs |
| Generated artifacts | Client stubs in four shapes: `Client` (unicast), `MulticastClient` (fan-out, e.g. service discovery), `MulticastTrigger` (source fires its own events), `Hook` (downstream registers as an event sink) |
| Cross-service type reuse | **Aliasing, not duplication** — `type X = Aapi.X` when service B needs a type service A owns |
| Parallel contract artifact | `manifest.yaml` — short per-service description kept "in lockstep" with code by the framework's codegen, specifically to prevent contract/implementation drift |
| External/LLM-facing contract | Every functional endpoint, web handler, and workflow is auto-published as an OpenAPI operation — doubles as an LLM tool-calling surface; auth and required claims flow through per call |

## Topic/Subject Contract

Every request and event is routed over a structured NATS subject:

```
<plane>.<trust>.<port>.<src>.<dest>.<id_or_locality>.<method>.<path...>
```

| Segment | Purpose |
|---|---|
| `plane` | Isolates tenants/test runs sharing one NATS cluster |
| `trust` | `safe` / `danger` (`:666` trust-root tier) / `reply` |
| `src` / `dest` | Flattened, validated hostnames |
| `id_or_locality` | Optionally pins to a replica (`id-XXXX`) or locality slot (`loc-region-zone-...`, broadest-first) |
| `method` / `path` | Mirror the HTTP request being carried |

ACLs for this subject space are **derived mechanically from source code at deploy time** by `gencreds`, scanning each microservice's actual publish/subscribe call sites and emitting a signed `.creds` file scoped to only those subjects — a capability allow-list, not a manually maintained one.

## Event Contract Shape

- An event is structurally a **multicast request published to a URL on the event source's own hostname**; sinks subscribe to that hostname.
- Events default to a **dedicated port `:417`** ("force eventing"), distinct from standard `:443` request traffic, enabling a NATS ACL where only the source can `PUB` and everyone else can only `SUB`.
- Because events are structurally requests, they **can return values to the publisher** (used in Microbus's own example to let subscribers grant/deny permission for an action) — a notable departure from typical fire-and-forget event semantics.

## What FamGo Should NOT Treat as a Literal Contract to Copy

- The `<plane>.<trust>.<port>.<src>.<dest>.<id_or_locality>.<method>.<path...>` subject syntax is **NATS-specific** — Kafka topic naming needs its own convention (see `patterns.md` for an adapted proposal).
- The `:417` dedicated event port mechanic is tied to Microbus's own transport layer — Kafka's equivalent boundary is topic-level ACLs, not a second port.
