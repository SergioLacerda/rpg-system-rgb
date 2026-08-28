# RGB PDF Authoring Toolchain

This directory contains the manual, external PDF authoring surface allowed by
ADR-016. It is not part of the Go module graph, the default test gate, or CI.

The canonical source remains `docs/core/**`. This toolchain renders a reviewed
PDF candidate with the editorial cover and R/G/B callout styling, then the
existing Go-owned publication command copies and validates that reviewed asset:

```bash
tools/pdfbuild/build-pdf.sh en v2.0
tools/pdfbuild/build-pdf.sh pt-br v2.0

make docs-pdf \
  PDF_SRC_EN=tools/pdfbuild/out/rgb-system-core-v2.0-en.pdf \
  PDF_SRC_PT_BR=tools/pdfbuild/out/rgb-system-core-v2.0-pt-br.pdf
```

The script expects Python, MkDocs, MkDocs Material, WeasyPrint, and Poppler.
Dependencies are pinned in `requirements.txt`, but installation is explicit and
local to `tools/pdfbuild/.venv`.

Fonts may be added under `fonts/` later. The stylesheet currently uses bundled
system fallbacks so the authoring step can be introduced without committing font
binaries.
