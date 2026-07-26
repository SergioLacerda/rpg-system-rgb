# Melee Weapons

Melee weapons represent **weapons used in close combat** in the RGB System.

Unlike firearms, melee weapons **use the character's R (Red) vector as the base for damage**.

This represents physical strength, impact force, and martial skill.

## Damage Calculation

The damage of a melee weapon is calculated by adding the character's **R value** to the weapon bonus.

```text
Damage = R + weapon bonus
```

Example:

```text
R = 3
Long Sword = +3

Total Damage = 6
```

## Melee Weapons Table

```text
Weapon          Damage Bonus
--------------  ----------Knife           +1
Dagger          +1
Short Sword     +2
Long Sword      +3
Axe             +3
```

## Interaction with the RGB System

Melee weapons are directly linked to the **R (Red) vector**.

```text
R → determines base damage
G → influences mobility and positioning
B → provides additional defense through shields
```

Characters with high **R** values tend to be more effective in melee combat.

## Interaction with the Damage Model

When a melee attack hits a target, the RGB damage resolution flow applies.

```text
Damage (R + weapon)
↓
Penetration (if applicable)
↓
Armor Reduction
↓
Shield Absorption
↓
Remaining Damage → Character
```

## Tactical Use

Melee weapons are particularly effective in:

- close combat situations
- confined spaces
- environments where firearms cannot be used

They also allow silent attacks and ambush tactics.

## Design Philosophy

Melee weapons in RGB follow three principles:

- **simplicity** — direct damage calculation based on R
- **balance** — damage proportional to the character's strength
- **system compatibility** — direct integration with the RGB damage model

## See Also

- [Combat](../../combat/attack_and_defense.md)
- [Damage Model](../../combat/damage_model.md)
- [Firearms](firearms.md)

← [Back to README](../README.md)
