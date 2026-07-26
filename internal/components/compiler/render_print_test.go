package compiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderPrintPageOrdersUnitsFromManifest(t *testing.T) {
	root := repoRoot(t)
	html, err := RenderPrintPage(root, "en")
	if err != nil {
		t.Fatalf("RenderPrintPage failed: %v", err)
	}
	// generated/pdf/core-v2-rules.manifest.json lists core.resource.health
	// before core.combat.attack-margin (KF-05) — confirm the health
	// heading's Markdown source text ("4 + R + B") appears before the
	// attack-margin section's heading text in the rendered output.
	healthIndex := strings.Index(html, "Character Creation")
	marginIndex := strings.Index(html, "Attack And Defense")
	if healthIndex == -1 || marginIndex == -1 {
		t.Fatalf("expected both source titles present in print page")
	}
	if healthIndex > marginIndex {
		t.Fatalf("expected Character Creation section before Attack And Defense section, got reverse order")
	}
}

func TestRenderPrintPagePTBR(t *testing.T) {
	root := repoRoot(t)
	html, err := RenderPrintPage(root, "PT-br")
	if err != nil {
		t.Fatalf("RenderPrintPage failed: %v", err)
	}
	if !strings.Contains(html, `lang="pt-BR"`) {
		t.Fatalf("expected pt-BR lang attribute:\n%s", html)
	}
}

func TestRenderPrintTreeWritesBothLocales(t *testing.T) {
	root := repoRoot(t)
	scratch := t.TempDir()
	mustCopyDocsTree(t, root, scratch)
	mustCopyGenerated(t, root, scratch)

	if err := RenderPrintTree(scratch); err != nil {
		t.Fatalf("RenderPrintTree failed: %v", err)
	}

	for _, locale := range []string{"en", "PT-br"} {
		path := filepath.Join(scratch, "generated", "library", "print", "core-v2-rules-"+locale+".html")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing print page for %s: %v", locale, err)
		}
	}
}

func mustCopyGenerated(t *testing.T, root, scratch string) {
	t.Helper()
	src := filepath.Join(root, "generated", "pdf")
	dst := filepath.Join(scratch, "generated", "pdf")
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		body, err := os.ReadFile(filepath.Join(src, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dst, entry.Name()), body, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
