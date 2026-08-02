package publication

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDescriptor(t *testing.T) {
	descriptor := Descriptor()
	if descriptor.ID != "publication" || descriptor.Name == "" || descriptor.Description == "" {
		t.Fatalf("unexpected descriptor: %+v", descriptor)
	}
}

func TestBuildLibraryValidationAndMarkdownEdges(t *testing.T) {
	if err := BuildLibrary(LibraryOptions{}); err == nil {
		t.Fatal("expected missing source/out options to fail")
	}
	root := t.TempDir()
	if err := BuildLibrary(LibraryOptions{SourceDir: filepath.Join(root, "missing"), OutDir: filepath.Join(root, "out")}); err == nil {
		t.Fatal("expected missing docs/core source to fail")
	}
	source := filepath.Join(root, "docs")
	out := filepath.Join(root, "public", "library")
	writeTestFile(t, filepath.Join(source, "core", "misc", "ignored.md"), "# Ignored\n")
	if err := BuildLibrary(LibraryOptions{SourceDir: source, OutDir: out}); err == nil {
		t.Fatal("expected no locale pages to fail")
	}

	writeTestFile(t, filepath.Join(source, "core", "en", "guide.md"), strings.Join([]string{
		"# Guide",
		"",
		"Text with **bold**, `code`, [external](https://example.com), [anchor](#local), and [doc](next.md#part).",
		"",
		"| Name | Value |",
		"| --- | --- |",
		"| R | Red |",
		"",
		"```",
		"<unsafe>",
		"```",
		"",
		"####### Too Deep",
	}, "\n"))
	writeTestFile(t, filepath.Join(source, "core", "PT-br", "sem-titulo.md"), "Sem titulo.\n")
	if err := BuildLibrary(LibraryOptions{SourceDir: source, OutDir: out}); err != nil {
		t.Fatal(err)
	}
	page := readTestFile(t, filepath.Join(out, "core", "en", "guide", "index.html"))
	for _, expected := range []string{
		"<strong>bold</strong>",
		"<code>code</code>",
		`href="https://example.com"`,
		`href="#local"`,
		`href="next/#part"`,
		"<td>Name</td>",
		"&lt;unsafe&gt;",
		"<h6",
	} {
		if !strings.Contains(page, expected) {
			t.Fatalf("rendered page missing %q:\n%s", expected, page)
		}
	}
	ptPage := readTestFile(t, filepath.Join(out, "core", "PT-br", "sem-titulo", "index.html"))
	if !strings.Contains(ptPage, "<title>sem-titulo</title>") {
		t.Fatalf("fallback title not rendered:\n%s", ptPage)
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
		filepath.Join(library, "core", "en", "index.html"),
		filepath.Join(library, "core", "PT-br", "index.html"),
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

func TestOutputHelpers(t *testing.T) {
	if outputURL("core/en/README.md") != "/core/en/" {
		t.Fatal("README URL should collapse to directory")
	}
	if outputFile("/out", "/") != filepath.Join("/out", "index.html") {
		t.Fatal("root output file should be index.html")
	}
	if relHref("/core/en/guide/", "/core/PT-br/") != "../../../core/PT-br/" {
		t.Fatalf("unexpected relative href: %s", relHref("/core/en/guide/", "/core/PT-br/"))
	}
	if headingLevel("plain") != 1 {
		t.Fatal("non-heading line should clamp to h1")
	}
	if slug("  RGB: Core V2!  ") != "rgb-core-v2" {
		t.Fatalf("unexpected slug: %s", slug("  RGB: Core V2!  "))
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
		filepath.Join(library, "core", "PT-br", "index.html"),
		filepath.Join(public, "rgb-latest-en.pdf"),
		filepath.Join(public, "rgb-v1-en.pdf"),
		filepath.Join(public, "rgb-latest-pt-br.pdf"),
		filepath.Join(public, "rgb-v1-pt-br.pdf"),
	} {
		writeTestFile(t, path, "content")
	}
	if err := os.MkdirAll(filepath.Join(library, "core", "en", "index.html"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Check(options); err == nil {
		t.Fatal("expected directory artifact to fail")
	}
}
