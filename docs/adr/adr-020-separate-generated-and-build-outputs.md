# ADR-020: Keep Tracked Output Under `generated/` And Disposable Output Under `build/`

## Status

Accepted.

Amends decisions 1 and 2 of
[ADR-018](adr-018-remove-in-repo-landing-and-relocate-publication-staging.md),
which put both the HTML Library and the release baseline under `build/`. Every
other decision of ADR-018 and of
[ADR-019](adr-019-consume-illusionista-surface-and-retire-landing-consumer.md)
is unchanged.

## Context

ADR-018 created a top-level `build/` that mixed two different kinds of output:

- `build/library/` is rebuilt on every run and ignored by git.
- `build/downloads/` is the tracked release baseline, consumed by
  `scripts/ci/compose-site.sh` and by the public-artifact gates.

`generated/` already holds the tracked, derived semantic projections
(`ai-context/`, `library/`, `pdf/`, `search/`), yet `.gitignore` listed it as a
whole while four of its files are tracked on purpose, so a new projection needed
`git add -f`.

## Decision

1. The two kinds of output get two roots, split by lifecycle:

   | Path | Content | Git |
   | --- | --- | --- |
   | `generated/{ai-context,library,pdf,search}/` | semantic projections (JSON) | tracked |
   | `generated/downloads/` | release baseline: PDFs, manifests, checksums, SBOM, provenance, skill packages | tracked |
   | `generated/client/` | written by the governance tooling | ignored |
   | `build/` | disposable output; the HTML Library is `build/site/library/` | ignored |

2. The HTML Library is written to `build/site/library/`. It must not be written
   to `generated/library/`: `BuildLibrary` removes its output directory before
   rendering, which would delete the tracked `core-v2-rules.json` projection.
3. The default paths of `cmd/rgb docs library|pdf|skill|check` and
   `cmd/rgb release ...` become `build/site/library` and `generated/downloads`.
   The `Makefile` variables are `BUILD_DIR` and `GENERATED_DIR`; CI artifact
   paths and `scripts/ci/compose-site.sh` follow.
4. `.gitignore` ignores `build/` and `generated/client/`, and no longer ignores
   `generated/` as a whole. Projections and the downloads baseline are tracked
   without `git add -f`.
5. `scripts/ci/check-generated-drift.sh` compares the projection outputs only:
   `generated/downloads` is excluded because it is rebuilt by the release gates,
   not by `make generate`.
6. The tracked files moved from `build/downloads/` to `generated/downloads/`
   without content changes. ADR-018 keeps its text; its `Status` points here.

## Consequences

- The rule is one sentence: `generated/` is versioned, `build/` is not.
- The published Library HTML is no longer under `generated/`, so
  `ValidateProjectPaths` does not scan it for leaked internal path markers. It is
  rendered from `docs/`, which that validator already scans.
- External callers relying on the `build/` defaults for downloads, or on
  `generated/` for the HTML, change paths. This is breaking and recorded in the
  CHANGELOG; ADR-018 is unreleased, so no tagged version shipped its defaults.
- Binary release files still live in git. Moving the baseline to GitHub Releases
  remains a separate decision.

## Rejected Alternatives

- **One root, `generated/site/` for the HTML.** It keeps ignored output inside a
  directory whose purpose is tracked artifacts, and needs a drift-script
  exclusion for it.
- **Keep `build/downloads/`.** A tracked baseline under an ignored, disposable
  root misleads about what may be deleted.
- **Write the HTML Library to `generated/library/`.** It collides with the
  projection of the same name and is erased on every build.
