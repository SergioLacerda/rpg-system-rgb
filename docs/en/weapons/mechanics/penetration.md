# Penetration

Penetration represents a weapon's ability to **overcome armor protection** in the RGB System.

Weapons with higher ballistic power or physical impact have greater ability to penetrate armor.

## Penetration Calculation

Penetration temporarily reduces the target's armor value during damage calculation.

```text
Effective Armor = Armor − Penetration
```

Example:

```text
Armor = 4
Weapon Penetration = 2

Effective Armor = 2
```

If the effective armor value reaches **0 or less**, the armor does not reduce damage.

## Weapon Penetration Table

```text
Weapon            Penetration
----------------  --------Pistol            1
SMG               2
Carbine           2
Assault Rifle     3
Sniper Rifle      4
Machine Gun       3
```

## Interaction with the Damage Model

Penetration is applied **before armor damage reduction**.

Full damage resolution flow:

```text
Weapon Damage
↓
Penetration
↓
Effective Armor
↓
Shield Absorption
↓
Remaining Damage → Character
```

## Interaction with the RGB System

Penetration is primarily related to weapon impact, but it indirectly interacts with RGB vectors.

```text
R → increases damage for melee weapons
G → allows characters to avoid attacks through mobility
B → provides additional protection through shields
```

Weapons with higher penetration are more effective against heavily armored targets.

## Design Philosophy

The penetration system in RGB follows three principles:

- **simplicity** — direct and easy-to-apply calculation
- **balance** — different weapons have different effectiveness against armor
- **compatibility** — direct integration with the RGB damage model

This allows weapons to be differentiated without adding excessive complexity to combat.

## See Also

- [Firearms](../categories/firearms.md)
- [Damage Model](../../combat/damage_model.md)
- [Armor](../../equipment/armor.md)

← [Back to Readme](../README.md)
