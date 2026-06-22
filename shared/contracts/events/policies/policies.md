# PHASE 9 — EVENT POLICIES

# STEP 9 — CREATE EVENT GOVERNANCE POLICIES


shared/contracts/events/policies/policies.md

# FILE CONTENT


# Enterprise Event Governance Policies

## Rules

- Events are immutable
- Never mutate existing payload contracts
- All events require schema versions
- All events require trace propagation
- All events require correlation IDs
- All consumers must support idempotency
- All retry exhaustion must go to DLQ
- All services must publish structured envelopes

## Naming

Topic:
domain.events.version

Event:
domain.action.state

Examples:
ride.requested
payment.completed
auth.login.failed

## Versioning

Breaking changes:
- create new version

Never:
- mutate existing version payloads

## Retention

Critical topics:
90 days

Audit topics:
365 days

DLQ:
30 days

Location telemetry:
24 hours