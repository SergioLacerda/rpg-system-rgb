# Prompt: Detect Conflicts

Record and classify disagreements when two or more sources (notes or images)
state different values for the same field of an entity.

## Steps

1. Isolate each disputed field into its own
   [`conflict-report`](../templates/conflict-report.md) — split multi-field
   disagreements (e.g. `name` vs. `loyalty`) into separate reports. When a
   disagreement reflects a whole-entity reboot or remake across many
   fields rather than an isolated field dispute, record alternative
   versions under `variants` on the `entity-package` per
   `variant-preserved` instead of flooding single-field conflict reports.
2. Record every competing claim with `value`, `source_ref`, and
   `source_kind` (`text` or `image`).
3. Apply `text-over-image` from [`../SKILL.md`](../SKILL.md) to set
   `resolution`:
   - If claims include both `text` and `image` and the disputed field is a
     **non-appearance** fact (name, loyalty, backstory, canon status,
     relationships, stats), set `resolution: resolved_by_authority_rule`
     and explain in `resolution_note` why the text source wins.
   - If all claims are `text`, or all claims are `image`, or the field
     genuinely addresses appearance, keep `resolution: unresolved` per
     `conflict-on-dispute` so the user decides (never break ties by recency
     or detail level).
4. Link the conflict report ID from every affected `entity-package`.
5. **Post-resolution archival:** keep only `unresolved` reports directly in
   `packages/conflicts/`. Once a report reaches `resolved_by_user` or
   `resolved_by_authority_rule`, archive it under `packages/conflicts/resolved/`
   and remove any stale copy from `packages/conflicts/`.

## Completion Check

- At least two claims are recorded on every conflict report.
- `resolution_note` is present whenever `resolution` is
  `resolved_by_authority_rule` or `resolved_by_user`.

## Worked Example

Read [`../examples/harbor-quartermaster.md`](../examples/harbor-quartermaster.md)
only when verifying how `conflict-report` interacts with an entity's
`unresolved_fields`.
