# richxcame/ride-hailing — Anti-Patterns

Source: Stage 1 inventories. These are explicitly documented problems or risks in richxcame's own repository — not speculation. Listed so FamGo does not accidentally reproduce them while adapting the patterns in `patterns.md`.

## Confirmed, Documented Risks

1. **Unauthenticated internal endpoints relying on network-layer trust only.** The two Realtime broadcast routes (`/internal/broadcast/ride`, `/internal/broadcast/user`) have no application-level auth middleware; the project's own docs recommend mTLS/network ACLs instead. This is a single point of failure if network segmentation is ever misconfigured — application-layer auth should be defense-in-depth, not optional. **FamGo audit item**: gap-analysis §2/§5 flags this exact pattern as something to verify is not present in `services/websocket-gateway/` or `gateway/`.
2. **Payment webhook signature verification skipped.** The Stripe webhook handler validates payload *shape* but not the cryptographic *signature*, meaning a forged request matching the expected JSON shape could be accepted as a genuine Stripe event. This is a textbook payment-fraud vector.
3. **A documented-but-unshipped service** (Negotiation) left in the README and port table with no corresponding code, compose entry, or k8s manifest. Creates a misleading impression of platform capability for anyone reading the docs without cross-checking the actual deployment.
4. **Two incompatible pagination conventions coexisting** rather than one being deprecated — increases client-integration cost and cognitive overhead for every new consumer of the API, with no documented plan to converge.
5. **A fully provisioned, fully infrastructure-ready message bus (NATS + JetStream) sitting disabled by default**, with zero application code wired to it. This is worse than not provisioning it at all in one specific way: it implies an event-driven capability to anyone reading the infrastructure that does not actually exist at the application layer — exactly the "scaffolded vs. delivered" gap this whole Stage 2/3/4 process is designed to catch in FamGo's own repo too.
6. **Conflating "admin user" and "trusted service caller" in one role system** via a seeded admin/service JWT used for service-to-service calls. If that seeded account's credentials ever leak, the blast radius is "anything an admin can do," not a narrowly-scoped service identity.

## Lower-Severity Hygiene Notes

- Read replica topology is mentioned in marketing copy (README tech-stack line) but never actually documented or configured anywhere else in the reviewed files — a gap between stated and demonstrated capability, smaller in stakes than the items above but the same category of risk.
