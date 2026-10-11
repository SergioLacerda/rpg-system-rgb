# Changelog

All notable project-level changes should be recorded here.

## Unreleased

- Removed the in-repository Astro landing page (`web/`) after its migration to an
  external surface; publication staging now lives under `build/`
  (`build/library`, `build/downloads`). **Breaking:** default output paths of
  `cmd/rgb docs library|pdf|skill|check` and `cmd/rgb release` changed, and the
  `landing-*`, `lint-web`, `test-web`, `audit-web`, `preview` and `web.*` make
  targets were removed. See ADR-018.
- Recorded ADR-019: the landing is consumed as a pinned external ILUSIONISTA
  surface; the `landing` consumer contract and projection surface were retired
  from the semantic vocabulary.
- Added public status taxonomy for implemented, experimental, planned, and
  deferred surfaces.
- Added minimal OSS governance files and templates.
- Added release and PDF quality gates for public downloads.
- Clarified that Maker and Specialist are planned surfaces, not available
  products.

## 2026-08-01

- Documented RGB Core V2 publication, landing, PDF, and release-gate handoffs.
