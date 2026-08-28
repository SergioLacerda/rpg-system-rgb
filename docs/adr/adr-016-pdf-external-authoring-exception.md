# ADR-016: Reinstate A Narrow External PDF Authoring Exception

## Status

Accepted.

## Context

A 2026-08-27 PDF editorial-quality review re-examined a pending critique
that proposed adding a headless-Chromium/Playwright HTML-to-PDF pipeline to
achieve "editorial quality" (cover, table of contents, RGB-vector rule
callout boxes) for the RGB Library PDF.

Two things emerged from direct verification, not just from reading the
critique:

1. **Most of the claimed gap does not exist.** Inspecting the currently
   shipped PDF (`web/landing/public/downloads/rgb-system-core-v2-*.pdf`)
   directly — via `pdftotext`, `pdftohtml -xml`, and a rasterized page
   render — shows it already has a correctly paginated table of contents,
   146 working internal hyperlinks on pages 1–3 alone, and running
   headers/footers (`page/total`). Only two things are actually missing: a
   designed cover (no RGB vector glyph, no version/license/date metadata —
   plain text only) and colored R/G/B rule callout boxes.
2. **The proposed fix conflicted with existing decisions.** The
   headless-Chromium pipeline directly conflicts with
   [ADR-005](adr-005-zero-external-dependency-posture.md) (which already
   named and rejected headless-Chromium PDF generation once, in mission
   `docs-compiler-html-pdf-scope-20260726`) and
   [ADR-010](adr-010-go-owned-publication-without-python.md) (which retired
   the interpreter-backed publication path and states PDF publication is "a
   Go-owned publishing step for reviewed PDF assets, not a full
   general-purpose PDF layout engine").

Deeper verification found the real problem underneath: the shipped PDF's own
metadata reports `Producer: WeasyPrint 62.3`, and `git log` traces it to a
pre-ADR-010 commit (`c2d2a5e`, "implement automated PDF generation workflow
with MkDocs"). A repository-wide search finds zero remaining trace of that
toolchain — no `mkdocs.yml`, no `requirements-docs*`, no `weasyprint` or
`mkdocs-pdf` reference anywhere. The tool that produced today's
good-quality PDF body was fully deleted when ADR-010 retired the automated
publication path, without a documented replacement for *authoring* PDF
content (as opposed to *publishing* it, which ADR-010's Go-owned pipeline
still does correctly). Nobody today has a reproducible way to produce the
next PDF release at all, cover/callouts aside.

ADR-005's original 2026-07-31 addendum already permitted exactly this class
of exception (Python documentation-build dependencies, OS rendering
libraries, a dedicated `docs-pdf`/`pdf-build` Make target, generated PDFs
under the landing download surface) before ADR-010 retired it on the premise
that the Go-owned replacement had reached parity. That parity claim holds for
copy/validate/manifest/checksums, but not for authoring the PDF body's
editorial content — nothing replaced MkDocs+WeasyPrint's role there.

## Decision

Reinstate a **narrow, external, non-Go, non-CI-blocking** PDF authoring
exception:

- A PDF authoring tool of this class (MkDocs + WeasyPrint, or an equivalent
  the maintainer chooses) may run **outside** the Go module graph and
  outside the default `make`/CI build, as a dedicated, manually-invoked
  release step — producing one reviewed PDF per release, per language.
- That authoring step's template must add what's actually missing: a
  designed cover (RGB vector glyph, version, language, license, release
  date) and CSS-based callout/admonition styling for R/G/B rule boxes. It
  does not need to reproduce anything else — the existing TOC pagination,
  internal links, and running headers/footers already work and must be
  preserved, not rebuilt.
- The resulting reviewed PDF is supplied to the existing, unchanged Go
  pipeline via `cmd/rgb docs pdf --source-en`/`--source-pt-br`
  (`make docs-pdf`), exactly as today. `internal/components/publication/
  pdf.go`'s copy-and-validate behavior does not change.
- `internal/components/tooling/release_artifacts.go`'s Poppler-based
  editorial smoke checks (`pdftotext`/`pdftohtml`/`pdftoppm`) remain the
  Go-owned validation gate over whatever the external tool produces.

This does **not** allow:

- any new Go module dependency (`go.sum` entry) for PDF generation or
  manipulation;
- headless-browser automation as part of the Go/CI build;
- treating the external authoring tool's output as canonical documentation
  — canonical source remains `docs/core/**` Markdown, unchanged (ADR-001);
- making any default build target (`make build`, `make test`, `make
  release-check`'s non-PDF gates) depend on the external tool's presence.

## Consequences

Positive:

- Closes a real reproducibility gap: the next PDF release has a documented
  path again, instead of depending on tooling that was silently deleted.
- Scopes the actual work correctly: template/cover/callout design, not a new
  pipeline or a new class of Go/CI dependency.
- Keeps ADR-010's substance intact — Go still owns copy, validate, manifest,
  and checksums; only PDF *authoring* (never publication) sits outside Go.
- Prevents re-litigating the headless-Chromium-vs-ADR-005 conflict from
  scratch the next time PDF editorial quality comes up.

Negative / accepted costs:

- Reintroduces an external, non-Go toolchain dependency for PDF authoring —
  the same tradeoff ADR-005's original addendum already accepted, now scoped
  more narrowly (cover + callout templates only, not full pipeline
  automation).
- Setting up the authoring tool, and designing the cover/callout templates,
  remain open implementation/authoring work — out of scope for this ADR and
  for the Strategist mission that produced it.

## Non-Goals

This ADR does not authorize any Go module dependency, any change to
`internal/components/publication/**` or `internal/components/tooling/**`, or
any change to the canonical Markdown authority model (ADR-001). It records
the authoring-exception decision only; the toolchain setup and the actual
next PDF release are separate, unauthorized-by-this-ADR implementation work.

## Supersession

This ADR narrows and partially reinstates the addendum ADR-010 retired —
scoped specifically to PDF *authoring* (cover/callout templates), not to the
full automated interpreter-backed publication pipeline ADR-010 removed.
ADR-010's Go-owned publication decision remains otherwise unchanged.
