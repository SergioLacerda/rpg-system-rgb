# RGB System Architecture – Layered Design Observation

This note documents an interesting structural property that appears when the RGB documentation
is analyzed as a whole.

Although the RGB System was designed as a simple tabletop RPG system, the documentation naturally
reveals a **layered architecture** similar to what is often found in game engines.

This structure improves clarity, modularity and long‑term expandability.

## Layered Architecture of the RGB System

When the documentation is organized by function, the system forms the following layers:

```text
System Layer
↓
Gameplay Layer
↓
Combat Layer
↓
Damage Layer
↓
Equipment Layer
```

Each layer answers a different design question.

## 1. System Layer

Documents:

- system_overview.md
- rgb_damage_interaction_model.md

Purpose:

This layer explains the **core philosophy and structure of the system**.

It defines the RGB vectors:

- **R (Red)** – damage and physical power
- **G (Green)** – mobility and reaction
- **B (Blue)** – shields and energy defense

It also shows how the system modules connect.

Equivalent concept in game engines:

Core system design.

## 2. Gameplay Layer

Documents:

- gameplay_loop.md
- combat_decision_model.md

Purpose:

This layer explains **how players interact with the system during play**.

Typical flow:

``````text
Character Creation
↓
Exploration
↓
Combat
↓
Recovery
```

During combat the player usually chooses between:

``````text
Attack
Move
Defend
```

Which map naturally to:

``````text
Attack → R
Move → G
Defend → B
```

Equivalent concept in game engines:

Gameplay framework.

# 3. Combat Layer

Documents:

- attack_and_defense.md
- movement.md

Purpose:

Defines how characters interact during encounters.

This includes:

- movement rules
- attack resolution
- defensive reactions
- positioning

Equivalent concept in game engines:

Combat system.

# 4. Damage Layer

Documents:

- damage_model.md
- rgb_damage_interaction_model.md

Purpose:

Handles the internal logic of damage calculation.

The RGB damage engine follows this structure:

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

This layer allows combat to remain tactical while keeping calculations simple.

Equivalent concept in game engines:

Damage engine.

# 5. Equipment Layer

Documents:

- armor.md
- shields.md
- firearms.md
- melee.md
- explosives.md

Purpose:

Defines the **data values used by the combat and damage systems**.

Examples:

- weapon damage
- armor reduction
- penetration values
- shield capacity

Equivalent concept in game engines:

Data layer.

# Complete RGB System Architecture

When combined, the RGB documentation forms the following structure:

```text
System Design
   │
   ├─ System Overview
   ├─ RGB Interaction Model
   │
Gameplay
   │
   ├─ Gameplay Loop
   ├─ Combat Decision Model
   │
Combat
   │
   ├─ Movement
   ├─ Attack & Defense
   │
Damage Engine
   │
   ├─ Damage Model
   │
Equipment Data
   │
   ├─ Armor
   ├─ Shields
   ├─ Weapons
```

# Design Consequences

This layered architecture produces several advantages.

### Modularity

New subsystems can be added without modifying core rules.

Examples:

- magic systems
- cybernetics
- vehicles
- superpowers

### Easier Balance

Most balancing occurs in:

```
Damage Layer
Equipment Layer
```

The rest of the system remains stable.

### Clear Documentation

Each layer answers a different question:

```text
| Layer | Question |
||-|
System | How the system is structured |
Gameplay | How players interact with the game |
Combat | How characters interact |
Damage | How damage is calculated |
Equipment | Which numerical values exist |
```

# Final Observation

The RGB System behaves less like a traditional RPG rule set and more like a **reusable RPG system engine**.

This makes the system particularly well suited for:

- modern campaigns
- fantasy worlds
- science fiction settings
- super‑powered universes

In this model:

```text
RGB System → rule engine
Setting (e.g., Aurora) → world built on top of the engine
```

This separation allows the RGB System to remain generic while different settings provide narrative context.

← [Back to README](README.md)
