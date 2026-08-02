# Example: Locate Authority For "Does armor or shield absorb damage first?"

**Procedure:** [`../procedures/locate-authority.md`](../procedures/locate-authority.md)

## Search

Matches in `generated/ai-context/core-specialist-pack.json`:

- `core.damage.armor-reduction` (`rule`)
- `core.damage.shield-absorption` (`rule`)
- `core.damage.flow` (`procedure`) — declares `ordered_steps` including
  both, in order.
- `core.example.combat-turn` (`example`) — illustrates the sequence but
  does not define it.

## Resolution

`core.damage.flow` outranks the two individual rule units for this
specific question (order), since it's the unit whose `ordered_steps`
actually settles "which comes first," while the two rule units each only
describe their own single step in isolation.

```text
Topic: does armor or shield absorb damage first?
Authority: core.damage.flow (procedure)
Source: docs/core/en/combat/damage_model.md
```

Follow-up: if the user wants the mechanics of each step individually
(not just the order), `core.damage.armor-reduction` and
`core.damage.shield-absorption` are the follow-up authorities — see
[`../procedures/compare-rules.md`](../procedures/compare-rules.md).
