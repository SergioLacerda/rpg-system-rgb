# RGB System Engine

The **RGB System Engine** describes the full architecture of the RGB role‑playing system.
It connects all system components into a layered structure similar to a game engine.

This model helps explain how the RGB system remains:

- modular
- scalable
- easy to balance
- adaptable to many settings

## RGB System Architecture

The RGB system can be understood as a layered engine:

```text
           Campaign Setting
                │
                ▼
        RGB Ability Engine
                │
                ▼
          Gameplay Layer
                │
                ▼
           Combat System
                │
                ▼
            Damage Engine
                │
                ▼
           Equipment Data
                │
                ▼
            RGB Vectors
```

Each layer builds on the one below it.

## Layer Description

## RGB Vectors (Core Layer)

The entire system is built on the three RGB vectors.

```text
        G
   mobility / positioning

R ---------------- B
power / damage      shield / defense
```

Vectors define the core mechanics:

- **R (Red)** → offensive power
- **G (Green)** → mobility and reaction
- **B (Blue)** → defense and energy

All other systems derive from these three values.

## Equipment Data

This layer defines numerical values used by the combat system.

Examples:

- weapon damage
- armor reduction
- penetration values
- shield capacity

Documents:

- firearms.md
- melee.md
- armor.md
- shields.md

## Damage Engine

The damage engine determines how attacks affect characters.

```text
Weapon
↓
Penetration
↓
Armor
↓
Shield
↓
Character
```

Documents:

- damage_model.md
- rgb_damage_interaction_model.md

## Combat System

The combat layer defines interaction between characters.

Examples:

- attack resolution
- movement
- positioning
- defense reactions

Documents:

- movement.md
- attack_and_defense.md

## Gameplay Layer

Defines how players interact with the system during play.

Typical gameplay loop:

```text
Character Creation
↓
Exploration
↓
Combat
↓
Recovery
```

Documents:

- character_creation.md
- gameplay_loop.md
- combat_decision_model.md

## RGB Ability Engine

This layer introduces abilities, powers, and special mechanics.

Abilities are built using the RGB vectors.

```text
Vector
Cost
Effect
Duration
Limit
```

Documents:

- skills_and_abilities.md
- rgb_ability_engine.md

## Campaign Setting

The top layer defines the narrative world.

Examples:

- modern campaigns
- fantasy worlds
- science fiction
- superhero universes

Example:

```text
RGB System → rule engine
Aurora → campaign setting built on top of the engine
```

## Complete System Flow

The full interaction between layers:

```text
Vectors
 ↓
Equipment Data
 ↓
Damage Engine
 ↓
Combat System
 ↓
Gameplay Rules
 ↓
Ability Engine
 ↓
Campaign Setting
```

## Design Advantages

The RGB System Engine provides several advantages.

## Modular Design

New subsystems can be added without rewriting core rules.

Examples:

- magic systems
- vehicles
- cybernetics
- superpowers

## Easier Balance

Balance occurs primarily in:

```text
Damage Engine
Equipment Data
```

The rest of the system remains stable.

## Clear Documentation

Each layer answers a specific question.

```text
Layer             Question
---------------------------------------------Vectors           What defines a character?
Equipment         What numbers exist?
Damage Engine     How damage works?
Combat System     How characters interact?
Gameplay          How players interact?
Abilities         What special powers exist?
Setting           What world the story happens in?
```

## Final Observation

The RGB System behaves less like a traditional RPG ruleset and more like a **modular RPG engine**.

Because of this architecture, the same core system can support many different genres
while remaining simple and consistent.

← [Back to README](README.md)
