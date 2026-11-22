# Prompt: Extract Entities

Extract every distinct story entity across all five scopes (`npc`, `location`,
`item`, `faction`, `other`) from raw text or images into packages grounded
exclusively in the story's own sources.

## Steps

1. Read the entire input batch before extracting so later passages can clarify
   or reclassify earlier mentions.
2. Open one [`entity-package`](../templates/entity-package.md) per distinct
   named actor, organization, place/zone, item/technology, or standalone event,
   project, phenomenon, or unclassified system; route ambiguous same-named
   mentions through [`detect-conflicts.md`](detect-conflicts.md) before merging.
3. Apply `layer-separated` from [`../SKILL.md`](../SKILL.md) to every candidate
   property, checking the entity type's normal field set deliberately:
   - **Explicit fact** — stated directly in a text source (`source_ref` declared
     in `source_refs`).
   - **Inferred suggestion** — plausible from story sources alone, tagged with
     `confidence` (`low|medium|high`) and explicit `rationale`.
   - **Unresolved field** — expected for the `entity_type` and absent from the
     sources (write `"None."` only after deliberate checking).
4. Route every image through [`visual-observation`](../templates/visual-observation.md)
   per `text-over-image` (grouping large collections into thematic packages plus
   `sources/images/IMAGE_CATALOG.md`).
5. Route single-field disagreements per `conflict-on-dispute` and multi-field
   reboots or drafts per `variant-preserved`.
6. Summarize the batch in one [`maker-report`](../templates/maker-report.md)
   per `non-canonical-status` and `prompt-user-decisions`.

## Completion Check

- Every distinct named entity, event, project, or unclassified system has its
  own dedicated `entity-package`.
- Every `entity-package` has at least one `explicit_facts` entry and disjoint
  field names across all three layers.
- Every reboot variant, `needs_review` entity, unlinked visual observation, and
  unresolved conflict is surfaced in `pending_decisions`.
- Every referenced conflict and visual observation ID matches an emitted package.

## Worked Example

Read [`../examples/harbor-quartermaster.md`](../examples/harbor-quartermaster.md)
only when checking how mixed text notes and an image sketch split across
packages.
