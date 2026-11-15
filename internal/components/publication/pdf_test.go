package publication

import (
	"os"
	"path/filepath"
	"testing"
)

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
