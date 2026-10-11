# ADR-018: Remove The In-Repository Landing Page And Relocate Publication Staging To `build/`

## Status

Accepted.

Supersedes the location-specific parts of
[ADR-004](adr-004-modular-monolith-monorepo.md),
[ADR-009](adr-009-retain-mkdocs-python-publication-exception.md),
[ADR-010](adr-010-go-owned-publication-without-python.md),
[ADR-013](adr-013-skill-distribution-via-existing-publication-pipeline.md) and
[ADR-016](adr-016-pdf-external-authoring-exception.md) that name `web/` or
`web/landing/public/**`. Their other decisions are unchanged.

Amended by [ADR-019](adr-019-consume-illusionista-surface-and-retire-landing-consumer.md), which resolves the `landing` consumer left open under
Consequences and records the contract for consuming the external surface.

## Context

The landing page is now published from an immutable external ILLUSIONISTA
surface that `scripts/ci/compose-site.sh` downloads and verifies at deploy
time. The Astro application under `web/landing` no longer fed the deployed
site, yet `make check`, two CI jobs, Dependabot, and the Node toolchain still
built and tested it.

`web/landing/public/{library,downloads}` was never landing-only. It was the
staging directory of the Go publication pipeline and the tracked release
baseline: `pdf-regression-check` uses `PDF_PUBLIC_DIR` as its baseline
directory, and `compose-site.sh` copies it into the site and enforces baseline
parity. Deleting `web/` therefore required relocating that role first.

## Decision

1. The publication staging root is `build/`, keeping the existing layout:
   - `build/library/` is generated HTML Library output and is ignored by git.
   - `build/downloads/` is the tracked release baseline (PDFs, manifests,
     checksums, SBOM, provenance, and skill packages). It was moved without
     content changes.
2. The default output paths of `cmd/rgb docs library|pdf|skill|check` and
   `cmd/rgb release ...` change from `web/landing/public/...` to `build/...`.
   Every flag is kept; the Makefile and CI always pass explicit paths.
3. All landing-page behavior is removed from this repository: `web/`, the
   `landing-*`, `preview`, `lint-web*`, `test-web`, `audit-web` and `web.*`
   make targets, the `lint-web` and `test-web` CI jobs, the Dependabot npm
   ecosystem, and the `projection.landing.core-summary.v0_1` projection
   (`web/landing/generated/core-v2-summary.json`), which had no consumer
   outside `web/`.
4. `make install` now installs the pinned local validation tools
   (`actionlint`, `shellcheck`); there are no Node dependencies.
5. The CI artifacts that carried release files are renamed `release-pdfs` and
   `release-skills`.

## Consequences

- `make check` and CI no longer require Node or npm.
- External callers that relied on the old `cmd/rgb` defaults write to `build/`
  instead; this is a breaking change recorded in the CHANGELOG.
- The semantic vocabulary still lists `landing` as a consumer
  (`docs/core/semantic/consumer-contracts.v0.1.json`). It is canonical
  documentation authority (ADR-001); ADR-019 retires it.
- Required status checks named `lint-web` and `test-web` must be removed from
  branch protection before merging the CI change.
- Historical ADRs and plans keep their original `web/landing` paths as records.
