<!--
Template: Conflict Report
Schema: ../schemas/conflict-report.schema.yaml
One report per disputed field. Default resolution is "unresolved" — do not
pick a winner unless the authority rule in SKILL.md actually applies, or
the user explicitly chose one.
-->

# Conflict Report: {{id}}

**Status:** {{status}}
**Field in dispute:** {{field}}
**Linked entity:** {{linked_entity, if any}}

## Competing Claims

| Value | Source | Source Kind |
|---|---|---|
| {{value}} | {{source_ref}} | {{text\|image}} |
| {{value}} | {{source_ref}} | {{text\|image}} |

## Resolution

**Status:** {{unresolved \| resolved_by_authority_rule \| resolved_by_user}}

{{resolution_note, if resolved}}
