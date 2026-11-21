# Procedure: Trace

Show the full provenance chain behind a previous Specialist answer.

## Input

A prior answer or the semantic ID(s) it cited.

## Steps

1. Look up each cited semantic ID in `docs/core/semantic/core-v2.index.json`.
2. Report each unit's provenance fields verbatim from the index entry (never
   paraphrase paths):
   - `kind` and `source_status` (`authority_type`)
   - `source_path` and `projection_paths` (EN and PT-br)
   - `relationships` (`depends_on`, `clarifies`, `illustrated_by`, etc.,
     such as `core.combat.attack-margin` → `core.vector.r`, `core.vector.g`)
   - `provenance.decision_refs`
3. If any cited ID is missing from the index, flag it explicitly as an
   ungrounded Specialist defect.

## Output Shape

```text
ID: <semantic id>
Kind: <kind> | Status: <source_status> (<authority_type>)
Source: <source_path>
Depends on: <ids, or "none">
Decision refs: <adr paths, or "none">
```

## Worked Example

Read [`../examples/trace-combat-turn-example.md`](../examples/trace-combat-turn-example.md)
only when tracing multi-hop `relationships` chains across vectors and rules.
