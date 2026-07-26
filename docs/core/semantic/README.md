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
go run scripts/validate_semantic_docs.go
```

Validate the documentation L10n manifest with:

```bash
go run scripts/validate_docs_l10n_manifest.go docs/core/semantic/l10n-manifest.v0.1.json
```

The explicit index validation command is:

```bash
go run scripts/validate_semantic_index.go docs/core/semantic/core-v2.index.json
```

Validate accepted semantic source slices with:

```bash
go run scripts/validate_semantic_source.go docs/core/semantic/source/core-v2-rules.v0.1.json docs/core/semantic/core-v2.index.json
```

The schema contract is documented in [schema-v0.1.md](schema-v0.1.md).

Consumer contracts are defined in
[consumer-contracts.v0.1.json](consumer-contracts.v0.1.json) and summarized in
the accepted documentation authority ADR
[ADR-001](../../adr/adr-001-ai-first-documentation-authority.md).

Validate the consumer contracts with:

```bash
go run scripts/validate_semantic_contracts.go docs/core/semantic/consumer-contracts.v0.1.json docs/core/semantic/core-v2.index.json
```

Derived projection manifests are defined in
[projection-manifest.v0.1.json](projection-manifest.v0.1.json) and summarized in
the accepted documentation authority ADR
[ADR-001](../../adr/adr-001-ai-first-documentation-authority.md).

Validate the projection manifest with:

```bash
go run scripts/validate_semantic_projections.go docs/core/semantic/projection-manifest.v0.1.json docs/core/semantic/core-v2.index.json docs/core/semantic/consumer-contracts.v0.1.json
```

Generate derived projection artifacts with:

```bash
go run scripts/generate_semantic_projections.go docs/core/semantic/projection-manifest.v0.1.json docs/core/semantic/core-v2.index.json docs/core/semantic/source/core-v2-rules.v0.1.json
```

Validate generated projection artifacts with:

```bash
go run scripts/validate_generated_projections.go docs/core/semantic/projection-manifest.v0.1.json
```
