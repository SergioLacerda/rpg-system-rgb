# Armor

Armor represents **physical protection** used to reduce incoming damage in the RGB System.

Armor interacts with the **Damage Model** as the per-hit reduction layer.
Energy shields and other Blue (B) effects are separate absorption or
preservation layers.

## Armor Calculation

Armor protection is a fixed reduction value:

```text
Effective Armor = Armor - Penetration
```

For example:

```text
Armor = 4
Penetration = 1
Effective Armor = 3
```

Armor reduces damage per attack. Shields absorb remaining damage after armor.

## Armor Types

```text
Type      Protection   Mobility Penalty
--------  ----------   -------------Light     2            −1 G
Medium    4            −2 G
Heavy     6            −3 G
```

### Light Armor

- Minimal protection
- Low mobility penalty
- Suitable for agile characters

### Medium Armor

- Balanced protection
- Moderate mobility reduction

### Heavy Armor

- High physical protection
- Significant mobility penalty

Heavy armor users rely more on **R (strength)** and **B (defense)** than on mobility.

## Armor Interaction with the RGB System

Armor interacts with the RGB vectors in the following way:

```text
R → contributes pressure and physical impact
G → affected by armor mobility penalties
B → may provide separate preservation or shield effects
```

This relationship maintains the RGB tactical balance:

```text
R → deal damage
G → avoid damage
B → absorb damage
```

## Interaction with the Damage Engine

Armor reduces damage **after weapon penetration is applied**.

The simplified damage flow:

```text
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

For full rules see:

- [Damage Model](../combat/damage_model.md)

## Design Philosophy

Armor in the RGB system follows three principles:

- **Simple values** – easy to calculate during combat
- **Tactical trade-offs** – higher protection reduces mobility
- **Compatibility with shields** – allows layered defenses

This ensures armor remains relevant without overpowering mobility or defensive builds.

## See Also

- [Shields](shields.md)
- [Damage Model](../combat/damage_model.md)

← [Back to README](README.md)
