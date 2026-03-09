# RGB System – Quick Start

This quick start guide explains the core ideas of the RGB System and how to begin playing in just a few minutes.

## 1. The RGB Vectors

Characters in the RGB System are defined by three vectors:

```text
| Vector | Name | Function |
|||--|
| R | Red | damage, strength, offensive power |
| G | Green | agility, movement, reaction |
| B | Blue | shields, energy, special defense |
```

These three values define how a character performs in combat and interaction.

## 2. Creating a Character

A starting character typically receives:

```text
7 points
```

Players distribute these points between **R, G and B**.

Example:

```text
R = 3
G = 2
B = 2
```

## 3. Character Durability

Durability is calculated using simple formulas.

```text
Health = 4 + (R × 2)
Shield = B × 3
```

Example:

```text
R = 3 → Health = 10
B = 2 → Shield = 6
```

## 4. Basic Turn Structure

During combat a character can typically perform:

```text
1 main action
+
1 minor adjustment (optional)
```

Examples:

Main actions:

- attack
- move
- defend
- use ability

Minor adjustments:

- reload weapon
- change position slightly
- interact with environment

## 5. Combat Decisions

Combat revolves around three tactical choices:

```text
Attack
Move
Defend
```

These map directly to the RGB vectors:

```text
Attack → R
Move → G
Defend → B
```

Players naturally emphasize different strategies depending on their vector distribution.

## 6. Damage Resolution

When an attack hits, damage follows this sequence:

```text
Attack Check
↓
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

Armor reduces damage **per attack**, while shields absorb **accumulated damage**.

## 7. Example Combat Situation

A character with:

```text
R = 3
G = 2
B = 2
```

faces an enemy.

The player must choose between:

- increasing damage output
- repositioning
- reinforcing defense

Each decision corresponds to one of the RGB vectors.

## 8. Gameplay Loop

A typical RGB session follows this structure:

```text
Character Creation
↓
Exploration
↓
Combat Encounter
↓
Recovery
↓
Story Progression
```

This loop repeats as the story develops.

## Learn More

For deeper understanding see:

- Core Rules → [Core Rules](../core/README.md)
- Combat System → [Combat System](../combat/README.md)
- Equipment → [Equipment](../equipment/README.md)
- Weapons → [Weapons](../weapons/README.md)
- Gameplay Loop → [Gameplay Loop](../reference/gameplay_loop.md)
