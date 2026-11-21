# Procedure: Validate Example

Verify whether a worked example's numbers and state transitions obey the
canonical rule(s) it illustrates.

## Input

A worked example — either a `core.example.*` semantic unit or custom
numbers — together with the rule(s) it claims to illustrate.

## Steps

1. Identify the governing rule(s)/procedure(s) via
   `relationships.illustrated_by` on the rule side or `relationships.depends_on`
   on the example side (e.g. `core.combat.attack-margin` →
   `illustrated_by: [core.example.combat-turn]`).
2. Independently re-derive every numerical step from the rule's `formula`
   or `ordered_steps` in `docs/core/semantic/source/core-v2-rules.v0.1.json` —
   verify the arithmetic directly rather than trusting narrative phrasing.
3. Compare the re-derived values to the example's stated values field by
   field.
4. Enforce every structural invariant implied by the rule (e.g.
   `core.damage.flow`: health never increases from a damage step; shield
   absorbs before health).
5. Report **every** mismatch found across all fields and invariants, then
   cite all governing rule IDs and `source_path`s.

## Output Shape

```text
Example: <id or description>
Rule(s): <ids>
Result: <CONSISTENT | INCONSISTENT>
Mismatches: <field-by-field list, or "none">
Cited: <rule ids and their source_paths>
```

## Worked Example

Read [`../examples/validate-example-combat-turn.md`](../examples/validate-example-combat-turn.md)
only when auditing multi-step combat turns involving attack margin, armor,
and shield absorption.
