# Prompt: Extract Entities

Extract campaign entities (`npc`, `location`, `item`, `faction`, and
`other`: pivotal events, projects, phenomena, or named systems of
uncertain nature) from raw GM notes, text, or images into validated
packages grounded strictly in the story's own sources.

## Steps

1. Read the entire input batch before extracting anything so late facts
   can reclassify or resolve earlier passages.
2. Open one [`entity-package`](../templates/entity-package.md) per
   distinct entity across all five scopes (`npc` characters/operators,
   `location` places/zones/territories, `item` resources/gear/artifacts,
   `faction` organizations/cells, `other` pivotal events, projects,
   world phenomena, or unclassified named systems); check same-named
   mentions carefully and route ambiguity through
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
5. When two mentions of the same entity disagree on a single field, keep
   the disputed field in `unresolved_fields`, open a
   [`conflict-report`](../templates/conflict-report.md) per
   `conflict-on-dispute`, and link its ID in `conflicts`. When sources
   reflect distinct reboots, remakes, or alternative drafts across multiple
   fields, record shared core facts once and preserve each version under
   `variants` per `variant-preserved` with `status: needs_review`.
6. Summarize the entire batch in one
   [`maker-report`](../templates/maker-report.md) listing every emitted
   package ID, set `status` to `draft` or `needs_review` per
   `non-canonical-status`, and emit `pending_decisions` per
   `prompt-user-decisions` so the user can resolve ambiguities (e.g. an
   unclassified system/character), assign unlinked images, or pick the
   final reboot version.

## Completion Check

- Every distinct named actor, faction, location/zone, item/resource, and
  pivotal event, project, or unclassified system in the source batch has
  its own `entity-package` — never fold a standalone event, project, or
  named entity into another package's description.
- Every reboot or multi-source alternative version is preserved under
  `variants`, and every `needs_review` entity, unlinked image, unresolved
  conflict, and variant entity is surfaced in `pending_decisions` for user
  orientation.
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
