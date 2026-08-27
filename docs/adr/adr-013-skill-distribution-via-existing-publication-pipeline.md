# ADR-013: Skill Distribution Via The Existing Publication Pipeline

## Status

Accepted.

## Context

A proposal to distribute the RGB Specialist skill as a downloadable `.zip`
(captured in an internal pending-analysis note) suggested a git-tag-triggered
GitHub Actions job publishing the archive as a GitHub Release asset, with a
dedicated install page and shell commands to unzip it into an
`~/.agents/skills`-style location.

This repository has no git-tag-triggered release job and no
`action-gh-release`-style step anywhere in `.github/workflows/ci.yml`.
Publication instead happens on every push to `main`: a `deploy` job bundles
generated documentation artifacts into a single static site
(`web/landing/dist`) and publishes it to GitHub Pages via
`actions/upload-pages-artifact` + `actions/deploy-pages`.

ADR-010 already established that topology for the two existing publication
outputs:

```text
docs/core/**
  |
  +-- cmd/rgb docs library --------------------> web/landing/public/library/**
  |
  +-- cmd/rgb docs pdf ------------------------> web/landing/public/downloads/**
                                                   |
                                                   v
                                            cmd/rgb release manifest/check
```

ADR-008 separately established that this project retires orphaned, parallel
publication pipelines rather than let them accumulate alongside the active
one. Adding a GitHub-Release-based flow for skill distribution would
reintroduce exactly that pattern: a second publication mechanism, with its
own trigger, permissions, and versioning scheme, running alongside the
Pages-based one that already serves the Library and PDFs.

Separately, `cmd/rgb` is this project's unified CLI (ADR-006), and it is
never invoked by end users directly — only by `make` targets and CI jobs.
Any skill distribution mechanism needs to preserve that: no new
user-facing installer CLI for this project.

`internal/components/tooling/release_artifacts.go`'s
`WriteReleaseArtifactManifest`/`CheckReleaseArtifacts` — the existing
manifest/checksum machinery for PDF downloads — hardcodes an expected file
list of exactly four PDF files (two locales, `latest` + versioned) and runs
PDF-specific editorial validation (`pdfinfo`/`pdftotext`/`pdftohtml`/
`pdftoppm`, page-count, TOC-link, and raster-luminance checks). It is not a
generic "checksum any artifact set" utility and cannot be pointed at a skill
`.zip` without change.

## Decision

Skill distribution joins the existing Go-owned publication topology
(ADR-010) as a third artifact type, published the same way the Library and
PDFs already are — not via GitHub Releases:

```text
skills/<skill>/**
  |
  +-- cmd/rgb docs skill -----------------------> web/landing/public/downloads/**
                                                   |
                                                   v
                                     a lightweight skill-artifact manifest/checksum step
                                     (separate from, not a reuse of, the PDF-specific
                                      release manifest/check machinery above)
```

- A skill package (e.g. `skills/specialist/`) is zipped by a new
  `cmd/rgb docs skill` subcommand, with the archive's top-level entry being
  the skill's own directory name — never a repo-root-relative or
  auto-generated wrapper folder.
- The `.zip` (`<name>-<version>.zip`, `<name>-latest.zip`) is written to
  `web/landing/public/downloads/`, the same directory the PDFs already use,
  and is published as a static file through the existing push-to-`main` →
  GitHub Pages `deploy` job. No new CI trigger, secret, or permission is
  introduced.
- Manifest/checksum generation for the skill `.zip` is a new, minimal step
  (SHA256 + byte count only), deliberately kept separate from
  `WriteReleaseArtifactManifest`/`CheckReleaseArtifacts` rather than
  generalizing those functions — the PDF checker's entire design center is
  PDF-specific editorial validation that does not apply to a skill archive,
  and coupling the two risks regressing PDF checks when the skill path
  changes.
- No new CLI is introduced for end users. `cmd/rgb docs skill` is invoked
  the same way `cmd/rgb docs library`/`cmd/rgb docs pdf` already are: through
  a Makefile target, from CI, never directly by a person installing the
  skill.

## Consequences

Positive:

- skill distribution reuses an already-reviewed publication boundary instead
  of introducing a second one;
- no new GitHub Actions triggers, permissions, or secrets;
- naming/versioning for the skill archive mirrors the existing PDF
  convention (`<name>-<version>`, `<name>-latest`), keeping the public
  `downloads/` directory self-consistent;
- the PDF-specific release-manifest/check machinery stays focused on PDF
  editorial validation and is not stretched to cover an unrelated artifact
  shape.

Accepted costs:

- a new, separate manifest/checksum code path is needed for the skill `.zip`
  rather than a single shared implementation across all downloadable
  artifacts;
- distribution here covers packaging/publication only — skill package
  *content* (what `skills/specialist/` actually contains, and whether it has
  a runnable contract) is defined by separate, independent work and is out
  of scope for this decision;
- multi-runtime skill packaging (for example a Codex- or Gemini-specific
  layout) is explicitly out of scope; this decision targets the Claude Code
  skill-loading convention already used by `skills/specialist/` and
  `skills/maker/` in this repository.

## Relationship To Prior ADRs

This ADR extends ADR-010's topology rather than amending it: ADR-010's
Library/PDF flow is unchanged, and this decision adds a third, independent
branch alongside it. It reaffirms ADR-008's rejection of parallel/orphaned
publication pipelines and ADR-006's unified-CLI boundary (`cmd/rgb` is never
an end-user-facing installer).
</content>
