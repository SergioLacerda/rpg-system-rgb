# Procedure: Validate

Check whether a proposed ability declaration satisfies `core.ability.contract`.

## Input

An ability description — structured (YAML/JSON) or prose — providing its
declaration fields.

## Steps

1. Load `core.ability.contract`'s `required_fields` directly from
   `docs/core/semantic/source/core-v2-rules.v0.1.json`.
2. Verify every required field is present and non-empty; collect **every**
   missing field rather than stopping at the first.
3. Apply the Go engine's value-level constraints from
   `internal/components/core/ability.go` (`Ability.Validate`) so Specialist
   matches the engine:
   - `vector`: valid RGB vector (`R`, `G`, or `B`)
   - `tier`: integer >= 1
   - `requirements`: every value >= 0
   - `action_type` (Go: `Timing`): `action` or `reaction`
   - `cost` and `range`: >= 0
   - `duration`: non-empty
   - `effect`/`effects` and `tags`: non-empty lists
   If `core.ability.contract` and `Ability.Validate` ever disagree on a rule,
   report both sides explicitly as a contract/engine mismatch rather than
   silently choosing one.
4. Report every missing field and every value violation found (`PASS` only
   when both lists are empty).
5. Cite `core.ability.contract` with its `source_path` and disclose any
   status/locale flags (`citation-gated`).

## Output Shape

```text
Ability: <id or name>
Result: <PASS | FAIL>
Missing fields: <list, or "none">
Violations: <list, or "none">
Cited: core.ability.contract (docs/core/en/core/skills_and_abilities.md)
Disclosures: <draft | stale | needs_review | localized | generated, or "none (canonical en)">
```

## Worked Example

Read [`../examples/validate-sample-ability.md`](../examples/validate-sample-ability.md)
only when validating an ability written in free-form prose or checking
field-mapping nuances.
