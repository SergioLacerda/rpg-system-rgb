# Shields

Shields are **additional defensive equipment** that increase a character's protection in the RGB System.

Shields work together with **Armor** and the **Damage Model**, creating layered defenses.

There are two main types of shields:

- Physical Shields
- Energy / Magical Shields

## Physical Shields

Physical shields represent defensive equipment carried by a character.

They provide additional protection but may impose mobility penalties.

Physical shields may be **destroyed, dropped, or lost** depending on the Game Master's decision.

```text
Shield Type     Protection   Mobility Penalty
-------------   ----------   -------------Light Shield    +1           0
Medium Shield   +2           −1 G
Heavy Shield    +3           −2 G
```

### Characteristics

- Adds protection to defense
- May affect mobility (G)
- Subject to equipment loss or damage

Physical shields are most commonly used in:

- medieval settings
- tactical combat environments
- close‑quarters combat

## Energy / Magical Shields

Energy or magical shields represent **protective energy fields**, advanced technology, or magical barriers.

Unlike armor, these shields **absorb damage rather than reducing it**.

Characteristics:

- do not regenerate during combat
- have a fixed value during a fight
- return to maximum after the combat ends

Shield value is determined by the Blue vector.

```text
Shield = B × 3
```

Example:

```text
B = 3
Shield = 9
```

Energy shields represent the **absorption layer** in the RGB damage system.

## Interaction with the RGB Damage Model

Shields are applied **after armor reduction**.

Damage resolution order:

```text
Weapon Damage
↓
Penetration
↓
Armor Reduction
↓
Shield Absorption
↓
Remaining Damage → Character
```

This layered defense system keeps combat simple while preserving tactical depth.

## Interaction with RGB Vectors

Shields interact with the three RGB vectors in different ways.

```text
R → determines resilience and health
G → may be reduced by heavy shields
B → determines energy shield capacity
```

This maintains the RGB combat philosophy:

```text
R → deal damage
G → avoid damage
B → absorb damage
```

## Design Philosophy

Shields in the RGB system follow three principles:

- **layered defense** – armor reduces damage, shields absorb damage
- **tactical trade-offs** – heavier shields reduce mobility
- **modular design** – shields adapt to different settings

Examples:

- medieval campaigns → physical shields
- science fiction → energy shields
- fantasy settings → magical shields

## See Also

- [Armor](armor.md)
- [Damage Model](../combat/damage_model.md)
- [Attack and Defense](../combat/attack_and_defense.md)

← [Back to README](README.md)
