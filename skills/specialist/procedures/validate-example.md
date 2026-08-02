# Procedure: Validate Example

Check whether a worked example's numbers are internally consistent with
the rule(s) it illustrates. This generalizes the fix behind
`.analysis/refined/20260801-core-rules-executable-examples-turn-contract`
(a published combat example that had shield absorption backwards and
health increasing after damage) into a repeatable procedure.

## Input

A worked example — either a `core.example.*` semantic unit, or an
arbitrary example text (e.g. a GM's own numbers) plus the rule(s) it
claims to follow.

## Steps

1. Identify the rule(s)/procedure(s) the example is illustrating, via
   `relationships.illustrated_by` on the rule side, or `relationships.depends_on`
   on the example side (see `core.combat.attack-margin`'s
   `illustrated_by: [core.example.combat-turn]`).
2. Re-derive the example's numbers from the rule's `formula` or
   `ordered_steps` independently — do not just re-read the example's
   stated result and assume it followed the formula.
3. Compare the re-derived numbers to the example's stated numbers, field
   by field.
4. Apply any structural invariants the rule implies even if the example
   doesn't state them (e.g. `core.damage.flow`: health must never increase
   from a damage step; shield absorbs before health).
5. Report every mismatch found — do not stop at the first one.

## Output Shape

```text
Example: <id or description>
Rule(s): <ids>
Result: <CONSISTENT | INCONSISTENT>
Mismatches: <field-by-field list, or "none">
Cited: <rule ids and their source_paths>
```

## Worked Example

See [`../examples/validate-example-combat-turn.md`](../examples/validate-example-combat-turn.md)
— includes the actual `tests/fixtures/combat_example_doc_parity_test.go`
regression test this procedure describes in prose form.

## Forbidden

- Declaring an example consistent because its narrative reads plausibly —
  the numbers must be independently re-derived and compared, not
  eyeballed.
