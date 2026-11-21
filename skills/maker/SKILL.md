---
name: rgb-maker
description: Story and campaign content structurer. Use when turning raw GM notes, lore documents, text, or images into organized sources, entity packages, visual observations, conflict reports, or session reports.
---

# Maker Skill

Structures raw campaign notes, lore documents, text, and images into
organized story artifacts (NPCs, locations/zones, items/resources,
factions, and pivotal events, projects, or named systems) with explicit
provenance separation — staying focused purely on the target story's own
universe and sources without projecting external rules or ungrounded facts.

Conforms to [`consumer.maker.v0_1`](../../docs/core/semantic/consumer-contracts.v0.1.json)
(mirrored in [`config.yaml`](config.yaml)).

## Workflow & Folder Organization

1. **Organize raw inputs by type** before extraction when reorganizing a
   story/campaign workspace:
   - `sources/text/` — notes, lore docs, session recaps, transcripts
   - `sources/images/` — portraits, concept art, sketches, maps
2. **Extract entities** from all text inputs → follow
   [`prompts/extract-entities.md`](prompts/extract-entities.md).
3. **Record appearance and resolve/surface disputes** for images and
   conflicting sources → follow
   [`prompts/detect-conflicts.md`](prompts/detect-conflicts.md).
4. **Emit schema + template pairs** under `packages/` (`maker-report.md`,
   `entities/`, `visual-observations/`, `conflicts/` for open disputes, and
   `conflicts/resolved/` for archived decisions):

| Kind | Schema | Template | Purpose |
|---|---|---|---|
| Entity Package | [`schemas/entity-package.schema.yaml`](schemas/entity-package.schema.yaml) | [`templates/entity-package.md`](templates/entity-package.md) | one extracted entity with facts/inferences/unresolved fields |
| Maker Report | [`schemas/maker-report.schema.yaml`](schemas/maker-report.schema.yaml) | [`templates/maker-report.md`](templates/maker-report.md) | session-level summary of one input batch |
| Visual Observation | [`schemas/visual-observation.schema.yaml`](schemas/visual-observation.schema.yaml) | [`templates/visual-observation.md`](templates/visual-observation.md) | appearance-only properties read from an image |
| Conflict Report | [`schemas/conflict-report.schema.yaml`](schemas/conflict-report.schema.yaml) | [`templates/conflict-report.md`](templates/conflict-report.md) | single disputed field awaiting resolution or authority-rule note |

## Core Operating Rules

- **`layer-separated`**: Classify every candidate property into exactly one
  layer — `explicit_facts`, `inferred_suggestions`, or `unresolved_fields` —
  and emit both structured metadata and rendered Markdown (`templates/*.md`).
  Ground every fact and inference strictly in the story's own sources; leave
  `canonical_dependencies` empty unless the user explicitly asks to link an
  external rule system.
- **`text-over-image`**: Text is authoritative for all facts (names, stats,
  relationships, canon status); images are authoritative for appearance only
  (`schemas/visual-observation.schema.yaml`). When text and image conflict on
  a non-appearance property, text wins outright. Never assign an image to a
  named entity unless a text source or the user explicitly links them.
- **`conflict-on-dispute`**: Whenever sources disagree on a field that
  `text-over-image` does not settle (text vs. text, or image vs. image on
  appearance), open a `conflict-report` with `resolution: unresolved` under
  `packages/conflicts/` so the user decides; keep the disputed field in
  `unresolved_fields` on the entity. Once resolved (`resolved_by_user` or
  `resolved_by_authority_rule`), archive the report under
  `packages/conflicts/resolved/` so `packages/conflicts/` holds open disputes
  only.
- **`variant-preserved`**: When multiple inputs represent reboots, remakes,
  or distinct timeline versions of the same entity (where multiple fields
  diverge structurally), keep shared core facts in `explicit_facts` and
  preserve each distinct version intact under `variants` as structured
  appendices with `status: needs_review` — never blend two reboots into a
  single muddled description.
- **`prompt-user-decisions`**: At the end of every extraction run, surface a
  numbered decision/guidance prompt for every ambiguity (`needs_review`
  entities of uncertain nature, unlinked images, unresolved conflicts, and
  multi-variant entities awaiting final canon choice) so the user can steer
  promotion directly.
- **`non-canonical-status`**: Disclose `draft`, `stale`, `needs_review`, or
  `generated` on every output. Maker is an authoring and organization aid; it
  never invents story canon or overwrites user decisions.

## The Three Truth Layers

Keep these three layers distinct across all artifacts:

- **canonical truth** — established facts within the story's primary canon.
- **public or scoped knowledge** — what characters or factions in the
  fiction are meant to know.
- **observed representation** — what a specific note or image shows, which
  may be incomplete, biased, or wrong in-fiction.

## Worked Example

Read [`examples/harbor-quartermaster.md`](examples/harbor-quartermaster.md)
only when checking end-to-end packaging across text notes, an image sketch,
and a text-vs-text conflict report.
