# Autofleet/locomotion — Migration Plan

Source: derived from `features.md`, `contracts.md`, `patterns.md`, `anti-patterns.md` in this folder.

## Honest Scope Statement

There is very little to migrate from this source. The repository is a thin rider client with no accessible backend. This migration plan reflects that reality rather than manufacturing work to fill the template.

## Phase 1 — The only concrete actions this source justifies

| Step | Action | Target | Depends on |
|---|---|---|---|
| 1.1 | Confirm `packages/payment-sdk/` remains a clean, swappable boundary (validates existing FamGo design; no new code required if already true) | `packages/payment-sdk/` | code review |
| 1.2 | Verify FamGo's client apps (`apps/mobile`, `apps/rider-web`, `apps/driver-web`) do not embed matching/pricing logic client-side — keep them as renderers of backend state, consistent with Autofleet's thin-client pattern | `apps/mobile/`, `apps/rider-web/`, `apps/driver-web/` | code review |
| 1.3 | Cross-check the "scheduled ride converges to ASAP at activation" modeling decision against FamGo's `ride-service` request handling and the planned `scheduler-service` (gap-analysis §1) for consistency | `services/ride-service/`, new `services/scheduler-service/` | scheduler-service build-out |

## Explicitly Not a Phase 2 or 3

There is no further phased work to extract from Autofleet — the source material is exhausted by Phase 1. Do not invent additional migration steps to make this package symmetrical with `richxcame/migration-plan.md` or `microbus/migration-plan.md`; their greater depth reflects their source repositories actually containing more to extract, not a template requirement.

## Where to Look Instead

For ride-lifecycle backend patterns (matching, dispatch, driver acceptance, support flows), `docs/adoption/richxcame/` is the relevant extraction package — Autofleet's own Stage 1 analysis confirms it has nothing usable in those areas.
