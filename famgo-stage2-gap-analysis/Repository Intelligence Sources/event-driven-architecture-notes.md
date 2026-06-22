# Event-Driven Architecture Notes — Microbus

Source: https://www.microbus.io/, https://docs.microbus.io/ (Microservice Substrate, Agentic Workflows, Agentic RAD, Security in Depth sections), https://github.com/microbus-io/fabric.

## Framing

Microbus is a Go microservice fabric whose stated purpose is running **agentic (LLM-driven) workflows** in production, not a general-purpose ride-hailing or e-commerce backend. Its event-driven mechanics exist to decouple microservices and to make multi-step, long-running workflows durable. The notes below map the requested categories (Events, Contracts, Topics, Sagas, Outbox, Retries, Dead Letters) onto what Microbus's own documentation actually describes — including categories where Microbus uses a different name or mechanism for a conceptually similar problem, and categories where no equivalent was found.

## Events

- Implemented as **first-class pub/sub on the message bus**, not as direct service-to-service calls. The documentation's own framing example: a `DeleteUser` handler that directly calls `filestoreapi`, `creditcardapi`, `groupmanagerapi`, etc. creates a cyclic dependency graph that grows every time a new consumer is added; publishing an `OnUserDeleted` event instead lets consumers subscribe without the producer knowing they exist.
- Mechanically, an event is **a multicast request published to a URL on the event source's own hostname**; event sinks subscribe to that hostname rather than their own. Because they are structurally just requests, events can return values to the publisher (used in the framework's own example to let subscribers grant/deny permission for an action).
- Events default to a **dedicated port, `:417`** ("force eventing"), separate from standard request traffic on `:443`. This lets operators write a NATS ACL that allows only the event source to `PUB` on `hostname:417` while every other microservice may only `SUB` — an authorization boundary expressed entirely through port-based ACLs rather than application logic.
- Consumers are dynamic: new subscribers can be deployed without releasing a new version of the producing microservice, which is the documentation's explicit goal — keeping the microservice dependency graph a DAG rather than a tangle of mutual dependencies.

## Contracts

- The unit of contract is the **microservice's `api` sub-package** (e.g., a `calculator` microservice's contract lives in `calculator/calculatorapi`). This package defines the typed request/response structs and generates **client stubs** — type-safe wrapper functions standing in for raw `GET`/`POST`/`Request`/`Publish` bus calls.
- Four stub shapes are generated depending on call semantics: `Client` (unicast), `MulticastClient` (multicast/fan-out calls such as service discovery), `MulticastTrigger` (used by the source microservice to fire its own events), and `Hook` (used by downstream microservices to register as event sinks).
- Cross-microservice type reuse is handled by **aliasing**, not duplication: if microservice B's endpoint accepts a type owned by microservice A, B declares `type X = Aapi.X` rather than redefining the shape, keeping a single source of truth for the contract.
- A separate, parallel contract artifact is **`manifest.yaml`** — a short per-microservice YAML description that the framework's codegen keeps "in lockstep" with the actual code, specifically so that contract and implementation cannot drift apart silently.
- Every functional endpoint, web handler, and workflow is additionally **auto-published as an OpenAPI operation**, which doubles as the contract surface for LLM tool-calling: passing a list of endpoint URLs to an LLM client is sufficient to expose them as callable tools, with auth and required claims flowing through per call.

## Topics

- There is no exposed "topic name" abstraction distinct from addressing — Microbus routes every request and event over a **structured NATS subject** that the framework constructs and parses automatically:
  `<plane>.<trust>.<port>.<src>.<dest>.<id_or_locality>.<method>.<path...>`
- Each segment is purpose-built: `plane` isolates tenants/test runs sharing one NATS cluster; `trust` collapses access tier into `safe` / `danger` (the `:666` trust-root tier) / `reply`; `src`/`dest` are flattened, validated hostnames; `id_or_locality` optionally pins a request to a specific replica (`id-XXXX`) or to a locality slot (`loc-region-zone-...`, broadest-first); `method` and `path` mirror the HTTP request being carried.
- **NATS ACLs are derived mechanically from source code at deploy time** by a `gencreds` tool, which scans each microservice's call patterns and emits a signed `.creds` file containing only the subject patterns that microservice's code actually publishes to or subscribes from — a capability-style allow-list rather than a manually maintained topic permission list.
- **Locality-aware routing** is itself implemented as a topic-addressing feature: a locality such as `us-west-b-1` registers nested subject slots (`loc-us`, `loc-us-west`, `loc-us-west-b`, `loc-us-west-b-1`), and callers narrow toward the most specific slot that still has live subscribers.
- A **short-circuit transport** bypasses NATS entirely when two connectors share a process, using a process-global trie that accepts the same subject syntax — so the "topic" addressing scheme is uniform whether or not a broker is actually involved, at the cost of the in-process path being publisher-attested rather than broker-verified.

## Sagas

- **The term "saga" does not appear in Microbus's documentation.** The closest structural equivalent is the **workflow graph** executed by the Foreman core service: a directed graph of task endpoints connected by typed transitions (unconditional, conditional, switch, goto, fan-out, dynamic fan-out/`forEach`, fan-in, error, timeout).
- Compensation-style logic is expressed through the **error transition**: when a task returns an error, the flow diverts to a designated handler task with the serialized error placed in that handler's state under the key `onErr`. The documentation explicitly frames this handler as the place to "compensate, log, or retry through an alternative path" — i.e., compensating-transaction logic is something the workflow author writes inside an error-handler task, not a built-in saga primitive.
- Long-running, multi-service coordination is durable by construction: each workflow step's input state and resulting state changes are persisted to a SQL database by the Foreman, so a crashed flow resumes exactly where it left off rather than restarting or silently dropping a step.
- Loop/retry-style control flow within a graph (e.g., a manual-review cycle) is expressed with **`goto`** transitions, which a task triggers imperatively rather than the Foreman evaluating them from state — distinguishing author-driven branching from state-driven (conditional/switch) branching.
- There is no orchestration-vs-choreography distinction documented as such; Microbus's model is **orchestration only** — a single Foreman owns graph execution, fan-out/fan-in, and state merging, in contrast to a choreographed saga where services react to each other's events with no central coordinator.

## Outbox

- **The term "outbox" does not appear in Microbus's documentation**, and there is no documented transactional-outbox-plus-relay pattern for guaranteeing atomic "write to DB + publish message."
- The closest analog is the Foreman's own durability mechanism: workflow **state is persisted to a SQL database** ("the Sequel library") on every step, recording both the input state and the changes each step produced — described as creating "a full history of the flow's execution." This guarantees workflow durability across process restarts but is scoped to the workflow engine's own step log, not a general-purpose mechanism for any microservice to atomically pair a local database write with an outbound bus message.
- Plain (non-workflow) event publication does not appear to go through any durable staging table — events are described purely as multicast pub/sub requests over NATS, with no documented at-least-once or outbox-backed delivery guarantee at that layer.

## Retries

- Retries are **explicit and bounded by construction**, not automatic. The single primitive is `flow.Retry(maxAttempts, initialDelay, multiplier, maxDelay)`, called from inside a task body; it returns `true` while attempts remain and `false` once the budget is exhausted. The task itself supplies the retryable condition (e.g., gating retry on an HTTP `408` status) — the framework does not infer what is retryable.
- The documentation explicitly forbids a simpler-looking alternative: `OnError`/`OnTimeout` transitions **cannot target the same task they originated from** ("the graph validator rejects self-targeted error transitions because they would retry indefinitely with no backoff budget"), forcing all retry logic through the bounded `flow.Retry` primitive rather than an unbounded loop in the graph itself.
- Separately, at the **transport layer**, every bus call uses an **ack-or-fail-fast** mechanic: the downstream connector sends an ack (`100 Continue`) before forwarding a request to the handler; if no ack arrives within a short ack timeout (250ms default), the upstream connector fails fast rather than waiting out the full request timeout (20s default) — a fast-failure mechanism, not itself a retry.
- At the **workflow-dispatch layer**, two automatic controllers react to downstream health without the workflow author writing retry code: a **rate-limit valve** (triggered by `429`) paces dispatch down and recovers it along a TCP-CUBIC-shaped curve; a **circuit breaker** (triggered by a `404` ack-timeout, `503`, or `529`) parks the task's pending backlog and probes the endpoint on an exponential schedule (100ms, doubling, capped at 1 minute) until it recovers.
- A noted limitation, stated directly in the docs: the circuit breaker has **no auto-give-up** — "the Foreman probes a tripped breaker indefinitely, so a task whose endpoint is permanently dead will park forever." Bounding that exposure (via flow lifetime or an explicit timeout) is left to the workflow author.

## Dead Letters

- **The term "dead letter" / "dead-letter queue" does not appear in Microbus's documentation**, and no DLQ-equivalent component was found in the package reference or transport documentation reviewed.
- The closest functional analog is the **circuit breaker's "parked" state**: when a task's downstream is unreachable, its pending backlog is removed from the active selection index and held (`parked` column on a `microbus_steps` table) rather than discarded, while a single probe per shard checks for recovery. On recovery, the parked backlog is released as a "rolling wave" rather than all at once.
- This differs from a conventional dead-letter queue in two material ways: (1) parked work is not redirected to a separate inspectable channel for manual triage — it stays attached to its original workflow step; and (2) there is **no automatic expiry or poison-message threshold** — a permanently broken downstream results in permanently parked work unless the workflow author has bounded the flow with a timeout. Microbus's design notes treat this as a deliberate scope boundary (bounding is "the workflow author's responsibility"), not an oversight to be patched.

## Summary Table

| Requested category | Native Microbus term/mechanism | Notes |
|---|---|---|
| Events | "Events" (first-class pub/sub, port `:417`) | Directly documented, dedicated page |
| Contracts | `api` sub-packages, client stubs, `manifest.yaml`, auto-generated OpenAPI | Directly documented across several pages |
| Topics | NATS subject structure (`<plane>.<trust>.<port>.<src>.<dest>.<id_or_locality>.<method>.<path...>`) | Directly documented, dedicated page; ACLs derived from source code |
| Sagas | Workflow graph + error transitions + `goto` loops (Foreman-orchestrated) | Conceptually present; term "saga" not used |
| Outbox | Foreman's per-step SQL state persistence | Conceptually adjacent (durability of workflow steps) but not a transactional-outbox-for-publishing pattern; term not used |
| Retries | `flow.Retry`, ack-or-fail-fast, rate-limit valve, circuit breaker | Directly documented, explicit and bounded by design |
| Dead Letters | Circuit-breaker "parked" backlog | Conceptually adjacent but explicitly has no auto-expiry; term not used |
