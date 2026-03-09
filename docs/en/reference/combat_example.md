# Combat Example

This example demonstrates how combat is resolved using the **RGB System**.

Two characters engage in combat and resolve actions using the core rules:

- RGB vectors
- movement
- attack and defense
- armor and shields
- damage calculation

## Characters

## Character A – Assault Fighter

```text
R = 3
G = 2
B = 1
```

Derived values:

```text
Health = 4 + (R × 2) = 10
Shield = B × 3 = 3
Movement = G × 2 = 4 meters
```

Equipment:

- Medium Armor (Protection 4)
- Rifle

## Character B – Mobile Operator

```text
R = 2
G = 3
B = 2
```

Derived values:

```text
Health = 8
Shield = 6
Movement = 6 meters
```

Equipment:

- Light Armor (Protection 2)
- SMG

## Combat Turn Example

## Step 1 — Positioning

Characters move according to:

```text
Movement = G × 2 meters
```

Character B moves first due to higher mobility.

## Step 2 — Attack

Character A fires the rifle.

Example weapon damage:

```text
Weapon Damage = 6
Penetration = 2
```

## Step 3 — Armor Interaction

Armor reduces incoming damage after penetration.

Example:

```text
Target Armor = 2
Penetration = 2
Effective Armor = 0
```

Armor is bypassed.

## Step 4 — Shield Absorption

Shield absorbs damage before health.

```text
Shield = 6
Damage = 6
```

Result:

```text
Remaining Shield = 0
Health Damage = 0
```

## Second Attack

Character B attacks with an SMG.

```text
Weapon Damage = 4
Penetration = 1
```

Character A armor:

```text
Armor = 4
Effective Armor = 3
```

Damage calculation:

```text
Damage = 4 - 3 = 1
```

Character A receives:

```text
Health Damage = 1
Remaining Health = 9
```

## Damage Flow Summary

The RGB system resolves damage using a layered structure.

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

## Tactical Interpretation

The example shows the role of the RGB vectors:

```text
R → determines offensive power
G → determines mobility and positioning
B → determines shield capacity
```

Each vector supports a different combat strategy.

## Design Goal

Combat resolution in the RGB system is designed to be:

- fast to calculate
- tactically meaningful
- consistent across settings

Because the same formulas apply to all characters, the system remains easy to balance.

## See Also

- [Attack and Defense](attack_and_defense.md)
- [Movement](movement.md)
- [Damage Model](damage_model.md)

← [Back to README](README.md)
