# Example: Trace `core.example.combat-turn`

**Procedure:** [`../procedures/trace.md`](../procedures/trace.md)

## Trace

```text
ID: core.example.combat-turn
Kind: example | Status: canonical (canonical_markdown_bridge)
Source: docs/core/en/reference/combat_example.md
Depends on: core.combat.attack-margin, core.damage.flow,
  core.damage.impact-source, core.damage.armor-reduction,
  core.damage.shield-absorption
Decision refs: docs/core/en/reference/combat_example.md
```

Following one dependency further:

```text
ID: core.damage.flow
Kind: procedure | Status: canonical (canonical_semantic)
Source: docs/core/en/combat/damage_model.md
Depends on: core.combat.attack-margin, core.damage.impact-source,
  core.damage.penetration, core.damage.armor-reduction,
  core.damage.shield-absorption
Decision refs: docs/adr/adr-001-ai-first-documentation-authority.md,
  docs/core/en/combat/damage_model.md
```

This is how a user gets from "here's a worked example" back to "here's the
rule it illustrates, and here's the ADR that established it as canonical."
