# MkDocs Library Integration Implementation Plan

**Goal:** Publish the canonical `docs/` tree as the landing page Library through MkDocs.

**Architecture:** Keep `docs/` as the source of truth, generate a static MkDocs site into
`web/landing/public/library/`, and let Astro publish it as a static asset under the
configured base path. Existing localized Astro Library routes become bridge pages so
older links do not 404.

**Tech Stack:** MkDocs, Make, Astro, Node/npm.

---

## Tasks

1. Add `requirements-docs.txt` with the MkDocs dependency.
2. Add root `mkdocs.yml` with `docs_dir: docs` and
   `site_dir: web/landing/public/library`.
3. Add `docs/index.md` as the MkDocs home page.
4. Add `docs-install`, `docs-build`, and `docs-preview` to the root Makefile.
5. Make `landing-build` depend on `docs-build`.
6. Point the landing header Library link to `/library/` through `landingPath`.
7. Replace localized Astro Library pages with bridge pages to `/library/`.
8. Document the Library install/build/preview flow in `web/landing/README.md`.
9. Verify `make docs-build` and `make landing-build`.
