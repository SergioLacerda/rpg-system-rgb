# ADR-005: Zero External Dependency Posture

## Status

Accepted.

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
