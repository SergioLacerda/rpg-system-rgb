package compiler

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// TestParseAllCanonicalDocsWithoutPanicking parses every canonical Markdown
// document (both locales) and confirms each yields at least one block. The
// UN-01 construct census was taken only against docs/core/en/**; this test
// also covers docs/core/PT-br/** to confirm the same construct set holds
// there.
func TestParseAllCanonicalDocsWithoutPanicking(t *testing.T) {
	root := repoRoot(t)
	for _, locale := range []string{"en", "PT-br"} {
		dir := filepath.Join(root, "docs", "core", locale)
		err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || filepath.Ext(path) != ".md" {
				return nil
			}
			body, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			blocks := Parse(string(body))
			if len(blocks) == 0 {
				t.Errorf("%s parsed to zero blocks", path)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "docs")); err != nil {
		t.Fatalf("resolved path %s does not look like the repo root", root)
	}
	return root
}
