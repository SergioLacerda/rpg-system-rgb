package publication

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDescriptor(t *testing.T) {
	descriptor := Descriptor()
	if descriptor.ID != "publication" || descriptor.Name == "" || descriptor.Description == "" {
		t.Fatalf("unexpected descriptor: %+v", descriptor)
	}
}

func TestCheckValidatesRequiredArtifacts(t *testing.T) {
	if err := Check(CheckOptions{}); err == nil {
		t.Fatal("expected missing check options to fail")
	}
	root := t.TempDir()
	library := filepath.Join(root, "library")
	public := filepath.Join(root, "downloads")
	options := CheckOptions{LibraryDir: library, PublicDir: public, Basename: "rgb", Version: "v1"}
	if err := Check(options); err == nil {
		t.Fatal("expected missing artifacts to fail")
	}
	for _, path := range []string{
		filepath.Join(library, "index.html"),
		filepath.Join(library, "en", "library", "index.html"),
		filepath.Join(library, "library", "index.html"),
		filepath.Join(public, "rgb-latest-en.pdf"),
		filepath.Join(public, "rgb-v1-en.pdf"),
		filepath.Join(public, "rgb-latest-pt-br.pdf"),
		filepath.Join(public, "rgb-v1-pt-br.pdf"),
	} {
		writeTestFile(t, path, "content")
	}
	if err := Check(options); err != nil {
		t.Fatal(err)
	}
	empty := filepath.Join(public, "rgb-v1-pt-br.pdf")
	writeTestFile(t, empty, "")
	if err := Check(options); err == nil {
		t.Fatal("expected empty artifact to fail")
	}
}

func TestPublishPDFsValidationAndFallbackSources(t *testing.T) {
	if err := PublishPDFs(PDFOptions{}); err == nil {
		t.Fatal("expected missing PDF options to fail")
	}
	root := t.TempDir()
	public := filepath.Join(root, "downloads")
	options := PDFOptions{PublicDir: public, Basename: "rgb", Version: "v1"}
	if err := PublishPDFs(options); err == nil {
		t.Fatal("expected missing fallback PDF source to fail")
	}
	versionedEN := filepath.Join(public, "rgb-v1-en.pdf")
	latestPT := filepath.Join(public, "rgb-latest-pt-br.pdf")
	writeTestFile(t, versionedEN, "%PDF-1.7\nEN\n")
	writeTestFile(t, latestPT, "%PDF-1.7\nPT\n")
	if err := PublishPDFs(options); err != nil {
		t.Fatal(err)
	}
	if !fileExists(filepath.Join(public, "rgb-latest-en.pdf")) || !fileExists(filepath.Join(public, "rgb-v1-pt-br.pdf")) {
		t.Fatal("fallback publication should create latest and versioned aliases")
	}

	bad := filepath.Join(root, "bad.pdf")
	writeTestFile(t, bad, "not a pdf")
	options.SourceEN = bad
	options.SourcePT = latestPT
	if err := PublishPDFs(options); err == nil {
		t.Fatal("expected non-PDF source to fail")
	}
}

func TestRequireFileRejectsDirectories(t *testing.T) {
	dir := t.TempDir()
	if err := requireFile(dir); err == nil {
		t.Fatal("expected directory artifact to fail")
	}
	if fileExists(dir) {
		t.Fatal("directories must not count as files")
	}
}

func TestCopyFileReadErrors(t *testing.T) {
	if err := copyFile(filepath.Join(t.TempDir(), "missing.pdf"), filepath.Join(t.TempDir(), "out.pdf")); err == nil {
		t.Fatal("expected missing source to fail")
	}
}

func TestCheckRejectsDirectoryArtifact(t *testing.T) {
	root := t.TempDir()
	library := filepath.Join(root, "library")
	public := filepath.Join(root, "downloads")
	options := CheckOptions{LibraryDir: library, PublicDir: public, Basename: "rgb", Version: "v1"}
	for _, path := range []string{
		filepath.Join(library, "index.html"),
		filepath.Join(library, "library", "index.html"),
		filepath.Join(public, "rgb-latest-en.pdf"),
		filepath.Join(public, "rgb-v1-en.pdf"),
		filepath.Join(public, "rgb-latest-pt-br.pdf"),
		filepath.Join(public, "rgb-v1-pt-br.pdf"),
	} {
		writeTestFile(t, path, "content")
	}
	if err := os.MkdirAll(filepath.Join(library, "en", "library", "index.html"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Check(options); err == nil {
		t.Fatal("expected directory artifact to fail")
	}
}
