
# Combat Decision Model

The RGB System encourages players to make tactical decisions during combat.

Each RGB vector represents a different strategic approach.

## RGB Tactical Triangle

Combat decisions can be visualized using the RGB triangle.

```text
        G
   mobility / positioning

R ---------------- B
power / damage    shield / defense
```

Each vector represents a different combat strategy.

## Combat Strategies

## R – Offensive Strategy

Focus on direct damage and aggressive combat.

Typical actions:

- powerful attacks
- heavy weapons
- breaking enemy defenses

Advantages:

- high damage output
- faster enemy elimination

Weakness:

- lower survivability

## G – Mobility Strategy

Focus on movement and positioning.

Typical actions:

- flanking enemies
- avoiding attacks
- controlling distance

Advantages:

- harder to hit
- tactical positioning

Weakness:

- lower direct damage

## B – Defensive Strategy

Focus on survivability and protection.

Typical actions:

- shield usage
- defensive abilities
- absorbing enemy attacks

Advantages:

- high durability
- strong defensive capacity

Weakness:

- slower combat resolution

## Tactical Interactions

The RGB system allows multiple combat styles.

```text
| Strategy | Vector Focus |
|----------|--------------|
Aggressive fighter | High R |
Mobile skirmisher | High G |
Defensive guardian | High B |
```

Hybrid builds are also possible.

Examples:

```text
| Build | Style |
|------|------|
R + G | fast striker |
R + B | heavy warrior |
G + B | evasive defender |
```

## Combat Decision Flow

During combat a player typically chooses between three tactical options:

```text
Attack
Move
Defend
```

These actions map naturally to the RGB vectors.

```text
Attack → R
Move → G
Defend → B
```

## Interaction with Damage Model

Combat decisions interact directly with the damage system.

```text
Attack
↓
Weapon Damage
↓
Penetration
↓
Armor
↓
Shield
↓
Character
```

Players must choose whether to:

- increase damage
- avoid damage
- absorb damage

## Design Philosophy

The RGB System creates a simple tactical model:

```text
R → Deal damage
G → Avoid damage
B → Absorb damage
```

These three strategies form the core tactical decisions in RGB combat.

← [Back to README](README.md)
