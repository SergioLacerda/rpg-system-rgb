package library

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildSiteLinksEveryGeneratedPageNoOrphans(t *testing.T) {
	scratch := t.TempDir()
	mustWriteFixturePages(t, scratch)

	if err := BuildSite(scratch); err != nil {
		t.Fatalf("BuildSite failed: %v", err)
	}

	for _, locale := range []string{"en", "PT-br"} {
		htmlDir := filepath.Join(scratch, "generated", "library", "html", locale)
		navBody, err := os.ReadFile(filepath.Join(htmlDir, "index.html"))
		if err != nil {
			t.Fatalf("missing nav page for %s: %v", locale, err)
		}
		nav := string(navBody)

		err = filepath.WalkDir(htmlDir, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || filepath.Ext(path) != ".html" {
				return nil
			}
			rel, err := filepath.Rel(htmlDir, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			if rel == "index.html" {
				return nil
			}
			if !strings.Contains(nav, `href="`+rel+`"`) {
				t.Errorf("nav page for %s missing link to %s", locale, rel)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func mustWriteFixturePages(t *testing.T, scratch string) {
	t.Helper()
	pages := map[string]string{
		"en/README.html":                       "<html><body>Overview</body></html>",
		"en/combat/attack_and_defense.html":    "<html><body>Attack</body></html>",
		"en/combat/movement.html":              "<html><body>Movement</body></html>",
		"en/core/character_creation.html":      "<html><body>Creation</body></html>",
		"PT-br/README.html":                    "<html><body>Visão geral</body></html>",
		"PT-br/combat/attack_and_defense.html": "<html><body>Ataque</body></html>",
		"PT-br/core/character_creation.html":   "<html><body>Criação</body></html>",
	}
	for rel, body := range pages {
		path := filepath.Join(scratch, "generated", "library", "html", filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
