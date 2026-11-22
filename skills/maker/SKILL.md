---
name: rgb-maker
description: Story and campaign structurer. Use when organizing raw story notes, lore, or images into entity packages and visual observations, resolving reboots or source conflicts, measuring workspace stats, or exporting a compact chat-loading pack.
---

# Maker Skill

Structures raw story notes, lore, and images into organized, source-grounded
artifacts (`npc`, `location`, `item`, `faction`, and `other`: pivotal events,
projects, phenomena, or unclassified named systems) strictly within the target
story's own universe.

Conforms to [`consumer.maker.v0_1`](../../docs/core/semantic/consumer-contracts.v0.1.json)
(mirrored in [`config.yaml`](config.yaml)).

## Workflow

1. **Organize raw inputs** by type before extraction:
   - `sources/text/` — notes, lore docs, recaps, transcripts
   - `sources/images/` — portraits, concept art, sketches, maps (preserve subfolders)
2. **Extract entities** from all text inputs → follow
   [`prompts/extract-entities.md`](prompts/extract-entities.md).
3. **Record visual observations and resolve or surface disputes** → follow
   [`prompts/detect-conflicts.md`](prompts/detect-conflicts.md). For large image
   sets (20+ files or nested folders), group related visuals into thematic
   `visual-observation` packages and emit a full file-to-package manifest at
   `sources/images/IMAGE_CATALOG.md`.
4. **Validate and render packages** under `packages/` (`maker-report.md`,
   `entities/`, `visual-observations/`, `conflicts/` for open disputes, and
   `conflicts/resolved/` for archived decisions):
   - `go run ./cmd/rgb maker process --batch <batch.json> --out-dir <packages-dir>`
     (validates schemas, archives resolved conflicts to `conflicts/resolved/`,
     and prunes stale `.md` files).
5. **Measure or export on demand**:
   - `go run ./cmd/rgb maker stats --dir <workspace-dir>` — raw vs. filtered line
     counts and cataloged image breakdown.
   - `go run ./cmd/rgb maker export --dir <workspace-dir> --out <chat-dir>` —
     condenses the workspace into four token-efficient Markdown files
     (`00_BOOT_PROMPT.md`, `01_CORE_AND_RULES.md`, `02_CAST_AND_FACTIONS.md`,
     `03_TECH_AND_LOCATIONS.md`) for chat loading.

| Kind | Schema | Template |
|---|---|---|
| Entity Package | [`schemas/entity-package.schema.yaml`](schemas/entity-package.schema.yaml) | [`templates/entity-package.md`](templates/entity-package.md) |
| Maker Report | [`schemas/maker-report.schema.yaml`](schemas/maker-report.schema.yaml) | [`templates/maker-report.md`](templates/maker-report.md) |
| Visual Observation | [`schemas/visual-observation.schema.yaml`](schemas/visual-observation.schema.yaml) | [`templates/visual-observation.md`](templates/visual-observation.md) |
| Conflict Report | [`schemas/conflict-report.schema.yaml`](schemas/conflict-report.schema.yaml) | [`templates/conflict-report.md`](templates/conflict-report.md) |

## Core Rules (Single Source of Truth)

- **`layer-separated`**: Place every candidate property into exactly one bucket
  — `explicit_facts`, `inferred_suggestions`, or `unresolved_fields` — grounded
  strictly in the story's own sources, while keeping **canonical truth**,
  **in-fiction/public knowledge**, and **observed representation** distinct.
  Keep `canonical_dependencies` empty unless the user explicitly requests linking
  an external rule system.
- **`text-over-image`**: Text governs all facts (names, stats, relationships,
  canon status); images govern appearance only. Link an image to a named entity
  only when a text source or the user explicitly confirms the match; otherwise
  leave `linked_entity` empty and mark `status: needs_review`.
- **`conflict-on-dispute`**: For any single-field disagreement not settled by
  `text-over-image` (text vs. text, or image vs. image on appearance), open a
  `conflict-report` with `resolution: unresolved` under `packages/conflicts/`
  and keep the field in `unresolved_fields` on the entity. Once resolved
  (`resolved_by_user` or `resolved_by_authority_rule`), archive the report under
  `packages/conflicts/resolved/` so `packages/conflicts/` holds open disputes
  only.
- **`variant-preserved`**: When multiple sources represent reboots, remakes, or
  distinct timeline drafts diverging across multiple fields, record shared core
  facts once in `explicit_facts` and keep each version intact as its own
  `variants` entry with `status: needs_review`.
- **`prompt-user-decisions`**: At the end of every extraction or image-mapping
  run, present a numbered decision prompt covering every `needs_review` entity,
  unlinked visual observation, unresolved conflict, and multi-variant entity so
  the user can steer final canon promotion.
- **`non-canonical-status`**: Mark every output with `draft`, `stale`,
  `needs_review`, or `generated`; preserve all user decisions verbatim.

## Completion Check

- `maker process` succeeds with zero validation errors and zero stale `.md` files.
- Every resolved conflict is archived in `packages/conflicts/resolved/`, leaving
  only `unresolved` reports in `packages/conflicts/`.
- Every remaining ambiguity is surfaced in `pending_decisions` and presented to
  the user.

## Worked Example

Read [`examples/harbor-quartermaster.md`](examples/harbor-quartermaster.md)
only when checking end-to-end packaging across text notes, an image sketch,
and a text-vs-text conflict report.
