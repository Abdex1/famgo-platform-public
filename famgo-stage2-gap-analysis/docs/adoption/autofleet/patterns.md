# Autofleet/locomotion — Patterns

Source: Stage 1 ride-lifecycle workflow analysis.

## Patterns Worth Adapting

1. **Thin client, opaque backend boundary.** The rider app treats matching/dispatch/pricing entirely as an opaque API response — it does not attempt to replicate or second-guess backend logic client-side. This is a clean separation-of-concerns pattern: the client's job is to render state transitions it receives, not to compute them. Relevant if FamGo's own mobile/web apps (`apps/mobile`, `apps/rider-web`, `apps/driver-web`) are tempted to embed any matching/pricing logic client-side — Autofleet's pattern argues against that.
2. **Payment/billing as an explicitly pluggable integration point**, stated as a design intent rather than discovered as an accident. Validates keeping `packages/payment-sdk/` as a clean boundary in FamGo rather than letting payment logic leak into `ride-service` or `wallet-service` internals.
3. **Two-mode ride request (ASAP / scheduled) that converges to one code path at activation time**, rather than maintaining scheduled rides as a permanently parallel flow. Worth checking against FamGo's own `ride-service` request handling and richxcame's `Scheduler` service (which polls for due scheduled rides) for consistency — three different sources (Autofleet, richxcame, and FamGo's own scheduler-shaped gap) all touch this same modeling question.

## Patterns Not Observable (and therefore not extractable)

- Matching/dispatch algorithm design
- Driver-side acceptance UX or backend trigger
- Any resilience, retry, or error-handling pattern — none documented at the client or backend level in the public repo
