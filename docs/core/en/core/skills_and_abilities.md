# Skills and Abilities

Skills represent special capabilities beyond the normal actions of the RGB system.
They allow characters to expand their abilities depending on the campaign setting.

Skills are **optional modules** and may represent:

- martial techniques
- supernatural powers
- advanced technology
- magic systems

The Game Master determines whether these abilities are available in the campaign.

## RGB Ability Structure

Abilities in the RGB system are organized around the **three core vectors**.

```text
        G
   mobility / positioning

R ---------------- B
power / damage      shield / defense
```

Each vector naturally forms a **branch of abilities**.

- **R (Red)** → pressure, impact, interruption, and transformation abilities
- **G (Green)** → movement, timing, evasion, range, and positioning abilities
- **B (Blue)** → preservation, stabilization, block, shield, and resistance
  abilities

This structure allows the system to scale naturally with different settings.

## Red Ability Tree (R)

Red abilities enhance **offensive power and physical force**.

Typical abilities:

- Double Power
- Delayed Strike
- Air Pressure Attack
- Power Strike
- Armor Break

These abilities increase a character's ability to deal damage.

Example progression:

```text
Level 3 → +1 R to ranged attacks (Precision)
Level 6 → Double power strike
Level 9 → Area impact attack
```

## Green Ability Tree (G)

Green abilities focus on **mobility, speed, and tactical positioning**.

Typical abilities:

- Enhanced Reflexes
- Momentary Acceleration
- Phantom Step (short teleport-like movement)
- Automatic Dodge
- Supernatural Sprint
- Enhanced Jump
- Silent Movement
- Movement Distortion
- Supernatural Evasion

These abilities allow characters to control space and avoid damage.

## Blue Ability Tree (B)

Blue abilities represent **energy manipulation, endurance, and defensive systems**.

Typical abilities:

- Energy Shield
- Stunning Shield
- Intelligent Shield / Living Armor
- Invisibility
- Clones
- Energy Absorption

These abilities improve survivability, continuity, and energy-based defenses.

## Ability Contract

Every ability must declare the following fields before it can be treated as
canonical, validated, bundled, or used by a Specialist:

| Field | Purpose |
| --- | --- |
| `id` | stable identifier |
| `name` | table-facing name |
| `vector` | primary RGB vector |
| `tier` | level, rank, or access tier |
| `requirements` | prerequisites |
| `action_type` | action, reaction, passive, stance, or special timing |
| `cost` | vector points, shield, callback, item charge, or other cost |
| `range` | target or area reach |
| `duration` | instant, round, scene, persistent, or conditional |
| `effect` | mechanical effect |
| `limits` | per turn, per combat, cooldown, or narrative restriction |
| `tags` | classification tags |
| `source_status` | proposed, optional, canonical, deprecated, or example |

Abilities that lack these fields are examples or design notes, not canonical
rules.

## Martial Factor (Optional Module)

The **Martial Factor** introduces advanced defensive techniques.

Examples:

- Complete Defense
- Special Absorption
- Defensive stance techniques

The Game Master decides whether these techniques require:

- an action
- a reaction
- a defensive stance

## Blue Factor (Optional Module)

In campaigns that allow **magic or meta abilities**, the Blue vector can manipulate
energy-based defenses.

These abilities usually interact with the shield system:

```text
Shield = B × 3
```

Examples:

- magical shields
- energy barriers
- energy absorption
- advanced defensive fields

## Callback System (Optional)

If the campaign includes powerful abilities, the Game Master may apply a **callback system**
to represent the cost or fatigue caused by these abilities.

If callback exceeds **100%**, the character faints.

```text
| Intensity | Recovery |
|-----------|----------|
| Light     | hours |
| Moderate  | days |
| Severe    | weeks |
| Extreme   | months |
```

The Game Master may also apply permanent consequences depending on the situation.

## Special Absorption

Special Absorption allows characters to reduce incoming damage using a declared
resource or vector expression.

```text
1 vector point spent → reduces 1 damage
```

Possible uses:

- **Red (R)** → physical impact absorption
- **Blue (B)** → energy or special damage absorption

The Game Master may adjust this ratio depending on the campaign.

## Complete Defense

Complete Defense is a legacy name for a block or guard stance where the
character prepares to receive pressure rather than evade it.

This stance also negates grapple attempts.

```text
Complete Defense = Block procedure + declared absorption effect
```

## Example Ability Categories

### Mobility

- flight
- levitation
- teleportation

### Perception

- night vision
- advanced sensors
- enhanced awareness

### Technology

- remote device control
- hacking abilities
- drone interface

## Creating RGB Abilities

Every ability should use the contract above. At minimum, it must define vector,
cost, action type, effect, duration, limits, and source status.

This keeps abilities consistent with the RGB system.

See also:

- [Character Creation](character_creation.md)
- [Attributes](attributes.md)

← [Back to README](README.md)
