# Example: Classify "The fighter charges in and swings a greatsword"

**Procedure:** [`../procedures/classify.md`](../procedures/classify.md)

## Classification

```text
Action: The fighter charges in and swings a greatsword.
Classification: Press (R)
Cited: core.action.press
```

The character is changing the source of pressure through direct force —
this is `core.action.press`, which `depends_on` `core.vector.r`.

## Contrast: An Ambiguous Case

**Action:** "The scout shoves a mercenary off the walkway."

```text
Action: The scout shoves a mercenary off the walkway.
Classification: ambiguous — plausibly Press (R, forceful contact) or
  Reposition (G, changing the target's position/relation to the fight)
Cited: core.action.press, core.action.reposition
```

Per `procedures/classify.md` step 3, Specialist reports both plausible
classifications instead of silently picking one — the GM decides which
resolution fits their table.
