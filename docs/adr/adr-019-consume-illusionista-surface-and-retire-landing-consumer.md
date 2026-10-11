# ADR-019: Consume The RGB Landing As A Pinned ILUSIONISTA Surface And Retire The `landing` Consumer

## Status

Accepted.

Amends [ADR-001](adr-001-ai-first-documentation-authority.md) and
[ADR-002](adr-002-rgb-core-v2-design-control.md) only where they name `landing`
as a consumer or derived projection of Core V2.

Resolves the point [ADR-018](adr-018-remove-in-repo-landing-and-relocate-publication-staging.md)
left for a separate decision (the `landing` entries of the semantic vocabulary)
and records the integration contract that ADR-018 only referenced.

## Context

[ADR-018](adr-018-remove-in-repo-landing-and-relocate-publication-staging.md)
removed the in-repository landing page. Two things stayed unresolved:

1. No ADR records how this repository consumes the landing that replaced it.
   The contract existed only in `.github/workflows/ci.yml` and
   `scripts/ci/compose-site.sh`.
2. The semantic layer still described `landing` as a component that consumes
   Core V2 projections (`consumer.landing.v0_1`, the `landing` projection
   surface, "landing summaries"). That described the old in-repository Astro
   application, not the current arrangement.

The current arrangement runs in the opposite direction. ILUSIONISTA
(`SergioLacerda/illusionist-web`) builds, tests and releases the RGB landing as
an immutable static archive (`illusionist-rgb-surface-vX.Y.Z.tar.gz`, tag
`rgb-vX.Y.Z`). RGB System is a consumer of that archive. The archive holds
presentation only: it carries no Library, no downloads, and nothing generated
from the Core V2 semantic index.

## Decision

1. **Ownership.**

   | Part | Owner |
   | --- | --- |
   | Landing presentation surface (`illusionist-rgb-surface-*`) | ILUSIONISTA |
   | Library generation (`cmd/rgb docs library`), PDFs, skill packages, manifests, SBOM, provenance | RGB System |
   | Composition of the final site tree and GitHub Pages deployment | RGB System |
   | Core V2 rules, semantic index, consumer contracts | RGB System |

2. **Dependency direction.** RGB System depends on ILUSIONISTA for the
   presentation surface. ILUSIONISTA does not read RGB projections, and this
   repository does not modify the archive it consumes.

3. **Consumption contract.** The deploy job composes the site from the surface
   under these rules, implemented by `scripts/ci/compose-site.sh` and the
   `ILLUSIONIST_*` values in `.github/workflows/ci.yml`:
   - the version is an exact release tag; `latest`, `main`, a commit and a
     workflow run are rejected;
   - integrity is pinned twice: the release `.sha256` asset must verify, and
     the archive digest must equal `ILLUSIONIST_RGB_SHA256` committed in
     `ci.yml`;
   - the archive must provide `index.html`, `pt-br/index.html`,
     `en/index.html` and `_astro/`;
   - collision guard: the archive must not carry non-empty `library/` or
     `downloads/`; the deploy fails instead of overwriting;
   - RGB mounts its generated Library at `library/` and its release downloads
     at `downloads/`, and every tracked file under `generated/downloads/` must
     still be published (baseline parity);
   - the surface is compiled for the base `/rpg-system-rgb` on
     `https://sergiolacerda.github.io`; serving it under another base needs a
     new surface release.

4. **Semantic vocabulary.** `landing` is retired as a component consumer and as
   a projection surface of Core V2:
   - `consumer.landing.v0_1` is removed from
     `docs/core/semantic/consumer-contracts.v0.1.json`;
   - `landing` is removed from `allowed_projection_surfaces`,
     `projection_surfaces` and `component_consumers`.

   The presentation surface remains non-canonical, as ADR-001 requires. It has
   no rules authority, and no RGB projection feeds it. If ILUSIONISTA ever needs
   RGB content at build time, that requires a written consumer contract first
   ([ADR-007](adr-007-bundle-consumer-contract-policy.md)) and a new ADR.

5. **Upgrading the surface.** A new surface version is one change that updates
   the `ILLUSIONIST_RGB_*` values together: version, artifact, checksum asset
   and `ILLUSIONIST_RGB_SHA256`. Rolling back is the same change in reverse;
   archived releases are immutable.

6. **Historical ADRs.** Their text is not rewritten. Their `Status` sections
   point to [ADR-018](adr-018-remove-in-repo-landing-and-relocate-publication-staging.md)
   and to this ADR where a decision was superseded or amended.

## Consequences

- One ADR now answers "who produces the landing, and how do we consume it".
- The semantic layer no longer promises a consumer that does not exist, so
  `make validate` and the generated indexes cannot be read as endorsing
  "landing summaries".
- Changing the surface version is a reviewed, single-purpose change with a
  pinned digest; a compromised or moved release fails closed at deploy.
- The deploy depends on a download from another repository's GitHub Release.
  A removed or renamed release breaks the deploy until the pin is updated.
- `docs/README-details.md` already lists `Landing` as implemented and hosted
  outside this repository; it now links here.

## Rejected Alternatives

- **Keep `consumer.landing.v0_1` as an external consumer.** The surface reads no
  RGB semantic unit at build time, so the contract would describe outputs that
  do not exist.
- **Vendor the surface source back into this repository.** It reintroduces the
  Node toolchain and the duplicated ownership that ADR-018 removed.
- **Track `latest` or `main` of the surface.** A prebuilt tree cannot be
  reviewed or rolled back reliably, and a silent surface change would ship
  straight to Pages.
