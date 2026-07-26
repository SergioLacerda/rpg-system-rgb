# Character Creation

Character creation in the **RGB System** is intentionally simple.  
Players distribute points between the three RGB vectors to define how their character performs in combat and interaction.

```text
        G
   relation / positioning

R ---------------- B
pressure            preservation
```

- **R (Red)** → pressure, force, attack, disruption, and impact
- **G (Green)** → movement, timing, evasion, range, and positioning
- **B (Blue)** → endurance, blocking, stabilization, shields, and protection

This distribution defines the character's tactical style.

## Starting Points

Characters begin with:

```text
7 points
```

These points must be distributed between:

```text
R (Red)
G (Green)
B (Blue)
```

Example:

```text
R = 3
G = 2
B = 2
```

Different distributions naturally create different combat roles.

```text
| Style | Vector Focus |
|------|--------------|
Heavy attacker | High R |
Mobile skirmisher | High G |
Defensive guardian | High B |
```

Hybrid builds are also possible.

## Character Progression

As characters gain experience they become stronger. See
[Progression](progression.md) for the full advancement rules, including
the per-advancement choice between vector growth and alternatives such as
new abilities, specialization, or new resources.

## Health

Character durability uses both pressure tolerance and preservation.

```text
Health = 4 + R + B
```

Example:

```text
R = 3
B = 2
Health = 9
```

This prevents **R** from owning both offense and durability by itself. Higher
**R** still helps characters withstand physical pressure, while higher **B**
represents preservation, stability, and defensive endurance.

## Quick Reference

```text
Attack    : margin = attacker R - defender G
Movement  : G × 2 meters
Health    : 4 + R + B
Shield    : B × 3
```

These formulas summarize the core mechanical interactions of the RGB system.

## Design Philosophy

Character creation in RGB follows three principles:

- **simplicity** — few numbers define the character
- **tactical diversity** — different vector distributions create different play styles
- **modularity** — abilities and equipment can expand the system without changing its core rules

Because every character distributes points between the same three vectors, the system naturally produces balanced and varied builds.

See also:

- [Attributes](attributes.md)
- [Progression](progression.md)
- [Combat](../combat/attack_and_defense.md)
- [Movement](../combat/movement.md)

← [Back to README](README.md)
