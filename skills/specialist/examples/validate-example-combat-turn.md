# Example: Validate `core.example.combat-turn`

**Procedure:** [`../procedures/validate-example.md`](../procedures/validate-example.md)

## Re-derivation (second attack, per `core.damage.flow`)

```text
Weapon Damage = 4, Penetration = 1, Defender Armor = 4, Defender Shield = 3

Effective Armor  = max(4 - 1, 0)        = 3
Armor Reduced    = min(3, 4)            = 3
After Armor      = 4 - 3               = 1
Shield Absorbed  = min(3, 1)            = 1
Health Damage    = 1 - 1               = 0
```

## Comparison

```text
Example: core.example.combat-turn (second attack)
Rule(s): core.damage.flow, core.damage.armor-reduction, core.damage.shield-absorption
Result: CONSISTENT
Mismatches: none
Cited: docs/core/en/reference/combat_example.md
```

This is the corrected state of the example — before
`.analysis/refined/20260801-core-rules-executable-examples-turn-contract`,
the published document stated `Health Damage = 1, Remaining Health = 9`
for this attack, which fails re-derivation twice over: the shield
(3 available) was never applied to absorb the 1 remaining damage, and
health cannot increase from a damage step under any re-derivation of
`core.damage.flow`.

The actual regression test for this validation is
`tests/fixtures/combat_example_doc_parity_test.go`
(`TestCombatExampleDocMatchesFixture`) — it re-parses the published
Markdown's numbers on every test run and fails if they ever drift from the
fixture again.
