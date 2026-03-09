# Damage Model

The RGB System uses a **layered damage model** that combines weapons, armor, shields and special abilities.
This model keeps combat simple while still allowing tactical decisions.

This document explains how damage flows through the system.

## Damage Resolution Flow

Damage in the RGB System follows a clear sequence:

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

This flow shows how different elements of the system interact during combat.

## Step-by-Step Damage Resolution

### 1. Attack Check

The attacker attempts to hit the defender.

Combat rules determine if the attack succeeds.

See:

- ../combat/attack_and_defense.md

If the attack fails, no damage is applied.

### 2. Weapon Damage

If the attack succeeds, the weapon determines the **base damage value**.

Weapons define how much damage is applied before defenses.

See:

- ../weapons/firearms.md
- ../weapons/melee.md
- ../weapons/explosives.md

### 3. Armor Reduction

Armor provides **physical protection**.

Armor reduces incoming damage depending on the armor type.

Typical armor categories include:

- Light Armor
- Medium Armor
- Heavy Armor

See:

- ../equipment/armor.md

### 4. Shield Absorption

Some characters may possess shields.

Shields absorb damage **before it reaches the character**.

Energy shields are typically calculated as:

Shield = B × 3

See:

- ../equipment/shields.md

### 5. Remaining Damage

After armor and shields are resolved, the remaining damage is applied to the character.

This damage affects the character's vitality.

## Penetration

Certain weapons may partially ignore armor protection.

Penetration represents the ability of a weapon to bypass armor defenses.

See:

- ../weapons/penetration.md

Penetration is applied **before armor reduction**.

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
Weapon Damage
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

- **R (Red)** influences damage and physical power
- **G (Green)** influences movement and reaction
- **B (Blue)** influences shields and special defenses

Together they create a balanced tactical system.

← [Back to README](README.md)
