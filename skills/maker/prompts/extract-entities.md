# Prompt: Extract Entities

Use when the input is raw GM notes, text, or images that may describe one
or more entities (NPCs, locations, items, factions).

## Instructions

1. Read the entire input before extracting anything — a fact stated late
   in the notes may change how an earlier passage should be classified.
2. For each distinct entity found, open one
   [`entity-package`](../templates/entity-package.md).
3. Classify every candidate fact into exactly one bucket:
   - **Explicit fact** — the source states it directly ("the quartermaster
     is missing three fingers on his left hand").
   - **Inferred suggestion** — plausible given the source, but not stated
     ("likely lost the fingers in the same accident mentioned earlier" —
     mark confidence and say why).
   - **Unresolved field** — the entity type normally has this field (e.g.
     a faction usually has a stated goal), the source doesn't cover it,
     and nothing licenses an inference.
4. If an image is part of the input, route its contribution through
   [`detect-conflicts.md`](detect-conflicts.md)'s authority-rule check
   before merging it into the entity package — appearance facts go through
   a [`visual-observation`](../templates/visual-observation.md), never
   directly into `explicit_facts`.
5. If the same entity appears more than once in the batch with differing
   details, do not merge silently — open a
   [`conflict-report`](../templates/conflict-report.md) for the disputed
   field(s) and reference it from both mentions.
6. When the batch is done, summarize it in one
   [`maker-report`](../templates/maker-report.md) listing every package
   produced.

## Do Not

- Invent a name, relationship, or stat not licensed by an explicit fact or
  a clearly-marked inference.
- Merge two same-named entities without checking whether they're actually
  the same entity (a conflict-report candidate, not an assumption).
- Skip the unresolved-fields section because nothing came to mind — check
  the entity type's normal field set deliberately.
