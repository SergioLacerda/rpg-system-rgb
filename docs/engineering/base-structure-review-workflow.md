# Base Structure Review Workflow

## Purpose

A repeatable workflow to review the RGB System V2 base structure whenever the
scaffold changes: new components, new entrypoints, new tooling, or new
documentation surfaces. It combines automated gates with a short manual
checklist so the review is cheap enough to run on every structural change.

This workflow enforces locally what a CI pipeline would enforce remotely. When
a CI pipeline is added later, each gate below maps one-to-one to a CI job.

## Automated Gates

Run the fast gate from the repository root during local development:

```bash
make check-fast
```

Run the full development gate before opening or merging PRs:

```bash
make check
```

Run the release gate before publishing downloads:

```bash
make release-check
```

`check-fast` runs:

| Gate | Command | Verifies |
| --- | --- | --- |
| Format | `make fmt-check` | Go formatting is already normalized |
| Static analysis | `make vet` | suspicious constructs (`go vet`) |
| Lint baseline | `make lint` | golangci-lint rules in `.golangci.yml` (falls back to `go vet` if not installed) |
| Unit tests | `make test` | all package tests |
| Architecture guardrails | `make test-arch` | import boundaries in `tests/architecture/` |
| Semantic docs | `make validate` | semantic documentation contracts |
| Coverage floor | `make cover-check` | total statement coverage stays at or above `COVER_THRESHOLD` (see `Makefile`) |

`check` adds:

| Gate | Command | Verifies |
| --- | --- | --- |
| Landing lint | `make lint-web` | Astro/TypeScript checks |
| Landing tests | `make test-web` | landing unit tests and coverage |
| Landing build | `make landing-build` | MkDocs strict build and Astro static output |
| Workflow lint | `make lint-yaml` | GitHub Actions workflow syntax via actionlint |
| Shell lint | `make lint-shell` | CI shell scripts via shellcheck |
| Generated drift | `make check-generated-drift` | `make generate` leaves tracked generated artifacts clean |

`release-check` adds:

| Gate | Command | Verifies |
| --- | --- | --- |
| PDF build | `make docs-pdf` | bilingual latest and versioned PDF downloads |
| Artifact manifest | `make release-artifact-manifest` | Go-owned release manifest and `SHA256SUMS` for PDF artifacts |
| Editorial validation | `make pdf-editorial-check` | Go-owned PDF headers, metadata, TOC sanity, raster smoke, manifest contents, and checksums |

The coverage floor starts deliberately low relative to the current total
(`make cover` shows the current number) — it is a regression guard, not a
target. Raise `COVER_THRESHOLD` in the `Makefile` over time as coverage
improves; never lower it without recording why (mirrors the `gocyclo`
threshold ratchet in `.golangci.yml`).

Supporting targets:

```bash
make cover                  # test coverage summary
make fmt                    # normalize formatting
make go-file-size-report    # informational: non-test .go files over 200 lines
```

`make review-structure` remains as the legacy base-structure gate and is covered
by `make check-fast`.

## Required Gate By Context

| Context | Required gate |
| --- | --- |
| Local development loop | `make check-fast` before handoff when Go or semantic sources changed |
| Pull request | `make check`; CI may split the same work across jobs |
| `main` push | `make check` plus deployment smoke checks |
| Release candidate | `make release-check` before publishing or announcing downloads |

## Supply Chain Update Policy

External tool versions are intentionally pinned where repository rebuilds would
otherwise float:

- CI installs `golangci-lint` through `GOLANGCI_LINT_VERSION`.
- GitHub Actions Python jobs use `PYTHON_VERSION`.
- The actionlint Docker image uses an explicit version tag.
- Node uses an explicit major version and npm dependencies are governed by
  `web/landing/package-lock.json`.
- PDF editorial checks require Poppler command line tools (`pdfinfo`,
  `pdftotext`, `pdftohtml`, and `pdftoppm`). CI installs them through
  `poppler-utils`; manifest/checksum/raster validation is implemented in Go
  and no longer depends on Python or Pillow.

GitHub Marketplace actions currently use conventional major tags
(`actions/checkout@v4`, `actions/setup-go@v5`, and related first-party actions).
Digest or full-SHA pinning is deferred because Dependabot already tracks these
surfaces and keeps security updates visible in small PRs. Revisit this exception
if untrusted third-party actions are added.

Dependabot owns update PRs for GitHub Actions, Go modules, the landing npm
workspace, and Python requirements under `docs-build/` while MkDocs remains
the documented renderer exception.

The relevant gate must pass before structural changes are proposed for commit.

## PDF Editorial Gate

The release PDF renderer remains `mkdocs-to-pdf` for this first quality gate.
The PDF configs intentionally use `docs/styles/rgb-pdf.css` instead of the dark
web Library stylesheet, so printed pages default to white backgrounds, dark
text, visible tables, wrapping code blocks, and explicit chapter breaks.

`make pdf-editorial-check` rejects common publication defects:

- missing PDF headers or suspiciously small files;
- missing title, subject, or producer metadata where Poppler exposes it;
- table-of-contents entries that resolve to page `0`;
- missing extractable TOC links on critical pages;
- latest and versioned PDFs that differ for the same locale;
- missing release manifest, aggregate checksums, or per-version `.sha256` files;
- rasterized cover, TOC, and first content pages that are blank, dark themed, or
  have excessive marks at page edges.

Principal bookmark validation, full PDF/UA tagging, and durable visual
regression baselines remain deferred. Those checks require a renderer or parser
decision beyond the first editorial smoke gate.

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
