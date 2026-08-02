# Example: Validate A Draft Ability

**Procedure:** [`../procedures/validate.md`](../procedures/validate.md)

## Input

```yaml
id: ability.suppressive-burst
name: Suppressive Burst
vector: R
tier: 2
requirements: { R: 4 }
action_type: action
cost: 2
range: 12
duration: instant
effect: ["target loses their minor action next turn if hit"]
tags: ["combat", "control"]
```

## Result

```text
Ability: ability.suppressive-burst
Result: PASS
Missing fields: none
Violations: none
Cited: core.ability.contract (docs/core/en/core/skills_and_abilities.md)
```

All thirteen `core.ability.contract` required fields are present, and the
value-level checks mirrored from `internal/components/core/ability.go`
(`Ability.Validate`) all hold: vector is a valid RGB letter, tier ≥ 1,
requirement values non-negative, `action_type` is `action`, cost/range
non-negative, duration non-empty, effect and tags non-empty.

## Contrast: A Failing Ability

```yaml
id: ability.mystery-strike
name: Mystery Strike
vector: R
tier: 0
action_type: instant
cost: -1
```

```text
Ability: ability.mystery-strike
Result: FAIL
Missing fields: requirements, range, duration, effect, tags, source_status
Violations: tier must be at least 1 (got 0); action_type "instant" is not
  "action" or "reaction"; cost must be non-negative (got -1)
Cited: core.ability.contract (docs/core/en/core/skills_and_abilities.md)
```
