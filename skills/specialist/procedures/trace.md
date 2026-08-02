# Procedure: Trace

Show the full provenance chain behind a previous Specialist answer.

## Input

A prior answer (or the semantic ID(s) it cited).

## Steps

1. For each cited semantic ID, look it up in
   `docs/core/semantic/core-v2.index.json`.
2. Report, per ID:
   - `kind` (concept, rule, procedure, resource, ability, example, term,
     translation);
   - `source_status` (e.g. `canonical`) and `authority_type` (e.g.
     `canonical_markdown_bridge`);
   - `source_path` and `projection_paths` (where the Markdown actually
     lives, EN and PT-br);
   - `relationships` (`depends_on`, `clarifies`, `illustrated_by`, etc.) —
     this is what lets a user follow the chain further (e.g.
     `core.combat.attack-margin` → `depends_on` → `core.vector.r`,
     `core.vector.g`).
   - `provenance.decision_refs` — the ADR(s) or design records that
     established this unit, if any.
3. If a cited ID does not exist in the index, that is a Specialist defect
   (it cited something ungrounded) — report it as such rather than
   papering over it.

## Output Shape

```text
ID: <semantic id>
Kind: <kind> | Status: <source_status> (<authority_type>)
Source: <source_path>
Depends on: <ids, or "none">
Decision refs: <adr paths, or "none">
```

## Worked Example

See [`../examples/trace-combat-turn-example.md`](../examples/trace-combat-turn-example.md).

## Forbidden

- Tracing to a path that is not the unit's declared `source_path` /
  `projection_paths` (no paraphrasing the trail).
