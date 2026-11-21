# Procedure: Explain

Answer a question about a rule, resource, or concept in the RGB System.

## Input

A natural-language question, or a semantic ID directly (e.g.
`core.combat.attack-margin`).

## Steps

1. If the input is not already a known semantic ID, resolve it via
   [`locate-authority.md`](locate-authority.md) first; if no unit matches,
   stop and report `"no canonical rule found for <question>"`.
2. Enforce `rule-over-example`: ground the explanation in `kind: rule` or
   `kind: procedure` units (only illustrate with `core.example.*` afterwards
   if helpful).
3. Read the unit's `retrieval_summary` and its canonical Markdown file at
   `projection_paths.markdown_en` (plus `markdown_pt_br` when answering in
   Portuguese) — read the source directly rather than relying on memory.
4. Compose the answer strictly from those verified sources.
5. Attach every semantic ID and `source_path` used, plus all applicable
   status/locale disclosures (`draft`, `stale`, `needs_review`, `localized`,
   `generated`) per `citation-gated`.

## Output Shape

```text
Answer: <explanation, grounded in the cited unit(s)>
Cited: <semantic id> (<source_path>)
Disclosures: <draft | stale | needs_review | localized | generated, or "none (canonical en)">
```

## Worked Example

Read [`../examples/explain-attack-margin.md`](../examples/explain-attack-margin.md)
only if the formatting of formula citations or margin examples is unclear.
