package tooling

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFixtureDoc writes a small docs/core/<locale>/<relPath> file under
// repoRoot for GenerateSkillSnapshot to concatenate.
func writeFixtureDoc(t *testing.T, repoRoot, locale, relPath, content string) {
	t.Helper()
	target := filepath.Join(repoRoot, "docs", "core", locale, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateSkillSnapshotFormat(t *testing.T) {
	root := t.TempDir()
	// GenerateSkillSnapshot has no per-input locale filtering: it
	// concatenates every declared (input, locale) pair, outer loop over
	// Inputs, inner loop over Locales. So every input needs a fixture for
	// every declared locale.
	writeFixtureDoc(t, root, "en", "topic-a.md", "# Topic A (EN)\n\nBody.\n")
	writeFixtureDoc(t, root, "pt", "topic-a.md", "# Topic A (PT)\n\nCorpo.\n")
	writeFixtureDoc(t, root, "en", "topic-b.md", "# Topic B (EN)\n\nNo trailing newline")
	writeFixtureDoc(t, root, "pt", "topic-b.md", "# Topic B (PT)\n\nSem nova linha final")

	cfg := SkillSnapshotConfig{
		OutputPath: "out/snapshot.md",
		Inputs:     []string{"topic-a.md", "topic-b.md"},
		Locales:    []string{"en", "pt"},
	}

	if err := GenerateSkillSnapshot(root, cfg); err != nil {
		t.Fatalf("GenerateSkillSnapshot failed: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(root, "out", "snapshot.md"))
	if err != nil {
		t.Fatal(err)
	}

	want := skillSnapshotHeader +
		"<!-- Source: docs/core/en/topic-a.md -->\n\n" +
		"# Topic A (EN)\n\nBody.\n" +
		"\n---\n\n" +
		"<!-- Source: docs/core/pt/topic-a.md -->\n\n" +
		"# Topic A (PT)\n\nCorpo.\n" +
		"\n---\n\n" +
		"<!-- Source: docs/core/en/topic-b.md -->\n\n" +
		"# Topic B (EN)\n\nNo trailing newline\n" +
		"\n---\n\n" +
		// topic-b's "pt" fixture has no trailing newline in its source —
		// GenerateSkillSnapshot must add exactly one, not double up.
		"<!-- Source: docs/core/pt/topic-b.md -->\n\n" +
		"# Topic B (PT)\n\nSem nova linha final\n"

	if string(got) != want {
		t.Errorf("snapshot format mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGenerateSkillSnapshotStripsRelativeLinksButKeepsExternalOnes(t *testing.T) {
	root := t.TempDir()
	writeFixtureDoc(t, root, "en", "topic.md",
		"See [Movement](../combat/movement.md) and [the repo](https://example.com/repo) "+
			"and [an email](mailto:team@example.com).\n")

	cfg := SkillSnapshotConfig{
		OutputPath: "out/snapshot.md",
		Inputs:     []string{"topic.md"},
		Locales:    []string{"en"},
	}
	if err := GenerateSkillSnapshot(root, cfg); err != nil {
		t.Fatalf("GenerateSkillSnapshot failed: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(root, "out", "snapshot.md"))
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(got), "](../combat/movement.md)") {
		t.Error("expected the relative link to be stripped to plain text, but it survived")
	}
	if !strings.Contains(string(got), "See Movement and") {
		t.Errorf("expected the relative link's label to remain as plain text, got:\n%s", got)
	}
	if !strings.Contains(string(got), "[the repo](https://example.com/repo)") {
		t.Errorf("expected the external link to survive unchanged, got:\n%s", got)
	}
	if !strings.Contains(string(got), "[an email](mailto:team@example.com)") {
		t.Errorf("expected the mailto link to survive unchanged, got:\n%s", got)
	}
}

func TestGenerateSkillSnapshotMissingFileErrors(t *testing.T) {
	root := t.TempDir()
	cfg := SkillSnapshotConfig{
		OutputPath: "out/snapshot.md",
		Inputs:     []string{"missing.md"},
		Locales:    []string{"en"},
	}

	err := GenerateSkillSnapshot(root, cfg)
	if err == nil {
		t.Fatal("expected an error for a missing input file, got nil")
	}
	if !strings.Contains(err.Error(), "missing.md") {
		t.Errorf("expected error to name the missing file, got: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(root, "out", "snapshot.md")); statErr == nil {
		t.Error("expected no output file to be written on error, but one exists")
	}
}
