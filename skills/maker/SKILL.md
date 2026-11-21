---
name: rgb-maker
description: RGB Maker content structurer. Use when turning raw GM notes, text, or images into entity packages, visual observations, conflict reports, or Maker session reports.
---

# RGB Maker Skill

Structures raw GM notes, text, and images into campaign/world artifacts
(NPCs, locations, items, factions) with explicit provenance separation.
Conforms to [`consumer.maker.v0_1`](../../docs/core/semantic/consumer-contracts.v0.1.json),
mirrored for quick lookup in [`config.yaml`](config.yaml) (the JSON file
wins if they disagree). See [ADR-001](../../docs/adr/adr-001-ai-first-documentation-authority.md)
for the AI-first documentation authority model.

## Workflow & Package Dispatch

1. **Extract entities** from notes, text, or images → follow
   [`prompts/extract-entities.md`](prompts/extract-entities.md).
2. **Resolve or record source disputes** when two or more inputs conflict
   on a field → follow
   [`prompts/detect-conflicts.md`](prompts/detect-conflicts.md).
3. **Emit schema + template pairs** for every package produced:

| Kind | Schema | Template | Purpose |
|---|---|---|---|
| Entity Package | [`schemas/entity-package.schema.yaml`](schemas/entity-package.schema.yaml) | [`templates/entity-package.md`](templates/entity-package.md) | one extracted entity with facts/inferences/unresolved fields |
| Maker Report | [`schemas/maker-report.schema.yaml`](schemas/maker-report.schema.yaml) | [`templates/maker-report.md`](templates/maker-report.md) | session-level summary of one input batch |
| Visual Observation | [`schemas/visual-observation.schema.yaml`](schemas/visual-observation.schema.yaml) | [`templates/visual-observation.md`](templates/visual-observation.md) | appearance-only properties read from an image |
| Conflict Report | [`schemas/conflict-report.schema.yaml`](schemas/conflict-report.schema.yaml) | [`templates/conflict-report.md`](templates/conflict-report.md) | single disputed field awaiting resolution or authority-rule note |

4. **Validate and render deterministically** through the Go component CLI
   when running inside this monorepo:

```bash
go run ./cmd/rgb maker process --batch <batch.json> --index docs/core/semantic/core-v2.index.json
```

## Core Operating Rules

- **`layer-separated`**: Classify every entity candidate property into
  exactly one layer — `explicit_facts`, `inferred_suggestions`, or
  `unresolved_fields` — and emit both machine-readable schema fields and
  rendered Markdown (`templates/*.md`).
- **`text-over-image`**: Text is authoritative for all facts (names, stats,
  relationships, canon status); images are authoritative for appearance only
  (`schemas/visual-observation.schema.yaml`). When text and image conflict on
  a non-appearance property, text wins outright.
- **`conflict-on-dispute`**: Whenever sources disagree on a field that
  `text-over-image` does not settle (text vs. text, or image vs. image on
  appearance), open a `conflict-report` with `resolution: unresolved` so the
  user decides; keep the disputed field in `unresolved_fields` on the entity.
- **`non-canonical-status`**: Disclose `draft`, `stale`, `needs_review`, and
  `generated` on every output. Consult `docs/core/semantic/core-v2.index.json`
  strictly as read-only reference for `canonical_dependencies` (using
  `canonical_markdown_bridge`, `canonical_semantic`, `translation`, or
  `design_record` units — never `projection` or `generated_artifact`).
  Maker is an optional authoring aid and never creates or overrides canon.

## The Three Truth Layers

Keep these three layers distinct across all artifacts:

- **canonical truth** — what RGB System rules define (`core-v2.index.json`).
- **public or scoped knowledge** — what characters or factions in the
  fiction are meant to know.
- **observed representation** — what a specific note or image shows, which
  may be incomplete, biased, or wrong in-fiction.

## Worked Example

Read [`examples/harbor-quartermaster.md`](examples/harbor-quartermaster.md)
only when checking end-to-end packaging across text notes, an image sketch,
and a text-vs-text conflict report.
