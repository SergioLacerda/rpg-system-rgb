# Base Structure Review Workflow

## Purpose

A repeatable workflow to review the RGB System V2 base structure whenever the
scaffold changes: new components, new entrypoints, new tooling, or new
documentation surfaces. It combines automated gates with a short manual
checklist so the review is cheap enough to run on every structural change.

This workflow enforces locally what a CI pipeline would enforce remotely. When
a CI pipeline is added later, each gate below maps one-to-one to a CI job.

## Automated Gates

Run the full gate from the repository root:

```bash
make review-structure
```

This runs, in order:

| Gate | Command | Verifies |
| --- | --- | --- |
| Static analysis | `make vet` | suspicious constructs (`go vet`) |
| Lint baseline | `make lint` | golangci-lint rules in `.golangci.yml` (falls back to `go vet` if not installed) |
| Unit tests | `make test` | all package tests |
| Architecture guardrails | `make test-arch` | import boundaries in `tests/architecture/` |
| Semantic docs | `make validate` | semantic documentation contracts |

Supporting targets:

```bash
make cover      # test coverage summary
make fmt        # normalize formatting
```

The gate must pass before structural changes are proposed for commit.

## Architecture Rules Under Test

`tests/architecture/architecture_test.go` enforces the dependency direction of
the scaffold. The rules mirror the depguard rules in `.golangci.yml`.

```text
cmd/*  ->  internal/app  ->  internal/components/<component>  ->  internal/components
```

1. Components must not import `internal/app` or `cmd` (dependencies point
   inward only).
2. Components must not import sibling components; shared contracts live in
   `internal/components`.
3. The shared contract package `internal/components` must depend only on the
   standard library.
4. Entrypoints in `cmd` reach project code only through `internal/app`.

Relaxing or changing any rule requires an ADR (see `docs/adr/`), per the
guardrail non-regression mandate (M016).

## Manual Review Checklist

Walk this list after the automated gates pass:

### Layout

- [ ] Every new top-level directory has a clear owner concept
      (`cmd`, `internal`, `docs`, `skills`, `web`, `generated`, `scripts`, `tests`).
- [ ] Generated artifacts live only under `generated/` and are reproducible
      from canonical sources.
- [ ] No canonical content was added under `generated/`.

### Components

- [ ] A new component boundary has a package under
      `internal/components/<name>` with a doc comment and a `Descriptor()`.
- [ ] The component is registered in `internal/app`.
- [ ] The component does not leak types from other components.

### Documentation

- [ ] English docs under `docs/core/en/**` remain the canonical source;
      PT-br files are localized projections (ADR-003).
- [ ] Structural decisions with lasting impact have an ADR in `docs/adr/`.
- [ ] Semantic manifests under `docs/core/semantic/` were regenerated if
      source documents changed (`make validate` confirms).

### Tests

- [ ] New behavior in `internal/` ships with tests in the same package
      (`*_test.go` next to the code).
- [ ] Cross-cutting or contract-level checks live under `tests/`
      (`tests/architecture`, `tests/semantic_docs`).
- [ ] Coverage did not silently drop (`make cover`).

## Test Layout Convention

```text
internal/<pkg>/x.go          unit tests: internal/<pkg>/x_test.go
tests/architecture/          import-boundary guardrails (whole repo)
tests/semantic_docs/         documentation contract guardrails
```

Unit tests stay next to the code they test. The `tests/` tree is reserved for
repo-wide guardrails that do not belong to a single package.

## Outcome

Record review findings that require follow-up in the analysis intake surface
(PT-br allowed) or as ADRs under `docs/adr/` (decision surface, English).
