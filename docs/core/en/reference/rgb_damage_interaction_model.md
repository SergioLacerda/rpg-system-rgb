# RGB System – Damage Model and Interaction Design Notes

This document consolidates the design observations about the RGB System damage model,
its mathematical balance, and how the RGB vectors interact with combat mechanics.

The goal of this document is to explain **why the system behaves consistently and remains balanced**,
while keeping the rules simple.

## RGB Interaction Model

The RGB System is built around three fundamental vectors:

- **R (Red)** — pressure, impact, disruption and physical force
- **G (Green)** — relation, mobility, timing and reaction
- **B (Blue)** — preservation, shields, energy and special resistance

These vectors interact to create tactical gameplay.

```text
        G
   relation / reaction

R - B
pressure          preservation
```

Each vector influences a different aspect of gameplay:

| Vector | Gameplay Role |
|||
R | Pressure, impact and physical force |
G | Mobility, positioning and evasion |
B | Preservation, energy shields and special defenses |

## System Interaction Flow

The RGB System connects character attributes, equipment and combat mechanics.

```text
RGB Vectors
↓
Character Creation
↓
Equipment Selection
↓
Combat Interaction
↓
Damage Resolution
```

This modular structure allows the system to adapt to multiple settings such as:

- modern campaigns
- fantasy worlds
- science fiction settings
- super‑powered universes

## RGB Damage Model

Combat damage in the RGB System follows a layered structure.

```text
Hit or Contact Check
↓
Impact Source
↓
Penetration
↓
Armor Reduction
↓
Shield Absorption
↓
Remaining Damage → Character
```

This layered approach ensures that weapons, armor and shields interact predictably.

## Core System Formulas

The RGB system defines character durability using simple formulas.

```text
Health = 4 + R + B
Shield = B × 3
```

Weapons, attributes, abilities or procedures declare the **impact source**,
while armor provides **fixed reduction values**.

Penetration modifies the effective armor value.

## Effective Armor

Penetration interacts with armor through a simple relationship:

```text
Effective Armor = Armor − Penetration
```

Final damage becomes:

```text
Damage After Armor = Impact Source − Effective Armor
```

This keeps calculations simple and predictable.

## Natural Damage Scale

If the following ranges are used:

| Element | Typical Scale |
|--||
Weapons | 3–10 |
Armor | 1–6 |
Penetration | 0–4 |

A natural balance relationship appears:

```text
Impact Source ≈ Armor + Penetration
```

This ensures that weapons and defenses remain balanced.

## Example

### Light Weapon

```text
Damage = 4
Penetration = 1
Armor = 4
```

```text
Effective Armor = 4 − 1 = 3
Final Damage = 4 − 3 = 1
```

### Heavy Weapon

```text
Damage = 8
Penetration = 3
Armor = 6
```

```text
Effective Armor = 6 − 3 = 3
Final Damage = 8 − 3 = 5
```

## Shield Interaction

After armor reduction, shields absorb remaining damage.

```text
Remaining Damage − Shield
```

Shields follow the RGB formula:

```text
Shield = B × 3
```

This creates two defensive layers:

```text
| Defense Type | Function |
|--|-|
Armor | Reduces damage per attack |
Shield | Absorbs accumulated damage |
```

This separation simplifies combat resolution.

## Balance Rule of Thumb

When designing weapons and armor, the following approximation keeps combat balanced:

```text
Armor ≈ Impact Source / 2
Penetration ≈ Impact Source / 3
```

This keeps most attacks meaningful while preserving defensive value.

## Unexpected Advantage: Multiple Enemy Combat

One interesting property of the RGB damage model is that it performs well when characters face **multiple opponents simultaneously**.

Many RPG systems struggle with this scenario because defensive values scale poorly against repeated attacks.

In RGB, the layered model distributes damage naturally:

```text
Impact Source → Penetration → Armor → Shield → Character
```

### Armor

Armor reduces damage **for each individual attack**, preventing weak enemies
from overwhelming armored characters too quickly.

### Shield

Shields absorb **total accumulated damage**, acting as a buffer against multiple small attacks.

## Tactical Consequences

This creates a stable combat dynamic.

```text
| Enemy Type | Result |
||--|
Many weak enemies | Armor reduces most attacks |
Few strong enemies | Penetration becomes important |
Energy-heavy enemies | Shields become critical |
```

Because armor works per hit and shields work cumulatively,
the system remains stable even with many attackers.

## Design Conclusion

The RGB damage structure forms a clear defensive hierarchy:

```text
Weapon
↓
Penetration
↓
Armor
↓
Shield
↓
Character
```

This structure:

- keeps calculations simple
- supports tactical combat
- scales naturally with equipment
- handles multi‑enemy encounters well

The interaction between **R, G and B vectors** ensures that different character builds
remain viable while encouraging tactical choices during combat.

← [Back to README](README.md)
