# RGB System — One Page Rules

A minimal reference to start playing the **RGB RPG System** immediately.

## Core Concept

Characters are defined by three vectors:

| Vector | Name | Function |
|||--|
| R | Red | damage, strength, offensive power |
| G | Green | agility, mobility, reaction |
| B | Blue | shields, energy, special defense |

Players distribute **7 starting points** between R, G and B.

Example:

R = 3  
G = 2  
B = 2

## Character Durability

```text
Health = 4 + (R × 2)
Shield = B × 3
```

Health represents physical endurance.  
Shield represents energy or special protection.

## Turn Structure

Each turn a character may perform:

```text
1 Action
+
1 Minor Adjustment
```

Examples:

Actions:

- attack
- move
- defend
- use ability

Minor adjustments:

- reload
- small reposition
- interact with object

## Tactical Choices

Most combat decisions fall into three categories:

```text
Attack → R
Move → G
Defend → B
```

This creates three combat strategies:

| Strategy | Vector |
|||
Striker | R |
Skirmisher | G |
Guardian | B |

Hybrid builds are possible.

## Attack Resolution

An attack succeeds if the attacker overcomes the defender according to the combat rules.

If the attack succeeds, damage is resolved.

## Damage Flow

```text
Attack Check
↓
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

Definitions:

Armor → reduces damage per hit  
Shield → absorbs accumulated damage  

## Example Damage

Weapon Damage: 7  
Penetration: 2  
Armor: 4

```text
Effective Armor = Armor − Penetration
Effective Armor = 2
Final Damage = 7 − 2 = 5
```

If a shield exists, damage is absorbed by the shield first.

## Gameplay Loop

A typical session follows:

```text
Character Creation
↓
Exploration
↓
Combat
↓
Recovery
↓
Story Progression
```

## RGB Tactical Model

```text
        G
   mobility / reaction

R - B
power / damage    shield / defense
```

R → deal damage  
G → avoid damage  
B → absorb damage  

## Design Philosophy

The RGB system focuses on:

- simple numeric relationships
- tactical combat decisions
- modular rules
- adaptable settings

The same rules can support:

- modern campaigns
- fantasy worlds
- science fiction
- superpowered settings
