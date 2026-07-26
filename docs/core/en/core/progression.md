# Progression

As characters gain experience, they grow stronger. Progression in the
**RGB System** is built on a small, fixed advancement budget per level,
which a player may spend either as raw vector growth or, when the
campaign uses the Skills and Abilities module, as one of several
alternative advancement choices.

This document is the canonical source for character progression. It
replaces the short summary previously kept inline in
[Character Creation](character_creation.md), which now links here.

## Advancement Budget

```text
+2 vector points per level
```

The Game Master determines when characters gain a level, depending on the
campaign. This baseline never changes — every advancement choice below
spends the same fixed budget, it does not add to it.

## Advancement Choices

By default, the full budget becomes vector points. When the campaign
allows the [Skills and Abilities](skills_and_abilities.md) module, the
player may instead spend a level's advancement on one alternative below.
Exactly one choice applies per level; choices do not stack within the
same level.

| Choice | Effect |
| --- | --- |
| **Vector growth** (default) | +2 points, distributed across R, G, and B as the player chooses |
| **New ability** | Gain one ability whose `tier` and `requirements` (see the Ability Contract in [Skills and Abilities](skills_and_abilities.md)) the character currently meets |
| **Ability improvement** | Reduce one known ability's `cost`, extend its `duration`, or raise its `limits` ceiling by one step, at the Game Master's discretion |
| **Specialization** | Commit to one vector's ability tree (Red, Green, or Blue): future new-ability choices in that tree ignore one tier of `requirements` |
| **New resource** | Gain a campaign-defined resource pool (e.g. an extra reaction charge, a callback-recovery buffer) sized by the Game Master |
| **New reaction** | Gain access to one reaction procedure not otherwise available at the character's current tier |
| **New state or maneuver access** | Gain the ability to declare one tactical state or maneuver (see [Attack and Defense](../combat/attack_and_defense.md), [Movement](../combat/movement.md)) that was previously restricted or unavailable |

## Specialization In Practice

Specialization is the one choice with a lasting structural effect: it
commits a character to a vector's identity rather than granting an
immediate benefit. A Red-specialized character does not gain a Red
ability from specializing alone — specialization changes how *future*
Red abilities are gated, per the Ability Contract's `requirements` field.
This keeps specialization meaningful without requiring new resource
tracking beyond what the Ability Contract already defines.

## Design Rationale

The RGB System V2 initial analysis (`base_project` intake, §7.5) left
this choice open: pure numeric vector growth, or a per-advancement choice
among several options. This document resolves it as a **hybrid**, not a
new invention:

- The pure-numeric baseline (`+2 vector points per level`) was already
  canonical, stated in Character Creation before this document existed.
  This document keeps it unchanged as the default — no existing example
  or rule that assumed pure numeric growth is invalidated.
- [Skills and Abilities](skills_and_abilities.md) already showed
  abilities gated by level in its worked examples (`Level 3 → +1 R to
  ranged attacks`, `Level 6 → Double power strike`, `Level 9 → Area
  impact attack`) — a per-advancement ability choice was already implied
  in practice, just never formalized as a documented option alongside
  vector growth. This document only makes that existing pattern explicit
  and consistent with the Ability Contract's `tier`/`requirements`
  fields, which already exist specifically to gate ability access by
  advancement.
- A hybrid keeps the system's stated design philosophy intact: character
  creation and progression both favor **simplicity** (one fixed budget,
  one choice per level) while still supporting **tactical diversity**
  (the campaign, via the optional Skills and Abilities module, can offer
  more than pure numbers without adding new subsystems).

## Example

A level-6 character with `R=5, G=3, B=4` and the Skills and Abilities
module active can choose:

```text
Vector growth  → R=7, G=3, B=4 (or any +2 split)
New ability    → gain one ability meeting tier <= character's current tier
Specialization → commit to Red, Green, or Blue for future ability gating
```

The player picks one; the choice does not carry over or accumulate.

See also:

- [Character Creation](character_creation.md)
- [Skills and Abilities](skills_and_abilities.md)
- [Attributes](attributes.md)

← [Back to README](README.md)
