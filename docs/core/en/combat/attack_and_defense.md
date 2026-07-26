# Attack And Defense

## Initiative

More Green (G) points → acts first.  
A surprise attack ignores initiative.

## Resolution Philosophy

Core V2 resolves attacks with a **margin**. This keeps deterministic play
possible while allowing richer outcomes than a flat hit or miss.

```text
Attack Margin = Attacker (R + modifiers) - Defender (G + modifiers)
```

Use the margin as follows:

| Margin | Result |
| ---: | --- |
| 3 or more | strong hit |
| 1 to 2 | hit |
| 0 | hit with cost or exposed follow-up |
| -1 to -2 | miss with opportunity |
| -3 or less | clear miss |

The older teaching shorthand remains valid:

```text
Attacker R ≥ Defender G
```

That shorthand means the attacker gets at least a basic hit when the margin is
zero or higher.

## Attack

Offensive pressure. The attacker uses Red (R) to change the target or the
source of danger.

```text
margin = Attacker (R + modifiers) - Target (G + modifiers)
```

If the result is a hit, continue to the damage model or to the declared
non-damage consequence.

## Defensive Procedures

Defense is no longer one combined number. Choose the procedure that describes
what the defender is doing.

| Procedure | Primary Owner | Purpose |
| --- | --- | --- |
| Evade | G | avoid or alter contact |
| Reposition | G | change range, cover, or engagement |
| Block | B or equipment | intentionally receive pressure |
| Armor reduction | equipment | reduce impact per hit |
| Shield absorption | B-derived resource or equipment | absorb remaining damage |
| Protect ally | B, equipment, or ability | redirect or contain pressure |
| Interrupt | R | stop an action by applying pressure first |
| Counter | R or G by procedure | answer pressure through force or timing |

Do not collapse these procedures into a generic `Defense = Armor + Shield`
formula. Armor and shield act later in the damage pipeline.

## Evade

```text
Evade uses G against the attacker's R.
```

Evade changes contact, timing, range, or position. If evasion succeeds, no
damage is applied unless an ability or area effect says otherwise.

## Block

Block uses B, equipment, or an ability to receive pressure intentionally. A
block may reduce the attack margin, protect another target, or send the attack
into armor and shield layers.

## Grapple

Allows a character to immobilize the opponent or force movement during a turn.

```text
A grapple action succeeds when:

Attacker (G) ≥ Target (G)
Attacker (R) ≥ Target (R)

The attacker may use meta‑abilities.
The target may use pre-declared bonuses.
```

## Damage

If an attack deals damage, use the damage model.

## Damage Resolution

1. check hit
2. establish impact source
3. apply penetration
4. apply armor reduction
5. apply shield absorption
6. apply remaining health or state consequence

## Meta Abilities or Magic

### Magical Attack during combat

```text
Each turn points are recovered to use spells or meta‑abilities.

Regeneration = 1 B point per level
```

### Magical Defense During Combat

Magical defense is a typed shield, block, resistance, or absorption effect. It
must state which defensive procedure it modifies and whether it acts before or
after armor.

See also:

- [Firearms](../weapons/categories/firearms.md)
- [Armor](../equipment/armor.md)

← [Back to README](README.md)
