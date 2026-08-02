# Procedure: Classify

Classify a described action by RGB vector and tactical decision procedure.

## Input

A short description of what a character does (e.g. "the fighter charges
in and swings a greatsword" or "the scout ducks behind cover").

## Steps

1. Match the described action against the three decision procedures:

   | Procedure | Vector | Semantic ID | Meaning |
   |---|---|---|---|
   | Press | R | `core.action.press` | changes the source of pressure — attack, force, disruption |
   | Reposition | G | `core.action.reposition` | changes relation to pressure — movement, evasion, positioning |
   | Sustain | B | `core.action.sustain` | preserves continuity under pressure — endurance, shields, stabilization |

2. If the action clearly maps to exactly one procedure, classify it and
   cite the matching `core.action.*` unit and its `depends_on` vector
   (`core.vector.r`, `core.vector.g`, or `core.vector.b`).
3. If the action plausibly maps to more than one procedure (e.g. "shove an
   enemy off a ledge" touches both Press and Reposition), report all
   plausible classifications instead of silently picking one — this is an
   uncertainty, not a rule gap, and must be surfaced per `SKILL.md`'s
   "surface uncertainty" requirement.
4. If the action does not describe a combat/tactical decision at all (pure
   roleplay, no mechanical resolution), classification does not apply —
   say so rather than forcing a fit.

## Output Shape

```text
Action: <restated action>
Classification: <Press|Reposition|Sustain> (<vector>)
Cited: <core.action.* id>
```

## Worked Example

See [`../examples/classify-charge-attack.md`](../examples/classify-charge-attack.md).

## Forbidden

- Forcing an ambiguous action into a single procedure to look decisive.
- Classifying without citing the `core.action.*` unit.
