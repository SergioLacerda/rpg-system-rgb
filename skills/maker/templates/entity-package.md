<!--
Template: Entity Package
Schema: ../schemas/entity-package.schema.yaml
Fill every section. Leave "Unresolved Fields" and "Conflicts" explicitly
empty ("None.") rather than deleting the section — an absent section reads
as "not checked," an empty one reads as "checked, nothing found."
-->

# Entity Package: {{id}}

**Type:** {{entity_type}} · **Status:** {{status}}
**Source(s):** {{source_refs}}

## Explicit Facts

Stated directly in the source.

| Field | Value | Source |
|---|---|---|
| {{field}} | {{value}} | {{source_ref}} |

## Inferred Suggestions

Plausible, but not stated directly — mark confidence honestly.

| Field | Value | Confidence | Rationale |
|---|---|---|---|
| {{field}} | {{value}} | {{low\|medium\|high}} | {{rationale}} |

## Unresolved Fields

Fields this entity type normally has, that the source didn't cover.

- {{field}} — {{reason}}

## Conflicts

Links to conflict reports touching this entity, if any.

- {{conflict-report id}}

## Canonical Dependencies

RGB System semantic IDs this entity's mechanics depend on, if any (e.g. an
ability referencing `core.ability.contract`). Read-only — this package does
not modify canonical units.

- {{semantic id}}
