# ADR-010: Go-Owned Publication Without Python

## Status

Accepted.

## Context

ADR-005 established the project's zero external dependency posture and later
allowed a narrow documentation publication exception for the former renderer
stack. ADR-008 removed the orphaned Go HTML/print/PDF compiler path. ADR-009
then kept the former exception only until a replacement topology proved
parity.

The remaining interpreter-backed surface was publication:

- documentation HTML rendering;
- PDF download publication;
- CI setup for the renderer;
- dependency update tracking for that renderer.

Release metadata and editorial validation were already Go-owned through
`cmd/rgb release manifest`, `cmd/rgb release check`, and the tooling package.

## Decision

Retire the interpreter-backed publication exception and make publication
Go-owned.

The active topology is now:

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

Astro remains the static landing package and hosting surface. It packages the
generated public Library and PDF downloads from `web/landing/public/**`.

This is a new publication boundary, not a revival of the ADR-008 deleted
compiler/library packages or the removed `rgb compile` command.

## Consequences

Positive:

- project publication no longer installs or runs the retired interpreter stack;
- CI and local Make targets use Go plus the existing Astro/Node landing
  workspace;
- dependency update tracking no longer needs the retired renderer ecosystem;
- the public Library is generated through the unified `cmd/rgb` CLI topology;
- release manifest and checksum behavior remains Go-owned.

Accepted costs:

- the Go HTML renderer intentionally supports a constrained Markdown subset
  based on the current public `docs/core/**` files;
- PDF publication is a Go-owned publishing step for reviewed PDF assets, not a
  full general-purpose PDF layout engine;
- Poppler remains an external command-line validation dependency for PDF
  editorial smoke checks. It is not part of publication rendering and does not
  reintroduce the retired interpreter runtime.

## Replacement Parity

The replacement preserves:

- HTML Library generation from canonical `docs/core/**`;
- Astro packaging of the Library;
- English and PT-br public documentation surfaces;
- stable `latest` PDF aliases;
- versioned PDF files;
- release manifest generation;
- aggregate `SHA256SUMS` and per-version `.sha256` files;
- PDF editorial smoke checks.

## Supersession

This ADR supersedes ADR-009 for active publication behavior. ADR-009 remains
as the historical record of why the exception was retained until this
replacement existed.
