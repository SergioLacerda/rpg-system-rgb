# MkDocs Library RGB Theme Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Theme the MkDocs-generated Library so it keeps the landing page header,
footer, RGB palette, spacing, and reading rhythm.

**Architecture:** Keep `docs/` as the documentation source and MkDocs as the
documentation renderer. Add one CSS file through `extra_css` and one minimal
MkDocs `main.html` override that replaces only the visible chrome while preserving
MkDocs content, navigation, search, and scripts.

**Tech Stack:** MkDocs, Jinja templates, CSS, Astro visual tokens.

---

### Task 1: Wire MkDocs Theme Inputs

**Files:**
- Modify: `mkdocs.yml`
- Create: `docs/styles/rgb-library.css`
- Create: `docs/overrides/main.html`

**Step 1:** Add `theme.name: mkdocs` and `theme.custom_dir: docs/overrides`.

**Step 2:** Add `extra_css: [styles/rgb-library.css]`.

**Step 3:** Create the CSS and template override files.

### Task 2: Implement Landing-Compatible Chrome

**Files:**
- Create: `docs/overrides/main.html`
- Reference: `web/landing/src/components/Header.astro`
- Reference: `web/landing/src/components/Footer.astro`

**Step 1:** Extend MkDocs `base.html`.

**Step 2:** Override `site_name`, `site_nav`, and `footer` blocks.

**Step 3:** Keep MkDocs search, previous/next, scripts, and content blocks untouched.

### Task 3: Implement RGB Visual Styling

**Files:**
- Create: `docs/styles/rgb-library.css`
- Reference: `web/landing/src/styles/tokens.css`

**Step 1:** Define RGB CSS variables.

**Step 2:** Override MkDocs Bootstrap navbar, dropdowns, content, TOC, tables, code,
search modal, buttons, and responsive layout.

**Step 3:** Keep selectors scoped to the generated Library page classes where practical.

### Task 4: Validate

**Files:**
- Generated: `web/landing/public/library/**`
- Generated: `web/landing/dist/library/**`

**Step 1:** Run `make docs-build`.

**Step 2:** Run `make landing-build LANDING_BASE=/rpg-system-rgb`.

**Step 3:** Preview and check `/rpg-system-rgb/library/`.
