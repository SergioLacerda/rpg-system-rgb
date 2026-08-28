# Advanced Test Coverage Context Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Resolve the current public-doc validation failure and mature the existing BDD context without adopting a Gherkin runner.

**Architecture:** Keep public documentation free of internal workspace path markers. Preserve the current BDD model: `.feature` files are readable specifications, Go mirror tests are executable enforcement, and `tests/features` validates the traceability contract.

**Tech Stack:** Go `testing`, Make targets, Markdown documentation, existing semantic docs validation.

---

### Task 1: Remove public internal path references

**Files:**
- Modify: `docs/adr/adr-015-skill-runtime-split.md`
- Modify: `skills/maker/README.md`

**Step 1:** Replace internal workspace-analysis path references with public descriptions of the accepted planning/review context.

**Step 2:** Run the public-path marker scan through the Go validation command and expect no internal path marker failures.

**Step 3:** Run `GOCACHE=/tmp/go-cache go test ./tests/semantic_docs ./internal/components/tooling`.

### Task 2: Document current BDD maturity model

**Files:**
- Modify: `docs/engineering/base-structure-review-workflow.md`

**Step 1:** Add a concise note explaining that the current Gherkin/BDD model is feature files plus Go mirror anchors.

**Step 2:** State that a real runner is deferred until it proves better drift reduction than the current lightweight mirror contract.

**Step 3:** Run `GOCACHE=/tmp/go-cache go test ./tests/features ./tests/core_behavior ./tests/properties ./tests/fixtures`.

### Task 3: Validate the integrated result

**Files:**
- No additional edits expected.

**Step 1:** Run `GOCACHE=/tmp/go-cache go test ./...`.

**Step 2:** Run `GOCACHE=/tmp/go-cache make cover-check`.

**Step 3:** Report any remaining failures as residuals rather than broadening scope.
