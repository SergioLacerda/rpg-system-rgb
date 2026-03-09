# Firearms

Firearms represent the **primary ranged weapons** used in modern and tactical RGB settings.

Unlike melee weapons, firearms **do not add the user's R (Red) vector to damage**.  
Their damage is determined by the weapon itself, representing ballistic force and ammunition power.

Firearms interact with the **RGB Combat System** and follow the standard **Damage Model** used throughout the system.

## Firearm Statistics

Each firearm is defined by four attributes:

```text
Damage      → base damage inflicted by the weapon
Capacity    → number of rounds in the magazine
Reload      → actions required to reload
Type        → category of weapon
```

## Firearms Table

```text
Type                Damage   Capacity   Reload
------------------  -------  ---------  -------------Compact Pistol      4        12         partial action
Pistol              4        15         partial action
Revolver            5        6          main action
SMG                 5        30         main action
Assault Rifle       7        30         main action
Heavy Rifle / DMR   8        20         main action
Sniper Rifle        9        10         main action
Machine Gun         7        100        2 actions
```

## Weapon Roles

Different firearms fill different tactical roles.

### Pistols

- lightweight sidearms
- quick to reload
- suitable for close combat

### SMGs

- high rate of fire
- effective at short‑to‑medium range

### Assault Rifles

- balanced damage and capacity
- effective in most combat situations

### Heavy / Precision Rifles

- higher damage
- lower ammunition capacity
- specialized roles

### Machine Guns

- sustained suppressive fire
- large ammunition capacity
- slower reload time

## Interaction with the RGB System

Firearms interact with RGB vectors indirectly.

```text
R → improves offensive abilities or weapon handling
G → influences positioning and tactical movement
B → provides defensive protection through shields
```

However, firearm damage **comes from the weapon**, not the character's R value.

## Interaction with the Damage Model

When a firearm attack hits a target, the RGB damage flow applies:

```text
Weapon Damage
↓
Penetration (if applicable)
↓
Armor Reduction
↓
Shield Absorption
↓
Remaining Damage → Character
```

This keeps firearm combat consistent with the rest of the RGB system.

## Design Philosophy

Firearms in RGB follow three principles:

- **clarity** — simple weapon statistics
- **tactical variety** — different weapons fulfill different roles
- **system compatibility** — weapons integrate with the RGB damage model

This allows firearms to remain powerful while keeping combat resolution fast and consistent.

## See Also

- [Combat](../../combat/attack_and_defense.md)
- [Weapon Ranges](../mechanics/weapon_ranges.md)
- [Damage Model](../../combat/damage_model.md)

← [Back to README](../README.md)
