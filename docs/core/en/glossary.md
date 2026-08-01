# Glossary

Terms used throughout the RGB System rules, gathered in one place. This page
is not normative on its own — each entry links to the rule document that
actually defines the term; the glossary only points there.

## R, G, B

The three vectors the system is built around. See
[Attributes](core/attributes.md).

- **R (Red)** — pressure and physical impact: attack, force, disruption.
- **G (Green)** — relation, timing, and positioning: movement, evasion,
  reach.
- **B (Blue)** — preservation under pressure: endurance, blocking, shields,
  stabilization.

## Margin

The difference between an attacker's effective `R` and a defender's
effective `G` in an attack, used to classify the outcome (strong hit, hit,
hit with cost, miss with opportunity, clear miss) instead of a flat
hit/miss. See [Attack and Defense](combat/attack_and_defense.md).

## Health

A character's vitality resource, `Health = 4 + R + B`. See
[Character Creation](core/character_creation.md).

## Impact

The pre-mitigation force of a hit, before penetration, armor reduction, or
shield absorption are applied. See [Damage Model](combat/damage_model.md).

## Penetration

How much of a weapon's impact bypasses armor before armor reduction is
applied. See [Penetration](weapons/mechanics/penetration.md).

## Armor Reduction

The mitigation layer that reduces impact per hit, applied after
penetration. See [Armor](equipment/armor.md).

## Shield Absorption

The mitigation layer that absorbs remaining impact after armor reduction,
drawn from a `B`-derived resource. See [Shields](equipment/shields.md).

## State

A tagged condition a character can carry (e.g. `injured`, `suppressed`,
`hidden`), created, changed, or cleared by specific procedures. Fifteen
states exist across physical, positional, tactical, defensive, and
informational categories.

## Procedure

A named defensive or tactical action with a stable identity and vector
ownership: Evade, Reposition, Block, Interrupt, Counterpressure. See
[Attack and Defense](combat/attack_and_defense.md).

## Ability

A discrete character capability declared with a minimal governed contract
(id, vector, tier, requirements, cost, range, duration, effect, limits,
tags, source status). See
[Skills and Abilities](core/skills_and_abilities.md).

## Tier

An ability's power/complexity band, used alongside `requirements` to gate
when a character can take it. See
[Skills and Abilities](core/skills_and_abilities.md).
