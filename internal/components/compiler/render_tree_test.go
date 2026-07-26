package compiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRenderHTMLTreeMirrorsSourceTree renders the real repo's docs into a
// scratch copy and confirms every source .md file produced exactly one
// matching .html file, 1:1, for both locales.
func TestRenderHTMLTreeMirrorsSourceTree(t *testing.T) {
	root := repoRoot(t)
	scratch := t.TempDir()
	mustCopyDocsTree(t, root, scratch)

	rendered, err := RenderHTMLTree(scratch)
	if err != nil {
		t.Fatalf("RenderHTMLTree failed: %v", err)
	}

	wantSourceCount := 0
	for _, locale := range []string{"en", "PT-br"} {
		sourceDir := filepath.Join(root, "docs", "core", locale)
		outputDir := filepath.Join(scratch, "generated", "library", "html", locale)
		err := filepath.WalkDir(sourceDir, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || filepath.Ext(path) != ".md" {
				return nil
			}
			wantSourceCount++
			rel, err := filepath.Rel(sourceDir, path)
			if err != nil {
				return err
			}
			htmlPath := filepath.Join(outputDir, strings.TrimSuffix(rel, ".md")+".html")
			if _, err := os.Stat(htmlPath); err != nil {
				t.Errorf("missing rendered HTML for %s: %v", rel, err)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if rendered != wantSourceCount {
		t.Fatalf("rendered %d files, want %d", rendered, wantSourceCount)
	}
}

func mustCopyDocsTree(t *testing.T, root, scratch string) {
	t.Helper()
	src := filepath.Join(root, "docs", "core")
	dst := filepath.Join(scratch, "docs", "core")
	err := filepath.WalkDir(src, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, body, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}
