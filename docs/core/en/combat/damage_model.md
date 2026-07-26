# Damage Model

The RGB System uses a **layered damage model** that combines impact sources,
penetration, armor, shields, and special abilities.
This model keeps combat simple while still allowing tactical decisions.

This document explains how damage flows through the system.

## Damage Resolution Flow

Damage in the RGB System follows a clear sequence:

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

This flow shows how different elements of the system interact during combat.

## Step-by-Step Damage Resolution

### 1. Attack Check

The attacker attempts to hit the defender.

Combat rules determine if the attack succeeds.

See:

- ../combat/attack_and_defense.md

If the attack fails, no damage is applied.

### 2. Impact Source

If the attack succeeds, identify the **impact source**.

Impact may come from a weapon, attribute, ability, procedure, or explicit
exception. The default is:

- melee impact uses weapon value plus the attacker's `R` when the weapon or
  procedure permits it;
- firearm impact uses weapon value by default;
- explosive impact uses the explosive profile;
- ability impact uses the ability's declared effect.

Examples cannot override the declared impact source.

See:

- ../weapons/categories/firearms.md
- ../weapons/categories/melee.md
- ../weapons/categories/explosives.md

### 3. Penetration

Apply penetration to armor before armor reduces damage.

```text
Effective Armor = Armor - Penetration
```

If effective armor is zero or less, armor does not reduce that hit.

See:

- ../weapons/mechanics/penetration.md

### 4. Armor Reduction

Armor provides **physical protection**.

Armor reduces incoming damage per hit after penetration.

Typical armor categories include:

- Light Armor
- Medium Armor
- Heavy Armor

See:

- ../equipment/armor.md

### 5. Shield Absorption

Some characters may possess shields.

Shields absorb remaining damage **after armor reduction** and before it reaches
health or a state consequence.

Energy shields are typically calculated as:

Shield = B × 3

See:

- ../equipment/shields.md

### 6. Remaining Damage

After armor and shields are resolved, the remaining damage is applied to the character.

This damage affects the character's health or creates a declared state
consequence.

## Special Abilities

Some abilities may modify how damage is applied.

Examples include:

- absorption techniques
- energy shields
- defensive martial techniques

See:

- ../core/skills_and_abilities.md

## Damage Interaction Overview

The RGB damage model can be summarized as:

```text
Impact Source
      ↓
Penetration
      ↓
Armor Reduction
      ↓
Shield Absorption
      ↓
Remaining Damage
```

This layered model keeps combat easy to understand while allowing
different equipment and abilities to interact in meaningful ways.

## RGB Interaction Model

The system can also be visualized as an interaction between the RGB vectors:

```text
        G
   mobility / reaction

R - B
power / damage    shield / energy
```

Each vector represents a different defensive or offensive strategy in combat.

- **R (Red)** changes the source of pressure through force, impact, or
  interruption.
- **G (Green)** changes relation to pressure through movement, timing, and
  positioning.
- **B (Blue)** preserves continuity through blocks, shields, stabilization, and
  resistance.

Together they create a balanced tactical system.

← [Back to README](README.md)
