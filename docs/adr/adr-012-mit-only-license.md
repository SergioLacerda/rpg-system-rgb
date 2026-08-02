# ADR-012: MIT-Only License

## Status

Accepted.

## Context

ADR-011 adopted a split licensing model: source code, scripts, CI, and
tooling were dual-licensed under MIT or Apache-2.0, at the user's option,
while rules text, examples, documentation, diagrams, web copy, and generated
Library/PDF content were licensed under Creative Commons Attribution 4.0
International (CC BY 4.0), which requires attribution.

The goal has since narrowed: a single, maximally permissive license across
the entire repository, with unrestricted open-source use, and with the
copyright notice itself naming the actual author rather than only the
project name.

## Decision

Adopt MIT as the sole license for the entire repository — code, rules text,
documentation, diagrams, web copy, generated content, and art alike. The
Apache-2.0 option and the CC BY 4.0 license are retired from active use.

The root `LICENSE` file holds the full MIT license text directly (not an
index to separate license files), with:

```
Copyright (c) 2026 Sergio Lacerda (RGB System)
```

`LICENSE-MIT`, `LICENSE-APACHE-2.0`, `LICENSE-CC-BY-4.0`, and `LICENSES.md`
are removed — there is no split model left to map or duplicate.

## Consequences

- The repository's licensing surface collapses to a single `LICENSE` file.
- Anyone may use, copy, modify, merge, publish, distribute, sublicense, and
  sell all repository content, subject only to MIT's copyright-notice
  retention condition — matching the request for unrestricted open-source
  access.
- The CC BY 4.0 attribution *requirement* for rules text, examples, docs,
  diagrams, and art is dropped. This is an explicit, accepted trade-off in
  favor of being fully unrestricted, not an oversight.
- The Apache-2.0 explicit patent grant is no longer offered for code.
- All prose references to the split model (`README.md`, `docs/README-details.md`)
  are updated to point at `LICENSE` directly.

## Supersession

This ADR supersedes ADR-011. ADR-011 remains as the historical record of
why the split model was originally adopted.
