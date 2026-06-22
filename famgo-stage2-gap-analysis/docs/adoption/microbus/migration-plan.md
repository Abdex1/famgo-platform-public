# Microbus (microbus-io/fabric) — Migration Plan

Source: derived from `features.md`, `contracts.md`, `patterns.md`, `anti-patterns.md` in this folder, cross-referenced with `docs/adoption/gap-analysis.md` §3 and `docs/adoption/migration-matrix.md`.

## Principle

Microbus is architecturally instructive, not code-compatible — different transport (NATS vs. FamGo's Kafka), different domain (agentic workflows vs. ride-hailing). Every step below re-implements a *concept* against FamGo's existing `packages/event-bus/`, `packages/kafka-sdk/`, and `platform/` scaffolding.

## Phase 1 — Populate currently-empty event-bus contract/topic scaffolding

| Step | Action | Target | Depends on |
|---|---|---|---|
| 1.1 | Define the adapted Kafka topic-naming convention (`<domain>.<entity>.<event-name>.v<version>`) as a documented standard | `standards/events/topics/` | none |
| 1.2 | Populate the 5 currently-empty contract domain folders (`dispatch`, `driver`, `notification`, `payment`, `safety`) following the existing `auth`/`ride` examples | `packages/event-bus/contracts/{dispatch,driver,notification,payment,safety}/` | 1.1 |
| 1.3 | Populate topic files for each domain beyond the existing `auth.go` | `packages/event-bus/topics/` | 1.1, 1.2 |
| 1.4 | Apply the type-aliasing rule (one owning service, others alias) as a lint/review checklist item across all contract folders | `packages/event-bus/contracts/` | 1.2 |

## Phase 2 — Saga / workflow engine

| Step | Action | Target | Depends on |
|---|---|---|---|
| 2.1 | Design a minimal typed-transition workflow step type (input state, output state, transition rules: unconditional/conditional/switch/fan-out/fan-in/error/timeout) in Go | `platform/saga/` | Phase 1 contracts (events drive transitions) |
| 2.2 | Persist per-step state to Postgres for crash-resumability | `platform/saga/`, `platform/database/postgres/` | 2.1 |
| 2.3 | Migrate `ride-service`'s existing inline `saga.go`/`events.go`/`event_subscribers.go` logic onto the shared engine | `services/ride-service/internal/application/saga.go`, `platform/saga/` | 2.1, 2.2 |
| 2.4 | Add a default-bounded circuit-breaker/parking lifetime (do **not** copy Microbus's documented unbounded-probe limitation — see anti-patterns.md #1) | `platform/saga/`, `platform/resilience/circuitbreaker/` | 2.1 |

## Phase 3 — Outbox, retries, DLQ

| Step | Action | Target | Depends on |
|---|---|---|---|
| 3.1 | Build a general-purpose transactional outbox (Microbus's own docs admit they don't have one of these — this is FamGo doing more than the reference, not less) | `platform/outbox/`, `packages/kafka-sdk/internal/outbox/` | Phase 1 |
| 3.2 | Implement `flow.Retry`-style bounded retry primitive, wired to the already-populated `consumer/retry_handler.go` | `packages/kafka-sdk/internal/retries/` | none |
| 3.3 | Define DLQ policy **with an explicit expiry/poison-message threshold** (the one thing Microbus's own parked-backlog model admits it lacks — see anti-patterns.md #2) | `standards/events/dlq/`, `packages/kafka-sdk/internal/dlq/` | 3.2 |

## Phase 4 — ACL automation (lower priority, larger effort)

| Step | Action | Target | Depends on |
|---|---|---|---|
| 4.1 | Build a CI-time tool that scans `packages/event-bus/contracts/` usage per service and generates scoped Kafka ACL/SASL bindings | `packages/event-bus/governance/ownership.go`, CI tooling | Phase 1 complete (contracts must exist to scan) |

## Explicitly Out of Scope

- Microbus's literal NATS subject syntax and `:417` dedicated event port — not applicable to Kafka.
- Auto-published OpenAPI-as-LLM-tool-calling surface — no current FamGo agentic use case identified.
- Locality-aware multi-region subject routing — no current FamGo multi-region requirement identified.
