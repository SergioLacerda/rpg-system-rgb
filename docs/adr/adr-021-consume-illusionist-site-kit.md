# ADR-021: Consume The Illusionist Site Kit

## Status

Accepted. Amends ADR-019's ownership and composition decisions for the site
and Library, and ADR-010's Go-owned HTML publication decision.

## Decision

RGB supplies its `docs/core` tree unchanged at the path required by the pinned
Illusionist RGB site kit. The kit owns the landing, Docs tab, navigation,
rendering and search. RGB verifies the release checksum and a separately pinned
SHA-256 before executing the kit. CI adds the RGB-owned PDF and Specialist
skill downloads to the generated site. PDF authoring and publication remain in
RGB; the Illusionist PDF output is not implemented by the provider.

The private release is downloaded by GitHub CLI using `GH_TOKEN`, supplied in
CI from the `ILLUSIONIST_RELEASE_TOKEN` secret. Local builds may use an
authenticated `GH_TOKEN` or an already downloaded kit in `ILLUSIONIST_KIT_DIR`;
both paths verify the same pinned digest. The kit contains the compiler built
for that site release, so its pinned digest covers the compiler as well.

The local HTML renderer and `rgb docs library` are retired. `make docs-build`
builds the entire site in an empty output directory, while `make docs-check`
builds in a temporary directory and checks the site and RGB PDF artifacts.
The default Portuguese Library is at `library/`, and English is at
`en/library/`. URLs from the former renderer are not preserved by the kit;
external links to those paths need updating. Source documents and their
content are RGB's responsibility.

The semantic Library consumer contract remains reserved because existing
semantic projections name it as owner. It is not the input contract for the
Illusionist Docs tab and does not describe the kit's search implementation.
