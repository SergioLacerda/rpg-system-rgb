# ADR-011: Split License Model

## Status

Accepted.

## Context

The repository contains two materially different kinds of work:

- executable software surfaces such as Go packages, shell scripts, CI, and
  Astro application code;
- creative and rules-authority surfaces such as RPG rules, examples,
  documentation, diagrams, public copy, generated Library/PDF content, and art.

The prior single CC BY 4.0 repository license was suitable for creative rules
text, but it was less conventional for software reuse and downstream tooling.

## Decision

Adopt a split model:

| Surface | License |
| --- | --- |
| Source code, scripts, CI, and tooling | MIT or Apache-2.0 |
| Rules text, examples, documentation, diagrams, web copy, generated content, and art | CC BY 4.0 |

`LICENSES.md` is the repository surface map. The root `LICENSE` file points to
the concrete license texts.

## Consequences

- Software consumers can use the code under a conventional permissive software
  license.
- Rules and documentation retain attribution requirements appropriate for
  creative RPG content.
- New files should be classified by surface ownership before publication.
- Ambiguous generated artifacts inherit the license of their source surface.
