# ADR-005: Zero External Dependency Posture

## Status

Accepted, with a documentation-PDF exception added on 2026-07-31.

## Context

RGB System V2 has never imported anything outside the Go standard
library, in any of its Go modules, at any point in its history — no
`go.sum` exists anywhere in the repository. This was never decided in a
single place; it emerged as a consistent pattern across independent
missions, each reasoning locally that a proposed dependency would be "the
first" of its kind and rejecting it on that basis:

- `docs-compiler-html-pdf-scope-20260726` rejected automating PDF
  generation via headless-Chromium (an external binary requirement) or a
  Go-native PDF library (a new Go module dependency), keeping PDF export
  a documented manual browser step instead.
- `guardrails-tooling-scope-20260726` rejected adopting `gocognit`
  (cognitive-complexity linting) specifically because it requires
  installing a new external Go tool, even though the project already had
  build-failing cyclomatic-complexity coverage via `gocyclo` (bundled
  with the already-required `golangci-lint` dev tool, not a module
  dependency).

Each decision was locally reasoned and undocumented as a standing rule.
This ADR formalizes the pattern those decisions were already following,
so future missions can cite one authoritative source instead of
re-deriving the same reasoning independently each time.

## Decision

**RGB System V2's Go modules do not take on external dependencies —
no `go.sum` entries, no external binaries, no runtimes beyond the Go
toolchain itself — unless a future ADR explicitly revisits this
decision.**

This applies to:

- Go module dependencies (anything that would produce a `go.sum` entry);
- external binaries or runtimes invoked by any `make` target or `cmd/*`
  entrypoint (e.g. a headless browser, a Java runtime, a database
  server);
- build-time or render-time dependencies embedded into generated output
  (e.g. a JavaScript rendering library bundled into
  `generated/library/html/**`).

It does **not** apply to:

- developer tooling that is not a dependency of the module itself
  (`golangci-lint` is already required to run `make lint`, the same way a
  compiler is required to build; it does not appear in `go.sum` and does
  not ship inside any built artifact);
- plain-text, markup-embeddable formats with no rendering dependency of
  their own — see Consequences below for the Mermaid application of this
  principle.

## Consequences

- Every future proposal to add a Go module, external binary, or runtime
  dependency must be weighed against this ADR, not re-argued from first
  principles — and requires a new ADR to override, per this repo's own
  `M016`-equivalent guardrail-non-regression pattern (see
  `tests/architecture/architecture_test.go`'s own header comment and
  `docs/engineering/base-structure-review-workflow.md`).
  `docs-readme-c4-adr-scope-20260726` applied this reasoning directly:
  C4 architecture diagrams (`docs/architecture/**`) use Mermaid syntax
  embedded in plain Markdown — no PlantUML (requires a Java runtime or a
  rendering server) and no Go diagrams-as-code library (a new module
  dependency) — specifically because both alternatives would violate this
  ADR.
- This keeps the project's build and review surface small: `go build
  ./...` and `go test ./...` never need network access or a lockfile
  audit.
- The project accepts the corresponding cost: some features (automated
  PDF generation, cognitive-complexity linting, richer diagram styling)
  stay deferred or manual rather than automated, until a future ADR
  explicitly decides the trade-off is worth it for a specific case.

## Addendum 2026-07-31: Documentation PDF Generation Exception

The project accepts a narrow exception for automated PDF generation in the
documentation publication pipeline.

This exception exists because the landing page and Library now expose stable
PDF download links, and keeping those links fresh by manual browser export is
operationally fragile. A generated PDF is a distribution artifact for readers,
not a source of rules authority.

The exception allows:

- Python documentation-build dependencies required by a MkDocs PDF plugin;
- operating-system rendering libraries required by that plugin, when documented
  as setup prerequisites;
- Makefile targets dedicated to documentation PDF generation, such as
  `docs-pdf` or `pdf-build`;
- generated PDF files under the landing static download surface.

The exception does not allow:

- Go module dependencies or any `go.sum` entry;
- runtime dependencies used by game/system code;
- headless-browser automation unless a later ADR explicitly authorizes that
  class of dependency;
- treating generated PDFs as canonical documentation sources;
- silently making every build depend on PDF rendering before the dependency
  profile is documented and validated.

Implementation constraints:

- Keep the PDF dependency profile separate from the base documentation profile,
  for example `requirements-docs-pdf.txt` in addition to
  `requirements-docs.txt`.
- Prefer an explicit `docs-pdf` or `pdf-build` target during stabilization.
  `landing-build` may depend on it only after the PDF dependency profile is
  reproducible in the supported development environment.
- Generated PDFs must be traceable to canonical Markdown or semantic source
  units through the existing projection/manifest model.
- The landing page should link to stable `latest` aliases while preserving
  versioned PDF files for release history.
- Release manifest generation, checksum verification, and PDF editorial
  validation are project-owned tooling and should run through Go. Python remains
  allowed here only for the MkDocs renderer stack until a later ADR retires or
  replaces that publication path.
