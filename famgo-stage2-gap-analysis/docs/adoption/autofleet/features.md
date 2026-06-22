# Autofleet/locomotion — Features

Source: Stage 1 ride-lifecycle workflow analysis (`ride-lifecycle-workflows.md`).

## Scope Warning (read first)

`Autofleet/locomotion` is a **rider-facing mobile client only**. The matching/dispatch engine, pricing engine, and driver-side app are not in the public repository — they live behind Autofleet's partner API, which Stage 1 could not access. This extraction package is therefore necessarily thin. It documents what the repo *actually shows*, not what a full ride-hailing backend would need; padding it out with invented backend detail would misrepresent the source.

## Features Actually Observable in the Repository

- **Four-screen rider flow**: "Where to?" (request) → "Get ride offers" (matching result display) → "Ride booked" (acceptance/assignment confirmation) → "Add payment method" (completion/payment).
- **ASAP vs. scheduled ride request modes** — confirmed both by Autofleet's product description and by a commit-history reference to "when future ride becomes active – treat as ASAP" logic, i.e., scheduled rides convert to ASAP semantics once their time arrives, rather than being a permanently distinct request type.
- **In-Ride capability bundle**: live vehicle tracking, ride status updates, and cancellation are grouped together in the app's advertised feature list, implying these are handled by one tracking/status subsystem in the client rather than separate modules.
- **Post-ride capability bundle**: ride summary, driver rating, driver comments, booking history with receipts.
- **Pluggable payment/billing integration point** — the README explicitly states payment/billing services "can be easily built into Locomotion," i.e., the app is designed with payment as a swappable integration rather than a single fixed implementation.
- **GitBook-style documentation skeleton** (`getting-started/`, `faq/`, `contributing/`) — a reasonable docs scaffold shape, even though the content itself is currently empty stubs.

## Features NOT in the Repository (do not assume these exist anywhere adoptable)

- Driver matching / dispatch algorithm
- Pricing / fare computation
- Driver-side application of any kind
- Trip-start trigger logic
- Cancellation policy enforcement (fees, time windows)
- Any support/FAQ/contact content (pages exist but are empty)
