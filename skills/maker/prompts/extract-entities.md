# Prompt: Extract Entities

Extract campaign entities (NPCs, locations, items, factions) from raw GM
notes, text, or images into validated packages grounded strictly in the
story's own sources.

## Steps

1. Read the entire input batch before extracting anything so late facts
   can reclassify or resolve earlier passages.
2. Open one [`entity-package`](../templates/entity-package.md) per
   distinct entity (`npc`, `location`, `item`, `faction`, `other`); check
   same-named mentions carefully and route ambiguity through
   [`detect-conflicts.md`](detect-conflicts.md) before merging.
3. Enforce `layer-separated` by placing every candidate property into
   exactly one bucket and checking the entity type's normal field set
   deliberately:
   - **Explicit fact** — stated directly in a text source (`source_ref`
     must match the package's `source_refs`).
   - **Inferred suggestion** — plausible given the story sources alone,
     tagged with `confidence` (`low|medium|high`) and explicit `rationale`.
   - **Unresolved field** — expected for the `entity_type`, missing from
     the source, and not licensed for inference (leave `"None."` only after
     checking deliberately).
4. Route every image input through a separate
   [`visual-observation`](../templates/visual-observation.md) per
   `text-over-image` (appearance only, never directly into `explicit_facts`),
   linking to an entity only when a text source or the user confirms the
   identity, and checking conflicts with
   [`detect-conflicts.md`](detect-conflicts.md).
5. When two mentions of the same entity disagree on a field, keep the
   disputed field in `unresolved_fields`, open a
   [`conflict-report`](../templates/conflict-report.md) per
   `conflict-on-dispute`, and link its ID in `conflicts`.
6. Summarize the entire batch in one
   [`maker-report`](../templates/maker-report.md) listing every emitted
   package ID, and set `status` to `draft` or `needs_review` per
   `non-canonical-status`.

## Completion Check

- Every emitted `EntityPackage` has at least one `explicit_facts` item and
  disjoint field names across all three layers.
- Every fact, inference, and open question stays strictly within the
  story's own universe without injecting external game mechanics.
- Every conflict and image observation ID referenced in the report matches
  an emitted package.

## Worked Example

Read [`../examples/harbor-quartermaster.md`](../examples/harbor-quartermaster.md)
only when checking how mixed text notes and an image sketch split across
packages.
