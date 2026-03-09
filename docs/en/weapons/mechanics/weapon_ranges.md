# Weapon Ranges

Weapon range defines **the effectiveness of ranged weapons at different combat distances** in the RGB System.

These rules represent the loss of precision and effectiveness as the target moves farther away.

## Distance Modifiers

Attacks with ranged weapons receive modifiers depending on the distance to the target.

```text
Distance   Modifier
---------  --------Short      0
Medium     −1 R
Long       −2 R
```

These modifiers represent the increasing difficulty of hitting distant targets.

## Range Table (meters)

```text
Weapon             Short   Medium   Long
-----------------  -----   ------   --Compact Pistol     10      20       35
Pistol             15      30       50
Revolver           20      40       60
SMG                20      50       80
Assault Rifle      40      120      250
Heavy Rifle        60      200      400
Sniper Rifle       150     500      1000
Machine Gun        60      200      400
```

## Extreme Range (Optional)

Attacks beyond long range may still be attempted, but they are significantly more difficult.

```text
Extreme Range → −4 R
```

The Game Master decides whether the target is still visible and whether the shot is plausible.

## Interaction with the RGB System

Range interacts mainly with the **R (Red) vector**, which represents offensive precision.

```text
R → accuracy and effectiveness of attacks
G → positioning to obtain better line of sight
B → defensive protection against attacks
```

Characters with higher **R** values tend to handle distance penalties more effectively.

## Interaction with the Combat System

Range modifiers are applied during the attack step.

Simplified resolution flow:

```text
Attack
↓
Apply distance modifier
↓
Resolve defense
↓
Apply damage model
```

This keeps the system consistent with:

- attack_and_defense.md
- damage_model.md

## Design Philosophy

The RGB range system follows three principles:

- **simplicity** — few modifiers that are easy to remember
- **consistency** — the same model for all ranged weapons
- **tactical balance** — different weapons have different effective ranges

This ensures that positioning has real impact during combat.

## See Also

- [Firearms](firearms.md)
- [Penetration](penetration.md)
- [Damage Model](../combat/damage_model.md)

← [Back to Firearms](firearms.md)
