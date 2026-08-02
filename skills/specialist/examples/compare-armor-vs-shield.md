# Example: Compare `core.damage.armor-reduction` vs `core.damage.shield-absorption`

**Procedure:** [`../procedures/compare-rules.md`](../procedures/compare-rules.md)

## Comparison

```text
A: core.damage.armor-reduction (rule) — Armor reduces incoming damage per
   hit, after penetration and before shield absorption.
B: core.damage.shield-absorption (rule) — Shields absorb remaining damage
   after armor reduction and before it reaches health or a state
   consequence.
Relationship: sequence — both are ordered_steps of core.damage.flow
   (procedure), with A immediately preceding B.
Cited: core.damage.armor-reduction, core.damage.shield-absorption
   (docs/core/en/combat/damage_model.md)
```

The two rules are not alternatives or overlapping mitigation — they are
sequential layers on the same incoming damage. Armor reduces the impact
value first; only the amount that survives armor reaches shield. Neither
unit's `relationships` declares `depends_on` on the other directly — the
sequencing authority is `core.damage.flow`, not a link between the two
rule units themselves (see
[`../procedures/locate-authority.md`](../procedures/locate-authority.md)'s
worked example for why `core.damage.flow` is cited for ordering
questions).
