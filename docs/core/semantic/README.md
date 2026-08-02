# Semantic Documentation Index

This folder contains the AI-first sidecar indexes for RGB System documentation.

The current pilot is:

- `core-v2.index.json`
- `consumer-contracts.v0.1.json`
- `projection-manifest.v0.1.json`
- `l10n-manifest.v0.1.json`
- `schema-v0.1.md`
- `source/core-v2-rules.v0.1.json`

The sidecar index is a bridge over current Markdown. It provides stable IDs,
authority classification, source status, projection paths, relationships,
retrieval summaries, and provenance for AI and future component consumers.

During the pilot:

- `docs/core/en/**` remains the Markdown review projection and is still the canonical
  bridge for units not promoted into `docs/core/semantic/source/**`;
- `docs/core/PT-br/**` remains a localized projection;
- the semantic index maps identity, source authority, and retrieval metadata;
- accepted source files under `docs/core/semantic/source/**` are canonical only for
  their promoted source IDs;
- generated surfaces such as Library, PDF, landing, bundles, search data, and
  AI context packs remain derived artifacts.

Validate the pilot with:

```bash
make validate
# equivalent: go run ./cmd/rgb validate
```

This runs, in order: project paths, semantic index
([schema-v0.1.md](schema-v0.1.md)), documentation L10n manifest, consumer
contracts (summarized in the accepted documentation authority ADR
[ADR-001](../../adr/adr-001-ai-first-documentation-authority.md)), accepted
semantic source slices, the projection manifest, and generated projection
artifacts.

Regenerate derived projection artifacts with:

```bash
make generate
# equivalent: go run ./cmd/rgb generate
```

`cmd/rgb` is the unified CLI for semantic validation and projection generation.
The deprecated `cmd/rgb-tooling` binary remains only as a compatibility path
while older callers migrate.
