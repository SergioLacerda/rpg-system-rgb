# Worked Example: Harbor Quartermaster

Demonstrates [`prompts/extract-entities.md`](../prompts/extract-entities.md)
and [`prompts/detect-conflicts.md`](../prompts/detect-conflicts.md) end to
end.

## Input

**Note 1** (session prep, text):

> The quartermaster at the harbor is a gruff woman named Yssa Kell. She's
> missing two fingers on her left hand — lost them to a rope winch years
> ago. She controls who gets dock space, so the party will need her favor.

**Note 2** (session recap, text, written later):

> Reminder: the harbor quartermaster's name is Yssa Voss, not Kell — I
> misheard it at the table.

**Image** (`harbor-npc-sketch.png`): a sketch of a stern woman in a
weathered coat, both hands visible and undamaged.

## Output

### `observation.harbor-npc-sketch`

```yaml
schema: rgb-maker-visual-observation/0.1
id: observation.harbor-npc-sketch
status: draft
source_ref: harbor-npc-sketch.png
appearance_facts:
  - field: clothing
    value: weathered coat
    confidence: high
  - field: hands
    value: both hands appear undamaged in the sketch
    confidence: medium
linked_entity: npc.harbor-quartermaster
non_authoritative_note: >
  This observation may describe appearance only. It does not establish or
  override any fact, name, relationship, or canon status — those come from
  text sources only, per the authority rule in SKILL.md.
```

The sketch's "undamaged hands" is an appearance observation. It is
**not** a conflict with Note 1's missing fingers — per the authority rule,
the image was never authoritative over that fact (physical injury history
is a fact, not appearance-only), so it's simply recorded, not disputed.

### `conflict.harbor-quartermaster-name`

```yaml
schema: rgb-maker-conflict-report/0.1
id: conflict.harbor-quartermaster-name
status: needs_review
field: name
linked_entity: npc.harbor-quartermaster
claims:
  - value: Yssa Kell
    source_ref: "Note 1 (session prep)"
    source_kind: text
  - value: Yssa Voss
    source_ref: "Note 2 (session recap)"
    source_kind: text
resolution: unresolved
```

Both claims are `text` — the authority rule doesn't resolve text-vs-text
disagreements. This stays `unresolved` for the GM to pick.

### `npc.harbor-quartermaster`

```yaml
schema: rgb-maker-entity-package/0.1
id: npc.harbor-quartermaster
entity_type: npc
status: needs_review
source_refs: ["Note 1 (session prep)", "Note 2 (session recap)", "harbor-npc-sketch.png"]
explicit_facts:
  - field: role
    value: harbor quartermaster, controls dock space assignments
    source_ref: "Note 1 (session prep)"
  - field: missing_fingers
    value: two fingers missing on left hand, lost to a rope winch
    source_ref: "Note 1 (session prep)"
inferred_suggestions:
  - field: attitude_toward_party
    value: likely transactional/guarded until the party earns favor
    confidence: medium
    rationale: "note 1 frames dock access as something the party must earn"
unresolved_fields:
  - field: name
    reason: "disputed — see conflict.harbor-quartermaster-name"
conflicts: ["conflict.harbor-quartermaster-name"]
canonical_dependencies: []
```

Note that `name` moved from what would have been an explicit fact into
`unresolved_fields`, with a pointer to the conflict report — Maker does
not pick "Kell" or "Voss" on its own.

### `report.2026-08-01-harbor-session`

```yaml
schema: rgb-maker-report/0.1
id: report.2026-08-01-harbor-session
status: draft
input_batch: ["Note 1 (session prep)", "Note 2 (session recap)", "harbor-npc-sketch.png"]
entity_packages: [npc.harbor-quartermaster]
conflict_reports: [conflict.harbor-quartermaster-name]
visual_observations: [observation.harbor-npc-sketch]
open_questions:
  - "Confirm the quartermaster's actual name (Kell vs Voss) before next session."
```
