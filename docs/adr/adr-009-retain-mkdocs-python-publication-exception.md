# ADR-009: Retain The MkDocs/Python Publication Exception Until Replacement Parity

## Status

Superseded by [ADR-010](adr-010-go-owned-publication-without-python.md).

Historical status before supersession: Accepted.

## Context

ADR-005 established the project's zero external dependency posture and then
added a narrow documentation-PDF exception on 2026-07-31. That exception
allows Python documentation-build dependencies required by the MkDocs PDF
plugin, operating-system rendering libraries documented for that plugin, and
generated PDFs under the landing download surface.

The project-owned release validation path has since been moved to Go:
release manifest generation, checksum verification, PDF metadata checks, TOC
checks, and raster smoke checks no longer depend on Python or Pillow. That
leaves the renderer/publication stack itself as the remaining Python surface:

- `make docs-build` uses `python3 -m mkdocs` for the HTML Library;
- `make docs-pdf` uses MkDocs plus `mkdocs-to-pdf` for bilingual PDF
  downloads;
- CI uses `actions/setup-python` for documentation and deployment jobs;
- `docs-build/requirements-docs.txt` and
  `docs-build/requirements-docs-pdf.txt` define the Python dependency
  profiles;
- Dependabot tracks the `/docs-build` `pip` ecosystem while the exception
  remains active.

ADR-008 is also binding context. It retired the previous Go HTML/print/PDF
pipeline because it had no active caller and would have preserved a parallel,
orphaned publication narrative. A future non-Python renderer must therefore
be a new justified architecture decision with real consumers and parity
evidence, not a silent resurrection of `internal/components/compiler`,
`internal/components/library`, `cmd/rgb-compiler`, or `rgb compile`.

## Decision

Keep the MkDocs/Python publication exception for now.

The current publication topology remains:

```text
docs/**
  |
  +-- MkDocs strict build ----------------> web/landing/public/library/**
  |
  +-- MkDocs + mkdocs-to-pdf + Poppler ---> web/landing/public/downloads/*.pdf
                                             |
                                             v
                                      Go-owned release artifact gates
```

Retiring Python from documentation publication requires a later ADR or an
explicit amendment to this ADR. That later decision must identify the
replacement topology and provide parity evidence before deleting the current
Python surfaces.

The replacement must preserve or deliberately supersede all current public
release behavior:

- HTML Library generation;
- Astro landing packaging of the Library;
- English and PT-br PDF generation;
- stable `latest` PDF aliases;
- versioned PDF files;
- release manifest generation;
- aggregate `SHA256SUMS` and per-version `.sha256` files;
- PDF editorial smoke checks.

Until that parity exists, do not remove or weaken:

- `PYTHON`, `MKDOCS`, `docs-install`, `docs-pdf-install`, `docs-build`, or
  `docs-pdf` in the Makefile;
- `docs-build/requirements-docs.txt`;
- `docs-build/requirements-docs-pdf.txt`;
- Python setup and pip caching in CI documentation/deployment jobs;
- the Dependabot `pip` entry for `/docs-build`.

## Consequences

Positive:

- keeps the current public documentation and PDF release path stable;
- avoids deleting a working renderer without a proven replacement;
- keeps ADR-005's exception narrow and explicit;
- prevents accidental reversal of ADR-008 by forbidding a silent revival of
  the deleted Go compiler/library pipeline;
- leaves release artifact validation Go-owned and independent from the
  Python renderer exception.

Negative / accepted costs:

- the repository still needs Python for documentation publication;
- `docs-build/` remains a dependency surface that CI and Dependabot must
  track;
- fully zero-Python publication remains deferred until a later ADR selects
  and proves a replacement topology.

## Replacement Requirements

Any future proposal to retire MkDocs/Python must include at least:

1. A renderer topology that does not reintroduce the ADR-008 deleted pipeline
   by implication.
2. Local validation equivalent to or stronger than:
   - `make docs-build`;
   - `make landing-build`;
   - `make docs-pdf`;
   - `make release-artifact-manifest`;
   - `make release-artifact-check`;
   - `make pdf-editorial-check`.
3. CI changes proving the replacement runs without `actions/setup-python`,
   pip cache, or `docs-build/requirements-*.txt`.
4. A deletion plan for obsolete Python publication files only after parity
   passes.
5. Updated public documentation describing the new publication path.

Full PDF/UA tagging, durable visual-regression baselines, and deeper
accessibility checks remain future quality work. They are desirable, but they
are not prerequisites for keeping the current MkDocs/Python exception active.
