# Combat

## Initiative

More Green (G) points → acts first.  
A surprise attack ignores initiative.

## Attack

Offensive stance. The attacker uses Red (R) points.

An attack hits when:

```text
Attacker R ≥ Defender G
```

Resolution:

```text
If Attacker (R + modifiers) ≥ Target (G + modifiers) → hit
Otherwise → dodge
```

## Defense

Maximizes evasion. The defender uses Green (G) points.

```text
Defense = Armor + Shield
```

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

FinalDamage = Attacker (R + modifiers) – Target (Armor + modifiers)  
RemainingHealth = CurrentHealth – FinalDamage

## Damage Resolution

1. check hit
2. calculate damage
3. apply penetration
4. reduce armor
5. reduce shield
6. apply to health

## Meta Abilities or Magic

### Magical Attack during combat

```text
Each turn points are recovered to use spells or meta‑abilities.

Regeneration = 1 B point per level
```

### Magical Defense during combat

```text
Defense = Armor + Physical Shield (if any) + Magical Shield
```

See also:

- [Firearms](../weapons/firearms.md)
- [Armor](../equipment/armor.md)

← [Back to README](README.md)
