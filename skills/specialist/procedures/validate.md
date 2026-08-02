# Procedure: Validate

Check whether a proposed ability declaration satisfies `core.ability.contract`.

## Input

An ability description — structured (a draft YAML/JSON block) or prose
that names the ability's id, name, vector, tier, requirements, action type,
cost, range, duration, effect, limits, and tags.

## Steps

1. Load `core.ability.contract`'s `required_fields` from
   `docs/core/semantic/source/core-v2-rules.v0.1.json`:

   ```text
   id, name, vector, tier, requirements, action_type, cost, range,
   duration, effect, limits, tags, source_status
   ```

2. Check every required field is present and non-empty in the input. List
   every missing field — do not stop at the first one.
3. Apply the same value-level checks the Go engine enforces
   (`internal/components/core/ability.go`, `Ability.Validate`), so a
   Specialist "pass" and an engine "pass" never disagree:
   - `vector` must be a valid RGB vector (R, G, or B);
   - `tier` must be at least 1;
   - each `requirements` value must be non-negative;
   - `action_type` (Go: `Timing`) must be `action` or `reaction`;
   - `cost` and `range` must be non-negative;
   - `duration` must be non-empty;
   - `effect`/`effects` and `tags` must each be non-empty lists.
4. Report every failure found (missing fields and value-level violations),
   not just the first.
5. If everything passes, report pass and cite `core.ability.contract`.

## Output Shape

```text
Ability: <id or name>
Result: <PASS | FAIL>
Missing fields: <list, or "none">
Violations: <list, or "none">
Cited: core.ability.contract (docs/core/en/core/skills_and_abilities.md)
```

## Worked Example

See [`../examples/validate-sample-ability.md`](../examples/validate-sample-ability.md).

## Forbidden

- Declaring an ability valid without checking every required field.
- Inventing a validation rule not present in `core.ability.contract` or
  `internal/components/core/ability.go` — if the two ever disagree, say so
  instead of picking one silently; that is a real product bug to report,
  not something Specialist resolves on its own.
