# Example: Identify Ambiguity — "Can a character Sprint and still use a reaction ability?"

**Procedure:** [`../procedures/identify-ambiguity.md`](../procedures/identify-ambiguity.md)

## Attempted Resolution

[`locate-authority.md`](../procedures/locate-authority.md) finds
`docs/core/en/combat/movement.md`'s "Optional Rule: Sprint," which states
restrictions "No attack during the same turn" and "Only movement and minor
actions allowed" — but Sprint is documented as prose, not as a
`core.term.*`/`core.action.*` semantic unit, so there is no semantic ID to
cite for it directly.

`internal/components/core/ability.go`'s `Timing` enum distinguishes
`action` and `reaction` abilities, but Sprint's restriction text predates
that distinction and does not mention reactions at all.

## Report

```text
Question: Can a character Sprint and still use a reaction ability the
  same turn?
Status: AMBIGUOUS
Consulted:
  - Sprint rule (docs/core/en/combat/movement.md) — restricts "attack" and
    limits to "movement and minor actions," but does not use the
    action/reaction vocabulary and does not mention reaction abilities.
  - core.ability.contract — defines action_type as action|reaction but
    does not cross-reference Sprint's restrictions.
Possible inference (non-authoritative): Sprint's restriction reads as
  aimed at the character's own turn economy (no attack, no full action),
  which a reaction by definition falls outside of — so a reaction ability
  plausibly remains usable. This is a GM call, not a rule; the text does
  not settle it either way.
```

This is a genuine gap, not a "no unit found" case — Sprint IS documented,
just not precisely enough to answer this specific question. Compare to
[`locate-authority.md`](../procedures/locate-authority.md) step 5's "zero
match" case, which is a different (simpler) situation.
