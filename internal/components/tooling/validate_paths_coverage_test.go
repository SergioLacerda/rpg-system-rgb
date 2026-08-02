package tooling

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateProjectPathsEdges(t *testing.T) {
	root := t.TempDir()
	if err := ValidateProjectPaths(root); err == nil {
		t.Fatal("expected missing required dirs to fail")
	}

	writePathFixture(t, filepath.Join(root, "README.md"), "public")
	if err := os.MkdirAll(filepath.Join(root, "docs", "core"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "docs", "adr"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "generated"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := ValidateProjectPaths(root); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(root, "docs", "generated"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := ValidateProjectPaths(root); err == nil {
		t.Fatal("expected forbidden legacy directory to fail")
	}
	if err := os.RemoveAll(filepath.Join(root, "docs", "generated")); err != nil {
		t.Fatal(err)
	}

	writePathFixture(t, filepath.Join(root, "docs", "core", "leak.md"), "see .sdd/runtime")
	if err := ValidateProjectPaths(root); err == nil {
		t.Fatal("expected public runtime marker to fail")
	}
	if err := os.Remove(filepath.Join(root, "docs", "core", "leak.md")); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(root, "generated", "client"), 0o750); err != nil {
		t.Fatal(err)
	}
	writePathFixture(t, filepath.Join(root, "generated", "client", "governance.txt"), "allowed .strategist/client marker")
	if err := ValidateProjectPaths(root); err != nil {
		t.Fatal(err)
	}
}

func TestRejectRuntimeRefsHelpers(t *testing.T) {
	dir := t.TempDir()
	if err := rejectRuntimeRefs(filepath.Join(dir, "missing"), nil); err == nil {
		t.Fatal("expected missing root to fail")
	}

	file := filepath.Join(dir, "file.txt")
	writePathFixture(t, file, "mentions .analysis/internal")
	if err := rejectRuntimeRefs(file, nil); err == nil {
		t.Fatal("expected file marker to fail")
	}
	writePathFixture(t, file, "public")
	if err := rejectRuntimeRefs(file, nil); err != nil {
		t.Fatal(err)
	}

	if isExcludedWalkDir("/no/such/root", filepath.Join(dir, "child"), []string{"child"}) {
		t.Fatal("unexpected excluded dir for unrelated root")
	}
	if !isExcludedWalkDir(dir, filepath.Join(dir, "child"), []string{"child"}) {
		t.Fatal("expected child directory exclusion")
	}
}

func writePathFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o640); err != nil {
		t.Fatal(err)
	}
}
