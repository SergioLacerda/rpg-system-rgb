# Core Terminology (EN / PT-br)

Canonical term pairs Specialist must use consistently in answers, sourced
from `core.term.*` and `core.translation.pt-br.*` semantic units in
`docs/core/semantic/core-v2.index.json`. Do not paraphrase these terms —
use them verbatim so answers stay traceable back to the semantic index.

| EN Term | Semantic ID | PT-br Term | Translation ID |
|---|---|---|---|
| Health | `core.term.health` | Vida | `core.translation.pt-br.health` |
| Impact Source | `core.term.impact-source` | Fonte de Impacto | `core.translation.pt-br.impact-source` |
| Attack Margin | `core.term.attack-margin` | Margem de Ataque | `core.translation.pt-br.attack-margin` |
| Armor Reduction | `core.term.armor-reduction` | Reducao por Armadura | `core.translation.pt-br.armor-reduction` |
| Shield Absorption | `core.term.shield-absorption` | Absorcao por Escudo | `core.translation.pt-br.shield-absorption` |

## Turn Economy (post-canonicalization)

Not yet backed by a dedicated `core.term.*` unit (see
`.analysis/refined/20260801-core-rules-executable-examples-turn-contract`
for the canonicalization work) — listed here directly from
`docs/core/en/combat/movement.md`'s "Combat Turn Structure" section:

| EN Term | PT-br Term |
|---|---|
| Movement | Movimento |
| Action | Ação |
| Minor Action | Ação Menor |

## Adding A Term

A term belongs here only once it has a `core.term.*` (and, if translated,
`core.translation.pt-br.*`) semantic unit — do not add ad hoc glossary
entries that aren't traceable to the index. Propose the semantic unit
first (`docs/core/semantic/core-v2.index.json`), then add the row here.
