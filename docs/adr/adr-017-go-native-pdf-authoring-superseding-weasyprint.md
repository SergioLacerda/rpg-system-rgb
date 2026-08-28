# ADR-017: Go-Native PDF Authoring, Superseding the WeasyPrint Exception

## Status

Accepted.

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

A Strategist evaluation mission on 2026-08-28 (mission
`20260828-pdf-python-toolchain-go-migration-evaluation`, tracked in the
project's internal analysis workspace, not part of this repository's public
surface) confirmed the Python toolchain at `tools/pdfbuild/` is real, current,
and intentional, and surveyed the realistic Go-native alternatives:

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

Route **(a)** is accepted, staged, and explicitly non-retroactive:

1. **The v2.0 PDF delivery already produced by the Python/MkDocs/WeasyPrint
   toolchain is not reworked or invalidated by this decision.** The
   ADR-016 exception remains in force for the current release; nothing about
   the already-shipped `rgb-system-core-v2-v2.0-*.pdf` assets or the
   `tools/pdfbuild/` toolchain that produced them changes as a result of this
   ADR by itself.
2. **The forward-looking direction is a Go-native paginated CSS/HTML→PDF
   renderer**, written in Go, to eliminate the project's last external,
   non-Go authoring dependency and its long-term maintenance burden (pinned
   Python interpreter, `.venv`, `pip`-installed MkDocs/WeasyPrint/pydyf
   versions, OS-level rendering library prerequisites) — not to reach parity
   with WeasyPrint by any means available, and specifically **not** routes
   (b) (swapping WeasyPrint for `wkhtmltopdf`) or headless Chromium, both
   rejected for the reasons in Context.
3. **This is gated, not immediate**, by the companion Strategist mission
   (`20260828-pdf-go-native-renderer-migration-planning`, same internal
   workspace as above): the regression test harness (Stage 2 / T2) must exist
   and pass against the
   *current* WeasyPrint output before any Go-native renderer code (Stage 3 /
   T3) is written, and the Python toolchain is only removed (Stage 4 / T4)
   after the Go-native renderer reaches parity per that harness.

In short: **keep Python for what has already shipped; migrate to Go going
forward to mitigate future maintenance risk**, not to redo work already
delivered.

## Consequences

Positive:

- Removes the last external, non-Go authoring dependency once Stage 3/4
  complete, satisfying ADR-005's zero-dependency goal in full rather than
  trading one external dependency for another.
- No rework of, or regression risk to, the already-delivered v2.0 PDFs —
  this decision only governs future authoring, not the current release.
- The regression harness (Stage 2) exists before renderer work starts,
  directly bounding the risk of a quality regression during the transition.

Accepted costs:

- Open-ended engineering effort to reach WeasyPrint's current quality bar
  (TOC pagination, running headers/footers, `@page`-driven cover/callout
  styling) in a from-scratch Go layout engine.
- Until Stage 3/4 land, the project continues to carry the Python toolchain
  exactly as ADR-016 scoped it — this ADR does not shrink that surface on
  its own, only commits to eventually retiring it.
- Real risk that the from-scratch engine never fully reaches parity; if that
  risk materializes, the companion mission's Stage 3/4 tasks stall and the
  ADR-016 exception remains the de facto steady state regardless of this
  ADR's `Accepted` status.

## Non-Goals

This ADR does not retroactively invalidate ADR-005's general zero-dependency
posture for anything other than PDF authoring, does not rework or invalidate
the already-shipped v2.0 PDF release, and does not by itself authorize
implementation work — renderer code (Stage 3) and toolchain removal (Stage 4)
in the companion planning mission remain gated on the Stage 2 regression
harness landing and passing first, and both stages are `implementation_handoff`
work that a Strategist mission drafts and hands off but does not execute
itself (see `05-approval-gate.md` / `06-execution.md`).

## Supersession

This ADR supersedes ADR-016 for the *future* PDF authoring direction. It does
not retroactively supersede ADR-016's authorization of the toolchain that
produced the already-shipped v2.0 PDFs — that authorization stands until
Stage 3/4 of the companion migration mission actually ship a validated
replacement. ADR-010's Go-owned publication decision (copy/validate/manifest/
checksums) remains unchanged.
