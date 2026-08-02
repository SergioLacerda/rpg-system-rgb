# Example: Explain "How does an attack succeed?"

**Question:** "How do I know if an attack hits?"

**Procedure:** [`../procedures/explain.md`](../procedures/explain.md)

**Resolution:** matches `core.combat.attack-margin` (kind: `rule`) —
preferred over `core.example.combat-turn` (kind: `example`), which only
illustrates one instance of this rule.

## Answer

```text
Answer: An attack succeeds when Attack Margin is zero or higher.

  Attack Margin = Attacker (R + modifiers) - Defender (G + modifiers)

Zero or higher is at least a basic hit; higher margins may unlock better
outcomes per the attack resolution rules.

Cited: core.combat.attack-margin (docs/core/en/combat/attack_and_defense.md)
```

## Trace (for follow-up)

`core.combat.attack-margin` `depends_on` `core.vector.r`, `core.vector.g`,
and is `illustrated_by` `core.example.combat-turn` — the worked combat
example, if the user wants numbers instead of the formula.
