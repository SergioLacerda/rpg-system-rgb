# Armor

Armor represents **physical protection** used to reduce incoming damage in the RGB System.

Armor interacts with the **Damage Model** and may also combine with **energy shields** if the Blue (B) vector is active.

## Armor Calculation

Armor protection can be understood as:

```text
Total Defense = Armor (physical protection) + Magical / Energy Shield (B points)
```

**Note:**  
The magical or energy shield is only applied if the **Blue Factor (B)** is active in the campaign setting.

For example:

```text
Armor = 4
Shield = B × 3
```

The armor reduces damage per attack, while the shield absorbs accumulated damage.

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
R → determines physical resilience and health
G → affected by armor mobility penalties
B → may provide additional shield protection
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

For full rules see:

- damage_model.md

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
