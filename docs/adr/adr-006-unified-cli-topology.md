# ADR-006: `cmd/rgb` As The Unified CLI

## Status

Accepted. Amended 2026-08-02 — see
[Amendment: publication pipeline retirement](#amendment-publication-pipeline-retirement-2026-08-02)
and
[Amendment: deprecated caller migration](#amendment-deprecated-caller-migration-2026-08-02)
below.

## Context

A base project architecture review found the CLI surface ambiguous:
`cmd/rgb` advertised scaffold readiness only (`app.Hello()`), while the
functional commands were split across `cmd/rgb-tooling` (`validate`,
`generate`, `bundle`) and `cmd/rgb-compiler` (`all`, `no-html`). Nothing
distinguished which binary a user or CI job was supposed to run, and the
review left the choice open: either promote `cmd/rgb` to a unified CLI, or
keep it as a status command and document the other two as the real CLIs.

## Decision

`cmd/rgb` is the unified CLI. It now dispatches to the same
`internal/app` use cases the other two binaries already called:

```
rgb                        -> scaffold/component status (unchanged)
rgb validate [repo-root]   -> internal/app.ValidateDocs
rgb generate [repo-root]   -> internal/app.GenerateProjections
rgb bundle   [repo-root]   -> internal/app.BuildBundle
rgb compile  all|no-html [repo-root] -> internal/app.Compile
```

`Makefile` targets (`validate`, `generate`, `bundle`, `compile`,
`compile-full`) now invoke `cmd/rgb` instead of the two split binaries.

`cmd/rgb-tooling` and `cmd/rgb-compiler` are **deprecated, not removed**.
They still build and behave exactly as before because real callers depended
on them at the time of this decision: `tests/semantic_docs/semantic_docs_test.go`
shelled out to `go run ./cmd/rgb-tooling validate`, and
`docs/core/semantic/README.md` documented `rgb-tooling` directly. Removing a
working entrypoint those depended on was a breaking change under M006 (RFC
Process for Breaking Changes) and was out of scope for the original decision.

## Consequences

Positive:

- one documented entrypoint for the product-facing CLI surface, matching
  the "one coherent CLI" gap called out by the architecture critique;
- `internal/app` stays the single seam between entrypoints and components,
  so `cmd/rgb`'s dispatch logic is a thin, independently testable layer
  (`cmd/rgb/main_test.go`) with no new component coupling.

Negative / accepted costs:

- two deprecated binaries remain in the tree until their callers migrate;
- `cmd/rgb`'s subcommand set must be kept in sync with `internal/app` by
  hand until there's enough surface to justify a flag-parsing library
  (see ADR-005 on the zero external dependency posture for why that's not
  a default reach).

## Follow-up

- Completed: `tests/semantic_docs/semantic_docs_test.go` and
  `docs/core/semantic/README.md` now call `cmd/rgb`.
- Remaining: propose `cmd/rgb-tooling` removal through the M006 RFC process if
  and when compatibility support is no longer needed. (`cmd/rgb-compiler` no
  longer applies here — see the publication-pipeline amendment.)

## Amendment: Publication Pipeline Retirement (2026-08-02)

A base project Library review found that,
unlike `cmd/rgb-tooling`, `cmd/rgb-compiler` and the `compile` subcommand
had **no real caller** — no test, CI workflow, or required Makefile target
used them; `internal/components/compiler` and `internal/components/library`
underneath them were equally uncalled outside that path. M006's
breaking-change concern (real callers depend on this) did not apply, so
[ADR-008](adr-008-retire-go-publication-pipeline.md) removed `compile`,
`cmd/rgb-compiler`, `internal/components/compiler`, and
`internal/components/library` outright instead of deprecating them. Every
mention of `compile`/`cmd/rgb-compiler` above describes this ADR's original
decision at the time it was written, not the current subcommand set — see
ADR-008 for what replaced it. The `cmd/rgb-tooling` decision above is
unaffected.

## Amendment: Deprecated Caller Migration (2026-08-02)

The repository callers named in the original decision have migrated to
`cmd/rgb`: `tests/semantic_docs/semantic_docs_test.go` runs
`go run ./cmd/rgb validate`, and `docs/core/semantic/README.md` documents
`cmd/rgb` as the equivalent command for `make validate` and `make generate`.

`cmd/rgb-tooling` remains present only as a compatibility path. Removing it is
still a breaking-change decision and requires the M006 RFC process.
