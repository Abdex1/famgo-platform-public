# Autofleet/locomotion — Anti-Patterns

Source: Stage 1 ride-lifecycle workflow analysis.

## Documented Risk

1. **Documentation skeleton mistaken for documentation.** The repository ships a full GitBook-style navigation structure (`SUMMARY.md` referencing `getting-started/`, `faq/faq.md`, `faq/contact-us.md`) where the actual pages are empty section-header stubs. Anyone scanning the file tree without opening the files would reasonably assume real onboarding/support documentation exists. This is the same category of risk flagged for richxcame's unshipped Negotiation service and is directly relevant to FamGo's own structure, which (per the gap analysis) has a comparable number of empty-but-present directories across `platform/`, `infra/`, and several `services/`. **Lesson**: a directory or file existing is not evidence that its content exists — apply this skepticism uniformly, including to FamGo's own repository, not just to third-party repos being evaluated.

## Lower-Severity Notes

- No support/contact workflow exists in the public repo at all, despite being advertised in the navigation. If FamGo is evaluating Autofleet as a reference for support flows (per the original Stage 1 prompt's "Support Flows" requirement), the honest finding is that there is nothing to reference — this should not be silently treated as "no gap found" but as "no source material available," a materially different conclusion that affects how confidently FamGo's own support-flow gap (gap-analysis §7) should be prioritized.
