# ADR-008: Retire The Go HTML/Print/PDF Publication Pipeline

## Status

Accepted.

## Context

The base project Library review found that RGB
Library publication actually happens through MkDocs (Library HTML) → Astro
(static hosting) → generated PDF downloads, and flagged that
`internal/components/library` still described an older, separate
publication concept that could mislead future work into treating it as the
active path.

Investigating the fate of `internal/components/library` surfaced a wider
finding: `internal/components/compiler` (a hand-written ~1,150-line Markdown
parser and HTML/print renderer, with its own tests) and the `cmd/rgb-compiler`
binary formed the same orphaned pipeline. A repository-wide search found zero
consumers of either package's output — no CI workflow, Makefile target
required by `check`/`review-structure`/`check-fast`, MkDocs config, or
Astro/TypeScript file referenced `generated/library/html/**`,
`generated/library/print/**`, `generated/library/PDF_EXPORT.md`, or any
`compiler`/`library` package function. `make compile` and `make
compile-full` were documented in `README.md` but not wired into any gate.
`cmd/rgb-compiler` (unlike `cmd/rgb-tooling`, see
[ADR-006](adr-006-unified-cli-topology.md)) had no real caller either — no
test shelled out to it.

This is a different situation from ADR-006's `cmd/rgb-tooling`/
`cmd/rgb-compiler` deprecation: there, real callers existed
(`tests/semantic_docs/semantic_docs_test.go`, `docs/core/semantic/README.md`)
so the binaries were kept working and only marked deprecated. Here, no real
caller existed for the compile path at all — keeping it working would have
preserved dead weight, not backward compatibility.

## Decision

Retire the Go HTML/print/PDF pipeline outright, not deprecate it:

- Delete `internal/components/library` (`library.go`, `site.go`,
  `pdf_export.go`, and their tests).
- Delete `internal/components/compiler` (Markdown parser, `block.go`,
  `render_html.go`, `render_print.go`, `render_tree.go`, and their tests).
- Delete `internal/app/compile.go` (`app.Compile`).
- Remove the `compile` subcommand from `cmd/rgb` (`cmd/rgb/main.go`); `rgb
  compile ...` now returns `unknown subcommand "compile"`.
- Delete `cmd/rgb-compiler` entirely (no ADR-006-style deprecation — it had
  no real caller to protect).
- Remove `compile`/`compile-full` from `Makefile` and their `README.md`
  documentation.
- Drop `library`/`compiler` from `internal/app.Components()` and the
  architecture diagrams (`docs/architecture/components.md`,
  `docs/architecture/containers.md`).

MkDocs + Astro (`make docs-build`, `make landing-build`, `make docs-pdf`) is
now the only documented publication path for the HTML Library and PDF
downloads — it already was, in practice, before this change.

## Consequences

Positive:

- removes a substantial parallel implementation (custom Markdown parser +
  two renderers + site/PDF-instruction writer) that had zero real
  consumers, eliminating the exact drift risk the Library review warned
  about ("two publication narratives that should be reconciled");
- `internal/app.Components()` and the architecture diagrams now describe
  only components with an active call path from `cmd/*`.

Negative / accepted costs:

- this is a real deletion of working, tested code (not a stub) — reversible
  only through git history, not through a runtime flag or config toggle;
- if a future need for programmatic (non-MkDocs) HTML/print rendering
  emerges, it will be rebuilt rather than resumed from this code.

## Relationship To ADR-006

ADR-006 originally described `cmd/rgb`'s `compile` subcommand and
`cmd/rgb-compiler`'s deprecation-not-removal alongside `rgb-tooling`. This
ADR supersedes that part of ADR-006: `compile` and `cmd/rgb-compiler` are
removed, not deprecated. ADR-006's decision about `cmd/rgb-tooling` (kept
deprecated, real callers) is unaffected and still stands.
