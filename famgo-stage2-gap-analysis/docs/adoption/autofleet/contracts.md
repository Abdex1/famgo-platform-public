# Autofleet/locomotion — Contracts

Source: Stage 1 ride-lifecycle workflow analysis.

## No API Contracts Are Documented in This Repository

The `locomotion` repository consumes Autofleet's backend ride engine through a **partner API that requires a separately-issued API key**. That API's request/response shapes, authentication scheme, and error model are not published in the public repository and were not accessible during Stage 1 research. Nothing in this folder can responsibly assert a concrete contract (endpoint paths, payload fields, status codes) for that backend.

## The Only Contract-Adjacent Fact Confirmed by the Repository

- The client supports **two request modes** for ride creation: immediate ("ASAP") and scheduled-for-future. This is the one shape-level detail confirmed both by Autofleet's own product description and corroborated independently by a commit-history log entry in the repository itself ("when future ride becomes active – treat as ASAP"). Any FamGo ride-request contract that wants to borrow this distinction should treat it as a **two-mode request type**, with scheduled requests converting to ASAP semantics at activation time rather than remaining a separate code path indefinitely.

## Recommendation

Do not write a `contracts.md` entry that resembles a real API spec for this source — there isn't one to extract. If FamGo needs a comparable contract reference for ride request/matching/offer/acceptance, `richxcame`'s API inventory (`docs/adoption/richxcame/contracts.md`) is the usable source; this folder exists for completeness, not because it has equivalent content to offer.
