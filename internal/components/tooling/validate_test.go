package tooling

import (
	"os"
	"path/filepath"
	"testing"
)

// repoRoot resolves the repository root relative to this test file
// (internal/components/tooling is two directories below the root).
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}
	if !dirExists(filepath.Join(root, "docs")) {
		t.Fatalf("resolved path %s does not look like the repo root", root)
	}
	return root
}

func TestValidatePassesAgainstRealRepoData(t *testing.T) {
	root := repoRoot(t)
	if err := Validate(root); err != nil {
		t.Fatalf("Validate failed against real repo data: %v", err)
	}
}

func TestValidateProjectPathsRejectsForbiddenMarker(t *testing.T) {
	scratch := t.TempDir()
	mustMakeMinimalProjectShape(t, scratch)
	if err := os.WriteFile(filepath.Join(scratch, "README.md"), []byte("references .analysis/ marker"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateProjectPaths(scratch); err == nil {
		t.Fatal("expected validation to fail when README.md references an internal path marker")
	}
}

// mustMakeMinimalProjectShape makes a minimal scratch copy sufficient for
// ValidateProjectPaths: the required/forbidden directory shape plus an
// empty README.md placeholder to be overwritten by the caller.
func mustMakeMinimalProjectShape(t *testing.T, scratch string) {
	t.Helper()
	for _, dir := range []string{"docs/core", "docs/adr", "generated"} {
		if err := os.MkdirAll(filepath.Join(scratch, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(scratch, "README.md"), []byte("placeholder"), 0o644); err != nil {
		t.Fatal(err)
	}
}
