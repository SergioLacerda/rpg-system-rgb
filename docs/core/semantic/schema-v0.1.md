# RGB Semantic Index Schema 0.1

## Purpose

This schema describes the first AI-first sidecar index for RGB Core V2. It is a
bridge contract, not the final semantic source-of-truth.

The schema is intentionally small. It gives each documentation unit a stable
identity, source authority, projection paths, relationships, retrieval metadata,
and provenance while current Markdown remains the human review surface.

## Top-Level Fields

| Field | Type | Required | Purpose |
| --- | --- | --- | --- |
| `schema` | string | yes | Schema identifier. Current value: `rgb-docs-semantic-index/0.1`. |
| `source_locale` | string | yes | Canonical bridge locale. Current value: `en`. |
| `default_localized_locale` | string | yes | Default localized projection locale. Current value: `PT-br`. |
| `authority_types` | string list | yes | Allowed authority types for units. |
| `source_statuses` | string list | yes | Allowed lifecycle and review statuses. |
| `kinds` | string list | yes | Allowed semantic unit kinds. |
| `projection_surfaces` | string list | yes | Allowed projection targets. |
| `component_consumers` | string list | yes | Allowed consumer tracks. |
| `units` | object list | yes | Semantic units in this index. |

## Unit Fields

| Field | Type | Required | Purpose |
| --- | --- | --- | --- |
| `id` | string | yes | Stable semantic identifier. |
| `kind` | enum | yes | Unit kind: concept, rule, procedure, example, term, translation, state, resource, ability, design_note, or generated_artifact. |
| `locale` | string | yes | Unit locale. Canonical bridge units start as `en`. |
| `authority_type` | enum | yes | Source authority class. |
| `source_status` | enum | yes | Lifecycle or review status. |
| `title` | string | yes | Human-readable title. |
| `source_path` | path | yes | Repository-relative source or projection path. |
| `projection_paths` | map | yes | Known projection paths keyed by projection surface. |
| `relationships` | map | yes | Relationship lists. Unit-ID relationships are validated. |
| `index` | object | yes | Retrieval fields. |
| `provenance` | object | yes | Revision and decision references. |
| `component_consumers` | string list | yes | Components expected to consume this unit. |
| `source_unit` | string | translation only | Source unit ID for translation units. |
| `translation_status` | enum | translation only | Translation review status. |

## Relationship Fields

The validator treats these fields as references to other semantic unit IDs:

- `depends_on`
- `clarifies`
- `implements`
- `illustrated_by`
- `translated_by`
- `supersedes`

Other relationship fields may be added later, but they are not considered stable
consumer contract until the validator recognizes them.

## Authority Types

| Value | Meaning |
| --- | --- |
| `canonical_markdown_bridge` | Current English Markdown remains canonical while indexed. |
| `canonical_semantic` | Future accepted source-of-truth unit. |
| `translation` | Localized surface linked to a source unit. |
| `projection` | Reader-facing derived output. |
| `generated_artifact` | Generated output with provenance, never source authority. |
| `design_record` | Governed decision or analysis evidence. |

## Projection Rules

- Canonical bridge units should expose at least `markdown_en`.
- Translated or localized Core units should expose `markdown_pt_br` when a PT-br
  document exists.
- Generated surfaces must remain `projection` or `generated_artifact`.
- AI context packs and search indexes are retrieval surfaces, not rule authority.

## Validation

Run:

```bash
go run scripts/validate_semantic_index.go docs/core/semantic/core-v2.index.json
```

The validator checks:

- required top-level fields;
- non-empty vocabularies;
- unique unit IDs;
- allowed `kind`, `authority_type`, `source_status`, projection surfaces, and
  component consumers;
- repository-relative source and projection paths;
- relationship target IDs;
- translation `source_unit` and `translation_status`;
- provenance `source_revision`;
- repository-relative `decision_refs` paths.
