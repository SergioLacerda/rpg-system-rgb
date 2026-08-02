# Contributing

Thanks for taking the time to improve RGB System.

## Project Status

Check the status table in [README.md](README.md) before proposing work. RGB Core
V2, semantic tooling, Library HTML, PDF downloads, and the landing page are
implemented. Maker and Specialist are planned surfaces and should not be
documented as available products.

## Useful Contributions

- Rule clarifications and examples for the canonical English docs.
- Portuguese localization updates that stay aligned with English source docs.
- Tests for implemented mechanics or semantic documentation contracts.
- Build, release, and documentation quality improvements.

## Workflow

1. Open an issue for material behavior, architecture, or public-doc changes.
2. Keep changes scoped to one concern.
3. Run the relevant gate before handoff:
   - `make check-fast` for Go or semantic-doc changes.
   - `make check` for public docs, landing, generated artifacts, or CI changes.
   - `make release-check` for PDF or release-asset changes.
4. Do not edit generated files by hand; run `make generate` or the documented
   publication target instead.

## Documentation Rules

- English technical docs are the canonical local documentation surface.
- Use the status taxonomy: Implemented, Experimental, Planned, Deferred.
- Avoid calling planned skill surfaces implemented.
- Avoid calling implemented publication surfaces placeholders.

## Conduct

Participation is governed by [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
