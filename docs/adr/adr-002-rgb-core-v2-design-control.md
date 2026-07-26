# ADR-002: RGB Core V2 Design Control

## Status

Accepted.

## Context

RGB Core V2 needed a controlled design surface before downstream tracks such as
Maker, Specialist, tooling, Library, PDF, landing, and bundles could consume the
model safely.

## Decision

Treat Core V2 as the authoritative source of the playable RGB model. Core owns:

- vector grammar;
- procedures;
- resources;
- tactical states;
- ability contract boundaries;
- documentation authority for rules, procedures, examples, and design notes.

Downstream tracks consume Core decisions. They do not define Core mechanics.

## Consequences

- Core documentation can be normalized before downstream product surfaces.
- Maker and Specialist must consume stable Core IDs and contracts.
- Library, PDF, generated bundles, search indexes, and AI context packs remain
  derived surfaces.
- Documentation authority must stay explicit when rules are promoted into
  semantic source units.
