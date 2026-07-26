# System Context

```mermaid
C4Context
    title RGB System V2 — System Context

    Person(gm, "Game Master", "Runs RGB tabletop sessions")
    Person(player, "Player", "Plays RGB characters")

    System(rgb, "RGB System V2", "Minimalist RPG rules, documentation, and tooling ecosystem")

    System_Ext(viewer, "Markdown/Git Viewer", "GitHub, GitLab, or a local Markdown renderer")

    Rel(gm, rgb, "Reads rules, runs encounters")
    Rel(player, rgb, "Reads rules, builds characters")
    Rel(rgb, viewer, "Renders docs, ADRs, and diagrams")
```

RGB System V2 is a single, self-contained system from the outside: no
external services, no network dependencies, no accounts. Game Masters and
Players interact with it entirely by reading documentation and (for
Tooling contributors) running local Go binaries.

← [Back to Architecture Index](README.md)
