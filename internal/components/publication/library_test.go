package publication

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildLibraryRendersBilingualCoreDocs(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "docs")
	out := filepath.Join(root, "public", "library")
	writeTestFile(t, filepath.Join(source, "core", "en", "README.md"), "# Core Overview\n\nSee [Combat](combat/attack.md).\n")
	writeTestFile(t, filepath.Join(source, "core", "en", "combat", "attack.md"), "# Attack\n\n- Roll R\n")
	writeTestFile(t, filepath.Join(source, "core", "PT-br", "README.md"), "# Visao Core\n\nTexto.\n")

	if err := BuildLibrary(LibraryOptions{SourceDir: source, OutDir: out}); err != nil {
		t.Fatalf("BuildLibrary returned error: %v", err)
	}

	index := readTestFile(t, filepath.Join(out, "index.html"))
	if !strings.Contains(index, "English Library") || !strings.Contains(index, "Biblioteca em Portuguese") {
		t.Fatalf("index did not include locale cards:\n%s", index)
	}
	page := readTestFile(t, filepath.Join(out, "core", "en", "index.html"))
	if !strings.Contains(page, "<h1 id=\"core-overview\">Core Overview</h1>") {
		t.Fatalf("expected rendered heading, got:\n%s", page)
	}
	if !strings.Contains(page, "href=\"combat/attack/\"") {
		t.Fatalf("expected rewritten markdown link, got:\n%s", page)
	}
}

func TestPublishPDFsPublishesLatestAndVersionedFiles(t *testing.T) {
	root := t.TempDir()
	public := filepath.Join(root, "downloads")
	en := filepath.Join(root, "en.pdf")
	pt := filepath.Join(root, "pt.pdf")
	writeTestFile(t, en, "%PDF-1.7\nEnglish\n")
	writeTestFile(t, pt, "%PDF-1.7\nPT-BR\n")

	err := PublishPDFs(PDFOptions{
		PublicDir: public,
		Basename:  "rgb-system-core-v2",
		Version:   "v0.2",
		SourceEN:  en,
		SourcePT:  pt,
	})
	if err != nil {
		t.Fatalf("PublishPDFs returned error: %v", err)
	}
	for _, name := range []string{
		"rgb-system-core-v2-latest-en.pdf",
		"rgb-system-core-v2-v0.2-en.pdf",
		"rgb-system-core-v2-latest-pt-br.pdf",
		"rgb-system-core-v2-v0.2-pt-br.pdf",
	} {
		if _, err := os.Stat(filepath.Join(public, name)); err != nil {
			t.Fatalf("expected %s: %v", name, err)
		}
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
