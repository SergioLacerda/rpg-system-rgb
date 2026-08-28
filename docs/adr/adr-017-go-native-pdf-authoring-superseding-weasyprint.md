# ADR-017: Go-Native PDF Authoring, Superseding the WeasyPrint Exception

## Status

Proposed.

## Context

[ADR-005](adr-005-zero-external-dependency-posture.md) established a zero
external-dependency posture for RGB System V2's Go modules and explicitly named
and rejected headless-Chromium PDF generation once, in mission
`docs-compiler-html-pdf-scope-20260726`.

[ADR-010](adr-010-go-owned-publication-without-python.md) retired the prior
interpreter-backed publication path, making copy/validate/manifest/checksums
Go-owned, and left PDF *authoring* — as opposed to *publication* — without a
replacement.

[ADR-016](adr-016-pdf-external-authoring-exception.md), accepted 2026-08-27,
reinstated a narrow, external, non-Go, non-CI-blocking PDF authoring exception
(MkDocs + WeasyPrint) specifically because no Go-native replacement existed for
WeasyPrint's CSS box-model/paginated layout engine (page breaks, running
headers/footers, generated TOC, `@page` rules), and because the only
Go-reachable alternative of comparable quality — headless Chromium — was the
exact approach ADR-005 had already rejected.

A Strategist evaluation mission on 2026-08-28
(`.analysis/refined/20260828-pdf-python-toolchain-go-migration-evaluation/`)
confirmed the Python toolchain at `tools/pdfbuild/` is real, current, and
intentional, and surveyed the realistic Go-native alternatives:

| Route | What it is | Problem |
|---|---|---|
| `gofpdf`/`fpdf`, `unipdf`, `maroto` | Low-level canvas/drawing APIs | No HTML/CSS input, no paginated flow layout — would mean reimplementing WeasyPrint's layout engine from scratch on top of a drawing API |
| `wkhtmltopdf` via `os/exec` | External binary (WebKit-based) | Still an external, non-Go dependency — same category of tradeoff as Python, and the upstream project is unmaintained |
| `chromedp` (headless Chromium) | Drive real Chrome to print-to-PDF | The exact approach ADR-005 already named and rejected once; ADR-016 implicitly reconfirmed that rejection by choosing WeasyPrint instead on 2026-08-27 |

No route reaches WeasyPrint's current output quality without either (a) writing
a paginated CSS layout engine from scratch in Go, a multi-week-to-months
undertaking with real risk of never reaching parity, or (b) trading the current
external non-Go dependency (Python) for a different one (wkhtmltopdf or
Chromium) — which does not satisfy ADR-005's actual goal of zero external
dependency, it only swaps which external dependency exists.

## Decision

**Not yet made.** This ADR is filed as `Proposed`, not `Accepted`, per the
Strategist mission's own approved scope: drafting this document is
documentation work; accepting it is a maintainer decision this mission does not
make on its own.

The concrete question this ADR must resolve before it can move to `Accepted`
is which route to commit to:

- **(a) Write a Go-native paginated CSS layout engine from scratch** —
  the only route that actually achieves ADR-005's zero-external-dependency
  goal, at the cost of a substantial, open-ended engineering effort with
  meaningful risk of never reaching WeasyPrint's current quality bar (TOC
  pagination, running headers/footers, `@page`-driven cover/callout styling).
- **(b) Accept a different external binary dependency** (e.g.
  `wkhtmltopdf`) in place of the Python/WeasyPrint stack — smaller
  engineering lift, but does not achieve zero-dependency; merely trades one
  external authoring dependency for another, arguably a worse-maintained one.
- **(c) Do not proceed** — leave the ADR-016 exception in place. Python
  authoring stays external, non-Go, non-CI-blocking, exactly as scoped today.

## Consequences

### If (a) is chosen

Positive: genuinely removes the last external-dependency exception in the
project; PDF authoring becomes buildable/testable the same way as the rest of
the Go module graph.

Accepted costs: open-ended implementation effort; regression risk during the
transition (mitigated by the regression test harness scoped in the companion
Strategist mission, `.analysis/refined/20260828-pdf-go-native-renderer-migration-planning/`,
which must be green against the *current* WeasyPrint output before any
renderer cutover); ongoing maintenance burden of an in-house layout engine.

### If (b) is chosen

Positive: smaller, bounded engineering effort; Go owns the authoring
orchestration even if not the rendering itself.

Accepted costs: does not resolve the zero-dependency goal ADR-005 exists for;
trades a well-maintained, actively-developed Python tool (WeasyPrint) for a
less-maintained external binary; would need ADR-005's scope explicitly
re-worded to permit this class of exception rather than only the current
Python one.

### If (c) is chosen

No change. ADR-016 continues to govern PDF authoring. This ADR would be marked
`Rejected` or `Withdrawn` and the companion regression-harness mission's Stage
3/4 tasks (renderer implementation, Python toolchain removal) do not proceed.

## Non-Goals

This ADR does not retroactively invalidate ADR-005's general zero-dependency
posture for anything other than PDF authoring, and does not authorize any
implementation work by itself — renderer code (Stage 3) and toolchain removal
(Stage 4) in the companion planning mission are explicitly blocked on this ADR
reaching `Accepted`, not merely `Proposed`.

## Supersession

If accepted with route (a) or (b), this ADR supersedes ADR-016 for PDF
authoring. ADR-010's Go-owned publication decision (copy/validate/manifest/
checksums) remains unchanged regardless of which route is chosen.
