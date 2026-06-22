# Ride Lifecycle Workflows — Autofleet/locomotion

Source: https://github.com/Autofleet/locomotion (branch: `master`), supplemented with publicly described platform behavior from Autofleet's own product pages (autofleet.io) where the repository itself defers to Autofleet's backend.

## Scope Note (read first)

`Autofleet/locomotion` is **a rider/passenger-facing mobile client only** (React Native + a thin Node.js layer), not a full ride-hailing backend. Per the repository's own README, the app is "pre-integrated to autofleet's backend ride engine that matches rides to available vehicles from partner fleets," and accessing that engine requires a partner API key obtained directly from Autofleet — it is not part of this open-source repository. Concretely:

- The repository's documentation directory (`getting-started/`, `faq/`, `contributing/`) is a GitBook-style skeleton; the requirements, configuration, and FAQ pages exist only as empty section headers in the current `master` branch — they contain no lifecycle detail.
- The dispatch/matching engine, pricing engine, and driver-side application are **not present in this repository** at all; they are Autofleet's proprietary, closed platform, reached only through an authenticated API.
- What follows is organized into (a) what the rider app itself does at each stage, based on the README's own description of its screens, and (b) what Autofleet's platform is described as doing at that stage, based on Autofleet's public product material — clearly separated so the boundary between "in this repo" and "outside this repo" stays visible.

## Workflow: Request Ride

**In the app (this repo):** The rider opens the app to a "Where to?" screen, enters a destination, and submits a ride request. This is the first of the four flow screenshots shown in the README (captioned "Where to?").

**Outside this repo (Autofleet platform):** The request is sent to Autofleet's backend ride engine via the partner API. Per Autofleet's own platform description, requests can be either immediate ("ASAP") or scheduled for a future time — a distinction also visible in the project's recent commit history (an "Actions" log entry references "when future ride becomes active – treat as ASAP" logic), confirming the client supports both modes even though the underlying scheduling logic lives server-side.

## Workflow: Driver Matching

**In the app (this repo):** Not implemented client-side. The rider app's role is to display the result of matching (ride offers), not to compute it.

**Outside this repo (Autofleet platform):** Autofleet describes this stage as "real time responsive ride dispatching for any order, vehicle, or driver constraint," using an optimization layer to select among available partner-fleet vehicles and compute pricing based on pickup/dropoff location and supply/demand conditions. None of this matching or pricing logic is visible in the `locomotion` repository — it is consumed by the app as an opaque API response.

## Workflow: Driver Acceptance

**In the app (this repo):** The second README screenshot ("Get ride offers") shows the app surfacing one or more ride offers to the rider after the backend has identified candidate vehicles. The rider-side flow is to receive and act on an offer, not to negotiate or directly contact a specific driver.

**Outside this repo (Autofleet platform):** Driver/vehicle acceptance of an assigned trip happens within Autofleet's separate fleet/driver-management tooling, which is not part of the `locomotion` codebase (the README explicitly scopes `locomotion` as the rider/passenger app only — there is no driver-app code in this repository).

## Workflow: Arrival

**In the app (this repo):** The README's roadmap-level feature list (carried over from the project's earlier description and still reflected in the current screen flow) includes live vehicle tracking and ride status updates, implying an in-app map view showing the assigned vehicle approaching the pickup point. This stage is not broken out as a separate documented screen in the current README, but is implied by the "In-Ride" tracking capability the project advertises.

**Outside this repo (Autofleet platform):** Real-time vehicle position is supplied by Autofleet's backend (the same engine responsible for matching), which the client polls or subscribes to for tracking data. The transport mechanism (REST polling vs. push) is not specified in the documentation reviewed.

## Workflow: Pickup

**In the app (this repo):** The third README screenshot ("Ride booked") corresponds to the point where a ride offer has been accepted and a vehicle is assigned/en route; the app's job at this boundary is to confirm the booking and transition the rider's view from "searching/offers" to "tracking an assigned ride."

**Outside this repo (Autofleet platform):** Confirmation that the rider has actually boarded the vehicle (trip-start trigger) is a backend/driver-app concern not represented in this repository.

## Workflow: Ride Start

**In the app (this repo):** Once a ride transitions to in-progress, the advertised capabilities are "ride status" and "ride updates" — i.e., the app continues to reflect backend-reported state changes rather than initiating the state transition itself.

**Outside this repo (Autofleet platform):** The actual trip-start event is generated by Autofleet's dispatch/trip-management system (or the driver-side app, which is outside this repository) and pushed or polled into the rider app's state.

## Workflow: Ride Completion

**In the app (this repo):** The README's feature list (and the fourth screenshot, "Add payment method") indicates a post-ride flow covering payment and receipt. The advertised capabilities include displaying a ride summary, prompting the rider for a driver rating, and supporting comments on the driver, along with a booking-history view of past trips and receipts.

**Outside this repo (Autofleet platform):** Fare finalization, payment capture, and receipt generation are processed through whichever payment/billing integration the operator has configured — the README states this integration point is intentionally left pluggable ("payment/billing services... can be easily built into Locomotion") rather than shipped as a single fixed implementation.

## Workflow: Cancellation

**In the app (this repo):** Ride cancellation is listed among the app's "In-Ride" capabilities (alongside tracking and status updates), implying a rider-initiated cancel action is exposed somewhere in the ride-tracking screen, though the specific point(s) in the flow at which cancellation is allowed (pre-acceptance, post-acceptance, en route) are not documented in the README or the GitBook stub pages.

**Outside this repo (Autofleet platform):** Cancellation policy (fees, time windows, driver-side notification) would be enforced by the backend engine and is not documented in this repository.

## Workflow: Support Flows

**In the app (this repo):** The repository's documentation skeleton includes a dedicated "Contact Us" page (`faq/contact-us.md`) and an "FAQ" page (`faq/faq.md`) as placeholders in the navigation structure (`SUMMARY.md`), but **both are currently empty stubs** containing only their section header — no actual support workflow, contact channel, or escalation path is documented in the repository as it stands.

**Outside this repo (Autofleet platform):** No public documentation of an in-app support/help flow (e.g., driver complaint reporting, lost-item flow, fare dispute) was found either in the repository or in the Autofleet product pages reviewed for this analysis.

## Summary Table

| Stage | Implemented in `locomotion` (rider app) | Implemented in Autofleet's backend (outside repo) | Documented in repo at a workflow level? |
|---|---|---|---|
| Request Ride | Yes — "Where to?" entry screen | Receives request, supports ASAP/scheduled modes | Partially (screenshot only) |
| Driver Matching | No | Yes — dispatch/optimization engine | No |
| Driver Acceptance | Displays resulting offer | Driver/fleet acceptance via separate tooling | No |
| Arrival | Implied via live tracking capability | Supplies real-time position data | No |
| Pickup | "Ride booked" confirmation screen | Trip-start trigger | No |
| Ride Start | Reflects backend-reported status | Generates the actual start event | No |
| Ride Completion | Summary, rating, comments, history; payment screen | Fare finalization, payment processing | Partially (screenshot + roadmap list) |
| Cancellation | Listed as an "In-Ride" capability | Cancellation policy enforcement | No (no detail beyond the capability name) |
| Support Flows | Placeholder FAQ/Contact pages (empty) | Not found | No |

## Overall Finding

The ride lifecycle for this product is **split across a public, thin client repository and a private, partner-gated backend platform**. The `locomotion` repository documents what the rider sees (four screen captures plus a short capability list) but does not document — and does not contain the code for — the matching, dispatch, pricing, or driver-acceptance logic that actually drives the lifecycle. Any adoption planning that assumes a self-contained, fully-documented ride lifecycle in this repository would be working from an incomplete picture; the authoritative source for stages 2–6 above would be Autofleet's partner API documentation, which is not publicly accessible without a granted API key.
